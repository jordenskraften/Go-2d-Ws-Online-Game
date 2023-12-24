package main

import (
	"log"
	"net/http"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/lobby"
)

func main() {
	MyHub := hub.NewHub("first Hub")
	//------------
	lobby1 := lobby.NewLobby("lobby#1", MyHub)
	lobby2 := lobby.NewLobby("lobby#2", MyHub)
	lobby3 := lobby.NewLobby("lobby#3", MyHub)
	log.Printf("%d, %d, %d", lobby1.Name, lobby2.Name, lobby3.Name)
	//-----------
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		hub.ServeWs(MyHub, w, r)
	})
	log.Fatal(http.ListenAndServe(":8080", nil))

}
