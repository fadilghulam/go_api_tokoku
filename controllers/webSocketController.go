package controllers

import (
	"encoding/json"
	"log"
	"time"

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
