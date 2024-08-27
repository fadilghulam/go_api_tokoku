package controllers

import (
	"encoding/json"
	"fmt"
	"go_api_tokoku/helpers"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"time"

	"github.com/gofiber/fiber/v2"
)

// func selectDetailByFlagId(flagId int) []map[string]interface{} {
// 	data, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT *
// 										FROM public.flag_detail
// 										WHERE flag_id= %v`, flagId))

// 	if err != nil {
// 		fmt.Println("Error:", err)
// 		return nil
// 	}

// 	if len(data) == 0 {
// 		return nil
// 	}
// 	return data
// }

// func queryTableView(input map[string]interface{}, valSourceColumnKey string, valSourceColumnKey2 map[string]interface{}) string {

//         qWhere := ""
//         if valSourceColumnKey == "GENERAL" {
//         } else {

// 			for sckKey, sckVal := range valSourceColumnKey2 {

//                 dataInputParams := nil
//                 if sckVal["param"] != "" {
//                     dataInputParams = input[sckVal["param"]];
//                 } else {
//                     if sckVal["dataType"] == "TIMESTAMP WITHOUT TIME ZONE"{
//                         dataInputParams = "2018-01-01"
//                     }
//                 }
//                 param = dataInputParams
//                 // $param = $dataInput[$sckVal->param];

//                 qWhereTemp = " AND " + sckVal["source"] + " IN (" + 99999 + ")"

//                 isArray = 1;
//                 if sckVal["isArray"] != "" {
//                     isArray = sckVal["isArray"]
//                 }

//                 isArraySource = 0;
//                 if sckVal["isArraySource"] != "" {
//                     isArraySource = sckVal["isArraySource"];
//                 }

//                 isArrayDestination = 0;
//                 if sckVal["isArrayDestination"] != "" {
//                     isArrayDestination = sckVal["isArrayDestination"];
//                 }

//                 if *int(sckVal["isArray"]) != nil {
//                     qWhereTemp = $this->getWhereDataType($sckVal->source, $param, $sckVal->dataType, $isArray, $isArraySource, $isArrayDestination);
//                 } else if (isset($sckVal->isDateRange)) {
//                     //percobaan
//                     qWhereTemp = " AND (" . $sckVal->source . "::DATE >= DATE('" . $param[0] . "') AND " . $sckVal->source . "::DATE <= DATE('" . $param[1] . "'))";
//                 } else {

//                     if (is_array($param)) {
//                         if (
//                             sizeof($param) > 0
//                         ) {
//                             $param = implode(',', $param);
//                             $qWhereTemp = " AND " . $sckVal->source . " IN (" . $param . ")";
//                         }
//                     } else {
//                         if ($param != '') {
//                             if ($sckVal->param == 'dateLatest') {
//                                 $qWhereTemp = " AND " . $sckVal->source . " > '" . $param . "'::" . $sckVal->dataType;
//                             } else {
//                                 $qWhereTemp = " AND " . $sckVal->source . " IN (" . $param . ")";
//                             }
//                         }
//                     }
//                 }

//                 $qWhere .= $qWhereTemp;
//             }
//         }

//         // print_r($qWhere);
//         return $qWhere;
//     }
// }

// func getData(flagDetail []map[string]interface{}, input map[string]interface{}) []map[string]interface{} {

// 	var qSelect, qWhere string

// 	for key, val := range flagDetail {
// 		if val["source_type"] == "TABLE" || val["source_type"] == "VIEW" || val["source_type"] == "FUNCTION" {
// 			qSelect = fmt.Sprintf("SELECT * FROM %v WHERE TRUE ", val["source_name"])

// 			switch val["source_type"] {
// 			case "TABLE":
// 			case "VIEW":
// 				qWhere = queryTableView(input, val["source_column_key"], val["source_column_key2"])
// 			}
// 		}
// 	}

// 		if ($val->source_type == "TABLE" || $val->source_type == "VIEW" || $val->source_type == "FUNCTION") {
// 			$qSelect = "SELECT * FROM {$val->source_name}
// 						WHERE TRUE ";

// 			$qWhere = "";

// 			switch ($val->source_type) {
// 				case 'TABLE':
// 				case 'VIEW':
// 					# code...
// 					$qWhere = $this->queryTableView($dataInput, $val->source_column_key, $val->source_column_key2);
// 					break;
// 				case 'FUNCTION':
// 					$qSelect = "SELECT * FROM {$val->source_name_data}(" . $this->queryFunction($dataInput, $val->source_column_key, $val->source_column_key2) . ")
// 						WHERE TRUE ";
// 					break;
// 				default:
// 					# code...
// 					break;
// 			}

// 			// print_r($qSelect . $qWhere);
// 			// $data[$val->source_alias] = $this->db->query(
// 			//     $qSelect . $qWhere

