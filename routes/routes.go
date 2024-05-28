package routes

import (
	"fmt"
	"go_api_tokoku/controllers"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
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

	testGroup := app.Group("go_api")
	testGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Landing Page 3!")
	})

	webSocketGroup := app.Group("ws")
	webSocketGroup.Get("/", func(c *fiber.Ctx) error {
		return c.SendFile("views/websockets.html")
	})
	webSocketGroup.Get("/echo", websocket.New(controllers.EchoHandler))
	webSocketGroup.Get("/print", websocket.New(controllers.EchoHandler2))
	webSocketGroup.Get("/addition", websocket.New(controllers.TestHandler))
	webSocketGroup.Get("/simpleSocket", websocket.New(controllers.SimpleSocketHandler))

	//authentication routes
	app.Post("/login", controllers.Login)
	app.Post("/sendOtp", controllers.SendOtp)
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

	//Cart routes
	authGroup.Post("/cart", controllers.InsertCart)
	authGroup.Get("/cart", controllers.GetCart)
	authGroup.Put("/updateCart", controllers.UpdateCart)
	authGroup.Delete("/deleteCart", controllers.DeleteCart)
	// authGroup.Post("/checkoutCart", controllers.CheckoutCart)
	authGroup.Post("/checkoutCart", controllers.CheckoutCart)

	//Voucher routes
	authGroup.Post("/addVoucher", controllers.InsertVoucher)
	authGroup.Get("/getVoucher", controllers.GetAllVoucher)
	authGroup.Post("/addVoucherCustomer", controllers.InsertVoucherCustomer)

	// //reports
	// app.Get("/revenues", controllers.GetRevenues)
	// app.Get("/solds", controllers.GetSolds)

}
