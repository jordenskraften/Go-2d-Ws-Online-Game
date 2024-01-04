// принимает жсон и юзера,
// решая в какое лобби отправить его месейдж
// будь то текст в чат или коорды в канвас
package exchanger

import (
	"log"
	"sync"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/entities"
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

func (ex *Exchanger) GetAllActiveLobbies() []*lobby.Lobby {
	ex.mu.RLock()
	defer ex.mu.RUnlock()

	return ex.Lobbies
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

	for _, val := range ex.Lobbies {
		if val.Name == name {
			log.Println("найдено лобби с именем " + name)
			return val
		}
	}
	return nil
}

func (ex *Exchanger) GetUserLobby(conn *hub.ConnItem) *lobby.Lobby {

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

func (ex *Exchanger) AddConnectionToLobby(conn *hub.ConnItem, lo *lobby.Lobby) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	lo.AddConnection(conn)
	//инфу всем в лобби о новичке
	canvasPos := lobby.NewPositionRandomCoords(conn.Name)
	msgPos := &entities.Position{
		Username: canvasPos.Username,
		X:        canvasPos.X,
		Y:        canvasPos.Y,
	}
	ex.BroadcastPositionMessage(conn, msgPos)
	ex.BroadcastChatUserStatus(lo, conn.Name, true)

}

func (ex *Exchanger) DeleteUserFromAllLobbies(conn *hub.ConnItem) {
	for {
		lobby := ex.GetUserLobby(conn)
		if lobby == nil {
			return
		} else {
			ex.RemoveСonnectionFromLobby(conn.Name, lobby)
		}
	}
}

func (ex *Exchanger) RemoveСonnectionFromLobby(connName string, lobby *lobby.Lobby) {
	ex.mu.Lock()
	defer ex.mu.Unlock()

	lobby.RemoveConnection(connName)
	//инфу всем в лобби о ливере
	ex.BroadcastCanvasInfo(lobby)
	ex.BroadcastChatUserStatus(lobby, connName, false)
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

	ex.SetUserLobby(conn, lobby.Name)
	log.Printf("user %s added to lobby %s, lobby have that user? %s\n", conn.Name, lobby.Name, lobby.Connections[conn.Name].Name)
}

func (ex *Exchanger) GetLobbiesNamesList() []string {
	lobbies := ex.GetAllActiveLobbies()
	lobbiesNames := make([]string, 0, len(lobbies))
	for _, val := range lobbies {
		lobbiesNames = append(lobbiesNames, val.Name)
	}
	return lobbiesNames
}

//метод который будет переключать юзеру лобби

func (ex *Exchanger) SetUserLobby(conn *hub.ConnItem, lobbyName string) {
	log.Printf("set lobby to user with name %s \n", conn.Name)
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
	} else {
		log.Println("lobby or connection is wrong")
	}

	ex.AddConnectionToLobby(conn, lobby)
	log.Printf("added user %s to lobby %s \n", conn.Name, lobby.Name)
	//вот здесь надо список лобби взять и отдать конекшну
	lobbiesNames := ex.GetLobbiesNamesList()
	log.Printf("active lobbies names list: %v", lobbiesNames)
	ex.Hub.SendLobbiesListToConnestion(conn, lobbiesNames, lobby.Name)
	log.Printf("finished changing lobby to user with name %s \n", conn.Name)

}

// --------------

func (ex *Exchanger) LobbyChangeCommand(conn *hub.ConnItem, lobbyName string) {
	curLobby := ex.GetUserLobby(conn)
	if curLobby == nil {
		log.Printf("no lobbies with name %s for user %s", lobbyName, conn.Name)
		return
	} else {
		if curLobby.Name == lobbyName {
			log.Printf("breaking change lobby, user %s already in lobby %s", conn.Name, lobbyName)
			return
		}
		ex.SetUserLobby(conn, lobbyName)
	}
}

