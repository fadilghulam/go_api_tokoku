package helpers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	db "go_api_tokoku/config"
	model "go_api_tokoku/models"
	"io"
	"log"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
	"github.com/gofiber/fiber/v2"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/api/option"
)

func ExecuteQuery(query string) ([]map[string]interface{}, error) {

	queries := fmt.Sprintf(`SELECT JSON_AGG(data) as data FROM (%s) AS data`, query)

	rows, err := db.DB.Raw(queries).Rows()
	if err != nil {
		return nil, err
	}

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// datas, err := SaveRowToDynamicStruct(rows, columns)
	// if err != nil {
	// 	return nil, err
	// }

	results, err := JsonDecode(rows, columns)
	if err != nil {
		return nil, err
	}

	if results[0]["data"] == nil {
		return nil, nil
	}

	return results[0]["data"].([]map[string]interface{}), nil
}

func InsertDataDynamic(ctx context.Context, data map[string]interface{}) (map[string]interface{}, error) {

	tx := db.DB.WithContext(ctx).Begin()
	defer tx.Rollback()

	transactionData := data["transactions"].(map[string]interface{})

	dataInsert := make(map[string]interface{})
	dataDelete := make(map[string]interface{})

	insertedIds := make(map[string][]interface{})

	for tableName, data := range transactionData {

		if IsDeletedIds(tableName) {
			pattern := regexp.MustCompile(`_deleted_ids`)

			// Perform the replacement
			tableName = pattern.ReplaceAllString(tableName, "")

			dataDelete[tableName] = data
		} else {
			dataSlice := data.([]interface{})

			for i := 0; i < len(dataSlice); i++ {
				for key, value := range dataSlice[i].(map[string]interface{}) {
					if b, ok := value.(bool); ok {
						if b {
							dataSlice[i].(map[string]interface{})[key] = 1
						} else {
							dataSlice[i].(map[string]interface{})[key] = 0
						}
					}
				}
			}

			dataInsert[tableName] = dataSlice
		}
	}

	for tablenames, data := range dataInsert {

		tempStruct, err := CreateStructInstance(tablenames)
		if err != nil {
			return nil, err
		}

		tempTable := tablenames
		tempSchema := ""
		whereSchema := ""
		tempTableName := strings.Split(tablenames, ".")
		if len(tempTableName) > 1 {
			tempTable = tempTableName[1]
			tempSchema = tempTableName[0]
			whereSchema = fmt.Sprintf(" AND table_schema = '%s'", tempSchema)
		}

		query := fmt.Sprintf(`SELECT 1
								FROM information_schema.columns
								WHERE table_name = '%s' AND column_name = 'dtm_upd' %s
								ORDER BY ordinal_position`, tempTable, whereSchema)

		var count int64

		if err := db.DB.Raw(query).Count(&count).Error; err != nil {
			panic(err.Error())
		}

		var columnTime string

		// Check if the query returned any rows
		if count == 0 {
			columnTime = "updated_at"
		} else {
			columnTime = "dtm_upd"
		}

		// fmt.Println(columnTime)
		param := fmt.Sprintf("id , (%s - '5 day'::interval) as %s ", columnTime, columnTime)
		dataSlice := data.([]interface{})

		// fmt.Println(dataSlice)
		for i := 0; i < len(dataSlice); i++ {

			element := dataSlice[i].(map[string]interface{})
			// fmt.Println(element["id"])
			temp, err := db.DB.Raw(fmt.Sprintf("SELECT %s FROM %s WHERE id = %v", param, tablenames, element["id"])).Rows()
			if err != nil {
				panic(err.Error())
			}

			columns, err := temp.Columns()
			if err != nil {
				panic(err.Error())
			}

			defer temp.Close()
			checkCol, err := SaveRowToDynamicStruct(temp, columns)
			if err != nil {
				panic(err.Error())
			}

			if len(checkCol) > 0 {
				for _, m := range checkCol {

					var tempRow map[string]interface{}
					id := m["id"]
					coltime := m[columnTime]

					var timeVal time.Time

					if dataSlice[i].(map[string]interface{})[columnTime] == nil {
						dataSlice[i].(map[string]interface{})[columnTime] = time.Now()
						timeVal = time.Now()
					} else {
						timeStr := dataSlice[i].(map[string]interface{})[columnTime].(string)
						if dataSlice[i].(map[string]interface{})[columnTime] == nil {
							dataSlice[i].(map[string]interface{})[columnTime] = time.Now()
						}

						timeVal, err = time.Parse("2006-01-02 15:04:05", timeStr)
						if err != nil {
							fmt.Println("Error parsing time:", err)
							continue
						}
					}

					if coltime != nil {

						// fmt.Println(timeVal, coltime.(time.Time))

						if timeVal.After(coltime.(time.Time)) {
							tempRow = dataSlice[i].(map[string]interface{})
						}
					} else {
						dataSlice[i].(map[string]interface{})[columnTime] = time.Now()
						tempRow = dataSlice[i].(map[string]interface{})
					}
					if len(tempRow) > 0 {

						idData := dataSlice[i].(map[string]interface{})["id"]
						insertedIds[tablenames] = append(insertedIds[tablenames], idData)

						structValuePtr := reflect.ValueOf(tempStruct)
						if structValuePtr.Kind() != reflect.Ptr || structValuePtr.Elem().Kind() != reflect.Struct {
							return nil, fmt.Errorf("tempStruct is not a pointer to a struct")
						}

						// Dereference the pointer to get the struct value
						structValue := structValuePtr.Elem()

						dataMap := dataSlice[i].(map[string]interface{})
						err := mapstructure.Decode(dataMap, structValue.Addr().Interface())
						if err != nil {
							return nil, err
						}

						// Validate the tempStruct
						err = ValidateStructInstance(tempStruct)
						if err != nil {
							return nil, err
						}

						delete(dataSlice[i].(map[string]interface{}), "id")
						db.DB.Model(tempStruct).Where("id = ?", id).Updates(dataSlice[i].(map[string]interface{}))
					}
				}
			} else {

				structValuePtr := reflect.ValueOf(tempStruct)
				if structValuePtr.Kind() != reflect.Ptr || structValuePtr.Elem().Kind() != reflect.Struct {
					return nil, fmt.Errorf("tempStruct is not a pointer to a struct")
				}

				// Dereference the pointer to get the struct value
				structValue := structValuePtr.Elem()

				dataMap := dataSlice[i].(map[string]interface{})
				err := mapstructure.Decode(dataMap, structValue.Addr().Interface())
				if err != nil {
					return nil, err
				}

				// Validate the tempStruct
				err = ValidateStructInstance(tempStruct)
				if err != nil {
					return nil, err
				}

				db.DB.Model(tempStruct).Create(dataSlice[i].(map[string]interface{}))
				idData := dataSlice[i].(map[string]interface{})["id"]
				insertedIds[tablenames] = append(insertedIds[tablenames], idData)
			}
		}
	}

	returnedData := make(map[string]interface{})
	for tablenames, ids := range insertedIds {
		tempIds := Implode(ids)
		rows, err := db.DB.Raw("SELECT JSON_AGG(data) as data FROM (SELECT * FROM "+tablenames+" WHERE id IN (?) ) data", tempIds).Rows()
		if err != nil {
			return nil, err
		}

		columns, err := rows.Columns()
		if err != nil {
			return nil, err
		}

		defer rows.Close()

		datas, err := JsonDecode(rows, columns)
		if err != nil {
			return nil, err
		}
		returnedData[tablenames] = datas
	}

	for tempname, value := range returnedData {
		for _, v := range value.([]map[string]interface{}) {
			for _, vdata := range v {
				returnedData[tempname] = vdata
			}
		}
	}

	for tablenames, data := range dataDelete {
		tempStruct, err := CreateStructInstance(tablenames)
		if err != nil {
			return nil, err
		}
		db.DB.Where("id IN (?)", Implode(data.([]interface{}))).Delete(tempStruct)
	}

	if err := tx.Commit().Error; err != nil {
		return nil, fiber.ErrInternalServerError
	}

	return returnedData, nil
}

