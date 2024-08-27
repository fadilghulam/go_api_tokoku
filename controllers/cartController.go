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

func InsertCart(c *fiber.Ctx) error {

	// inputBody, err := helpers.PostBody(c.Body())
	// if err != nil {
	// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{Message: err.Error(), Success: false})
	// }

	// customerId := inputBody["customer_id"].(float64)
	// produkId := inputBody["produk_id"].(float64)
	// qty := inputBody["qty"].(float64)

	customerId := c.FormValue("customer_id")
	produkId := c.FormValue("produk_id")
	qty := c.FormValue("qty")
	storeId := c.FormValue("store_id")
	harga := c.FormValue("harga")

	result := db.DB.Exec(
		fmt.Sprintf(
			`INSERT INTO tk.cart (customer_id, produk_id, qty, date_cart, store_id, harga) 
				VALUES (%v, %v, %v, CURRENT_DATE, %v, %v)
				ON CONFLICT (customer_id, produk_id, date_cart, store_id) 
				DO UPDATE SET qty = (cart.qty + excluded.qty), updated_at = now()`, customerId, produkId, qty, storeId, harga))

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Insert failed",
			Success: true,
		})
	} else {
		rowsAffected := result.RowsAffected
		if rowsAffected > 0 {
			return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
				Message: "Cart has been added",
				Success: true,
			})
		} else {
			return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
				Message: "Insert failed",
				Success: true,
			})
		}
	}
}

func UpdateCart(c *fiber.Ctx) error {
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

	newQty64 := int64(newQty)
	var query = ""
	if newQty64 != 0 {
		var cart model.Cart

		result := tx.Model(&model.Cart{}).Where("id IN (?)", cartId).Updates(model.Cart{Qty: newQty64, UpdatedAt: time.Now()})
		if result.Error != nil {
			tx.Rollback()
			log.Println(result.Error.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		// Fetch store_id and customer_id from the updated record
		if err := tx.Model(&model.Cart{}).Select("store_id, customer_id").First(&cart, "id IN (?)", cartId).Error; err != nil {
			tx.Rollback()
			log.Println(result.Error.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		storeID := cart.StoreId
		customerID := cart.CustomerId

		// Execute the second update query within the new transaction
		query2 := fmt.Sprintf(`UPDATE tk.cart SET updated_at = now() WHERE store_id = %v AND customer_id = %v`, storeID, customerID)
		if err := tx.Exec(query2).Error; err != nil {
			tx.Rollback()
			log.Println(result.Error.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		// Handle store_id and customer_id

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
						Message: "Cart has been updated",
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

	} else {
		query = fmt.Sprintf(
			`DELETE FROM tk.cart WHERE id IN(%v)`, cartId)

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
						Message: "Cart has been updated",
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
}

// func InsertCart(c *fiber.Ctx) error {
// 	cart := new(model.Cart)
// 	// Store the body in the user and return error if encountered
// 	err := c.BodyParser(cart)
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Something's wrong with your input"})
// 	}

// 	err = db.DB.Create(&cart).Error
// 	if err != nil {
// 		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Could not add cart"})
// 	}
// 	// Return the created cart
// 	return c.Status(201).JSON(fiber.Map{"status": "success", "message": "Cart has been added"})
// }

func GetCart(c *fiber.Ctx) error {
	customerId := c.Query("customerId")

	results, err := helpers.ExecuteQuery(fmt.Sprintf(
		`SELECT cart.store_id,
				cart.customer_id,
				JSONB_AGG(
					JSONB_BUILD_OBJECT(
						'id', cart.id,
						'qty', cart.qty,
						'stock', 0,
						'product', JSONB_BUILD_OBJECT(
										'id', p.id,
										'code', p.code,
										'name', p.name,
										'discount', COALESCE(dis.nominal,0),
										'harga', cart.harga,
										'image', p.foto,
										'point', COALESCE(pt.value,0),
										'promo', JSONB_BUILD_OBJECT(
													'id', 1,
													'name', 'Promo 123'
												),
										'produk_satuan', JSONB_BUILD_OBJECT(
											'carton', ps.carton,
											'ball', ps.ball,
											'slof', ps.slof,
											'pack', ps.pack
										)
									)
					) ORDER BY cart.id DESC
				) as items
	FROM tk.cart cart
	JOIN produk p
		ON cart.produk_id = p.id
	JOIN produk_satuan ps
		ON p.satuan_id = ps.id
	JOIN customer c
		ON c.id = %s
		AND cart.customer_id = c.id
	JOIN salesman s
		ON c.salesman_id = s.id
	--LEFT JOIN ref_harga_master rhm
	--	ON c.branch_id = rhm.branch_id
	--AND CURRENT_DATE BETWEEN rhm.date_start AND COALESCE(rhm.date_end, CURRENT_DATE)
	--LEFT JOIN ref_harga rh
	--	ON rhm.id = rh.ref_harga_master_id
	--	AND c.tipe = rh.customer_tipe
	--	AND s.tipe_salesman = rh.salesman_tipe
	--	AND p.id = rh.produk_id
	LEFT JOIN tk.discount dis
		ON p.id = dis.produk_id
		AND CURRENT_DATE BETWEEN dis.date_start AND COALESCE(dis.date_end, CURRENT_DATE)
		AND c.branch_id = dis.branch_id
		AND c.tipe = dis.customer_type_id
	LEFT JOIN tk.points pt
		ON p.id = pt.produk_id
		AND CURRENT_DATE BETWEEN pt.date_start AND COALESCE(pt.date_end, CURRENT_DATE)
		AND c.branch_id = pt.branch_id
	GROUP BY cart.store_id, cart.customer_id
	ORDER BY MAX(cart.updated_at) DESC `, customerId))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: "Something went wrong",
			Success: false,
			Data:    nil,
		})
	}
	if len(results) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Cart not found",
			Success: true,
			Data:    nil,
		})
	}
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Cart has been loaded",
		Success: true,
		Data:    results,
	})
}

