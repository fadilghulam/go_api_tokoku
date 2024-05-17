package controllers

import (
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
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

	result := db.DB.Exec(
		fmt.Sprintf(
			`INSERT INTO tk.cart (customer_id, produk_id, qty, date_cart) 
				VALUES (%v, %v, %v, CURRENT_DATE)
				ON CONFLICT (customer_id, produk_id, date_cart) 
				DO UPDATE SET qty = (cart.qty + excluded.qty)`, customerId, produkId, qty))

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
			`DELETE FROM tk.cart WHERE id = %v`, cartId)
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
		`SELECT cart.customer_id,
					cart.id,
					cart.qty,
					0 as stock,
					JSONB_AGG(
						JSONB_BUILD_OBJECT(
							'id', p.id,
							'code', p.code,
							'name', p.name,
							'discount', COALESCE(dis.nominal,0),
							'harga', rh.harga,
							'image', p.foto,
							'point', COALESCE(pt.value,0),
							'promo', JSONB_BUILD_OBJECT(
										'id', 1,
										'name', 'Promo 123'
									)
						) ORDER BY cart.id
					) as product
			FROM tk.cart cart
			JOIN produk p
				ON cart.produk_id = p.id
			JOIN customer c
				ON c.id = %s
			JOIN salesman s
				ON c.salesman_id = s.id
			LEFT JOIN ref_harga_master rhm
				ON c.branch_id = rhm.branch_id
			AND CURRENT_DATE BETWEEN rhm.date_start AND COALESCE(rhm.date_end, CURRENT_DATE)
			LEFT JOIN ref_harga rh
				ON rhm.id = rh.ref_harga_master_id
				AND c.tipe = rh.customer_tipe
				AND s.tipe_salesman = rh.salesman_tipe
				AND p.id = rh.produk_id
			LEFT JOIN tk.discount dis
				ON p.id = dis.produk_id
				AND CURRENT_DATE BETWEEN dis.date_start AND COALESCE(dis.date_end, CURRENT_DATE)
				AND c.branch_id = dis.branch_id
				AND c.tipe = dis.customer_type_id
			LEFT JOIN tk.points pt
				ON p.id = pt.produk_id
				AND CURRENT_DATE BETWEEN pt.date_start AND COALESCE(pt.date_end, CURRENT_DATE)
				AND c.branch_id = pt.branch_id
			GROUP BY cart.id`, customerId))

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseDataMultiple{
			Message: err.Error(),
			Success: false,
			Data:    nil,
		})
	}
	if len(results) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseDataMultiple{
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