func JsonDecode(rows *sql.Rows, columns []string) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for rows.Next() {
		rowData := make(map[string]interface{})

		// Create a slice of interface{} to hold the values for Scan
		values := make([]interface{}, len(columns))
		for i := range columns {
			var value interface{}
			values[i] = &value
		}

		// Scan the row into the slice of interface{}
		if err := rows.Scan(values...); err != nil {
			log.Fatal(err)
		}

		// Transfer values from slice to map
		for i, col := range columns {
			rowData[col] = *values[i].(*interface{})
		}

		// Append the map to the result slice
		result = append(result, rowData)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	for i, m := range result {
		for key, value := range m {
			if bytes, isBytes := value.([]byte); isBytes {
				// fmt.Println(isBytes)
				var decodedData []map[string]interface{}
				if err := json.Unmarshal(bytes, &decodedData); err != nil {
					log.Fatal(err)
				} else {
					result[i][key] = decodedData
					// fmt.Println(decodedData)
				}
			}
		}
	}

	return result, nil
}

func JsonDecodeMap(input map[string]interface{}) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	// Create a map to hold the data
	rowData := make(map[string]interface{})

	// Iterate over the input map
	for key, value := range input {
		if bytes, isBytes := value.([]byte); isBytes {
			// Unmarshal bytes into a map
			var decodedData []map[string]interface{}
			if err := json.Unmarshal(bytes, &decodedData); err != nil {
				return nil, err
			}
			rowData[key] = decodedData
		} else {
			rowData[key] = value
		}
	}

	// Append the map to the result slice
	result = append(result, rowData)

	return result, nil
}

