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
	MyConnectionsManages := transport.NewConnectionsManager(MyHub, MyExchanger)
	MyExchanger.CreateLobby("lobby#1")
	// MyExchanger.CreateLobby("lobby#2")
	// MyExchanger.CreateLobby("lobby#3")
	// log.Printf("%d, %d, %d", MyExchanger.Lobbies[0].Name, MyExchanger.Lobbies[1].Name, MyExchanger.Lobbies[2].Name)
	go testing(MyExchanger)
	//-----------
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		transport.ServeWs(MyHub, MyConnectionsManages, w, r)
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
		log.Printf("удалось добавить конект %s в любое %s", conn.Name, lobby.Name)
		log.Println(lobby.Connections)
		curLob := ex.GetUserLobby(conn)
		log.Println("текущее лобби клиента " + conn.Name + " " + curLob.Name)
		//--------
		//сменим лобби
		ex.ChangeUserLobby(conn, "lobby#1")
		curLob = ex.GetUserLobby(conn)
		log.Println(lobby.Connections)
		log.Println(curLob.Connections)
		log.Println("текущее лобби клиента " + conn.Name + " " + curLob.Name)
		//
		log.Println("теперь удалим этот конект с лобби")
		//ex.DeleteUserFromAllLobbies(conn)
		curLob.RemoveConnection(conn.Name)
		log.Println(lobby.Connections)
		//
		curLob = ex.GetUserLobby(conn)
		if curLob != nil {
			log.Println(curLob.Name)
		}
		log.Println("================================")
	} else {
		log.Println("не удалось добавить любой конект в любое лобби")
	}
}
