package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// NewServer initializes a new WebSocket server.
func NewServer() *Server {
	// Return a new Server instance with initialized fields
	return &Server{
		Clients:   make(map[string]*Client),
		Usernames: make(map[string]bool),
		Upgrader: websocket.Upgrader{
			// The Upgrader is responsible for upgrading HTTP connections to WebSocket
			// CheckOrigin allows requests from any origin (we can modify this for security)
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}
