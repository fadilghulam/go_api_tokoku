package controllers

import (
	"go_api_tokoku/helpers"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
	"gorm.io/gorm"
	// "github.com/gorilla/websocket"
)

func WebsocketHandler(c *fiber.Ctx) error {
	if websocket.IsWebSocketUpgrade(c) {
		return c.Next()
	}
	return fiber.ErrUpgradeRequired
}

func WsSearchUserByID(db *gorm.DB) fiber.Handler {
	return websocket.New(func(c *websocket.Conn) {
		defer c.Close()

		for {
			// Reading the message from the WebSocket client
			var userID string
			err := c.ReadJSON(&userID)
			if err != nil {
				log.Println("Error reading JSON:", err)
				return
			}

			// Query the user by ID
			// result := db.First(&user, userID)
			result, err := helpers.RefreshUser(userID)
			if err != nil {
				log.Println("User not found:", err)
				c.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}

			// Send user details back to the WebSocket client
			err = c.WriteJSON(result)
			if err != nil {
				log.Println("Error writing JSON:", err)
				return
			}
		}
	})
}
