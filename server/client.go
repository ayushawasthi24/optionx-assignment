package server

import (
	"github.com/gorilla/websocket"
)

// Client struct represents a connected user.
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	Active   chan struct{}
}
