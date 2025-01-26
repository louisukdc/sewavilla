package controller

import (
	"log"

	"github.com/gofiber/websocket/v2"
)

// Definisikan struct Message
type Message struct {
	Message string `json:"message"`
}

var clients = make(map[*websocket.Conn]bool) // Connected clients
var broadcast = make(chan string)            // Broadcast channel

// WebSocket handler
func WebSocketHandler(c *websocket.Conn) {
	clients[c] = true
	log.Println("New WebSocket connection established")

	defer func() {
		delete(clients, c)
		c.Close()
		log.Println("WebSocket connection closed")
	}()

	for {
		var msg Message // Menggunakan struct Message yang telah didefinisikan
		if err := c.ReadJSON(&msg); err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		// Menampilkan isi pesan yang diterima
		log.Printf("Received message: %s", msg.Message)

		// Kirim kembali respons jika perlu
		response := "Received: " + msg.Message
		if err := c.WriteJSON(response); err != nil {
			log.Printf("Error sending message: %v", err)
			break
		}
	}
}

// Goroutine to broadcast messages to all connected clients
func HandleMessages() {
	for {
		// Receive message from the broadcast channel
		message := <-broadcast

		// Send the message to all connected clients
		for client := range clients {
			if err := client.WriteJSON(message); err != nil {
				log.Printf("Error writing message: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
