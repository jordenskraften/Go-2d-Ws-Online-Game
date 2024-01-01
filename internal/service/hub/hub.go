package hub

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/entities"
)

type ConnItem struct {
	Name string
	conn *websocket.Conn
	send chan []byte
}
type Hub struct {
	Name        string
	Connections map[string]*ConnItem
	mu          sync.RWMutex
}

func NewHub(name string) *Hub {
	return &Hub{
		Name:        name,
		Connections: make(map[string]*ConnItem),
	}
}

func (h *Hub) AddConnection(name string, conn *websocket.Conn, send chan []byte) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Connections[name] = &ConnItem{
		Name: name,
		conn: conn,
		send: send,
	}
}

func (h *Hub) RemoveConnection(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Connections, name)
}

func (h *Hub) GetConnectionByName(name string) *ConnItem {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if conn, ok := h.Connections[name]; ok {
		return conn
	}
	return nil
}

func (h *Hub) IsConnectionInHub(name string) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.Connections[name]
	return ok
}

func (h *Hub) GetAnyConnection() *ConnItem {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, val := range h.Connections {
		return val // Возвращаем первое доступное соединение
	}

	return nil // Возвращаем nil, если соединений нет
}

// ----------
func (h *Hub) BroadcastChatMessageToUserList(userlist []*ConnItem, username string, text string, date string) {
	//надо собрать json объект для отправки с полем type ChatMessage из аргументов стрингов
	messagePayload := entities.ChatMessageData{
		Type:     "ChatMessageData",
		Username: username,
		Text:     text,
		Date:     date,
	}
	log.Println(messagePayload)

	messageJSON, err := json.Marshal(messagePayload)
	if err != nil {
		log.Println("Error marshalling JSON in hub broadcast chat message to userlist:", err)
		return
	}

	for _, user := range userlist {
		for _, conn := range h.Connections {
			if user.Name == conn.Name {
				log.Printf("found user %s and sending him chatMessage from hub broadcast:", conn.Name)
				//ну а тут происходит отправка если юзер в листе есть в конектах хаба с этим ником
				//conn.conn.WriteJSON(messageJSON)
				err := conn.conn.WriteMessage(websocket.TextMessage, messageJSON)
				if err != nil {
					log.Println("Error sending message to user:", err)
					continue
				}
			}
		}
	}
}

// ----------
func (h *Hub) BroadcastCanvasDataToUserList(userlist []*ConnItem, messagePayload entities.CanvasMessageData) {
	//надо собрать json объект для отправки с полем type ChatMessage из аргументов стрингов
	log.Println(messagePayload)

	messageJSON, err := json.Marshal(messagePayload)
	if err != nil {
		log.Println("Error marshalling JSON in hub broadcast chat message to userlist:", err)
		return
	}

	for _, user := range userlist {
		for _, conn := range h.Connections {
			if user.Name == conn.Name {
				log.Printf("found user %s and sending him chatMessage from hub broadcast:", conn.Name)
				//ну а тут происходит отправка если юзер в листе есть в конектах хаба с этим ником
				//conn.conn.WriteJSON(messageJSON)
				err := conn.conn.WriteMessage(websocket.TextMessage, messageJSON)
				if err != nil {
					log.Println("Error sending message to user:", err)
					continue
				}
			}
		}
	}
}

// -----------
func (h *Hub) SendLobbiesListToConnestion(conn *ConnItem, lobbies []string) {
	msg := entities.LobbiesNamesData{
		Type:  "LobbiesNamesData",
		Names: lobbies,
	}
	log.Println(msg)
	messageJSON, err := json.Marshal(msg)
	if err != nil {
		log.Println("Error marshalling JSON in hub SendLobbiesListToConnestion:", err)
		return
	}

	err = conn.conn.WriteMessage(websocket.TextMessage, messageJSON)
	if err != nil {
		log.Println("Error sending message to user:", err)
	}
	log.Println("SendLobbiesListToConnestion is done:", err)
}
