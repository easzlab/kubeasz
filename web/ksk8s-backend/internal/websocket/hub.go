package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Hub manages WebSocket connections per task ID.
type Hub struct {
	mu         sync.RWMutex
	clients    map[int64]map[*Client]bool
	broadcast  chan TaskMessage
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	TaskID int64
}

type TaskMessage struct {
	TaskID  int64
	Payload []byte
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]map[*Client]bool),
		broadcast:  make(chan TaskMessage, 256),
		register:   make(chan *Client, 64),
		unregister: make(chan *Client, 64),
	}
}

// Run starts the hub event loop.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.TaskID] == nil {
				h.clients[client.TaskID] = make(map[*Client]bool)
			}
			h.clients[client.TaskID][client] = true
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.TaskID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.clients, client.TaskID)
				}
				close(client.Send)
			}
			h.mu.Unlock()

		case msg := <-h.broadcast:
			h.mu.RLock()
			clients := h.clients[msg.TaskID]
			h.mu.RUnlock()
			for client := range clients {
				select {
				case client.Send <- msg.Payload:
				default:
					// Client is slow, drop message
				}
			}
		}
	}
}

// Register registers a client.
func (h *Hub) Register(c *Client) {
	h.register <- c
}

// Unregister unregisters a client.
func (h *Hub) Unregister(c *Client) {
	h.unregister <- c
}

// Broadcast sends a message to all clients for a task.
func (h *Hub) Broadcast(taskID int64, payload []byte) {
	select {
	case h.broadcast <- TaskMessage{TaskID: taskID, Payload: payload}:
	case <-time.After(100 * time.Millisecond):
		// Drop if broadcast channel is backed up
	}
}

// ClientCount returns the number of connected clients for a task.
func (h *Hub) ClientCount(taskID int64) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients[taskID])
}

// WritePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.Conn.WriteMessage(websocket.TextMessage, message)

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ReadPump pumps messages from the websocket to the hub (mostly to detect disconnect).
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// log unexpected close
			}
			break
		}
	}
}
