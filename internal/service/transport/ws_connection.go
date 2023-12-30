// ws это транспортный уровень по сути между клиентом и бизнес логикой
// но пусть побудет пока тут :(
package transport

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
)

type WsConnection struct {
	Name string
	hub  *hub.Hub
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
func ServeWs(hub *hub.Hub, cm *ConnectionsManager, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Ошибка установки соединения WebSocket:", err)
		return
	}
	defer conn.Close()

	wsConn := &WsConnection{
		Name: generateRandomName(),
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
	}
	//вот этот метод переделать и структурку в хабе и гуд будет
	wsConn.hub.AddConnection(wsConn.Name, wsConn.conn, wsConn.send)
	defer wsConn.hub.RemoveConnection(wsConn.Name)
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

	//добавим юзера в стартовое лобби
	cm.SetupUser(wsConn.Name, *wsConn.hub)
	defer cm.RemoveUser(wsConn.Name, *wsConn.hub)

	//цикл прослушки
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		wsConn.ReadLoop(cm)
	}()
	wg.Wait() // Ждем завершения работы горутины ReadLoop перед закрытием соединения
}

func (wsCon *WsConnection) ReadLoop(cm *ConnectionsManager) {
	defer wsCon.conn.Close()

	for {
		var msgData map[string]interface{}
		err := wsCon.conn.ReadJSON(&msgData)
		if err != nil {
			log.Printf("Error reading JSON: %v", err)
			break
		} else {
			cm.DistributeMessage(wsCon.Name, msgData)
		}
	}
}

func generateRandomName() string {
	rand.Seed(time.Now().UnixNano())

	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	nameLength := 7
	nameBytes := make([]byte, nameLength)

	for i := range nameBytes {
		nameBytes[i] = letters[rand.Intn(len(letters))]
	}

	return string(nameBytes)
}
