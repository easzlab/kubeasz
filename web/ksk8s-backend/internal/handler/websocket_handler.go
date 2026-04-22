package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/easzlab/ksk8s/internal/websocket"
	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var upgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub      *websocket.Hub
	ringMap  *websocket.LogRingMap
}

func NewWebSocketHandler(hub *websocket.Hub, ringMap *websocket.LogRingMap) *WebSocketHandler {
	return &WebSocketHandler{
		hub:     hub,
		ringMap: ringMap,
	}
}

// ServeWS handles WebSocket connections for task log streaming.
func (h *WebSocketHandler) ServeWS(c *gin.Context) {
	taskID, _ := strconv.ParseInt(c.Param("id"), 10, 64)
	offsetStr := c.Query("offset")
	offset := 0
	if offsetStr != "" {
		offset, _ = strconv.Atoi(offsetStr)
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &websocket.Client{
		Hub:    h.hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		TaskID: taskID,
	}

	h.hub.Register(client)

	// Send historical lines from ring buffer
	ring := h.ringMap.Get(taskID)
	lines, total := ring.Since(offset)
	for _, line := range lines {
		payload, _ := json.Marshal(line)
		select {
		case client.Send <- payload:
		default:
		}
	}

	// Send a "ready" marker so client knows replay is done
	ready, _ := json.Marshal(map[string]interface{}{
		"type":        "ready",
		"total_lines": total,
	})
	select {
	case client.Send <- ready:
	default:
	}

	go client.WritePump()
	go client.ReadPump()
}
