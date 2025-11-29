package socket

import (
	"sync"

	"github.com/dibyajyoti-mandal/code-exec-engine/models"
	"github.com/gorilla/websocket"
)

type ClientMap struct {
	clients map[string]*websocket.Conn
	mu      sync.RWMutex
}

var Pool = ClientMap{
	clients: make(map[string]*websocket.Conn),
}

func (c *ClientMap) Add(id string, conn *websocket.Conn) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.clients[id] = conn
}

func (c *ClientMap) Remove(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if conn, ok := c.clients[id]; ok {
		conn.Close()
		delete(c.clients, id)
	}
}

func (c *ClientMap) SendResult(clientID string, result models.Result) {
	c.mu.RLock()
	conn, ok := c.clients[clientID]
	c.mu.RUnlock()

	if !ok {
		return
	}
	conn.WriteJSON(result)
}

// Send to EVERYONE connected
func (c *ClientMap) Broadcast(result models.Result) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, conn := range c.clients {
		conn.WriteJSON(result)
	}
}
