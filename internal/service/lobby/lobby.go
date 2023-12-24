package lobby

import (
	"sync"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
)

type Lobby struct {
	Name        string
	mu          sync.Mutex
	hub         *hub.Hub
	Connections map[string]*hub.WsConnection
	Canvas      *Canvas
	Chat        *Chat
}

func NewLobby(name string, currHub *hub.Hub) *Lobby {
	return &Lobby{
		Name:        name,
		Connections: make(map[string]*hub.WsConnection),
		hub:         currHub,
		Canvas:      NewCanvas(name),
		Chat:        NewChat(name),
	}
}

func (lo *Lobby) AddConnection(conn *hub.WsConnection) {
	lo.mu.Lock()
	defer lo.mu.Unlock()

	lo.Connections[conn.Name] = conn
}

//вот здесь в чат и в лобби будут добавляться месейджи
//и при добавлении еще будет отправляться в хаб список имен или конекшнов хз
//по которым надо разослать через ws месейджи
