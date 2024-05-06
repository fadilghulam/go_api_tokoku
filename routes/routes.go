package routes

import (
	"fmt"
	"go_api_tokoku/controllers"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	tokenString := c.Get("Authorization")
	if tokenString == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Missing Authorization header",
		})
	}

	// Extract the token from the Authorization header
	secretKey := os.Getenv("JWTKEY")
	tokenString = tokenString[7:] // Remove "Bearer " prefix
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check the token signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// Provide the key to validate the token
		return []byte(secretKey), nil
	})
	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid or expired token",
		})
	}

	// Token is valid, set user information in context
	claims, _ := token.Claims.(jwt.MapClaims)
	c.Locals("user", claims["user"])

	return c.Next()
}

func Setup(app *fiber.App) {
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page!")
	})

	//authentication routes
	app.Post("/login", controllers.Login)
	app.Post("/loginLegacy", controllers.LoginOrigin)

	authGroup := app.Group("")
	authGroup.Use(AuthMiddleware)

	//Products routes
	authGroup.Get("/produkTerkini", controllers.GetProdukTerkini)
	authGroup.Get("/produkDetail", controllers.GetProdukDetail)
	authGroup.Get("/TestRoute", controllers.TestRoute)

	//Transaction routes
	authGroup.Get("/transactions", controllers.GetTransactions)
	authGroup.Get("/points", controllers.GetPointsCustomer)
	authGroup.Get("/pointsHistory", controllers.GetPointsHistory)
	// app.Get("/payments/:paymentId", controllers.GetPaymentDetails)
	// app.Post("/payments", controllers.CreatePayment)
	// app.Delete("/payments/:paymentId", controllers.DeletePayment)
	// app.Put("/payments/:paymentId", controllers.UpdatePayment)

	// //Order routes
	// app.Get("/orders", controllers.OrdersList)
	// app.Get("/orders/:orderId", controllers.OrderDetail)
	// app.Post("/orders", controllers.CreateOrder)
	// app.Post("/orders/subtotal", controllers.SubTotalOrder)
	// app.Get("/orders/:orderId/download", controllers.DownloadOrder)
	// app.Get("/orders/:orderId/check-download", controllers.CheckOrder)

	// //reports
	// app.Get("/revenues", controllers.GetRevenues)
	// app.Get("/solds", controllers.GetSolds)

}
