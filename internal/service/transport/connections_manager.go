package transport

import (
	"log"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/exchanger"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
)

type ConnectionsManager struct {
	hub       *hub.Hub
	exchanger *exchanger.Exchanger
}

func NewConnectionsManager(hub *hub.Hub, ex *exchanger.Exchanger) *ConnectionsManager {
	log.Println("создаем конект менежер")
	return &ConnectionsManager{
		hub:       hub,
		exchanger: ex,
	}
}

// находит юзера в хабе и добавляет через эксченжер в старт лобби
func (cm *ConnectionsManager) SetupUser(name string, hub hub.Hub) {
	conn := cm.hub.GetConnectionByName(name)
	if conn == nil {
		return
	}
	log.Printf("конект менежер добавляет юзера в пул %s", conn.Name)
	cm.exchanger.SetupConnection(conn)
	log.Printf("конект менежер закончил добавление юзера в пул %s", conn.Name)
}

func (cm *ConnectionsManager) RemoveUser(name string, hub hub.Hub) {
	conn := cm.hub.GetConnectionByName(name)
	if conn == nil {
		return
	}
	log.Printf("конект менежер удаляет юзера из пула %s", conn.Name)
	cm.exchanger.DeleteUserFromAllLobbies(conn)
	log.Printf("конект менежер закончил удаление юзера из всех лобби %s", conn.Name)
}
