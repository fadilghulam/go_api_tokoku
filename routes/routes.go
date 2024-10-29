package routes

import (
	"fmt"
	"go_api_tokoku/controllers"
	"go_api_tokoku/helpers"
	"os"
	"strconv"

	db "go_api_tokoku/config"

	"github.com/dgrijalva/jwt-go"
	// fsio "github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"

	// socketio "github.com/googollee/go-socket.io"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/zishang520/socket.io/v2/socket"
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

func testFunc(userId string) string {
	return userId + " exist"
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
	webSocketGroup.Get("/v2", func(c *fiber.Ctx) error {
		return c.SendFile("views/websockets2.html")
	})
	webSocketGroup.Get("/echo", websocket.New(controllers.EchoHandler))
	webSocketGroup.Get("/print", websocket.New(controllers.EchoHandler2))
	webSocketGroup.Get("/addition", websocket.New(controllers.TestHandler))
	webSocketGroup.Get("/simpleSocket", websocket.New(controllers.SimpleSocketHandler))
	webSocketGroup.Get("/testNewSocket", websocket.New(controllers.HandleWebSocket))
	// webSocketGroup.Get("/testNewSocket", websocket.New(controllers.HandleWebSocket))
	webSocketGroup.Get("/testNewSocket2", controllers.WebsocketHandler, controllers.WsSearchUserByID(db.DB))

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

	app.Get("/getFlagToday", controllers.GetFlagToday)
	app.Get("/generateTransaction", controllers.GenerateTransactionsUserId)

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

	//Points routes
	authGroup.Get("/getPointsRule", controllers.GetPointsRule)

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
	authGroup.Get("/getCountTransactionItem", controllers.GetCountTransactionsItem)

	//Advertisement routes
	authGroup.Get("/getAdvertisement", controllers.GetAdvertisement)

	//Reviews routes
	authGroup.Get("/getReview", controllers.GetReview)
	authGroup.Post("/insertReview", controllers.InsertReview)
	authGroup.Get("/getComplaints", controllers.GetComplaints)
	authGroup.Post("/insertComplaints", controllers.InsertComplaints)
	authGroup.Get("/callCenter", controllers.GetCallCenter)
	authGroup.Get("/getCountReviews", controllers.GetCountCustomerReviews)
	authGroup.Get("/getCustomerReviews", controllers.GetCustomerReviews)

	//Membership routes
	authGroup.Get("/getMembership", controllers.GetMembership)

	//Notification routes
	authGroup.Get("/getNotification", controllers.GetNotificationsV2)
	authGroup.Put("/saveNotification", controllers.SaveNotification)
	authGroup.Post("/setNotifSetting", controllers.SetNotificationSettings)
	authGroup.Get("/getNotifSetting", controllers.GetNotificationSettings)
	// authGroup.Post("/insertReview", controllers.InsertReview)

	// //reports
	// app.Get("/revenues", controllers.GetRevenues)
	// app.Get("/solds", controllers.GetSolds)
}
func SocketIoSetup(app *fiber.App) {
	// Socket.IO setup
	socketio := socket.NewServer(nil, nil)
	socketio.On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		fmt.Println("Client connected: ", client.Id())

		// for i := range clients {
		// 	client := clients[i].(*socket.Socket)
		// 	client.Emit("message", "Hello World!")
		// }

		// client.On("message", func(args ...interface{}) {
		// 	for i := range clients {
		// 		tempClient := clients[i].(*socket.Socket)
		// 		tempClient.Emit("message", "Hello World from client message from")
		// 	}
		// })
		// socketio.FetchSocket()
		// fmt.Println(client.FetchSockets())

		client.On("requestData", func(args ...interface{}) {
			// fmt.Println(args)
			if len(args) == 0 {
				client.Emit("returnData", "No user id")
			} else {

				var userIdClient string

				switch args[0].(type) {
				case string:
					userIdClient = args[0].(string)
				case float64:
					userIdClient = strconv.Itoa(int(args[0].(float64)))
				}

				client.Join(socket.Room("room " + userIdClient))
				// fmt.Println(client.Rooms())

				// returnMap := make(map[string]interface{})

				// result, err := helpers.RefreshUser(userIdClient)
				// if err != nil {
				// 	fmt.Println(err)
				// 	client.Emit("returnData", "Failed to get user data")
				// }

				// if len(result) > 0 {

				// 	tokenString, expired, _ := controllers.CreateJWT(userIdClient)
				// 	jwtMap := map[string]interface{}{
				// 		"expired": expired,
				// 		"token":   tokenString,
				// 	}

				// 	returnMap["auth"] = true
				// 	returnMap["data"] = result[0]
				// 	returnMap["jwt"] = jwtMap
				// } else {

				// 	returnMap["auth"] = false
				// 	returnMap["data"] = nil
				// 	returnMap["jwt"] = nil
				// }

				// // fmt.Println(jwtMap)

				// client.Emit("returnData", returnMap)
				// socketio.Send(returnMap)

				// fmt.Println(returnMap)
			}
		})

		client.On("checkData", func(args ...interface{}) {
			// fmt.Println(args)
			if len(args) == 0 {
				client.Emit("returnData", "No user id")
			} else {

				var userIdClient string

				switch args[0].(type) {
				case string:
					userIdClient = args[0].(string)
				case float64:
					userIdClient = strconv.Itoa(int(args[0].(float64)))
				}

				// client.Join(socket.Room("room " + userIdClient))
				// fmt.Println(client.Rooms())

				returnMap := make(map[string]interface{})

				result, err := helpers.RefreshUser(userIdClient)
				if err != nil {
					fmt.Println(err)
					client.Emit("returnData", "Failed to get user data")
				}

				if len(result) > 0 {

					tokenString, expired, _ := controllers.CreateJWT(userIdClient)
					jwtMap := map[string]interface{}{
						"expired": expired,
						"token":   tokenString,
					}

					returnMap["auth"] = true
					returnMap["data"] = result[0]
					returnMap["jwt"] = jwtMap
				} else {

					returnMap["auth"] = false
					returnMap["data"] = nil
					returnMap["jwt"] = nil
				}

				// fmt.Println(jwtMap)

				client.Emit("returnData", returnMap)
				// socketio.Send(returnMap)

				// fmt.Println(returnMap)
			}
		})

		client.On("receivedUserId", func(args ...interface{}) {
			var userIdClient string

			switch args[0].(type) {
			case string:
				userIdClient = args[0].(string)
			case float64:
				userIdClient = strconv.Itoa(int(args[0].(float64)))
			}

			returnMap := make(map[string]interface{})

			result, err := helpers.RefreshUser(userIdClient)
			if err != nil {
				fmt.Println(err)
				client.Emit("returnData", "Failed to get user data")
			}
			tokenString, expired, err := controllers.CreateJWT(userIdClient)
			jwtMap := map[string]interface{}{
				"expired": expired,
				"token":   tokenString,
			}
			returnMap["auth"] = true
			returnMap["data"] = result[0]
			returnMap["jwt"] = jwtMap

			// client.Emit("returnData", returnMap)
			socketio.To(socket.Room("room "+userIdClient)).Emit("returnData", returnMap)
		})

		client.On("disconnect", func(args ...interface{}) {
			// client.Disconnect(true)
			fmt.Println("Client disconnected: ", client.Id())
			client.Disconnect(true)
		})
	})

	socketio.On("disconnect", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		fmt.Println("Client disconnected: ", client.Id())
		client.Disconnect(true)
	})

	socketio.Of("/connectCustomer", nil).On("connection", func(clients ...interface{}) {
		client := clients[0].(*socket.Socket)
		client.On("userId", func(args ...interface{}) {
			userIdClient := args[0].(string)
			result, err := helpers.RefreshUser(userIdClient)
			if err != nil {
				fmt.Println(err)
				client.Emit("returnData", "Failed to get user data")
			}
			client.Emit("returnData", result)
		})
	})
	app.Get("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))
	app.Post("/socket.io", adaptor.HTTPHandler(socketio.ServeHandler(nil)))
}