func JoinStrings(strings []string, separator string) string {
	result := ""
	for i, s := range strings {
		if i > 0 {
			result += separator
		}
		result += s
	}
	return result
}

func SplitToString(a []int, sep string) string {
	if len(a) == 0 {
		return ""
	}

	b := make([]string, len(a))
	for i, v := range a {
		b[i] = strconv.Itoa(v)
	}
	return strings.Join(b, sep)
}

func Implode(interfaceSlice []interface{}) []int64 {
	intSlice := make([]int64, len(interfaceSlice))
	for i, v := range interfaceSlice {

		if val, ok := v.(float64); ok {
			intSlice[i] = int64(val)
		} else if val, ok := v.(string); ok {
			temp, _ := strconv.Atoi(val)
			intSlice[i] = int64(temp)
		} else if val, ok := v.(int64); ok {
			intSlice[i] = int64(val)
		} else {
			fmt.Println("unknown type", v, reflect.TypeOf(v))
		}
	}
	return intSlice
}

func IsDeletedIds(s string) bool {
	pattern := `.*_deleted_ids`
	matched, _ := regexp.MatchString(pattern, s)
	return matched
}

func SaveRowToDynamicStruct(rows *sql.Rows, columns []string) ([]map[string]interface{}, error) {

	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range values {
		valuePtrs[i] = &values[i]
	}

	var results []map[string]interface{}

	// Iterate through the rows and store in the slice
	for rows.Next() {
		// Scan the values into the value pointers
		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create a map for the row
		rowMap := make(map[string]interface{})

		// Fill the map with column name and corresponding value
		for i, col := range columns {
			val := values[i]
			rowMap[col] = val
		}

		// Append the row map to the results slice
		results = append(results, rowMap)
	}

	return results, nil
}

func ByteResponse(responseBody []byte) (map[string]interface{}, error) {
	var responses map[string]interface{}
	if err := json.Unmarshal(responseBody, &responses); err != nil {
		// fmt.Println("Error:", err)
		return nil, err
	}

	responseData := make(map[string]interface{})
	for key, val := range responses {
		responseData[key] = val
	}

	return responseData, nil
}

func PostBody(body []byte) (map[string]interface{}, error) {
	// bodyBytes := c.Body()

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(body), &data); err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	return data, nil
}

func ConvertStringToInt64(i string) int64 {

	temp, _ := strconv.Atoi(i)
	temp2 := int64(temp)
	return temp2
}
func FloatToString(input_num float64) string {
	// to convert a float number to a string
	return strconv.FormatFloat(input_num, 'f', 6, 64)
}

func TrimLeftChar(s string) string {
	for i := range s {
		if i > 0 {
			return s[i:]
		}
	}
	return s[:0]
}

