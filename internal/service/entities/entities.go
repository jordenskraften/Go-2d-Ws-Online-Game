package entities

import (
	"math/rand"
	"time"
)

type ChatMessage struct {
	Date string `json:"date"`
	Text string `json:"text"`
}
type Position struct {
	Username string  `json:"username"`
	X        float32 `json:"x"`
	Y        float32 `json:"y"`
}

func NewPositionRandomCoords(name string) *Position {
	rand.Seed(time.Now().UnixNano())

	minX := 21
	maxX := 370
	randomX := float32(rand.Intn(maxX-minX+1) + minX)

	minY := 21
	maxY := 270
	randomY := float32(rand.Intn(maxY-minY+1) + minY)
	return &Position{
		Username: name,
		X:        randomX,
		Y:        randomY,
	}
}

type LobbyCommand struct {
	LobbyName string `json:"lobby_name"`
}

type LobbiesNamesData struct {
	Type         string   `json:"type"`
	Names        []string `json:"lobby_names"`
	CurrentLobby string   `json:"current_lobby"`
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
