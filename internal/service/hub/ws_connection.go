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
	for {
		messageType, p, err := wsCon.conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		// Конвертировать байты в текст
		messageText := string(p)

		log.Printf("Тип сообщения: %d, Сообщение: %s\n", messageType, messageText)
		// Обработка полученного сообщения (messageType, messageText)
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
