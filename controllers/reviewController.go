package controllers

import (
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"
	"math"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func GetReview(c *fiber.Ctx) error {

	customerId := c.Query("customerId")

	reviews := []model.Review{}
	err := db.DB.Where("customer_id = ?", customerId).Find(&reviews).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(reviews) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "No reviews found",
			Success: true,
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Success",
		Success: true,
		Data:    reviews,
	})
}

// func InsertReview(c *fiber.Ctx) error {

// 	review := new(model.Review)
// 	if err := c.BodyParser(review); err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	}
// 	err := db.DB.Create(&review).Error
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
// 		Message: "Success",
// 		Success: true,
// 	})
// }

func InsertReview(c *fiber.Ctx) error {

	review := new(model.Review)
	if err := c.BodyParser(review); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	err := tx.Create(&review).Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	getInsertedReview := tx.First(&review).Where("id = ? AND salesman_id IS NULL", review.ID)
	if getInsertedReview.Error != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	pointsRule := []model.PointsRule{}

	err = tx.Find(&pointsRule, "type = ?", "review").Error
	if err != nil {
		tx.Rollback()
		fmt.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	var tempPoints int16
	if len(pointsRule) >= 0 {
		reviewValue := reflect.ValueOf(review).Elem()

		for _, rule := range pointsRule {
			field := reviewValue.FieldByName(rule.SubType)
			// fmt.Println(field)
			if field.IsValid() && field.String() != "" {
				tempPoints += rule.Point
			}
		}

		customerPointHistory := new(model.CustomerPointHistory)
		customerPointHistory.CustomerId = review.CustomerId
		customerPointHistory.Point = tempPoints
		customerPointHistory.TransactionId = review.OrderId
		customerPointHistory.Type = "REVIEW"
		customerPointHistory.DateTime = time.Now()

		err = tx.Create(&customerPointHistory).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		err = tx.Model(&model.CustomerPoint{}).Where("customer_id = ?", review.CustomerId).Update("point", gorm.Expr("point + ?", tempPoints)).Error
		if err != nil {
			tx.Rollback()
			fmt.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

	}

	tx.Commit()

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    tempPoints,
	})
}

func GetCallCenter(c *fiber.Ctx) error {

	type datas struct {
		Whatsapp  string `json:"whatsapp"`
		Instagram string `json:"instagram"`
		Email     string `json:"email"`
		Linkedin  string `json:"linkedin"`
		Phone     string `json:"phone"`
	}

	// datas = map[string]string{

	data := datas{
		Whatsapp:  "623415088783",
		Instagram: "bksinside",
		Email:     "info@pt-bks.com",
		Linkedin:  "bercakawansejati",
		Phone:     "623415088783",
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    data,
	})
}

func GetComplaints(c *fiber.Ctx) error {

	complaints := []model.Complaints{}

	customerId := c.Query("customerId")
	if customerId == "" {
		log.Println("no customerId")
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "No complaints found",
			Success: true,
		})
	}
	// err := db.DB.Find(&complaints).Error
	err := db.DB.Where("customer_id = ?", customerId).Find(&complaints).Order("id ASC").Error
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(complaints) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "No complaints found",
			Success: true,
			Data:    nil,
		})
	}

	for i := range complaints {
		if complaints[i].Other != nil && *complaints[i].Other == "" {
			complaints[i].Other = nil
		}
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Success",
		Success: true,
		Data:    complaints,
	})
}

func InsertComplaints(c *fiber.Ctx) error {

	complaints := new(model.Complaints)
	if err := c.BodyParser(complaints); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	err := db.DB.Create(&complaints).Error
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Complaints has been added",
		Success: true,
	})
}