// 			// )->result_array();
// 			// print_r("SELECT JSON_AGG(data) AS data FROM (" . $qSelect . $qWhere . ") data");
// 			// exit();
// 			// print_r("SELECT JSON_AGG(data) AS data FROM (" . $qSelect . $qWhere . ") data");
// 			$tmpData = $this->db->query(
// 				"SELECT JSON_AGG(data) AS data FROM (" . $qSelect . $qWhere . ") data"
// 			)->row_object();
// 			$data[$val->source_alias] = json_decode($tmpData->data, TRUE);
// 		}
// 	}
// }

// func GetFlagToday(c *fiber.Ctx) error {
// 	dataSelectFlag, err := helpers.ExecuteQuery(`SELECT
// 												f.*
// 												FROM flag f
// 												JOIN flag_mapping fm
// 												ON fm.flag_id = f.id
// 												WHERE fm.flag_table = 'md.flag_merchandiser' AND f.is_today =1`)

// 	flagDetailTemp := []map[string]interface{}{}

// 	for _, val := range dataSelectFlag {
// 		// fmt.Println(val["id"])
// 		flagDetailTemp = append(flagDetailTemp, selectDetailByFlagId(int(val["id"].(float64)))...)
// 	}

// 	fmt.Println(flagDetailTemp)

// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(fiber.Map{
// 		"Message":    "Success",
// 		"Success":    true,
// 		"Data":       flagDetailTemp,
// 		"LengthData": len(flagDetailTemp),
// 	})
// }

func executeGORMQueryIndexString(query string, resultsChan chan<- map[string][]map[string]interface{}, index string, wg *sync.WaitGroup) {
	defer wg.Done()

	results, _ := helpers.ExecuteQuery(query)

	resultsChan <- map[string][]map[string]interface{}{index: results}
}

func GetFlagToday(c *fiber.Ctx) error {

	start := time.Now()
	type GetFlagRequest struct {
		MerchandiserId int    `json:"merchandiserId"`
		BranchId       int    `json:"branchId"`
		AreaId         int    `json:"areaId"`
		LoginAt        string `json:"loginAt"`
	}
	var flagReq GetFlagRequest
	if err := c.QueryParser(&flagReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	dataSend, err := json.Marshal(flagReq)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	responseData, err := helpers.SendCurl(dataSend, "GET", "https://rest.pt-bks.com/pluto-mobile/md/getDataTodayMD2")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	queries := []string{}
	keyQuery := []string{}
	// 	`SELECT * FROM public.get_customer_today_by_branch_merchandiser( 111) WHERE TRUE `,
	// 	`SELECT * FROM public.get_piutang_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_pengembalian_detail_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_stok_salesman_riwayat_today_by_merchandiser( 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_stok_merchandiser_riwayat_today_by_merchandiser( 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_stok_salesman_riwayat_today_by_merchandiser( 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_penjualan_detail_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_penjualan_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_pembayaran_piutang_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_kunjungan_log_customer_today_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_kunjungan_customer_today_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_pengembalian_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_pembayaran_piutang_detail_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_payment_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_transaction_detail_today_customer_by_merchandiser( 111, 38) WHERE TRUE `,
	// 	`SELECT * FROM public.get_transaction_today_customer_by_merchandiser( 111, 38) WHERE TRUE`,
	// }

	// fmt.Println(responseData)
	tempMerchandiserId := strconv.Itoa(flagReq.MerchandiserId)
	tempBranchId := strconv.Itoa(flagReq.BranchId)
	for key, val := range responseData["data"].(map[string]interface{}) {
		r := strings.NewReplacer("merchandiserId", tempMerchandiserId, "branchId", tempBranchId)
		tempString := r.Replace(val.(string))
		queries = append(queries, tempString)

		keyQuery = append(keyQuery, key)
	}

	var wg sync.WaitGroup
	resultsChan := make(chan map[string][]map[string]interface{}, len(queries))
	// tempResults := make([][]map[string]interface{}, len(queries))
	tempResults := make([]map[string]interface{}, 0, len(queries))

	for i, query := range queries {
		wg.Add(1)
		// go executeGORMQuery(query, resultsChan, i, &wg)
		go executeGORMQueryIndexString(query, resultsChan, keyQuery[i], &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for key, res := range result {
			tempResults = append(tempResults, map[string]interface{}{
				key: res,
			})
		}
	}
	finalResult := make(map[string]interface{})
	for _, val := range tempResults {
		for key, res := range val {
			// fmt.Println(key) //nama table
			// fmt.Println(res) //data each table
			finalResult[key] = res
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Function took %s", elapsed)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    finalResult,
		"keyData": keyQuery,
		"elapsed": elapsed,
	})
}
