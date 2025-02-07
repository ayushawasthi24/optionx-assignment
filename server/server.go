package server

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// WebSocket upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// NewServer initializes a new WebSocket server.
func NewServer() *Server {
	return &Server{
		Clients:   make(map[string]*Client),
		Usernames: make(map[string]bool),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}
