package controllers

import (
	"bytes"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"text/template"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func GetNotifications(c *fiber.Ctx) error {

	customerId := c.Query("customerId")

	notification := []model.TkNotification{}

	err := db.DB.Where("customer_id = ?", customerId).Find(&notification).Error
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(notification) == 0 {
		log.Println("Data not found")
		return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
			Message: "Data not found",
			Success: true,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseDataMultiple{
		Message: "Data notification has been loaded",
		Success: true,
		Data:    notification,
	})
}

func InsertTokenFCM(c *fiber.Ctx) error {

	token := new(model.TokenFcm)
	if err := c.BodyParser(token); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	tx := db.DB.Begin()

	if token.AppName == "" {
		token.AppName = "tokoku"
	}

	checkData := []model.TokenFcm{}
	err := tx.Where("app_name = ? AND user_id = ?", token.AppName, token.UserID).Find(&checkData).Error
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(checkData) == 0 {
		err := tx.Create(&token).Error
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
				Message: "Data Successfully Inserted",
				Success: true,
			})
		}
	} else {
		err := tx.Where("app_name = ? AND user_id = ?", token.AppName, token.UserID).Updates(token).Error
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
				Message: "Data Successfully Updated",
				Success: true,
			})
		}
	}
}

func SendNotificationFCM(c *fiber.Ctx) error {

	userIdsQuery := c.Query("userId")

	if userIdsQuery == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "userId query parameter is required",
		})
	}

	userIdInt, err := strconv.Atoi(userIdsQuery)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid userId value",
		})
	}

	// Convert comma-separated user IDs to a slice of integers
	// userIdsStr := strings.Split(userIdsQuery, ",")
	// var userIds []int
	// for _, idStr := range userIdsStr {
	// 	id, err := strconv.Atoi(idStr)
	// 	if err != nil {
	// 		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
	// 			"message": "Invalid userId value",
	// 		})
	// 	}
	// 	userIds = append(userIds, id)
	// }

	// fmt.Println(tokens)

	// var dataFcm map[interface{}]interface{}

	dataFcm := make(map[string]string)
	dataFcm["url"] = "md.transaction"
	dataFcm["dataId"] = "9999999999"
	dataFcm["popUp"] = "1"
	dataFcm["title"] = "test title"
	dataFcm["body"] = "test body"

	err = helpers.SendNotification("test title", "test body", userIdInt, dataFcm, c)

	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
		Message: "Success send notification",
		Success: true,
	})
}

func send(to, subject, title, body string) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", "no-reply@pt-bks.com")
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)

	// Add custom headers for the hidden data
	// for key, value := range data {
	// 	m.SetHeader("X-Custom-"+key, value)
	// }

	// Load and parse the HTML template
	tmpl, err := template.ParseFiles("views/email_template.html")
	if err != nil {
		return err
	}

	// Create a buffer to store the executed template
	var tpl bytes.Buffer
	templateData := struct {
		Subject string
		Title   string
		Body    string
	}{
		Subject: subject,
		Title:   title,
		Body:    body,
	}

	if err := tmpl.Execute(&tpl, templateData); err != nil {
		return err
	}

	// Set the body of the email as HTML
	m.SetBody("text/html", tpl.String())

	// Attach files
	// for _, attachment := range attachments {
	// 	content, err := base64.StdEncoding.DecodeString(attachment.Content)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	m.Attach(attachment.Filename, gomail.SetCopyFunc(func(w io.Writer) error {
	// 		_, err := w.Write(content)
	// 		return err
	// 	}))
	// }

	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	strconvPort, _ := strconv.Atoi(smtpPort)

	smtpUser := os.Getenv("SMTP_USERNAME")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	d := gomail.NewDialer(smtpHost, strconvPort, smtpUser, smtpPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmail(c *fiber.Ctx) error {
	type EmailRequest struct {
		To      string `json:"to"`
		Subject string `json:"subject"`
		Title   string `json:"title"`
		Body    string `json:"body"`
	}

	var req EmailRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println(err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Unable to parse request body",
		})
	}

	if err := send(req.To, req.Subject, req.Title, req.Body); err != nil {
		log.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to send email",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Email sent successfully",
	})
}
