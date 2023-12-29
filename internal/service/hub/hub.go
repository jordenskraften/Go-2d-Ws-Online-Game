package hub

import (
	"sync"

	"github.com/gorilla/websocket"
)

type ConnItem struct {
	Name string
	conn *websocket.Conn
	send chan []byte
}
type Hub struct {
	Name        string
	Connections map[string]*ConnItem
	mu          sync.RWMutex
}

func NewHub(name string) *Hub {
	return &Hub{
		Name:        name,
		Connections: make(map[string]*ConnItem),
	}
}

func (h *Hub) AddConnection(name string, conn *websocket.Conn, send chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Connections[name] = &ConnItem{
		Name: name,
		conn: conn,
		send: send,
	}
}

func (h *Hub) RemoveConnection(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Connections, name)
}

func (h *Hub) GetConnectionByName(name string) *ConnItem {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.Connections[name]; ok {
		return conn
	}
	return nil
}

func (h *Hub) IsConnectionInHub(name string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.Connections[name]
	return ok
}

func (h *Hub) GetAnyConnection() *ConnItem {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, val := range h.Connections {
		return val // Возвращаем первое доступное соединение
	}

	return nil // Возвращаем nil, если соединений нет
}
