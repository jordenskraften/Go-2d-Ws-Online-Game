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
	mu      sync.RWMutex
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

func (ex *Exchanger) AddConnectionToLobby(conn *hub.ConnItem, lobby *lobby.Lobby) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	lobby.AddConnection(conn)

}

func (ex *Exchanger) RemoveСonnectionFromLobby(connName string, lobbyName string) {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

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

func (ex *Exchanger) GetAnyLobby() *lobby.Lobby {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

	if len(ex.Lobbies) > 0 {
		return ex.Lobbies[0] // Возвращаем первый доступный лобби (можно выбирать другой логикой)
	}

	return nil // Возвращаем nil, если лобби нет
}

func (ex *Exchanger) GetLobbyByName(name string) *lobby.Lobby {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

	for _, val := range ex.Lobbies {
		if val.Name == name {
			log.Println("найдено лобби с именем " + name)
			return val
		}
	}
	return nil
}

func (ex *Exchanger) GetUserLobby(conn *hub.ConnItem) *lobby.Lobby {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

	//найти лобби в котором сидит юзер
	//ключевой метод, пригодится понимать куда отправлять меседж юзера

	for _, lobby := range ex.Lobbies {
		for _, user := range lobby.Connections {
			if conn == user {
				log.Printf("нашли лобби %s с юзером %s", conn.Name, lobby.Name)
				return lobby
			}
		}
	}
	log.Printf("не нашли лобби с юзером %s", conn.Name)
	return nil
}

func (ex *Exchanger) DeleteUserFromAllLobbies(conn *hub.ConnItem) {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

	for {
		lobby := ex.GetUserLobby(conn)
		if lobby == nil {
			return
		}
		lobby.RemoveConnection(conn.Name)
	}
}

func (ex *Exchanger) SetupConnection(conn *hub.ConnItem) {
	/*
		теперь надо что сделать
		функция которая в эксченжер сует коннект с вебсокета
		и записывает его в стартовое лобби если оно есть
		если этого лобби нет, оно создается
	*/
	log.Println("SetupConnection " + conn.Name)

	lobby := ex.GetAnyLobby()
	if lobby == nil {
		// Создаем новое лобби и используем его
		ex.CreateLobby("start lobby")
		lobby = ex.GetAnyLobby() // Получаем только что созданное лобби
		if lobby == nil {
			log.Println("Ошибка при создании лобби")
			return
		}
	}

	lobby.AddConnection(conn)
	log.Printf("user %s added to lobby %s, lobby have that user? = %s\n", conn.Name, lobby.Name, lobby.Connections[conn.Name].Name)

}

//метод который будет переключать юзеру лобби

func (ex *Exchanger) ChangeUserLobby(conn *hub.ConnItem, lobbyName string) {
	log.Printf("changing lobby to user with name %s \n", conn.Name)
	// проверить активен ли конект в хабе
	inHub := ex.Hub.IsConnectionInHub(conn.Name)
	if !inHub {
		log.Printf("user with name %s dont exist in hub \n", conn.Name)
		return
	}
	// проверить существует ли в пуле лобби с таким именем
	lobby := ex.GetLobbyByName(lobbyName)
	if lobby == nil {
		log.Printf("lobby with name %s dont exist \n", lobbyName)
		return
	}
	// удалить юзера из всех лобби
	// добавить юзера в лобби
	if inHub && lobby != nil {
		ex.DeleteUserFromAllLobbies(conn)
		ex.AddConnectionToLobby(conn, lobby)
		log.Printf("added user %s to lobby \n", conn.Name, lobby.Name)
	} else {
		log.Println("lobby or connection is wrong")
	}
	log.Printf("finished changing lobby to user with name %s \n", conn.Name)
}

//метод выбирающий в какое лобби писать месейдж
//тут надо наверно задействовать интерфейсы и рефлексию
//всовывать сюда один тип Messageable по интерфейсу
//а далее уже смотреть какой конкретно меседж
//и реализовать в свиче его бизнес логику
//чат месейдж, канвас месейдж, команд лобби меседж
// получается я должен сперва типы месейдж и интерфейс к ним прибабахать
