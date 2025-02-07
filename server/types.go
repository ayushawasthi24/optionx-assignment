package server

import (
	"sync"
	"github.com/gorilla/websocket"
)

// Message struct defines the format of messages exchanged over WebSocket.
type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver,omitempty"`
	Content  string `json:"content"`
}

// WelcomeMessage is a specialized message sent when a client joins the server.
type WelcomeMessage struct {
	Type         string   `json:"type"`
	ClientID     string   `json:"client_id"`
	YourUsername string   `json:"your_username"`
	ClientList   []string `json:"client_list"`
}

// Server struct defines the WebSocket server's state and operations.
type Server struct {
	Clients   map[string]*Client
	Usernames map[string]bool
	Upgrader  websocket.Upgrader
	ClientsMu sync.RWMutex
}
