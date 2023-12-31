package transport

import (
	"log"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/entities"
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

func (cm *ConnectionsManager) DistributeMessage(username string, msgData map[string]interface{}) {
	// определяет что за месейдж пришел от клиента
	// и решает какую логику к нему применить
	conn := cm.hub.GetConnectionByName(username)
	if conn == nil {
		log.Println("No connection with that nickname. Dirstibute message is stopped")
		return
	}
	log.Printf("New Message from username %s to server: \n", conn.Name)

	if msgType, ok := msgData["type"]; ok {
		switch msgType {
		case "LobbyCommand":
			lobbyName, lobbyNameExists := msgData["LobbyName"].(string)

			if !lobbyNameExists {
				log.Println("Incomplete LobbyCommand data")
			}

			message := entities.LobbyCommand{
				LobbyName: lobbyName,
			}
			log.Println("Received LobbyCommand object:", message)

		case "ChatMessage":
			text, textExists := msgData["text"].(string)

			if !textExists {
				log.Println("Incomplete ChatMessage data")
			}

			message := entities.ChatMessage{
				Text: text,
			}
			log.Println("Received ChatMessage object:", message)
			//тута в чат запись
			cm.exchanger.BroadcastChatMessage(conn, &message)

		case "Position":
			x, xOk := msgData["x"].(float64)
			y, yOk := msgData["y"].(float64)
			if !xOk || !yOk {
				log.Println("Error extracting coordinates for Position")
			}

			position := entities.Position{
				X: float32(x),
				Y: float32(y),
			}
			log.Println("Received Position object:", position)

			cm.exchanger.BroadcastPositionMessage(conn, &position)
		}
	}
}
