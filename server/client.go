package server

import (
	"github.com/gorilla/websocket"
	"github.com/google/uuid"
	"github.com/brianvoe/gofakeit/v6"
)

// Client struct represents a connected user.
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
	Active   chan struct{}
}

// NewClient initializes a new client with a WebSocket connection.
func NewClient(conn *websocket.Conn) *Client {
	return &Client{
		ID:       uuid.New().String(),
		Username: gofakeit.Username(), // Generate random username
		Conn:     conn,
	}
}
