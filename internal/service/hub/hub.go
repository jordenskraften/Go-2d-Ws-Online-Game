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
	//чтобы избавиться от импорта зависимости ws
	//можно сюда в аргументы вписать поля из примитивов ws connect
	//например канал send
	//и собрать новый конект тут из этих примитивов, добавив его в слайс
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Connections[conn.Name] = conn
}

func (h *Hub) RemoveConnection(conn *WsConnection) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Connections, conn.Name)
}
