package hub

import (
	"sync"

	ws "github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/server"
)

type Hub struct {
	Name        string
	Connections map[string]*ws.WsConnection
	mu          sync.Mutex
}

func NewHub(name string) *Hub {
	return &Hub{
		Name:        name,
		Connections: make(map[string]*ws.WsConnection),
	}
}

func (h *Hub) AddConnection(conn *ws.WsConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Connections[conn.Name] = conn
}
