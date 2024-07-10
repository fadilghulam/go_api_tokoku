package controllers

import (
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func GetItemsExchange(c *fiber.Ctx) error {

	customerId := c.Query("customerId")

	var query string
	if customerId == "" {
		query = fmt.Sprintf(`SELECT ie.id,
									ie.point,
									ie.date_start,
									ie.date_end,	
									ie.max_exchange,
									ie.about,
									ie.detail,
									ie.term_condition,	
									JSONB_BUILD_OBJECT(
										'id', i.id,
										'name', i.name,
										'description', i.description,
										'image', i.image
									) as item		
								FROM tk.item_exchange ie
								JOIN tk.item i
									ON ie.item_id = i.id	
								WHERE now() BETWEEN ie.date_start AND COALESCE(ie.date_end, now())		
		`)
	} else {
		query = fmt.Sprintf(`WITH exchanged as (SELECT exchange_id, coalesce(count(id),0) as counts 
										FROM tk.customer_point_history 
										WHERE customer_id = %v 
										AND exchange_id IS NOT NULL
										AND type = 'EXCHANGE'
										GROUP BY customer_id, exchange_id)
								SELECT ie.id,
									ie.point,
									ie.date_start,
									ie.date_end,
									ie.max_exchange,
									ie.about,
									ie.detail,
									ie.term_condition,
									COALESCE(ex.counts,0) as exchanged_count,
									JSONB_BUILD_OBJECT(
										'id', i.id,
										'name', i.name,
										'description', i.description,
										'image', i.image
									) as item
								FROM tk.item_exchange ie
								JOIN tk.item i
									ON ie.item_id = i.id
								LEFT JOIN exchanged ex
									ON ie.id = ex.exchange_id
								WHERE now() BETWEEN ie.date_start AND COALESCE(ie.date_end, now())`, customerId)
	}

	result, err := helpers.ExecuteQuery(query)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Item has been loaded",
		Success: true,
		Data:    result,
	})
}

func InsertCartItem(c *fiber.Ctx) error {

	cartItem := new(model.CartItem)

	err := c.BodyParser(cartItem)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	query := fmt.Sprintf(`INSERT INTO tk.cart_item (customer_id, item_exchange_id, point, qty) VALUES (%v, %v, %v, %v)
							ON CONFLICT (customer_id, item_exchange_id) 
							DO UPDATE SET qty = (cart_item.qty + excluded.qty), updated_at = now()`, cartItem.CustomerId, cartItem.ItemExchangeId, cartItem.Point, cartItem.Qty)

	err = tx.Exec(query).Error
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		tx.Commit()
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Cart item has been added",
			Success: true,
		})
	}
}

func GetCartItem(c *fiber.Ctx) error {
	customerId := c.Query("customerId")

	query := fmt.Sprintf(`SELECT ci.id,
								ci.customer_id,
								ci.item_exchange_id,
								ci.point,
								ci.qty,
								JSONB_BUILD_OBJECT(
									'id', i.id,
									'name', i.name,
									'description', i.description,
									'image', i.image
								) as item
							FROM tk.cart_item ci
							JOIN tk.item_exchange ie
								ON ci.item_exchange_id = ie.id
							JOIN tk.item i
								ON ie.item_id = i.id
							WHERE ci.customer_id = %v
							ORDER BY ci.updated_at DESC`, customerId)

	result, err := helpers.ExecuteQuery(query)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Cart item has been loaded",
		Success: true,
		Data:    result,
	})
}

func UpdateCartItem(c *fiber.Ctx) error {
	cartId := c.FormValue("cart_id")
	qty := c.FormValue("qty")

	newQty, err := strconv.Atoi(qty)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Quantity is not valid",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	var query = ""
	if newQty != 0 {
		// query = fmt.Sprintf(
		// 	`UPDATE tk.cart SET qty = %v, updated_at = now() WHERE id = %v`, newQty, cartId)

		query = fmt.Sprintf(`UPDATE tk.cart_item SET qty = %v, updated_at = now() WHERE id IN(%v)`, newQty, cartId)

	} else {
		query = fmt.Sprintf(
			`DELETE FROM tk.cart_item WHERE id IN(%v)`, cartId)
	}

	result := tx.Exec(query)

	if result.Error != nil {
		tx.Rollback()
		log.Println(result.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		rowsAffected := result.RowsAffected
		if rowsAffected > 0 {
			if err := tx.Commit().Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Something went wrong",
					Success: false,
				})
			} else {
				return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
					Message: "Cart item has been updated",
					Success: true,
				})
			}
		} else {
			return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
				Message: "Update failed",
				Success: true,
			})
		}
	}
}

