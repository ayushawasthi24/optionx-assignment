package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// BroadcastWelcomeMessage sends a welcome message to the newly connected client
func (s *Server) BroadcastWelcomeMessage(client *Client) {
	clientList := s.getClientList(client.ID)
	welcomeMsg := WelcomeMessage{
		Type:         "welcome",
		ClientID:     client.ID,
		YourUsername: client.Username,
		ClientList:   clientList,
	}
	msg, _ := json.Marshal(welcomeMsg)
	client.Send <- msg
}

// AddClient adds a new client to the server's client map
func (s *Server) AddClient(client *Client) {
	s.ClientsMu.Lock()
	s.Clients[client.ID] = client
	s.ClientsMu.Unlock()
}

// RemoveClient removes a client from the server and closes their connection
func (s *Server) RemoveClient(client *Client) {
	s.ClientsMu.Lock()
	delete(s.Clients, client.ID)
	delete(s.Usernames, client.Username)
	s.ClientsMu.Unlock()
	client.Conn.Close()
}

// HandleConnections handles incoming WebSocket connections and sets up client communication
func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP request to a WebSocket connection
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	clientID := uuid.New().String()
	username := s.generateUsername()
	client := &Client{
		ID:       clientID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte),
		Active:   make(chan struct{}, 1),
	}

	s.AddClient(client)
	go s.writePump(client)            // Start the writePump in a separate goroutine to send messages
	go s.readPump(client)             // Start the readPump in a separate goroutine to listen for incoming messages
	s.BroadcastWelcomeMessage(client)

	// Broadcast a join message
	joinMsg := fmt.Sprintf("%s has joined the server", client.Username)
	for _, c := range s.Clients {
		if c.ID != client.ID {
			c.Send <- []byte(joinMsg)
		}
	}
}

// readPump reads incoming WebSocket messages from the client
func (s *Server) readPump(c *Client) {
	defer s.RemoveClient(c)

	c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second)) // Set a read deadline for the WebSocket connection
	c.Conn.SetPongHandler(func(string) error {               // Set a pong handler to handle ping messages
		c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
		c.Active <- struct{}{}
		return nil
	})

	for {
		_, msg, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			log.Println("JSON decode error:", err)
			continue
		}

		// Handle private and broadcast messages
		s.ClientsMu.Lock()
		if message.Type == "private" {
			for _, recipient := range s.Clients {
				if recipient.ID == message.Receiver {
					recipient.Send <- msg
					break
				}
			}
		} else if message.Type == "broadcast" {
			for _, recipient := range s.Clients {
				if recipient.ID != c.ID {
					recipient.Send <- msg
				}
			}
		}
		s.ClientsMu.Unlock()
	}
}

// writePump writes WebSocket messages to the client
func (s *Server) writePump(c *Client) {
	ticker := time.NewTicker(10 * time.Second) // Set up a ticker to send a ping message every 10 seconds
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case msg, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, msg)

		case <-ticker.C: // Send a ping message every 10 seconds
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				s.ClientsMu.Lock()
				log.Printf("Ping error: %v", err)
				// Broadcast a leave message
				leaveMsg := fmt.Sprintf("%s has left the server", c.Username)
				for _, client := range s.Clients {
					if client.ID != c.ID {
						client.Send <- []byte(leaveMsg)
					}
				}
				delete(s.Clients, c.ID)
				delete(s.Usernames, c.Username)

				s.ClientsMu.Unlock()

				return
			}

		case <-c.Active:
		}
	}
}

func (s *Server) getClientList(excludeID string) []string {
	s.ClientsMu.RLock()
	defer s.ClientsMu.RUnlock()
	var list []string
	for id, client := range s.Clients {
		if id != excludeID {
			list = append(list, client.Username)
		}
	}
	return list
}

func (s *Server) generateUsername() string {
	return gofakeit.Username()
}
