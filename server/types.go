package server

import (
	"sync"
	"github.com/gorilla/websocket"
)

// Message struct defines the format of messages exchanged over WebSocket.
type Message struct {
	Type     string `json:"type"`     // Type of message (e.g., "private", "broadcast")
	Sender   string `json:"sender"`   // Sender's unique identifier (client ID)
	Receiver string `json:"receiver,omitempty"` // Receiver's unique identifier (only for private messages)
	Content  string `json:"content"`  // The content of the message
}

// WelcomeMessage is a specialized message sent when a client joins the server.
type WelcomeMessage struct {
	Type         string   `json:"type"`         // Type of message (e.g., "welcome")
	ClientID     string   `json:"client_id"`     // Unique ID of the client receiving the welcome message
	YourUsername string   `json:"your_username"` // Username of the client
	ClientList   []string `json:"client_list"`   // List of usernames of all connected clients
}

// Server struct defines the WebSocket server's state and operations.
type Server struct {
	Clients   map[string]*Client      // Map to store connected clients by their unique IDs
	Usernames map[string]bool         // Map to track used usernames (prevents duplicates)
	Upgrader  websocket.Upgrader      // WebSocket upgrader to upgrade HTTP connections to WebSocket
	ClientsMu sync.RWMutex            // Mutex to handle concurrent access to the Clients and Usernames maps
}
