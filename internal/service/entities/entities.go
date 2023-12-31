package entities

type ChatMessage struct {
	Date string `json:"date"`
	Text string `json:"text"`
}
type Position struct {
	Username string  `json:"username"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
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

type CanvasMessageData struct {
	Type      string     `json:"type"`
	Positions []Position `json:"positions"`
}
