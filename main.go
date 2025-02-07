package main

import (
	"log"
	"net/http"
	"optionx-assignment/server"
)

func main() {
	s := server.NewServer()

	// Handle WebSocket connections at the /ws endpoint
	http.HandleFunc("/ws", s.HandleConnections)

	// Log a message indicating the server has started
	log.Println("WebSockets server started on :8080")

	// Start the HTTP server on port 8080 and handle incoming requests
	http.ListenAndServe(":8080", nil)
}