func GetCountCustomerReviews(c *fiber.Ctx) error {
	customerId := c.Query("customerId")

	data, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT COUNT(sq.id) as count_data, 
													CASE WHEN sq.is_reviewed = 1 THEN 'Reviewed' ELSE 'NotReviewed' END as tag
									FROM (
									SELECT tr.id, 
											tr.total_transaction, 
											tr.transaction_date,
											JSONB_AGG(
												JSONB_BUILD_OBJECT(
													'id', trd.id,
													'produk_id', trd.produk_id,
													'point', trd.point,
													'code', p.code,
													'name', p.name,
													'image', p.foto,
													'qty', trd.qty,
													'unit', iu.name,
													'harga', trd.harga,
													'note', trd.note,
													'produk_satuan', JSONB_BUILD_OBJECT(
																		'carton', ps.carton,
																		'ball', ps.ball,
																		'slof', ps.slof,
																		'pack', ps.pack
																	)
												) ORDER BY p.order
											) as details,
											CASE WHEN rv.id IS NOT NULL THEN 1 ELSE 0 END as is_reviewed
									FROM tk.transaction tr
									JOIN tk.transaction_detail trd
										ON tr.id = trd.transaction_id
									JOIN produk p
										ON trd.produk_id = p.id
									LEFT JOIN produk_satuan ps
										ON p.satuan_id = ps.id
									LEFT JOIN tk.review rv
										ON tr.id = rv.order_id
									LEFT JOIN tk.unit_mapping um
										ON p.id = um.produk_id
									LEFT JOIN tk.item_unit_tk iu
										ON um.item_unit_id = iu.id
									WHERE tr.customer_id = %s 
									GROUP BY tr.id, rv.id

									UNION ALL

									SELECT tr.id, 
											0 as total_transaction, 
											tr.transaction_date,
											JSONB_AGG(
												JSONB_BUILD_OBJECT(
													'id', trd.id,
													'produk_id', trd.item_id,
													'point', trd.point,
													'code', '',
													'name', i.name,
													'image', i.image,
													'qty', trd.qty,
													'unit', iu.name,
													'harga', 0,
													'note', trd.note,
													'produk_satuan', JSONB_BUILD_OBJECT(
																		'buah', 1
																	)
												) ORDER BY i.name
											) as details,
											CASE WHEN rv.id IS NOT NULL THEN 1 ELSE 0 END as is_reviewed
									FROM tk.transaction_item tr
									JOIN tk.transaction_item_detail trd
										ON tr.id = trd.transaction_item_id
									JOIN tk.item i
										ON trd.item_id = i.id
									LEFT JOIN tk.review rv
										ON tr.id = rv.order_item_id
									LEFT JOIN tk.unit_mapping um
										ON i.id = um.item_id
									LEFT JOIN tk.item_unit_tk iu
										ON um.item_unit_id = iu.id
									WHERE tr.customer_id = %s 
									GROUP BY tr.id, rv.id
									) sq
									GROUP BY sq.is_reviewed`, customerId, customerId))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	var tempData []map[string]interface{}

	if len(data) == 0 {
		tempData = append(tempData, map[string]interface{}{
			"Reviewed":    0,
			"NotReviewed": 0,
		})
	} else {
		for _, v := range data {
			tempData = append(tempData, map[string]interface{}{
				v["tag"].(string): v["count_data"].(float64),
			})
		}

		flattenedResult := make(map[string]interface{})

		for _, v := range tempData {
			for key, value := range v {
				// fmt.Println(key, value)
				flattenedResult[key] = int(value.(float64))
			}
		}

		tempData = []map[string]interface{}{flattenedResult}
	}

	data = tempData

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    data[0],
	})
}

func GetCustomerReviews(c *fiber.Ctx) error {

	customerId := c.Query("customerId")
	mode := c.Query("mode")
	page := c.Query("page")
	pageSize := c.Query("pageSize")

	qWhere := "AND rv.id IS NULL"
	if strings.ToLower(mode) == "reviewed" {
		qWhere = " AND rv.id IS NOT NULL"
	}

	var qLimit, qPage string
	iPage, _ := strconv.Atoi(page)
	iPageSize, _ := strconv.Atoi(pageSize)

	if pageSize != "" {
		qLimit = " LIMIT " + pageSize
	} else {
		qLimit = " LIMIT 20"
		iPageSize = 20
	}

	if page == "" {
		iPage = 0
	} else {
		iPage = iPage - 1
	}

	tempQ := strconv.Itoa(iPage * iPageSize)
	qPage = " OFFSET " + tempQ

	query := fmt.Sprintf(`SELECT sq.* FROM (
		SELECT tr.id, 
						tr.total_transaction, 
						tr.transaction_date,
						'transaction' as table,
						DATE(tr.transaction_date),
						JSONB_BUILD_OBJECT(
							'id', tr.reference_id,
							'name', MAX(s.name),
							'phone', MAX(s.phone),
							'type', tr.reference_name
						) as courier,
						JSONB_AGG(
							JSONB_BUILD_OBJECT(
								'type', 'produk',
								'id', trd.id,
								'produk_id', trd.produk_id,
								'point', trd.point,
								'code', p.code,
								'name', p.name,
								'image', p.foto,
								'qty', trd.qty,
								'unit', iu.name,
								'harga', trd.harga,
								'note', trd.note,
								'produk_satuan', JSONB_BUILD_OBJECT(
													'carton', ps.carton,
													'ball', ps.ball,
													'slof', ps.slof,
													'pack', ps.pack
												)
							) ORDER BY p.order
						) as products
		FROM tk.transaction tr
		JOIN tk.transaction_detail trd
			ON tr.id = trd.transaction_id
		JOIN produk p
			ON trd.produk_id = p.id
		LEFT JOIN produk_satuan ps
			ON p.satuan_id = ps.id
		LEFT JOIN tk.review rv
			ON tr.id = rv.order_id
		LEFT JOIN tk.unit_mapping um
			ON p.id = um.produk_id
		LEFT JOIN tk.item_unit_tk iu
			ON um.item_unit_id = iu.id
		LEFT JOIN salesman s
			ON tr.reference_id = s.id
			AND tr.reference_name = 'SALESMAN'
		WHERE tr.customer_id = %s AND tr.transaction_state_id = 3
			%s
		GROUP BY tr.id

		UNION ALL

		SELECT tr.id, 
				0 as total_transaction, 
				tr.transaction_date,
				'transaction_item' as table,
				DATE(tr.transaction_date),
				JSONB_BUILD_OBJECT(
							'id', tr.reference_id,
							'name', MAX(s.name),
							'phone', MAX(s.phone),
							'type', tr.reference_name
						) as courier,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'type', 'item',
						'id', trd.id,
						'produk_id', trd.item_id,
						'point', trd.point,
						'code', '',
						'name', i.name,
						'image', i.image,
						'qty', trd.qty,
						'unit', iu.name,
						'harga', 0,
						'note', trd.note,
						'produk_satuan', JSONB_BUILD_OBJECT(
											'buah', 1
										)
					) ORDER BY i.name
				) as details
		FROM tk.transaction_item tr
		JOIN tk.transaction_item_detail trd
			ON tr.id = trd.transaction_item_id
		JOIN tk.item i
			ON trd.item_id = i.id
		LEFT JOIN tk.review rv
			ON tr.id = rv.order_item_id
		LEFT JOIN tk.unit_mapping um
			ON i.id = um.item_id
		LEFT JOIN tk.item_unit_tk iu
			ON um.item_unit_id = iu.id
		LEFT JOIN salesman s
			ON tr.reference_id = s.id
			AND tr.reference_name = 'SALESMAN'
		WHERE tr.customer_id = %s AND tr.transaction_state_id = 3
			%s
		GROUP BY tr.id
		) sq
		ORDER BY sq.transaction_date DESC`, customerId, qWhere, customerId, qWhere)

	// fmt.Println(query)

	var wg sync.WaitGroup
	resultsChan := make(chan map[int][]map[string]interface{}, 2)

	queries := []string{
		query,
		query + qPage + qLimit,
	}

	tempResults := make([][]map[string]interface{}, len(queries))

	// Launch concurrent Goroutines
	for i, query := range queries {
		wg.Add(1)
		go executeGORMQuery(query, resultsChan, i, &wg)
	}

	// Wait for all Goroutines to finish
	wg.Wait()
	close(resultsChan)

	for result := range resultsChan {
		for index, res := range result {
			tempResults[index] = res
		}
	}

	if len(tempResults) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	type newResponseDataMultiple struct {
		Message    string      `json:"message"`
		Success    bool        `json:"success"`
		Data       interface{} `json:"datas"`
		TotalPages int         `json:"total_pages"`
	}

	var tempTotalPages int
	if len(tempResults[0]) < iPageSize {
		tempTotalPages = 1
	} else {
		tempTotalPages = int(math.Ceil(float64(len(tempResults[0])) / float64(iPageSize)))
	}

	return c.Status(fiber.StatusOK).JSON(newResponseDataMultiple{
		Message:    "Data has been loaded successfully",
		Success:    true,
		Data:       tempResults[1],
		TotalPages: tempTotalPages,
	})
}
