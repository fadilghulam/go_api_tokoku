package controllers

import (
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"

	"github.com/gofiber/fiber/v2"
)

func GetPointsRule(c *fiber.Ctx) error {

	pointsRule := []model.PointsRule{}

	err := db.DB.Where("CURRENT_DATE BETWEEN date_start AND COALESCE(date_end, CURRENT_DATE)").Order("id asc").Find(&pointsRule).Error
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(pointsRule) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(helpers.ResponseWithoutData{
			Message: "Points rule not found",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Data has been loaded successfully",
		Success: true,
		Data:    pointsRule,
	})
}
