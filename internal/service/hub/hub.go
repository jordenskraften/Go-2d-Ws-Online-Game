package hub

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type ConnItem struct {
	Name string
	conn *websocket.Conn
	send chan []byte
}
type Hub struct {
	Name        string
	Connections map[string]*ConnItem
	mu          sync.Mutex
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
	log.Println("Список подключений на сервере после подключения:")
	str := ""
	for name := range h.Connections {
		str += name + ", "
	}
	log.Println(str)
	log.Println("-----------------------")

	//трайнем подключить челика к лобби
}

func (h *Hub) RemoveConnection(name string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.Connections, name)
	log.Println("Список подключений на сервере после отсоединения:")
	str := ""
	for name := range h.Connections {
		str += name + ", "
	}
	log.Println(str)
	log.Println("-----------------------")
}

func (h *Hub) GetConnectionByName(name string) *ConnItem {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, val := range h.Connections {
		if val.Name == name {
			log.Println("найден конект в хабе с именем " + name)
			return val
		}
	}
	return nil
}

func (h *Hub) GetAnyConnection() *ConnItem {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, val := range h.Connections {
		return val // Возвращаем первое доступное соединение
	}

	return nil // Возвращаем nil, если соединений нет
}
