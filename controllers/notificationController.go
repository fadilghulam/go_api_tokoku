package controllers

import (
	"bytes"
	"crypto/rand"
	"fmt"
	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
	"time"

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

func send(to, subject, title, otp, body string) error {

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
	tmpl, err := template.ParseFiles("views/email_template2.html")
	if err != nil {
		return err
	}

	// Create a buffer to store the executed template
	var tpl bytes.Buffer
	templateData := struct {
		Subject  string
		Title    string
		Body     string
		Otp      string
		OtpSlice []string
	}{
		Subject:  subject,
		Title:    title,
		Body:     body,
		Otp:      otp,
		OtpSlice: strings.Split(otp, ""),
	}

	// fmt.Println(templateData)

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

func generateOTP(length int) string {
	otp := ""
	for i := 0; i < length; i++ {
		// Generate a random digit between 0 and 9
		num := make([]byte, 1)
		_, err := rand.Read(num)
		if err != nil {
			log.Fatal(err)
		}
		digit := strconv.Itoa(int(num[0] % 10))
		otp += digit
	}
	return otp
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
			"message": "Unable to parse request body",
			"success": false,
		})
	}

	checkUser, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT u.*, u.full_name as person_name, u.id::integer as new_user_id FROM public.user u
														JOIN hr.person p
															ON u.id = p.user_id
														WHERE TRUE AND p.email = '%s'
														ORDER BY dtm_crt DESC`, req.To))

	if err != nil {
		log.Println(err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get user info",
			"success": false,
		})
	}

	if len(checkUser) == 0 {
		log.Println("user not found " + req.To)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "User not found",
			"success": true,
		})
	} else {
		checkApp, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT * 
															FROM public.user u
															JOIN user_level ul
															ON ul.id = ANY(u.level_id)
															LEFT JOIN app a
															ON a.id = ANY(ul.access_app_id)
															WHERE (UPPER(a.name) = UPPER('%s') OR -1 = ANY(ul.access_app_id)) AND u.id = %v`, "TOKOKU", checkUser[0]["id"]))

		if err != nil {
			log.Println(err.Error())
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Failed to get app permission",
				"success": false,
			})
		}

		if len(checkApp) == 0 {
			log.Println("app permission not found" + req.To)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"message": "Anda tidak memiliki akses untuk aplikasi ini",
				"success": true,
			})
		} else {

			checkOtp, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT * FROM public.otp o
																WHERE UPPER(o.app_name) = UPPER('%s') AND o.user_id = %v AND UPPER(o.type) = '%s'
																AND expired_at >=NOW() AND confirmed_at IS NULL
																`, "TOKOKU", checkUser[0]["id"], "EMAIL"))

			if err != nil {
				log.Println(err.Error())
				return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
					"message": "Failed to get log otp",
					"success": false,
				})
			}

			if len(checkOtp) == 0 {
				otp := generateOTP(5)

				dataOtp := map[string]string{
					"otp": otp,
				}

				hashOTP := helpers.NewCurl(dataOtp, "POST", "https://rest.pt-bks.com/olympus/generateHashed", c)

				// fmt.Println("Otp: " + otp)

				if err := send(req.To, req.Subject, req.Title, otp, req.Body); err != nil {
					log.Println(err.Error())
					return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
						"message": "Failed to send email otp",
						"success": false,
					})
				} else {

					tableOtp := new(model.Otp)

					tableOtp.AppName = "TOKOKU"
					tableOtp.Type = "EMAIL"
					tableOtp.SendTo = req.To
					tableOtp.Otp = hashOTP["otpHash"].(string)
					tableOtp.UserID = int32(checkUser[0]["new_user_id"].(float64))
					tableOtp.CreatedAt = time.Now()
					tableOtp.UpdatedAt = time.Now()
					tableOtp.ExpiredAt = time.Now().Add(time.Minute * 3)

					if err := db.DB.Create(tableOtp).Error; err != nil {
						log.Println(err.Error())
						return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
							"message": "Failed to insert otp",
							"success": false,
						})
					}
				}

				// fmt.Println(otpHash)
			} else {
				return c.Status(http.StatusOK).JSON(fiber.Map{
					"message": "Otp already sent",
					"success": true,
				})
			}

		}
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Email sent successfully",
		"success": true,
	})
}
