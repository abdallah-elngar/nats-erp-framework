package websocket

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client يمثل عميل WebSocket
type Client struct {
	conn *websocket.Conn
	hub  *Hub
	send chan []byte
	id   string
	mu   sync.Mutex
}

// NewClient ينشئ عميلاً جديداً
func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		hub:  hub,
		send: make(chan []byte, 256),
		id:   generateID(),
	}
}

// Read يقرأ الرسائل من العميل
func (c *Client) Read() {
	defer func() {
		c.hub.Unregister(c)
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		c.hub.Broadcast(message)
	}
}

// Write يكتب الرسائل إلى العميل
func (c *Client) Write() {
	defer c.conn.Close()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.mu.Lock()
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			c.mu.Unlock()

			if err != nil {
				return
			}
		}
	}
}

// Send يرسل رسالة إلى العميل
func (c *Client) Send(message []byte) {
	c.send <- message
}

// Close يغلق العميل
func (c *Client) Close() {
	close(c.send)
	c.conn.Close()
}

// generateID يولد معرفاً فريداً
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
