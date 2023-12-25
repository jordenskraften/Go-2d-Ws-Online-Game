package entities

type ChatMessage struct {
	Username string `json:"username"`
	Date     string `json:"date"`
	Text     string `json:"text"`
}

type Canvas struct {
	Positions map[string]Position `json:"positions"`
}

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
