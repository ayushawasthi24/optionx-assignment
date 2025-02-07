package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// HandleConnections manages incoming WebSocket connections.
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

func (s *Server) AddClient(client *Client) {
	s.ClientsMu.Lock()
	s.Clients[client.ID] = client
	s.ClientsMu.Unlock()
}

func (s *Server) RemoveClient(client *Client) {
	s.ClientsMu.Lock()
	delete(s.Clients, client.ID)
	delete(s.Usernames, client.Username)
	s.ClientsMu.Unlock()
	client.Conn.Close()
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	conn, err := s.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error:", err)
		return
	}
	sessionID := uuid.New().String()
	username := s.generateUniqueUsername()
	client := &Client{
		ID:       sessionID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte),
		Active:   make(chan struct{}, 1),
	}

	s.AddClient(client)
	go s.writePump(client)
	go s.readPump(client)
	s.BroadcastWelcomeMessage(client)
}

func (s *Server) readPump(c *Client) {
	defer s.RemoveClient(c)

	c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
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

		s.ClientsMu.Lock()
		if message.Type == "private" {
			for _, recipient := range s.Clients {
				if recipient.Username == message.Receiver {
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

func (s *Server) writePump(c *Client) {
	ticker := time.NewTicker(10 * time.Second)
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

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
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

func (s *Server) generateUniqueUsername() string {
	for {
		username := gofakeit.Username()
		s.ClientsMu.RLock()
		if !s.Usernames[username] {
			s.ClientsMu.RUnlock()
			s.ClientsMu.Lock()
			s.Usernames[username] = true
			s.ClientsMu.Unlock()
			return username
		}
		s.ClientsMu.RUnlock()
	}
}
