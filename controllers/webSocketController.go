package controllers

import (
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

// func ServeHTML(w http.ResponseWriter, r *http.Request) {
// 	http.ServeFile(w, r, "websockets.html")
// }