func DeleteCartItem(c *fiber.Ctx) error {

	// bodyBytes := c.Body()

	// var data map[string]interface{}
	// if err := json.Unmarshal([]byte(bodyBytes), &data); err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
	// 		Message: err.Error(),
	// 		Success: false,
	// 	})
	// }

	// cartIds := data["id_cart"].(string)

	cartIds := c.FormValue("cart_id")

	var carts []model.CartItem

	var ids = []int{}

	for _, id := range strings.Split(cartIds, ",") {
		tempId, _ := strconv.Atoi(id)
		ids = append(ids, tempId)
	}

	result := db.DB.Delete(&carts, ids)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		rowsAffected := result.RowsAffected
		if rowsAffected > 0 {
			return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
				Message: "Data has been deleted",
				Success: true,
			})
		} else {
			return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
				Message: "No data has been deleted",
				Success: true,
			})
		}
	}
}

// func CheckOutCartItem(c *fiber.Ctx) error {

// 	cartIds := c.FormValue("cartItemId")
// 	customerId := c.FormValue("customerId")
// 	prov := c.FormValue("prov")
// 	kab := c.FormValue("kab")
// 	kec := c.FormValue("kec")
// 	kel := c.FormValue("kel")
// 	note := c.FormValue("note")
// 	sr_id := c.FormValue("sr_id")
// 	var querySr_id, valueSr_id string
// 	if sr_id != "" {
// 		querySr_id = ", sr_id"
// 		valueSr_id = "," + sr_id
// 	}
// 	var queryRayon_id, valueRayon_id string
// 	rayon_id := c.FormValue("rayon_id")
// 	if rayon_id != "" {
// 		queryRayon_id = ", rayon_id"
// 		valueRayon_id = "," + rayon_id
// 	}
// 	var queryBranch_id, valueBranch_id string
// 	branch_id := c.FormValue("branch_id")
// 	if branch_id != "" {
// 		queryBranch_id = ", branch_id"
// 		valueBranch_id = "," + branch_id
// 	}

// 	tx := db.DB.Begin()

// 	// var storeIds []int
// 	// var storeIdsStr string
// 	getCartItems := fmt.Sprintf(`SELECT ci.*, ie.item_id
// 								FROM tk.cart_item ci
// 								JOIN tk.item_exchange ie
// 									ON ci.item_exchange_id = ie.id
// 								WHERE ci.id IN (%v)`, cartIds)

// 	// var results [][]interface{}
// 	result, err := helpers.ExecuteQuery(getCartItems)
// 	if err != nil {
// 		tx.Rollback()
// 		log.Println("failed to get store id: ", err.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	}

// 	for i := 0; i < len(result); i++ {

// 		var transactionID int

// 		query := fmt.Sprintf(`INSERT INTO tk.transaction_item (transaction_state_id, customer_id, transaction_date, provinsi, kabupaten, kecamatan, kelurahan, note %s %s %s)
// 								VALUES (1, %v, NOW(), '%s', '%s', '%s', '%s', '%s' %v %v %v)
// 								RETURNING id||''`, querySr_id, queryRayon_id, queryBranch_id, customerId, prov, kab, kec, kel, note, valueSr_id, valueRayon_id, valueBranch_id)

// 		// fmt.Println(query)

// 		firstInsert := tx.Raw(query).Scan(&transactionID)
// 		if firstInsert.Error != nil {
// 			tx.Rollback()
// 			log.Println("failed to insert transaction: ", firstInsert.Error.Error())
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Something went wrong",
// 				Success: false,
// 			})
// 		}

// 		tempCartIds := strings.Split(result[i]["id"].(string), ",")
// 		// fmt.Println(tempCartIds)

// 		for j := 0; j < len(tempCartIds); j++ {
// 			anotherInsertQuery := fmt.Sprintf(`INSERT INTO tk.transaction_item_detail (transaction_item_id, item_id, qty, note, point)
// 								SELECT %v, ie.item_id, ci.qty, null, ci.point
// 								FROM tk.cart_item ci
// 								JOIN tk.item_exchange ieON ci.item_exchange_id = ie.id
// 								WHERE ci.id = %v`, transactionID, tempCartIds[j])

