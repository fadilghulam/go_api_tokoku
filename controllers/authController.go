package controllers

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"go_api_tokoku/helpers"

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
