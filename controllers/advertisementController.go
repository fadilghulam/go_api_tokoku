package controllers

import (
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"

	"github.com/gofiber/fiber/v2"
)

func GetAdvertisement(c *fiber.Ctx) error {

	advertisements := []model.Advertisement{}
	err := db.DB.Find(&advertisements).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Success",
		Success: true,
		Data:    advertisements,
	})
}