func ArrayColumn(data []any, column string) []any {

	var result []any
	for _, v := range data {
		result = append(result, v.(map[string]interface{})[column])
	}

	return result
}

func SendNotification(title string, body string, userIds int, customerId int64, dataSend map[string]string, c *fiber.Ctx) error {

	tokenFCM := new([]model.TokenFcm)

	// fmt.Println(userIds)

	err := db.DB.Where("user_id IN (?) AND app_name = ? AND customer_id = ?", userIds, "tokoku", customerId).Find(&tokenFCM).Error

	if err != nil {
		fmt.Println("error executing query ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	var tokens []string
	for _, token := range *tokenFCM {
		tokens = append(tokens, token.Token)
	}

	type NotificationRequest struct {
		Title  string            `json:"title"`
		Body   string            `json:"body"`
		Tokens []string          `json:"tokens"`
		Data   map[string]string `json:"data"` // additional data
	}

	var req NotificationRequest

	req.Title = title
	req.Body = body
	req.Tokens = tokens
	req.Data = dataSend

	// Load the service account key JSON file
	opt := option.WithCredentialsFile("middleware/tokoku.json")

	// Initialize the Firebase app
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		fmt.Printf("error initializing app: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize Firebase app",
		})
	}

	// Initialize the FCM client
	client, err := app.Messaging(ctx)
	if err != nil {
		fmt.Printf("error getting Messaging client: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to initialize FCM client",
		})
	}

	// Create the multicast message to send
	message := &messaging.MulticastMessage{
		Tokens: req.Tokens,
		Notification: &messaging.Notification{
			Title: req.Title,
			Body:  req.Body,
		},
		Data: req.Data,
	}

	// Send the message
	response, err := client.SendMulticast(ctx, message)
	if err != nil {
		fmt.Printf("error sending FCM message: %v\n", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send FCM notification",
		})
	}

	fmt.Println("Successfully sent FCM message to ", req.Tokens)
	return c.JSON(fiber.Map{
		"message": "Notification sent successfully",
		"success": response.SuccessCount,
		"failure": response.FailureCount,
		"errors":  response.Responses,
	})
}

func NewCurl(data map[string]string, method string, url string, c *fiber.Ctx) map[string]interface{} {

	client := &http.Client{}

	dataSend, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	// Create a POST request with a JSON payload
	req, err := http.NewRequest(method, url, bytes.NewReader(dataSend))
	if err != nil {
		fmt.Println("Error creating request:", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	responseData, err := ByteResponse(responseBody)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	return responseData
}

func RefreshUser(userId string) ([]map[string]interface{}, error) {

	datas, err := ExecuteQuery(fmt.Sprintf(`WITH data_customer AS (SELECT c.user_id, JSONB_AGG(c.*) as datas FROM tk.get_customer_by_userid(%v) c GROUP BY c.user_id)

											SELECT NULL as employee,
												u.id,
												u.full_name as name,
												u.username,
												u.profile_photo,
												p.email,
												p.phone,
												p.ktp,
												ARRAY[]::varchar[] as permission,
												dc.datas as "userInfo"
											FROM public.user u
											LEFT JOIN data_customer dc
												ON u.id = dc.user_id
											LEFT JOIN hr.person p
												ON u.id = p.user_id
											WHERE u.id = %v
											GROUP BY u.id, p.id, dc.datas`, userId, userId))

	if err != nil {
		return nil, err
	}

	return datas, nil
}

func SendCurl(data []byte, method string, url string) (map[string]interface{}, error) {

	client := &http.Client{}

	// req, err := http.NewRequest("GET", "https://rest.pt-bks.com/pluto-mobile/md/getDataTodayMD2", bytes.NewReader(dataSend))
	req, err := http.NewRequest(method, url, bytes.NewReader(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	responseData, err := ByteResponse(responseBody)
	if err != nil {
		fmt.Println("Error reading response:", err)
	}

	switch responseData["data"].(type) {
	case map[string]interface{}:
		if len(responseData["data"].(map[string]interface{})) == 0 {
			responseData["data"] = nil
		}
	case []interface{}:
		if len(responseData["data"].([]interface{})) == 0 {
			responseData["data"] = nil
		}
	default:
		responseData["data"] = nil

	}

	return responseData, nil
}
