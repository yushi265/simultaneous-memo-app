package websocket

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients by page ID
	rooms map[string]map[*Client]bool

	// Inbound messages from the clients
	broadcast chan *Message

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	mu sync.RWMutex
}

// Message represents a WebSocket message
type Message struct {
	PageID  string `json:"pageId"`
	Type    string `json:"type"`
	Content []byte `json:"content"`
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[string]map[*Client]bool),
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastToRoom(message)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.rooms[client.pageID]; !ok {
		h.rooms[client.pageID] = make(map[*Client]bool)
	}
	h.rooms[client.pageID][client] = true
	log.Printf("Client registered to page %s", client.pageID)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if clients, ok := h.rooms[client.pageID]; ok {
		if _, ok := clients[client]; ok {
			delete(clients, client)
			close(client.send)
			log.Printf("Client unregistered from page %s", client.pageID)

			// Clean up empty rooms
			if len(clients) == 0 {
				delete(h.rooms, client.pageID)
			}
		}
	}
}

func (h *Hub) broadcastToRoom(message *Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if clients, ok := h.rooms[message.PageID]; ok {
		for client := range clients {
			select {
			case client.send <- message.Content:
			default:
				// Client's send channel is full, close it
				close(client.send)
				delete(clients, client)
			}
		}
	}
}