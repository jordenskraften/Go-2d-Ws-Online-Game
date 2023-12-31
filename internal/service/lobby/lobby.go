package lobby

import (
	"log"
	"sync"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
)

type Lobby struct {
	Name        string
	mu          sync.RWMutex
	hub         *hub.Hub
	Connections map[string]*hub.ConnItem
	Canvas      *Canvas
	Chat        *Chat
}

func NewLobby(name string, currHub *hub.Hub) *Lobby {
	return &Lobby{
		Name:        name,
		Connections: make(map[string]*hub.ConnItem),
		hub:         currHub,
		Canvas:      NewCanvas(name),
		Chat:        NewChat(name),
	}
}

func (lo *Lobby) AddConnection(conn *hub.ConnItem) {
	lo.mu.Lock()
	defer lo.mu.Unlock()

	lo.Connections[conn.Name] = conn
}

func (lo *Lobby) RemoveConnection(name string) {

	lo.mu.Lock()
	defer lo.mu.Unlock()

	lo.Canvas.RemoveUser(name)
	delete(lo.Connections, name)
	log.Printf("удалили юзера из %s лобби %s", lo.Name, name)
}

func (lo *Lobby) GetActiveConnectionsList() []*hub.ConnItem {
	lo.mu.RLock()
	defer lo.mu.RUnlock()

	list := make([]*hub.ConnItem, 0, len(lo.Connections))

	for _, conn := range lo.Connections {
		list = append(list, conn)
	}

	return list
}
