package controllers

import (
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"

	"github.com/gofiber/fiber/v2"
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

func InsertReview(c *fiber.Ctx) error {

	review := new(model.Review)
	if err := c.BodyParser(review); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	err := db.DB.Create(&review).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Success",
		Success: true,
	})
}
