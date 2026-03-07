package ws

import (
	"github.com/gorilla/websocket"
)

type Hub struct {
	Client       map[*Client]bool
	Broadcast    chan []byte
	AddClient    chan *Client
	removeClient chan *Client
}

func NewHub() *Hub {
	return &Hub{
		Client:       make(map[*Client]bool),
		Broadcast:    make(chan []byte),
		AddClient:    make(chan *Client),
		removeClient: make(chan *Client),
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
				select {
				case client.send <- msg:
					go func() {
						for msg := range client.send {
							client.conn.WriteMessage(websocket.TextMessage, msg)
						}
					}()
				default:
					close(client.send)
					delete(h.Client, client)
				}
			}
		}
	}
}
