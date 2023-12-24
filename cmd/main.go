package main

import (
	"log"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/lobby"
)

func main() {
	MyHub := hub.NewHub("first Hub")
	log.Println(MyHub.Name)
	conn1 := &hub.WsConnection{Name: "client#1"}
	conn2 := &hub.WsConnection{Name: "client#2"}
	MyHub.AddConnection(conn1)
	MyHub.AddConnection(conn2)
	log.Println(MyHub.Connections[conn1.Name].Name)
	log.Println(MyHub.Connections[conn2.Name].Name)
	//------------
	lobby1 := lobby.NewLobby("lobby#1", MyHub)
	lobby1.AddConnection(conn1)
	lobby1.AddConnection(conn2)
	lobby2 := lobby.NewLobby("lobby#2", MyHub)
	lobby2.AddConnection(conn1)
	log.Println("Лобби 1:", lobby1.Name)
	connections1 := ""
	for name := range lobby1.Connections {
		connections1 += name + ", "
	}
	log.Println("Имена соединений в мапе Лобби 1:", connections1)

	log.Println("Лобби 2:", lobby2.Name)
	connections2 := ""
	for name := range lobby2.Connections {
		connections2 += name + ", "
	}
	log.Println("Имена соединений в мапе Лобби 2:", connections2)

}
