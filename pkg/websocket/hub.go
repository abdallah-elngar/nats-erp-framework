package websocket

import (
	"net/http"
	"sync"
)

// Hub يدير اتصالات WebSocket
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// NewHub ينشئ مركز WebSocket جديداً
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run يقوم بتشغيل المركز
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

// Register يسجل عميلاً
func (h *Hub) Register(client *Client) {
	h.register <- client
}

// Unregister يلغي تسجيل عميل
func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

// Broadcast يبث رسالة
func (h *Hub) Broadcast(message []byte) {
	h.broadcast <- message
}

// GetClients يعيد عدد العملاء
func (h *Hub) GetClients() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// ServeHTTP يعالج طلبات WebSocket
func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(conn, h)
	h.Register(client)

	go client.Write()
	go client.Read()
}
