package hub

import (
	"sync"
)

type Hub struct {
	Name        string
	Connections map[string]*WsConnection
	mu          sync.Mutex
}

func NewHub(name string) *Hub {
	return &Hub{
		Name:        name,
		Connections: make(map[string]*WsConnection),
	}
}

func (h *Hub) AddConnection(conn *WsConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Connections[conn.Name] = conn
}

func (h *Hub) RemoveConnection(conn *WsConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Connections, conn.Name)
}
