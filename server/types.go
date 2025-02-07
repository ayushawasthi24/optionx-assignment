package server

import (
	"sync"

	"github.com/gorilla/websocket"
)

// Message struct defines the WebSocket message format.
type Message struct {
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver,omitempty"`
	Content  string `json:"content"`
}

type WelcomeMessage struct {
	Type         string   `json:"type"`
	ClientID     string   `json:"client_id"`
	YourUsername string   `json:"your_username"`
	ClientList   []string `json:"client_list"`
}

type Server struct {
	Clients   map[string]*Client
	Usernames map[string]bool
	Upgrader  websocket.Upgrader
	ClientsMu sync.RWMutex
}
