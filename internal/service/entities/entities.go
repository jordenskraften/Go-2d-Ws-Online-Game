package entities

type ChatMessage struct {
	Date string `json:"date"`
	Text string `json:"text"`
}
type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type LobbyCommand struct {
	LobbyName string `json:"lobby_name"`
}

type ChatMessageData struct {
	Type     string `json:"type"`
	Username string `json:"username"`
	Date     string `json:"date"`
	Text     string `json:"text"`
}
