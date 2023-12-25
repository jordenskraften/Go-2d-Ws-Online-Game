package main

import (
	"log"
	"net/http"
	"time"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/exchanger"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/transport"
)

func main() {
	MyHub := hub.NewHub("first Hub")
	//------------
	MyExchanger := exchanger.NewExchanger(MyHub)
	MyExchanger.CreateLobby("lobby#1")
	MyExchanger.CreateLobby("lobby#2")
	MyExchanger.CreateLobby("lobby#3")
	log.Printf("%d, %d, %d", MyExchanger.Lobbies[0].Name, MyExchanger.Lobbies[1].Name, MyExchanger.Lobbies[2].Name)
	go testing(MyExchanger)
	//-----------
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		transport.ServeWs(MyHub, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}

func testing(ex *exchanger.Exchanger) {
	time.Sleep(5 * time.Second)

	lobby := ex.GetAnyLobby()
	conn := ex.Hub.GetAnyConnection()
	log.Println("================================")
	log.Println("лобби и конект " + lobby.Name + " " + conn.Name)
	if lobby != nil && conn != nil {
		ex.AddConnectionToLobby(conn, lobby)
		log.Println("удалось добавить любой конект в любое лобби")
		log.Println(lobby.Connections)
		log.Println("теперь удалим этот конект с лобби")
		lobby.RemoveConnection(conn.Name)
		log.Println(lobby.Connections)
		log.Println("================================")
	} else {
		log.Println("не удалось добавить любой конект в любое лобби")
	}
}
