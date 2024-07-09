package controllers

import (
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func InsertVoucher(c *fiber.Ctx) error {

	voucher := new(model.Voucher)

	err := c.BodyParser(voucher)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	err = tx.Create(&voucher).Error
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
			Message: "Voucher has been added",
			Success: true,
		})
	}
}

func GetAllVoucher(c *fiber.Ctx) error {
	customerId := c.Query("customerId")
	data, err := helpers.ExecuteQuery(
		fmt.Sprintf(`SELECT vc.id,
							vc.customer_id,
							v.id as voucher_id,
							v.icon,
							v.code,
							v.date_start,
							v.date_end,
							v.diskon,
							v.is_percentage,
							v.min_cost,
							v.max_diskon,
							v.note,
							vc.amount_left,
							-- JSONB_AGG(DISTINCT CASE WHEN -1 = ANY(v.produk_id) THEN null ELSE p.id END) as product_ids,
							CASE WHEN -1 = ANY(v.produk_id) THEN null ELSE JSONB_AGG(p.id) END as product_ids,
							JSONB_AGG(
								JSONB_BUILD_OBJECT(
									--'id', v.id,
									--'icon', v.icon,
									--'code', v.code,
									--'date_start', v.date_start,
									--'date_end', v.date_end,
									--'diskon', v.diskon,
									--'is_percentage', v.is_percentage,
									--'min_cost', v.min_cost,
									--'max_diskon', v.max_diskon,
									--'note', v.note,
									'id', p.id,
									'name', p.name,
									'code', p.code,
									'produk_satuan', JSONB_BUILD_OBJECT(
														'carton', ps.carton,
														'ball', ps.ball,
														'slof', ps.slof,
														'pack', ps.pack
													)
								) ORDER BY p.order 
							) as products
					FROM tk.voucher_customer vc
					JOIN tk.voucher v 
						ON vc.voucher_id = v.id 
					JOIN produk p
						ON CASE WHEN -1 = ANY(v.produk_id) THEN TRUE ELSE p.id = ANY(v.produk_id) END
						AND p.is_aktif = 1
					JOIN produk_satuan ps
						ON p.satuan_id = ps.id
					WHERE now() BETWEEN v.date_start AND COALESCE(v.date_end, now()) AND vc.customer_id = %s
					GROUP BY vc.id, v.id
					ORDER BY v.date_end`, customerId))

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "Data Has Been Loaded",
			Success: true,
			Data:    data,
		})
	}
}

func InsertVoucherCustomer(c *fiber.Ctx) error {
	voucher := new(model.Voucher)
	voucherCustomer := new(model.VoucherCustomer)

	err := c.BodyParser(voucherCustomer)
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	customerId := c.FormValue("customerId")
	code := c.FormValue("code")
	if code != "" {
		// fmt.Println(code)
		if err := tx.Where("code = ?", code).First(&voucher).Error; err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}

		if voucher.ID == 0 {
			tx.Rollback()
			return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
				Message: "Voucher Not Found",
				Success: false,
			})
		}

		customerId, err := strconv.ParseInt(customerId, 10, 64)
		if err != nil {
			tx.Rollback()
			log.Println(err.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
				Message: "Something went wrong",
				Success: false,
			})
		}
		voucherCustomer.CustomerId = customerId
		voucherCustomer.VoucherId = voucher.ID
		voucherCustomer.AmountLeft = voucher.Amount
	}

	tx.Where("customer_id = ? AND voucher_id = ?", voucherCustomer.CustomerId, voucherCustomer.VoucherId).Find(&voucherCustomer)
	if voucherCustomer.ID != 0 {
		tx.Rollback()
		return c.Status(fiber.StatusConflict).JSON(helpers.ResponseWithoutData{
			Message: "Voucher Already Added",
			Success: false,
		})
	}

	err = tx.Create(&voucherCustomer).Error
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	} else {
		voucherAdded, _ := helpers.ExecuteQuery(fmt.Sprintf("SELECT * FROM tk.voucher WHERE id = %v", voucherCustomer.VoucherId))
		tx.Commit()
		return c.Status(fiber.StatusOK).JSON(helpers.Response{
			Message: "Voucher has been added",
			Success: true,
			Data:    voucherAdded[0],
		})
	}
}
