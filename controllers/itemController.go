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

func GetItemsExchange(c *fiber.Ctx) error {

	customerId := c.Query("customerId")

	result, err := helpers.ExecuteQuery(fmt.Sprintf(`WITH exchanged as (SELECT exchange_id, coalesce(count(id),0) as counts 
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
													WHERE now() BETWEEN ie.date_start AND COALESCE(ie.date_end, now())`, customerId))
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
