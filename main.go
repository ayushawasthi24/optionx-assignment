package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	ID         string
	Connection *websocket.Conn
	LastActive time.Time
	Mu         sync.Mutex
}

// Server manages all client connections
type Server struct {
	Clients   map[string]*Client
	Upgrader  websocket.Upgrader
	ClientsMu sync.RWMutex
}

// NewServer initializes a new WebSocket server
func NewServer() *Server {
	return &Server{
		Clients: make(map[string]*Client),
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

// AddClient registers a new client connection
func (s *Server) AddClient(client *Client) {
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()
	s.Clients[client.ID] = client
}

// RemoveClient removes a client from the server
func (s *Server) RemoveClient(clientID string) {
	s.ClientsMu.Lock()
	defer s.ClientsMu.Unlock()
	delete(s.Clients, clientID)
}

// BroadcastWelcomeMessage sends welcome message to new client
func (s *Server) BroadcastWelcomeMessage(client *Client) error {
	s.ClientsMu.RLock()
	defer s.ClientsMu.RUnlock()

	existingClientIDs := make([]string, 0, len(s.Clients)-1)
	for id := range s.Clients {
		if id != client.ID {
			existingClientIDs = append(existingClientIDs, id)
		}
	}

	welcomeMsg := map[string]interface{}{
		"message":   "Welcome to the WebSocket server!",
		"client_id": client.ID,
		"clients":   existingClientIDs,
	}

	return client.Connection.WriteJSON(welcomeMsg)
}

// SendMessageToClient sends a message to a specific client
func (s *Server) SendMessageToClient(senderID, targetID, message string) error {
	s.ClientsMu.RLock()
	defer s.ClientsMu.RUnlock()

	targetClient, exists := s.Clients[targetID]
	if !exists {
		return fmt.Errorf("client %s not found", targetID)
	}

	msgPayload := map[string]string{
		"from":    senderID,
		"message": message,
	}

	return targetClient.Connection.WriteJSON(msgPayload)
}

// HandlePing handles ping/pong messages to check client connectivity
func (s *Server) HandlePing(client *Client, done chan struct{}) {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			client.Mu.Lock()
			err := client.Connection.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(10*time.Second))
			client.Mu.Unlock()

			if err != nil {
				log.Printf("Ping error for client %s: %v", client.ID, err)
				client.Connection.Close()
				s.RemoveClient(client.ID)
				return
			}
		case <-done:
			return
		}
	}
}

// HandleClient manages individual client connection lifecycle
func (s *Server) HandleClient(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}

	clientID := uuid.New().String()
	client := &Client{
		ID:         clientID,
		Connection: conn,
		LastActive: time.Now(),
	}

	s.AddClient(client)
	log.Printf("New client connected: %s", clientID)

	// Send welcome message with existing clients
	err = s.BroadcastWelcomeMessage(client)
	if err != nil {
		log.Printf("Welcome message error: %v", err)
		conn.Close()
		return
	}

	done := make(chan struct{})
	go s.HandlePing(client, done)
	defer func() {
		done <- struct{}{}
		conn.Close()
		s.RemoveClient(clientID)
	}()

	for {
		var msg map[string]interface{}
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Printf("Read error from %s: %v", clientID, err)
			break
		}

		client.LastActive = time.Now()
		log.Printf("Received message from %s: %v", clientID, msg)

		// Handle targeted message
		if targetID, ok := msg["id"].(string); ok {
			message, _ := msg["message"].(string)
			err := s.SendMessageToClient(clientID, targetID, message)
			if err != nil {
				log.Printf("Send message error: %v", err)
			}
		}
	}
}

func main() {
	server := NewServer()

	http.HandleFunc("/ws", server.HandleClient)
	log.Println("WebSocket server starting on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