// 			// fmt.Println(anotherInsertQuery)
// 			secondInsert := tx.Exec(anotherInsertQuery)
// 			if secondInsert.Error != nil {
// 				tx.Rollback()
// 				log.Println("failed to insert transaction: ", secondInsert.Error.Error())
// 				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 					Message: "Something went wrong",
// 					Success: false,
// 				})
// 			}
// 		}
// 	}

// 	deleteQuery := fmt.Sprintf("DELETE FROM tk.cart_item WHERE id IN (%v)", cartIds)
// 	deleteProc := tx.Exec(deleteQuery)
// 	if deleteProc.Error != nil {
// 		tx.Rollback()
// 		log.Println("failed to delete cart: ", deleteProc.Error.Error())
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	}

// 	if err := tx.Commit().Error; err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 			Message: "Something went wrong",
// 			Success: false,
// 		})
// 	} else {
// 		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
// 			Message: "Checkout success",
// 			Success: true,
// 		})
// 	}
// }

func CheckOutCartItem(c *fiber.Ctx) error {
	cartIds := c.FormValue("cartItemId")
	customerId := c.FormValue("customerId")
	prov := c.FormValue("prov")
	kab := c.FormValue("kab")
	kec := c.FormValue("kec")
	kel := c.FormValue("kel")
	note := c.FormValue("note")
	sr_id := c.FormValue("sr_id")
	rayon_id := c.FormValue("rayon_id")
	branch_id := c.FormValue("branch_id")

	tx := db.DB.Begin()

	getCartItems := fmt.Sprintf(`SELECT ci.id||'' as id, ie.item_id 
								FROM tk.cart_item ci
								JOIN tk.item_exchange ie
									ON ci.item_exchange_id = ie.id
								WHERE ci.id IN (%v)`, cartIds)

	// var results [][]interface{}
	result, err := helpers.ExecuteQuery(getCartItems)
	if err != nil {
		tx.Rollback()
		log.Println("failed to get cart item id: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	// fmt.Println(result[0]["id"].(string))

	for i := 0; i < len(result); i++ {

		var transactionID int64

		transactionItem := model.TkTransactionItem{
			TransactionStateID: 1,
			CustomerID:         helpers.ConvertStringToInt64(customerId),
			TransactionDate:    time.Now(),
			Provinsi:           prov,
			Kabupaten:          kab,
			Kecamatan:          kec,
			Kelurahan:          kel,
			Note:               note,
		}

		if sr_id != "" {
			transactionItem.SrID = helpers.ConvertStringToInt64(sr_id)
		}
		if rayon_id != "" {
			transactionItem.RayonID = helpers.ConvertStringToInt64(rayon_id)
		}
		if branch_id != "" {
			transactionItem.BranchID = helpers.ConvertStringToInt64(branch_id)
		}

		if err := tx.Create(&transactionItem).Error; err != nil {
			tx.Rollback()
			log.Println("failed to insert transaction: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		transactionID = transactionItem.ID

		tempCartIds := strings.Split(result[i]["id"].(string), ",")

		for j := 0; j < len(tempCartIds); j++ {
			anotherInsertQuery := fmt.Sprintf(`INSERT INTO tk.transaction_item_detail (transaction_item_id, item_id, qty, note, point)
 								SELECT %v, ie.item_id, ci.qty, null, ci.point
 								FROM tk.cart_item ci
 								JOIN tk.item_exchange ie ON ci.item_exchange_id = ie.id
 								WHERE ci.id = %v`, transactionID, tempCartIds[j])

			// fmt.Println(anotherInsertQuery)
			secondInsert := tx.Exec(anotherInsertQuery)
			if secondInsert.Error != nil {
				tx.Rollback()
				log.Println("failed to insert transaction: ", secondInsert.Error.Error())
				return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
					Message: "Something went wrong",
					Success: false,
				})
			}
		}

		customerExchangedPoints := fmt.Sprintf(`SELECT SUM(point) as point FROM tk.transaction_item_detail WHERE transaction_item_id = %v`, transactionID)

		point, err := helpers.ExecuteQuery(customerExchangedPoints)
		if err != nil {
			tx.Rollback()
			log.Println("failed to get customer point: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		customerPointHistory := model.CustomerPointHistory{
			CustomerId: helpers.ConvertStringToInt64(customerId),
			ExchangeId: transactionID,
			Point:      point[0]["point"].(int16),
			Type:       "EXCHANGE",
		}

		if err := tx.Create(&customerPointHistory).Error; err != nil {
			tx.Rollback()
			log.Println("failed to insert customer point history: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}
	}

	deleteQuery := fmt.Sprintf("DELETE FROM tk.cart_item WHERE id IN (%v)", cartIds)
	deleteProc := tx.Exec(deleteQuery)
	if deleteProc.Error != nil {
		tx.Rollback()
		log.Println("failed to delete cart: ", deleteProc.Error.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Cart has been processed",
			Success: true,
		})
	}
}

func GetTransactionsItem(c *fiber.Ctx) error {

	customerId := c.Query("id")
	types := c.Query("type")

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`WITH data_details as (
			SELECT 
				trd.transaction_item_id,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'id', trd.id,
						'item_id', trd.item_id,
						'point', trd.point,
						'name', i.name,
						'description', i.description,
						'image', i.image,
						'qty', trd.qty,
						'point', trd.point,
						'note', trd.note
					) ORDER BY i.name
				) as details,
				JSONB_BUILD_OBJECT(
					'id', tr.reference_id,
					'name', s.name,
					'phone', s.phone,
					'type', tr.reference_name
				) as courier,
				JSONB_BUILD_OBJECT(
					'id', c.id,
					'name', c.name,
					'outlet_name', c.outlet_name,
					'phone', c.phone,
					'address', c.alamat,
					'kelurahan', c.kelurahan,
					'kabupaten', c.kabupaten,
					'kecamatan', c.kecamatan,
					'provinsi', c.provinsi
				) as customer,
				JSONB_BUILD_OBJECT(
					'no_transaksi', tr.id,
					'jumlah_pesanan', SUM(trd.qty),
					'jumlah_point', SUM(trd.point),
					'voucher_tokoku', 0
				) as invoice,
				 rev.rating,
				JSONB_BUILD_OBJECT(
					'rating', rev.rating,
					'description', rev.description,
					'photo', rev.photo
				) as review
			FROM tk.transaction_item_detail trd
			JOIN tk.transaction_item tr
				ON trd.transaction_item_id = tr.id
			JOIN tk.transaction_state ts
				ON tr.transaction_state_id = ts.id
			JOIN customer c
				ON tr.customer_id = c.id
			JOIN tk.item i
				ON trd.item_id = i.id
			LEFT JOIN salesman s
				ON (tr.reference_id = s.id
				AND tr.reference_name = 'SALESMAN')
			LEFT JOIN tk.review rev
				ON tr.id = rev.order_item_id
				AND rev.order_id IS NULL
			WHERE tr.customer_id = %s AND ts.name = UPPER('%s')
			GROUP BY trd.transaction_item_id, tr.id, c.id, s.id, rev.id
		)
		
		SELECT LOWER(name) as name,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'id', tr.id,
						'type', ts.name,
						'transaction_date', tr.transaction_date,
						'items', trd.details,
						'courier', trd.courier,
						'customer', trd.customer,
						'invoice', trd.invoice,
						'note', tr.note,
						'review', CASE WHEN trd.rating IS NULL THEN NULL ELSE trd.review END
					) ORDER BY tr.transaction_date DESC
				) FILTER (WHERE tr.id IS NOT NULL) as datas
		FROM tk.transaction_state ts
		LEFT JOIN tk.transaction_item tr
			ON ts.id = tr.transaction_state_id
			AND tr.customer_id = %s
		LEFT JOIN data_details trd
			ON tr.id = trd.transaction_item_id
		WHERE ts.name = UPPER('%s')
		GROUP BY ts.name`, customerId, types, customerId, types))

	if results == nil {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Data not found",
			Success: true,
			Data:    nil,
		})
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}

	// newResults := make(map[string]interface{})
	// for _, val := range results {
	// 	// fmt.Println(i, val)
	// 	newResults[val["name"].(string)] = val["datas"]
	// }

	var newResults interface{}
	for _, val := range results {
		// fmt.Println(i, val)
		newResults = val["datas"]
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    newResults,
	})
}