// --------------
func (ex *Exchanger) BroadcastChatMessage(conn *hub.ConnItem, msg *entities.ChatMessage) {
	//определяем лобби юзера
	lobby := ex.GetUserLobby(conn)
	if lobby == nil {
		log.Printf("user %s has no lobby for broadcast chat message \n", conn.Name)
		return
	} else {
		//достаем чат лобби
		chat := lobby.Chat
		// записываем в него меседж
		chatMsg := chat.AddChatMessage(conn.Name, msg.Text)
		// достаем список активных конектов в лобби
		lobbyUsers := lobby.GetActiveConnectionsList()
		//эт логирование списка юзеров лобби
		// userNames := make([]string, len(lobbyUsers))
		// for i, userConn := range lobbyUsers {
		// 	userNames[i] = userConn.Name
		// }
		// userListStr := strings.Join(userNames, ", ")
		// log.Printf("users in lobby %s list are: %s", lobby.Name, userListStr)
		// всему списку по вебсокету отправляет
		ex.Hub.BroadcastChatMessageToUserList(lobbyUsers, chatMsg.Username, chatMsg.Text, chatMsg.Date)
	}
}

func (ex *Exchanger) BroadcastPositionMessage(conn *hub.ConnItem, msg *entities.Position) {
	//определяем лобби юзера
	lobby := ex.GetUserLobby(conn)
	if lobby == nil {
		log.Printf("user %s has no lobby for broadcast chat message \n", conn.Name)
		return
	} else {
		canvas := lobby.Canvas
		if !canvas.IsUserInCanvas(conn.Name) {
			canvas.AddUser(conn.Name, float32(msg.X), float32(msg.Y))
		} else {
			canvas.ChangeUserCoords(conn.Name, float32(msg.X), float32(msg.Y))
		}
		//теперь надо дернуть хаб чтобы он всем распространил месейдж
		lobbyUsers := lobby.GetActiveConnectionsList()
		canvasInfo := canvas.GetCanvasInfo()
		canvasMsg := entities.CanvasMessageData{
			Type:      "CanvasMessageData",
			Positions: []entities.Position{},
		}
		for _, val := range canvasInfo {
			canvasMsg.Positions = append(canvasMsg.Positions,
				entities.Position{
					Username: val.Username,
					X:        val.X,
					Y:        val.Y,
				})
		}

		ex.Hub.BroadcastCanvasDataToUserList(lobbyUsers, canvasMsg)

	}
}
func (ex *Exchanger) BroadcastCanvasInfo(lobby *lobby.Lobby) {
	canvas := lobby.Canvas
	//теперь надо дернуть хаб чтобы он всем распространил месейдж
	lobbyUsers := lobby.GetActiveConnectionsList()
	canvasInfo := canvas.GetCanvasInfo()
	canvasMsg := entities.CanvasMessageData{
		Type:      "CanvasMessageData",
		Positions: []entities.Position{},
	}
	for _, val := range canvasInfo {
		canvasMsg.Positions = append(canvasMsg.Positions,
			entities.Position{
				Username: val.Username,
				X:        val.X,
				Y:        val.Y,
			})
	}

	ex.Hub.BroadcastCanvasDataToUserList(lobbyUsers, canvasMsg)

}

func (ex *Exchanger) BroadcastChatUserStatus(lobby *lobby.Lobby, name string, status bool) {
	//определяем лобби юзера
	//достаем чат лобби
	chat := lobby.Chat
	// записываем в него меседж
	text := "has entered chat " + lobby.Name
	if !status {
		text = "has left chat" + lobby.Name
	}
	chatMsg := chat.AddChatMessage(name, text)
	// достаем список активных конектов в лобби
	lobbyUsers := lobby.GetActiveConnectionsList()
	ex.Hub.BroadcastChatMessageToUserList(lobbyUsers, chatMsg.Username, chatMsg.Text, "")

}
