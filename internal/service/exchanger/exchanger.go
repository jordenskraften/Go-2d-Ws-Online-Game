// принимает жсон и юзера,
// решая в какое лобби отправить его месейдж
// будь то текст в чат или коорды в канвас
package exchanger

import (
	"log"
	"sync"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/lobby"
)

type Exchanger struct {
	Lobbies []*lobby.Lobby
	Hub     *hub.Hub
	mu      sync.Mutex
}

func NewExchanger(hub *hub.Hub) *Exchanger {
	ex := &Exchanger{
		Hub:     hub,
		Lobbies: make([]*lobby.Lobby, 0),
	}
	return ex
}

func (ex *Exchanger) AddLobby(lobby *lobby.Lobby) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	ex.Lobbies = append(ex.Lobbies, lobby)
}

func (ex *Exchanger) CreateLobby(name string) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	lobby := lobby.NewLobby(name, ex.Hub)
	ex.Lobbies = append(ex.Lobbies, lobby)
}

func (ex *Exchanger) GetAnyLobby() *lobby.Lobby {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	if len(ex.Lobbies) > 0 {
		return ex.Lobbies[0] // Возвращаем первый доступный лобби (можно выбирать другой логикой)
	}

	return nil // Возвращаем nil, если лобби нет
}

func (ex *Exchanger) AddConnectionToLobby(conn *hub.ConnItem, lobby *lobby.Lobby) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	lobby.AddConnection(conn)

}

func (ex *Exchanger) RemoveСonnectionFromLobby(connName string, lobbyName string) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	//тут в слайсе находим сперва лобби с именем
	lobby := ex.GetLobbyByName(lobbyName)
	if lobby == nil {
		return
	}
	//log.Println(lobby)

	//в хабе аналогично ищем по имени есть ли такой конект
	conn := ex.Hub.GetConnectionByName(connName)
	if conn == nil {
		return
	}
	//log.Println(conn)

	lobby.RemoveConnection(connName)
}

func (ex *Exchanger) GetLobbyByName(name string) *lobby.Lobby {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	for _, val := range ex.Lobbies {
		if val.Name == name {
			log.Println("найдено лобби с именем " + name)
			return val
		}
	}
	return nil
}

func (ex *Exchanger) GetUserLobby(conn *hub.ConnItem) *lobby.Lobby {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	//найти лобби в котором сидит юзер
	//ключевой метод, пригодится понимать куда отправлять меседж юзера

	for _, lobby := range ex.Lobbies {
		for _, user := range lobby.Connections {
			if conn == user {
				log.Printf("нашли юзеру %d лобби %d", conn.Name, lobby.Name)
				return lobby
			}
		}
	}
	log.Printf("не нашли юзеру %d лобби", conn.Name)
	return nil
}

func (ex *Exchanger) DeleteUserFromAllLobbies(conn *hub.ConnItem) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	for {
		lobby := ex.GetUserLobby(conn)
		if lobby == nil {
			return
		}
		lobby.RemoveConnection(conn.Name)
	}
}

//метод выбирающий в какое лобби писать месейдж
