package controllers

import (
	"encoding/json"
	"fmt"
	"go_api_tokoku/helpers"
	"log"
	"time"

	"strconv"

	fsio "github.com/gofiber/contrib/socketio"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

// var upgrader = websocket.Upgrader{
// 	ReadBufferSize:  1024,
// 	WriteBufferSize: 1024,
// }

// func main() {
// 	http.HandleFunc("/echo", echoHandler)
// 	http.HandleFunc("/print", echoHandler2)
// 	http.HandleFunc("/", serveHTML)

// 	// go startWebSocketServer()

// 	fmt.Println("Server started at :8080")
// 	select {}
// }

// func startWebSocketServer() {
// 	log.Fatal(http.ListenAndServe(":8080", nil))
// }

type MessageObject struct {
	Data  string `json:"data"`
	From  string `json:"from"`
	Event string `json:"event"`
	To    string `json:"to"`
}

func EchoHandler(c *websocket.Conn) {
	defer c.Close()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Printf("Received message: %s\n", message)

		// Simulate some processing time
		time.Sleep(1 * time.Second)

		err = c.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

func EchoHandler2(c *websocket.Conn) {
	defer c.Close()

	for {
		messageType, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		log.Printf("Received message from client: %s\n", message)

		// Simulate some processing time
		time.Sleep(1 * time.Second)

		err = c.WriteMessage(messageType, message)
		if err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

// type Message struct {
// 	Num1   float64 `json:"num1"`
// 	Num2   float64 `json:"num2"`
// 	Result float64 `json:"result,omitempty"`
// }

func TestHandler(c *websocket.Conn) {
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		// var message Message
		// if err := json.Unmarshal(msg, &message); err != nil {
		// 	log.Println("Error unmarshalling JSON:", err)
		// 	continue
		// }

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		// Perform addition
		num1, ok1 := message["num1"].(float64)
		num2, ok2 := message["num2"].(float64)
		if !ok1 || !ok2 {
			log.Println("Invalid numbers")
			continue
		}

		// Perform addition
		result := num1 + num2
		// message.Result = message.Num1 + message.Num2

		message["result"] = result
		// Send result back to client
		response, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			continue
		}

		if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

func SimpleSocketHandler(c *websocket.Conn) {
	defer c.Close()
	for {
		messageType, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		if message["num1"] != nil || message["num2"] != nil {
			num1, _ := message["num1"].(float64)
			num2, _ := message["num2"].(float64)
			// if !ok1 || !ok2 {
			// 	// log.Println("Invalid numbers")
			// 	// continue
			// }

			if num1 != 0 && num2 == 0 {
				num2 = 0
			}

			if num1 == 0 && num2 != 0 {
				num1 = 0
			}

			// fmt.Println(num1, num2)

			result := num1 + num2
			message["result"] = result
			message["message"] = "You are performing addition with input: " + fmt.Sprint(strconv.ParseFloat(fmt.Sprint(num1), 64)) + "and " + fmt.Sprint(strconv.ParseFloat(fmt.Sprint(num2), 64))
		} else {
			err = c.WriteMessage(messageType, msg)
			if err != nil {
				log.Println("Error writing message:", err)
				return
			}
			// if message["str"] != nil {
			message["result"] = message["str"].(string)
			// } else {
			// 	message["result"] = message["num1"].(float64) + message["num2"].(float64)
			// }

			message["message"] = "Your data is sent, but no operation is performed"
		}

		response, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			continue
		}

		if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

func HandleWebSocket(c *websocket.Conn) {
	defer c.Close()
	for {
		_, msg, err := c.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			return
		}

		var message map[string]interface{}
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("Error unmarshalling JSON:", err)
			continue
		}

		fmt.Println(message)
		userID := message["userId"]
		var data []map[string]interface{}
		if userID != nil {
			// if err := db.DB.First(&user, userID).Error; err != nil {
			// 	if len(user.Username) == 0 {
			// 		log.Println("user not found")
			// 		c.WriteMessage(websocket.TextMessage, []byte("User not found"))
			// 		continue
			// 	}
			// 	log.Println("database query failed:", err)
			// 	break
			// }

			data, err = helpers.RefreshUser(userID.(string))
			if err != nil {
				log.Println("Error refreshing user:", err)
				break
			}

			if len(data) == 0 {
				log.Println("User not found")
				c.WriteMessage(websocket.TextMessage, []byte("User not found"))
				continue
			}
		}

		message["data"] = data
		message["message"] = "Data successfully loaded"

		response, err := json.Marshal(message)
		if err != nil {
			log.Println("Error marshalling JSON:", err)
			continue
		}

		if err := c.WriteMessage(websocket.TextMessage, response); err != nil {
			log.Println("Error writing message:", err)
			return
		}
	}
}

// func SetupSocketIOServer() *socketio.Server {
// 	server := socketio.NewServer(nil)

// 	server.OnConnect("/", func(so socketio.Conn) error {
// 		log.Println("Client connected:", so.ID())
// 		return nil
// 	})

// 	server.OnEvent("/", "message", func(so socketio.Conn, msg string) {
// 		var message map[string]interface{}
// 		if err := json.Unmarshal([]byte(msg), &message); err != nil {
// 			log.Println("Error unmarshalling JSON:", err)
// 			so.Emit("error", "Invalid message format")
// 			return
// 		}

// 		userID := message["userId"]
// 		var data []map[string]interface{}
// 		if userID != nil {
// 			data, err := helpers.RefreshUser(userID.(string))
// 			if err != nil {
// 				log.Println("Error refreshing user:", err)
// 				so.Emit("error", "Failed to refresh user data")
// 				return
// 			}

// 			if len(data) == 0 {
// 				log.Println("User not found")
// 				so.Emit("error", "User not found")
// 				return
// 			}
// 		}

// 		message["data"] = data
// 		message["message"] = "Data successfully loaded"

// 		response, err := json.Marshal(message)
// 		if err != nil {
// 			log.Println("Error marshalling JSON:", err)
// 			so.Emit("error", "Failed to create response")
// 			return
// 		}

// 		so.Emit("response", string(response))
// 	})

// 	server.OnDisconnect("/", func(so socketio.Conn, reason string) {
// 		log.Println("Client disconnected:", so.ID(), "Reason:", reason)
// 	})

// 	go server.Serve()
// 	return server
// }

func SocketIOHandler(c *fiber.Ctx) error {

	fsio.New(func(kws *fsio.Websocket) {

		// Retrieve the user id from endpoint
		userId := kws.Params("id")

		// Add the connection to the list of the connected clients
		// The UUID is generated randomly and is the key that allow
		// fsio to manage Emit/EmitTo/Broadcast
		// clients[userId] = kws.UUID

		// Every websocket connection has an optional session key => value storage
		kws.SetAttribute("user_id", userId)

		//Broadcast to all the connected users the newcomer
		kws.Broadcast([]byte(fmt.Sprintf("New user connected: %s and UUID: %s", userId, kws.UUID)), true, fsio.TextMessage)
		//Write welcome message
		kws.Emit([]byte(fmt.Sprintf("Hello user: %s with UUID: %s", userId, kws.UUID)), fsio.TextMessage)
	})

	fmt.Println("a")
	fmt.Println(fsio.EventConnect)
	fsio.On(fsio.EventConnect, func(ep *fsio.EventPayload) {
		fmt.Println("b")
		fmt.Printf("Connection event 1 - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	fsio.On(fsio.EventMessage, func(ep *fsio.EventPayload) {

		fmt.Println("c")
		fmt.Printf("Message event - User: %s - Message: %s", ep.Kws.GetStringAttribute("user_id"), string(ep.Data))

	})

	fsio.On(fsio.EventDisconnect, func(ep *fsio.EventPayload) {
		// Remove the user from the local clients
		// delete(clients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Disconnection event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	// On close event
	// This event is called when the server disconnects the user actively with .Close() method
	fsio.On(fsio.EventClose, func(ep *fsio.EventPayload) {
		// Remove the user from the local clients
		// delete(clients, ep.Kws.GetStringAttribute("user_id"))
		fmt.Printf("Close event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	// On error event
	fsio.On(fsio.EventError, func(ep *fsio.EventPayload) {
		fmt.Printf("Error event - User: %s", ep.Kws.GetStringAttribute("user_id"))
	})

	return nil
}
