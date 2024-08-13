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
	app.Post("/login2", controllers.Login2)
	app.Post("/sendOtp", controllers.SendOtp)
	app.Post("/loginLegacy", controllers.LoginOrigin)
	app.Get("/auth", controllers.Auth)
	app.Get("/oauth/callback", controllers.OAuthCallback)
	app.Post("/register", controllers.RegisterUser)
	app.Get("/sendNotif", controllers.SendNotificationFCM)
	app.Post("/sendEmail", controllers.SendEmail)
	app.Get("/loginOauth", controllers.LoginOauth)
	app.Post("/quickCheckout", controllers.QuickCheckout)

	// app.Post("/TestHash", controllers.TestHash)

	authGroup := app.Group("")
	authGroup.Use(AuthMiddleware)

	//Customer routes
	authGroup.Post("/insertTokenFCM", controllers.InsertTokenFCM)
	authGroup.Get("/getDataCustomer", controllers.RefreshUser)
	authGroup.Post("/uploadFileS3", controllers.DoUpload)
	authGroup.Post("/deleteFileS3", controllers.DoDelete)
	authGroup.Put("/updateProfile", controllers.UpdateProfile)

	//Salesman routes
	authGroup.Get("/getSalesman", controllers.GetSalesmanBy)

	//Products routes
	authGroup.Get("/produkTerkini", controllers.GetProdukTerkini)
	authGroup.Get("/getHargaProduk", controllers.GetHargaProduk)
	authGroup.Get("/produkDetail", controllers.GetProdukDetail)
	authGroup.Get("/TestRoute", controllers.TestRoute)

	//Transaction routes
	authGroup.Get("/transactions", controllers.GetTransactions)
	authGroup.Get("/getCountTransactions", controllers.GetCountTransactions)
	authGroup.Get("/points", controllers.GetPointsCustomer)
	authGroup.Get("/pointsHistory", controllers.GetPointsHistory)

	//Cart routes
	authGroup.Post("/cart", controllers.InsertCart)
	authGroup.Get("/cart", controllers.GetCart)
	authGroup.Put("/updateCart", controllers.UpdateCart)
	authGroup.Delete("/deleteCart", controllers.DeleteCart)
	authGroup.Post("/checkoutCart", controllers.CheckoutCart)
	// authGroup.Post("/quickCheckout", controllers.QuickCheckout)

	//Voucher routes
	authGroup.Post("/addVoucher", controllers.InsertVoucher)
	authGroup.Get("/getVoucher", controllers.GetAllVoucher)
	authGroup.Post("/addVoucherCustomer", controllers.InsertVoucherCustomer)

	//Item routes
	authGroup.Post("/insertExchange", controllers.InsertCartItem)
	authGroup.Get("/getCartItem", controllers.GetCartItem)
	authGroup.Put("/updateCartItem", controllers.UpdateCartItem)
	authGroup.Delete("/deleteCartItem", controllers.DeleteCartItem)
	authGroup.Get("/itemExchange", controllers.GetItemsExchange)
	authGroup.Post("/checkoutCartItem", controllers.CheckOutCartItem)

	//Transaction Item routes
	authGroup.Get("/getTransactionItem", controllers.GetTransactionsItem)
	authGroup.Get("/getCountTransactionsItem", controllers.GetCountTransactionsItem)

	//Advertisement routes
	authGroup.Get("/getAdvertisement", controllers.GetAdvertisement)

	//Reviews routes
	authGroup.Get("/getReview", controllers.GetReview)
	authGroup.Post("/insertReview", controllers.InsertReview)
	authGroup.Get("/getComplaints", controllers.GetComplaints)
	authGroup.Post("/insertComplaints", controllers.InsertComplaints)
	authGroup.Get("/callCenter", controllers.GetCallCenter)

	//Membership routes
	authGroup.Get("/getMembership", controllers.GetMembership)

	//Notification routes
	authGroup.Get("/getNotification", controllers.GetNotifications)
	// authGroup.Post("/insertReview", controllers.InsertReview)

	// //reports
	// app.Get("/revenues", controllers.GetRevenues)
	// app.Get("/solds", controllers.GetSolds)

}