// func GetCart(c *fiber.Ctx) error {
// 	customerId := c.Query("customerId")
// 	var carts []model.Cart

// 	result := db.DB.Find(&carts, "customer_id = ?", customerId)

// 	if result.Error != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
// 			Message: result.Error.Error(),
// 			Success: false,
// 			Data:    nil,
// 		})
// 	}
// 	if result.RowsAffected == 0 {
// 		return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseDataMultiple{
// 			Message: "Cart not found",
// 			Success: true,
// 			Data:    nil,
// 		})
// 	}
// 	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
// 		Message: "Cart has been loaded",
// 		Success: true,
// 		Data:    carts,
// 	})
// }

func DeleteCart(c *fiber.Ctx) error {

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

	var carts []model.Cart

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

// func CheckoutCart(c *fiber.Ctx) error {
// 	cartIds := c.FormValue("cart_id")
// 	prov := c.FormValue("prov")
// 	customer_id := c.FormValue("customer_id")
// 	kab := c.FormValue("kab")
// 	kec := c.FormValue("kec")
// 	kel := c.FormValue("kel")
// 	note := c.FormValue("note")
// 	sr_id := c.FormValue("sr_id")
// 	voucher_id := c.FormValue("voucher_id")
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

// 	// tempCartIds := strings.Split(cartIds, ",")

// 	tx := db.DB.Begin()

// 	// var storeIds []int
// 	// var storeIdsStr string
// 	getStoreId := fmt.Sprintf(`SELECT store_id, MAX(id) as max_id, string_agg(id||'',',') as ids FROM tk.cart WHERE id IN (%v) GROUP BY store_id`, cartIds)

// 	// var results [][]interface{}
// 	result, err := helpers.ExecuteQuery(getStoreId)
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

// 		query := fmt.Sprintf(
// 			`INSERT INTO tk.transaction (transaction_state_id, customer_id, transaction_date, provinsi, kabupaten, kecamatan, kelurahan, store_id, note %s %s %s)
// 	         SELECT 1, %v, NOW(), '%s', '%s', '%s', '%s', %v, '%s' %v %v %v
// 	         FROM tk.cart WHERE id = %v
// 	         RETURNING id`, querySr_id, queryRayon_id, queryBranch_id, customer_id, prov, kab, kec, kel, result[i]["store_id"], note, valueSr_id, valueRayon_id, valueBranch_id, result[i]["max_id"])

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

// 		tempCartIds := strings.Split(result[i]["ids"].(string), ",")
// 		// fmt.Println(tempCartIds)

// 		for j := 0; j < len(tempCartIds); j++ {
// 			anotherInsertQuery := fmt.Sprintf(`INSERT INTO tk.transaction_detail (transaction_id, produk_id, qty, harga, diskon, condition, pita, note)
// 								SELECT %v, produk_id, qty, harga, 0, null, null, null
// 								FROM tk.cart WHERE id = %v`, transactionID, tempCartIds[j])

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

// 	if voucher_id != "" {
// 		querySubVoucher := fmt.Sprintf("UPDATE tk.voucher_customer SET amount_left = amount_left - 1 WHERE voucher_id IN (%s) AND customer_id = %v", voucher_id, customer_id)
// 		subtractVoucher := tx.Exec(querySubVoucher)
// 		if subtractVoucher.Error != nil {
// 			tx.Rollback()
// 			log.Println("failed to subtract voucher: ", subtractVoucher.Error.Error())
// 			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
// 				Message: "Something went wrong",
// 				Success: false,
// 			})
// 		}
// 	}

// 	deleteQuery := fmt.Sprintf("DELETE FROM tk.cart WHERE id IN (%v)", cartIds)
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
// 			Message: "Cart has been processed",
// 			Success: true,
// 		})
// 	}
// }

func CheckoutCart(c *fiber.Ctx) error {
	cartIds := c.FormValue("cart_id")
	prov := c.FormValue("prov")
	customer_id := c.FormValue("customer_id")
	kab := c.FormValue("kab")
	kec := c.FormValue("kec")
	kel := c.FormValue("kel")
	note := c.FormValue("note")
	sr_id := c.FormValue("sr_id")
	voucher_id := c.FormValue("voucher_id")
	rayon_id := c.FormValue("rayon_id")
	branch_id := c.FormValue("branch_id")

	tx := db.DB.Begin()

	getStoreId := fmt.Sprintf(`SELECT store_id||'' as store_id, MAX(id) as max_id, string_agg(id||'',',') as ids FROM tk.cart WHERE id IN (%v) GROUP BY store_id`, cartIds)

	result, err := helpers.ExecuteQuery(getStoreId)
	if err != nil {
		tx.Rollback()
		log.Println("failed to get store id: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	getEstimateDate, err := helpers.ExecuteQuery(fmt.Sprintf(`WITH get_estimated as (
																SELECT DISTINCT ON (sq.estimated_date) DATE(sq.estimated_date) as estimated_date--, CURRENT_DATE + '2 week'::interval
																FROM (
																SELECT r.id, r.day, r.mode, r.weekly, 
																CURRENT_DATE, 
																DATE_PART('week', CURRENT_DATE) AS woy, 
																DATE(date_trunc('week', CURRENT_DATE)) + (r.day-1||' day')::INTERVAL  AS estimated_date
																FROM rute_customer rc
																JOIN rute r
																	ON rc.rute_id = r.id
																WHERE rc.customer_id = %s AND r.is_aktif = 1 
																AND r.mode = (CASE WHEN MOD(DATE_PART('week', CURRENT_DATE)::INTEGER,2) = 0 THEN 'genap' ELSE 'ganjil' END)
																AND MOD(DATE_PART('week', CURRENT_DATE)::INTEGER,r.weekly) = 0
																ORDER BY r.day
																) sq
																ORDER BY sq.estimated_date
																LIMIT 1
																)

																SELECT DATE(estimated_date) as estimated_date FROM get_estimated
																UNION
																SELECT DATE(CURRENT_DATE + '2 week'::interval) as estimated_date`, customer_id))

	if err != nil {
		tx.Rollback()
		log.Println("failed to get estimate date, ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	for i := 0; i < len(result); i++ {

		var transactionID int

		transaction := model.TkTransaction{
			TransactionStateID: 1,
			CustomerID:         helpers.ConvertStringToInt64(customer_id),
			TransactionDate:    time.Now(),
			Provinsi:           prov,
			Kabupaten:          kab,
			Kecamatan:          kec,
			StoreID:            helpers.ConvertStringToInt64(result[i]["store_id"].(string)),
			Note:               note,
			Kelurahan:          kel,
			EstimateDate:       getEstimateDate[0]["estimated_date"].(time.Time),
			// StoreID:            result[i]["store_id"].(int64),
		}

		if sr_id != "" {
			transaction.SrID = helpers.ConvertStringToInt64(sr_id)
		}
		if rayon_id != "" {
			transaction.RayonID = helpers.ConvertStringToInt64(rayon_id)
		}
		if branch_id != "" {
			transaction.BranchID = helpers.ConvertStringToInt64(branch_id)
		}

		if err := tx.Create(&transaction).Error; err != nil {
			tx.Rollback()
			log.Println("failed to insert transaction: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		transactionID = transaction.ID

		tempCartIds := strings.Split(result[i]["ids"].(string), ",")

		for j := 0; j < len(tempCartIds); j++ {
			anotherInsertQuery := fmt.Sprintf(`INSERT INTO tk.transaction_detail (transaction_id, produk_id, qty, harga, diskon, condition, pita, note)
								SELECT %v, produk_id, qty, harga, 0, null, null, null
								FROM tk.cart WHERE id = %v`, transactionID, tempCartIds[j])

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
	}

	if voucher_id != "" {
		querySubVoucher := fmt.Sprintf("UPDATE tk.voucher_customer SET amount_left = amount_left - 1 WHERE voucher_id IN (%s) AND customer_id = %v", voucher_id, customer_id)
		subtractVoucher := tx.Exec(querySubVoucher)
		if subtractVoucher.Error != nil {
			tx.Rollback()
			log.Println("failed to subtract voucher: ", subtractVoucher.Error.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}
	}

	deleteQuery := fmt.Sprintf("DELETE FROM tk.cart WHERE id IN (%v)", cartIds)
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

func QuickCheckout(c *fiber.Ctx) error {

	type product struct {
		Quantity int64   `json:"quantity"`
		Point    int16   `json:"point"`
		ProdukId int32   `json:"produkId"`
		Harga    float64 `json:"harga"`
		Diskon   float64 `json:"diskon"`
	}

	type templateQuickCheckout struct {
		CustomerID int64     `json:"customerId"`
		Provinsi   string    `json:"prov"`
		Kabupaten  string    `json:"kab"`
		Kecamatan  string    `json:"kec"`
		Kelurahan  string    `json:"kel"`
		Note       string    `json:"note"`
		SrID       int64     `json:"srId"`
		RayonID    int64     `json:"rayonId"`
		BranchID   int64     `json:"branchId"`
		VoucherID  int64     `json:"voucherId"`
		TotalPrice int64     `json:"totalPrice"`
		Products   []product `json:"product"`
	}

	requestBody := new(templateQuickCheckout)

	err := c.BodyParser(requestBody)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something's wrong with your input",
			Success: false,
		})
	}

	// fmt.Println(requestBody.Products)

	// for i := 0; i < len(requestBody.Products); i++ {
	// 	fmt.Println(requestBody.Products[i].Point)
	// 	// transactionDetail := model.TkTransactionDetail{
	// 	// 	TransactionID: int64(transactionID),
	// 	// 	ProdukID:      int64(requestBody.Products[i].ProdukId),
	// 	// 	Qty:           requestBody.Products[i].Quantity,
	// 	// 	Harga:         requestBody.Products[i].Harga,
	// 	// 	Diskon:        requestBody.Products[i].Diskon,
	// 	// 	Point:         int64(requestBody.Products[i].Point),
	// 	// }

	// 	// if err := tx.Create(&transactionDetail).Error; err != nil {
	// 	// 	tx.Rollback()
	// 	// 	log.Println("failed to insert transaction: ", err.Error())
	// 	// 	return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
	// 	// 		Message: "Something went wrong",
	// 	// 		Success: false,
	// 	// 	})
	// 	// }
	// }

	tx := db.DB.Begin()

	var transactionID int

	transaction := model.TkTransaction{
		TransactionStateID: 1,
		CustomerID:         requestBody.CustomerID,
		TransactionDate:    time.Now(),
		Provinsi:           requestBody.Provinsi,
		Kabupaten:          requestBody.Kabupaten,
		Kecamatan:          requestBody.Kecamatan,
		Kelurahan:          requestBody.Kelurahan,
		Note:               requestBody.Note,
		TotalTransaction:   requestBody.TotalPrice,
	}

	if requestBody.SrID != 0 {
		transaction.SrID = requestBody.SrID
	}
	if requestBody.RayonID != 0 {
		transaction.RayonID = requestBody.RayonID
	}
	if requestBody.BranchID != 0 {
		transaction.BranchID = requestBody.BranchID
	}

	if err := tx.Create(&transaction).Error; err != nil {
		tx.Rollback()
		log.Println("failed to insert transaction: ", err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	transactionID = transaction.ID

	for i := 0; i < len(requestBody.Products); i++ {
		// fmt.Println(requestBody.Products[i].Quantity)
		transactionDetail := model.TkTransactionDetail{
			TransactionID: int64(transactionID),
			ProdukID:      int64(requestBody.Products[i].ProdukId),
			Qty:           requestBody.Products[i].Quantity,
			Harga:         requestBody.Products[i].Harga,
			Diskon:        requestBody.Products[i].Diskon,
			Point:         int64(requestBody.Products[i].Point),
		}

		if err := tx.Create(&transactionDetail).Error; err != nil {
			tx.Rollback()
			log.Println("failed to insert transaction: ", err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}
	}

	if requestBody.VoucherID != 0 {
		querySubVoucher := fmt.Sprintf("UPDATE tk.voucher_customer SET amount_left = amount_left - 1 WHERE voucher_id IN (%v) AND customer_id = %v", requestBody.VoucherID, requestBody.CustomerID)
		subtractVoucher := tx.Exec(querySubVoucher)
		if subtractVoucher.Error != nil {
			tx.Rollback()
			log.Println("failed to subtract voucher: ", subtractVoucher.Error.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Checkout success",
			Success: true,
		})
	}
	// return nil
}
