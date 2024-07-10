package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	db "go_api_tokoku/config"
	"go_api_tokoku/helpers"
	model "go_api_tokoku/models"

	"github.com/dgrijalva/jwt-go"
	jwt2 "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gofiber/fiber/v2"
)

type InputLogin struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ResponseData struct {
	Success bool                   `json:"success"`
	Auth    bool                   `json:"auth"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
	Jwt     map[string]interface{} `json:"jwt"`
}

func createJWT(userID string) (string, error) {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	claims := jwt.MapClaims{
		"id":  userID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	secretKey := []byte(os.Getenv("JWTKEY"))

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	result, err := token.SignedString(secretKey)

	if err != nil {
		// return "", err
		fmt.Println("Error:", err.Error())
	}
	return result, nil
}

func LoginOrigin(c *fiber.Ctx) error {

	bodyBytes := c.Body()

	client := &http.Client{}

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(bodyBytes), &data); err != nil {
		fmt.Println("Error:", err)
		return err
	}

	if _, ok := data["password"]; ok {
		password := []byte(data["password"].(string))

		data["password"] = fmt.Sprintf("%x", md5.Sum(password))
	}

	data["appName"] = "TOKOKU"

	var user InputLogin
	err := json.Unmarshal(bodyBytes, &user)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}

	dataSend, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	// Create a POST request with a JSON payload
	req, err := http.NewRequest("POST", "https://rest.pt-bks.com/olympus/login", bytes.NewReader(dataSend))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// fmt.Println(resp.Body)

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// fmt.Println(string(responseBody))

	responseData, err := helpers.ByteResponse(responseBody)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	return c.Status(resp.StatusCode).JSON(responseData)
	// return c.Status(resp.StatusCode).JSON(responseData2)
}

func Login(c *fiber.Ctx) error {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Otp      string `json:"otp"`
		SendTo   string `json:"sendTo"`
	}
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	client := &http.Client{}

	// Data to send in the POST request
	data := map[string]interface{}{
		"username": loginReq.Username,
		"password": fmt.Sprintf("%x", md5.Sum([]byte(loginReq.Password))),
		"appName":  "TOKOKU",
	}

	if loginReq.Otp != "" {
		data["otp"] = loginReq.Otp
		data["sendTo"] = loginReq.SendTo
	}

	dataSend, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	// Create a POST request with a JSON payload
	req, err := http.NewRequest("POST", "https://rest.pt-bks.com/olympus/login", bytes.NewReader(dataSend))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	responseData, err := helpers.ByteResponse(responseBody)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// switch responseData["data"].(type) {
	// case map[string]interface{}:
	if len(responseData["data"].(map[string]interface{})) == 0 {
		responseData["data"] = nil
	}
	// case []interface{}:
	// 	if len(responseData["data"].([]interface{})) == 0 {
	// 		responseData["data"] = nil
	// 	}

	// }

	return c.Status(resp.StatusCode).JSON(responseData)
}

func SendOtp(c *fiber.Ctx) error {
	type OtpRequest struct {
		Phone string `json:"phone"`
	}
	var otpReq OtpRequest
	if err := c.BodyParser(&otpReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	client := &http.Client{}

	// Data to send in the POST request
	data := map[string]interface{}{
		"sendTo":  otpReq.Phone,
		"appName": "TOKOKU",
	}

	dataSend, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
	}

	// Create a POST request with a JSON payload
	req, err := http.NewRequest("POST", "https://rest.pt-bks.com/olympus/sendOtp", bytes.NewReader(dataSend))
	if err != nil {
		log.Fatal("Error creating request:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending request:", err)
		return c.SendStatus(http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	// Read the response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	responseData, err := helpers.ByteResponse(responseBody)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	// fmt.Println("responseData :", responseData)

	// if len(responseData["data"].(map[string]interface{})) == 0 {
	// 	responseData["data"] = nil
	// }

	if responseData["success"] == true {
		return c.Status(fiber.StatusOK).JSON(responseData)
	}

	return nil
}

func LoginGPT(c *fiber.Ctx) error {
	// Parse request body
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var loginReq LoginRequest
	if err := c.BodyParser(&loginReq); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	// Example validation (replace with your own logic)
	// if loginReq.Username != "example_user" || loginReq.Password != "example_password" {
	// 	return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid username or password"})
	// }

	// Data to send in the POST request
	data := map[string]interface{}{
		"username": loginReq.Username,
		"password": fmt.Sprintf("%x", md5.Sum([]byte(loginReq.Password))),
		"appName":  "TOKOKU",
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to marshal data"})
	}

	// URL to send the POST request
	url := "https://rest.pt-bks.com/olympus/login"

	// Send POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send HTTP request", "details": err.Error()})
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	responseData, err := helpers.ByteResponse(responseBody)
	if err != nil {
		log.Fatal("Error reading response:", err)
	}

	return c.Status(resp.StatusCode).JSON(responseData)
}

func Auth(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	oauthClientID := os.Getenv("OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")

	// fmt.Println("oauthClientID: ", oauthClientID)

	if oauthClientID == "" || oauthClientSecret == "" {
		log.Fatal("Missing OAuth credentials")
	}

	var oauth2Config = &oauth2.Config{
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,
		RedirectURL:  "http://localhost:4001/oauth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	url := oauth2Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url)
}

func OAuthCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	if code == "" {
		return c.Status(http.StatusBadRequest).SendString("Missing code")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	oauthClientID := os.Getenv("OAUTH_CLIENT_ID")
	oauthClientSecret := os.Getenv("OAUTH_CLIENT_SECRET")

	// fmt.Println("oauthClientID: ", oauthClientID)

	if oauthClientID == "" || oauthClientSecret == "" {
		log.Fatal("Missing OAuth credentials")
	}

	var oauth2Config = &oauth2.Config{
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,
		RedirectURL:  "http://localhost:4001/oauth/callback",
		Scopes:       []string{"profile", "email"},
		Endpoint:     google.Endpoint,
	}

	token, err := oauth2Config.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to exchange token")
	}

	client := oauth2Config.Client(c.Context(), token)
	userInfo, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to get user info")
	}

	// Parse user info
	var user map[string]interface{}
	if err := json.NewDecoder(userInfo.Body).Decode(&user); err != nil {
		return c.Status(http.StatusInternalServerError).SendString("Failed to parse user info")
	}

	// Generate JWT
	// fmt.Println(user)
	username := user["email"].(string)
	tokenString, err := createJWT(username)
	if err != nil {
		fmt.Println("Error:", err.Error())
		return c.Status(http.StatusInternalServerError).SendString("Could not generate token")
	}

	// return c.JSON(fiber.Map{"token": tokenString})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"user":    user,
		"token":   tokenString,
	})
}

type TemplateInputRegister struct {
	Nik          string `json:"nik"`
	FullName     string `json:"fullName"`
	UserName     string `json:"userName"`
	EmailAddress string `json:"emailAddress"`
	PhoneNumber  string `json:"phoneNumber"`
	Password     string `json:"password"`
}

type Claims struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	jwt2.RegisteredClaims
}

func RegisterUser(c *fiber.Ctx) error {

	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	inputRegister := new(TemplateInputRegister)

	if err := c.BodyParser(inputRegister); err != nil {
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	checkUser := new(model.User)
	user := new(model.User)

	tx := db.DB.Begin()

	err = tx.Table("public.user").Where("username = ?", inputRegister.UserName).Find(checkUser).Error
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	if len(checkUser.Username) != 0 {
		tx.Rollback()
		return c.Status(fiber.StatusBadRequest).JSON(helpers.ResponseWithoutData{
			Message: "Username already exists",
			Success: false,
		})
	}

	user.FullName = inputRegister.FullName
	user.Username = inputRegister.UserName

	password := []byte(inputRegister.Password)
	user.Password = fmt.Sprintf("%x", md5.Sum(password))
	user.LevelID = model.Int32Array{32}
	user.IsMultipleLogin = 0
	user.IsAktif = "y"
	user.IsVerified = 0

	err = tx.Table("public.user").Create(&user).Error
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	userID := user.ID

	var person model.HrPerson

	person.Ktp = inputRegister.Nik
	person.FullName = inputRegister.FullName
	person.Email = inputRegister.EmailAddress
	person.Phone = inputRegister.PhoneNumber
	person.UserID = userID
	person.AppId = 17

	err = tx.Table("hr.person").Create(&person).Error
	if err != nil {
		tx.Rollback()
		log.Println(err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	tx.Commit()

	datas, err := helpers.ExecuteQuery(fmt.Sprintf(`SELECT NULL as employee,
														u.id,
														u.full_name as name,
														ARRAY[]::varchar[] as permission,
														NULL as userinfo
													FROM public.user u
													WHERE u.id = %v
													GROUP BY u.id`, userID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	// Calculate expiration time
	addTime := 60 * 60 * 6
	addTime = addTime * 4 * 30

	// Create a new JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:       int64(user.ID),
		Username: user.Username,
		RegisteredClaims: jwt2.RegisteredClaims{
			IssuedAt:  jwt2.NewNumericDate(time.Now()),
			ExpiresAt: jwt2.NewNumericDate(time.Now().Add(time.Duration(addTime) * time.Second)),
		},
	})

	// Sign the token with a secret key
	// kunci := "your-secret-key"
	kunci := os.Getenv("JWTKEY")
	tokenString, err := token.SignedString([]byte(kunci))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(helpers.ResponseWithoutData{
			Message: "Something went wrong",
			Success: false,
		})
	}

	type jwtData struct {
		Token   string `json:"token"`
		Expired int    `json:"expired"`
	}

	returnJwt := jwtData{
		Token:   tokenString,
		Expired: addTime,
	}

	// return c.JSON(fiber.Map{"token": tokenString})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Success",
		"data":    datas[0],
		"jwt":     returnJwt,
	})

	// return c.Status(fiber.StatusOK).JSON(helpers.ResponseWithoutData{
	// 	Message: "Register Success",
	// 	Success: true,
	// })
}
