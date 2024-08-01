package controllers

import (
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"

	"github.com/gofiber/fiber/v2"
)

func GetSalesmanBy(c *fiber.Ctx) error {

	salesman := []model.Salesman{}

	err := db.DB.Where("id = ?", c.Query("salesmanId")).Find(&salesman).Error

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(salesman) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
			Message: "Salesman not found",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    salesman,
	})
}
