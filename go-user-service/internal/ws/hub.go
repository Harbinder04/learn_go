package ws

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Client       map[*websocket.Conn]bool
	Broadcast    chan []byte
	AddClient    chan *websocket.Conn
	removeClient chan *websocket.Conn
}

func NewHub() *Hub {
	return &Hub{
		Client:       make(map[*websocket.Conn]bool),
		Broadcast:    make(chan []byte),
		AddClient:    make(chan *websocket.Conn),
		removeClient: make(chan *websocket.Conn),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.AddClient:
			h.Client[client] = true
		case client := <-h.removeClient:
			delete(h.Client, client)
		case msg := <-h.Broadcast:
			for client := range h.Client {
				client.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}
