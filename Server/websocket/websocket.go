package websocket

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Server struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	upgrader  websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		clients:   make(map[*websocket.Conn]bool),
		broadcast: make(chan []byte),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // Allow all connections regardless of origin
			},
		},
	}
}

func (s *Server) HandleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer ws.Close()

	s.clients[ws] = true

	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			delete(s.clients, ws)
			break
		}
	}
}

func (s *Server) BroadcastToClients(message []byte) {
	s.broadcast <- message
}

func (s *Server) Run() {
	for {
		select {
		case message := <-s.broadcast:
			for client := range s.clients {
				err := client.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error writing to Websocket: %v", err)
					client.Close()
					delete(s.clients, client)
				}
			}
		}
	}
}
