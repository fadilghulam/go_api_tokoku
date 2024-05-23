package controllers

import (
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"
	"strconv"
	"strings"

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

	newQty, _ := strconv.Atoi(qty)

	var query = ""
	if newQty != 0 {
		query = fmt.Sprintf(
			`UPDATE tk.cart SET qty = %v, updated_at = now() WHERE id = %v`, newQty, cartId)
	} else {
		query = fmt.Sprintf(
			`DELETE FROM tk.cart WHERE id IN(%v)`, cartId)
	}

	result := db.DB.Exec(query)

	if result.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Update failed",
			Success: true,
		})
	} else {
		rowsAffected := result.RowsAffected
		if rowsAffected > 0 {
			return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
				Message: "Cart has been updated",
				Success: true,
			})
		} else {
			return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
				Message: "Update failed",
				Success: true,
			})
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
												)
									)
					) ORDER BY cart.id DESC
				) as items
	FROM tk.cart cart
	JOIN produk p
		ON cart.produk_id = p.id
	JOIN customer c
		ON c.id = %s
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
			Message: err.Error(),
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
			Message: result.Error.Error(),
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

func CheckoutCart(c *fiber.Ctx) error {
	cartIds := c.FormValue("cart_id")
	prov := c.FormValue("prov")
	kab := c.FormValue("kab")
	kec := c.FormValue("kec")
	kel := c.FormValue("kel")
	note := c.FormValue("note")
	sr_id := c.FormValue("sr_id")
	var querySr_id, valueSr_id string
	if sr_id != "" {
		querySr_id = ", sr_id"
		valueSr_id = "," + sr_id
	}
	var queryRayon_id, valueRayon_id string
	rayon_id := c.FormValue("rayon_id")
	if rayon_id != "" {
		queryRayon_id = ", rayon_id"
		valueRayon_id = "," + rayon_id
	}
	var queryBranch_id, valueBranch_id string
	branch_id := c.FormValue("branch_id")
	if branch_id != "" {
		queryBranch_id = ", branch_id"
		valueBranch_id = "," + branch_id
	}

	// tempCartIds := strings.Split(cartIds, ",")

	tx := db.DB.Begin()

	// var storeIds []int
	// var storeIdsStr string
	getStoreId := fmt.Sprintf(`SELECT store_id, MAX(id) as max_id, string_agg(id||'',',') as ids FROM tk.cart WHERE id IN (%v) GROUP BY store_id`, cartIds)

	// var results [][]interface{}
	result, err := helpers.ExecuteQuery(getStoreId)
	if err != nil {
		tx.Rollback()
		log.Fatal("failed to get store id: ", err.Error())
	}

	for i := 0; i < len(result); i++ {

		var transactionID int

		query := fmt.Sprintf(
			`INSERT INTO tk.transaction (transaction_state_id, customer_id, transaction_date, provinsi, kabupaten, kecamatan, kelurahan, store_id, note %s %s %s)
	         SELECT 1, customer_id, NOW(), '%s', '%s', '%s', '%s', %v, '%s' %v %v %v
	         FROM tk.cart WHERE id = %v
	         RETURNING id`, querySr_id, queryRayon_id, queryBranch_id, prov, kab, kec, kel, result[i]["store_id"], note, valueSr_id, valueRayon_id, valueBranch_id, result[i]["max_id"])

		// fmt.Println(query)

		firstInsert := tx.Raw(query).Scan(&transactionID)
		if firstInsert.Error != nil {
			tx.Rollback()
			log.Fatal("failed to insert transaction: ", firstInsert.Error)
		}

		tempCartIds := strings.Split(result[i]["ids"].(string), ",")
		// fmt.Println(tempCartIds)

		for j := 0; j < len(tempCartIds); j++ {
			anotherInsertQuery := fmt.Sprintf(`INSERT INTO tk.transaction_detail (transaction_id, produk_id, qty, harga, diskon, condition, pita, note)
								SELECT %v, produk_id, qty, harga, 0, null, null, null
								FROM tk.cart WHERE id = %v`, transactionID, tempCartIds[j])

			// fmt.Println(anotherInsertQuery)
			secondInsert := tx.Exec(anotherInsertQuery)
			if secondInsert.Error != nil {
				tx.Rollback()
				log.Fatal("failed to insert transaction: ", secondInsert.Error)
			}
		}
	}

	deleteQuery := fmt.Sprintf("DELETE FROM tk.cart WHERE id IN (%v)", cartIds)
	deleteProc := tx.Exec(deleteQuery)
	if deleteProc.Error != nil {
		tx.Rollback()
		log.Fatal("failed to delete cart: ", deleteProc.Error)
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Insert failed",
			Success: true,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Cart has been processed",
			Success: true,
		})
	}
}
