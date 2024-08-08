package controllers

import (
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"

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

func GetCallCenter(c *fiber.Ctx) error {

	type datas struct {
		Whatsapp  string `json:"whatsapp"`
		Instagram string `json:"instagram"`
		Email     string `json:"email"`
		Phone     string `json:"phone"`
	}

	// datas = map[string]string{

	data := datas{
		Whatsapp:  "6281359613831",
		Instagram: "@pt-bks.com",
		Email:     "armour.retail.family@pt-bks.com",
		Phone:     "087741135521",
	}

	return c.Status(fiber.StatusOK).JSON(helpers.Response{
		Message: "Success",
		Success: true,
		Data:    data,
	})
}

func GetComplaints(c *fiber.Ctx) error {

	complaints := []model.Complaints{}
	// err := db.DB.Find(&complaints).Error

	customerId := c.Query("customerId")
	if customerId == "" {
		log.Println("no customerId")
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "No complaints found",
			Success: true,
		})
	}
	// err := db.DB.Find(&complaints).Error
	err := db.DB.Where("customer_id = ?", customerId).Find(&complaints).Order("id ASC").Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(complaints) == 0 {
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
			Message: "No complaints found",
			Success: true,
			Data:    nil,
		})
	}

	for i := range complaints {
		if complaints[i].Other != nil && *complaints[i].Other == "" {
			complaints[i].Other = nil
		}
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Success",
		Success: true,
		Data:    complaints,
	})
}

func InsertComplaints(c *fiber.Ctx) error {

	complaints := new(model.Complaints)
	if err := c.BodyParser(complaints); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	err := db.DB.Create(&complaints).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}
	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Complaints has been added",
		Success: true,
	})
}
