package main

import (
	"log"

	ws "github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/server"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/lobby"
)

func main() {
	hub := hub.NewHub("first Hub")
	log.Println(hub.Name)
	conn1 := &ws.WsConnection{Name: "client#1"}
	conn2 := &ws.WsConnection{Name: "client#2"}
	hub.AddConnection(conn1)
	hub.AddConnection(conn2)
	log.Println(hub.Connections[conn1.Name].Name)
	log.Println(hub.Connections[conn2.Name].Name)
	//------------
	lobby1 := lobby.NewLobby("lobby#1", hub)
	lobby1.AddConnection(conn1)
	lobby1.AddConnection(conn2)
	lobby2 := lobby.NewLobby("lobby#2", hub)
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
