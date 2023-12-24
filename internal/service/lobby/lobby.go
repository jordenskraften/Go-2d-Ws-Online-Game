package lobby

import (
	"sync"

	ws "github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/server"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
)

type Lobby struct {
	Name        string
	mu          sync.Mutex
	hub         *hub.Hub
	Connections map[string]*ws.WsConnection
	Canvas      *Canvas
	Chat        *Chat
}

func NewLobby(name string, hub *hub.Hub) *Lobby {
	return &Lobby{
		Name:        name,
		Connections: make(map[string]*ws.WsConnection), // <-- добавленная запятая после этой строки
		hub:         hub,
		Canvas:      NewCanvas(name),
		Chat:        NewChat(name),
	}
}

func (lo *Lobby) AddConnection(conn *ws.WsConnection) {
	lo.mu.Lock()
	defer lo.mu.Unlock()

	lo.Connections[conn.Name] = conn
}

//вот здесь в чат и в лобби будут добавляться месейджи
//и при добавлении еще будет отправляться в хаб список имен или конекшнов хз
//по которым надо разослать через ws месейджи
