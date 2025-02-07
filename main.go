package main

import (
	"log"
	"net/http"
	"optionx-assignment/server"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not set
	}
	s := server.NewServer()

	// Handle WebSocket connections at the /ws endpoint
	http.HandleFunc("/ws", s.HandleConnections)

	// Log a message indicating the server has started
	log.Println("WebSockets server started on :8080")

	// Start the HTTP server on port 8080 and handle incoming requests
	http.ListenAndServe(":"+port, nil)
}
