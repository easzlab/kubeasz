package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/easzlab/ksk8s/internal/ssh"
	"github.com/gin-gonic/gin"
	gorilla "github.com/gorilla/websocket"
)

var sshUpgrader = gorilla.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type SSHHandler struct{}

func NewSSHHandler() *SSHHandler {
	return &SSHHandler{}
}

// ServeWS handles WebSocket connections for WebSSH.
func (h *SSHHandler) ServeWS(c *gin.Context) {
	addr := c.Query("addr")
	if addr == "" {
		addr = "127.0.0.1:22"
	}
	user := c.Query("user")
	if user == "" {
		user = "root"
	}
	password := c.Query("password")
	keyPath := c.Query("key_path")
	if keyPath == "" {
		keyPath = os.Getenv("KSK8S_SSH_KEY")
	}
	if keyPath == "" {
		keyPath = "/root/.ssh/id_rsa"
	}

	conn, err := sshUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	client, err := ssh.Connect(addr, user, password, keyPath)
	if err != nil {
		conn.WriteMessage(gorilla.TextMessage, []byte("SSH connection failed: "+err.Error()))
		return
	}
	defer client.Close()

	// Send connected message
	conn.WriteMessage(gorilla.TextMessage, []byte("\r\nConnected to "+addr+"\r\n"))

	// Goroutine: SSH stdout -> WebSocket
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := client.Read(buf)
			if err != nil {
				if err != io.EOF {
					conn.WriteMessage(gorilla.TextMessage, []byte("\r\nSSH read error: "+err.Error()+"\r\n"))
				}
				conn.Close()
				return
			}
			if err := conn.WriteMessage(gorilla.BinaryMessage, buf[:n]); err != nil {
				return
			}
		}
	}()

	// Main loop: WebSocket -> SSH stdin
	for {
		msgType, data, err := conn.ReadMessage()
		if err != nil {
			return
		}

		if msgType == gorilla.TextMessage {
			// Check for resize command
			var resizeMsg struct {
				Type string `json:"type"`
				Cols int    `json:"cols"`
				Rows int    `json:"rows"`
			}
			if json.Unmarshal(data, &resizeMsg) == nil && resizeMsg.Type == "resize" {
				client.Resize(resizeMsg.Cols, resizeMsg.Rows)
				continue
			}
		}

		if _, err := client.Write(data); err != nil {
			return
		}
	}
}
