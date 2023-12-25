// В пакете hub
package hub

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/entities"
)

type WsConnection struct {
	Name string
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Client struct {
	ID        string
	IPAddress string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Здесь можно реализовать свою логику проверки origin
		// Например, разрешить все запросы:
		return true
	},
}

// serveWs handles websocket requests from the peer.
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка установки соединения WebSocket:", err)
		return
	}
	defer conn.Close()

	wsConn := &WsConnection{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	wsConn.hub.AddConnection(wsConn)
	//надо добавить в лобби
	//потом для теста в лобби отправить месейдж чату и канвасу

	clientID := fmt.Sprintf("%d", time.Now().Unix())
	ipAddress := r.RemoteAddr
	client := Client{
		ID:        clientID,
		IPAddress: ipAddress,
	}
	log.Printf("Клиент с ID: %s, IP: %s, подключился в %s\n", client.ID, client.IPAddress, time.Now().Format("2006-01-02 15:04:05"))

	clientInfoJSON, err := json.Marshal(client)
	if err != nil {
		log.Println("Ошибка кодирования JSON:", err)
		return
	}

	err = conn.WriteMessage(websocket.TextMessage, clientInfoJSON)
	if err != nil {
		log.Println("Ошибка отправки информации о клиенте через WebSocket:", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		wsConn.ReadLoop()
	}()
	wg.Wait() // Ждем завершения работы горутины ReadLoop перед закрытием соединения
}

func (wsCon *WsConnection) ReadLoop() {
	defer wsCon.conn.Close()

	for {
		var msgData map[string]interface{}
		err := wsCon.conn.ReadJSON(&msgData)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			break
		}

		if msgType, ok := msgData["type"]; ok {
			switch msgType {
			case "ChatMessage":
				username, usernameExists := msgData["username"].(string)
				text, textExists := msgData["text"].(string)
				date, dateExists := msgData["date"].(string)

				if !usernameExists || !textExists || !dateExists {
					log.Println("Incomplete ChatMessage data")
					continue
				}

				message := entities.ChatMessage{
					Username: username,
					Text:     text,
					Date:     date,
				}
				log.Println("Received ChatMessage object:", message)

			case "Canvas":
				positionsData, ok := msgData["positions"].(map[string]interface{})
				if !ok {
					log.Println("Error obtaining position data for Canvas")
					continue
				}

				positions := make(map[string]entities.Position)
				for key, val := range positionsData {
					posData, ok := val.(map[string]interface{})
					if !ok {
						log.Println("Error processing position data for Canvas")
						continue
					}

					x, xOk := posData["x"].(float64)
					y, yOk := posData["y"].(float64)
					if !xOk || !yOk {
						log.Println("Error extracting coordinates for Canvas")
						continue
					}

					position := entities.Position{
						X: int(x),
						Y: int(y),
					}
					positions[key] = position
				}

				canvas := entities.Canvas{
					Positions: positions,
				}
				log.Println("Received Canvas object:", canvas)

			default:
				log.Println("Unknown message type")
			}
		}
	}
}

// Серверное приложение вызывает метод Upgrader.Upgrade из обработчика HTTP-запросов, чтобы получить *Conn:
//методы WriteMessage и ReadMessage соединения, чтобы отправлять и получать сообщения в виде фрагментов байтов
/*
for {
    messageType, p, err := conn.ReadMessage()
    if err != nil {
        log.Println(err)
        return
    }
    if err := conn.WriteMessage(messageType, p); err != nil {
        log.Println(err)
        return
    }
}
*/
