package server

import (
	"github.com/gorilla/websocket"
)

// Client struct represents a connected user.
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn // WebSocket connection instance associated with the client
	Send     chan []byte     // Channel for sending messages to the client
	Active   chan struct{}   // Channel to track the client's activity
}
