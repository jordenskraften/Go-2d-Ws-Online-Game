package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/exchanger"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/hub"
	"github.com/jordenskraften/Go-2d-Ws-Online-Game/internal/service/transport"
)

func main() {
	MyHub := hub.NewHub("first Hub")
	//------------
	MyExchanger := exchanger.NewExchanger(MyHub)
	MyConnectionsManages := transport.NewConnectionsManager(MyHub, MyExchanger)
	MyExchanger.CreateLobby("start at #1 lobby")
	MyExchanger.CreateLobby("#2 lobby another")
	MyExchanger.CreateLobby("new #3 lobby")
	//-----------

	//-----------
	mux := http.NewServeMux()

	// Обработчик для веб-сокетов
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		transport.ServeWs(MyHub, MyConnectionsManages, w, r)
	})

	// Путь до папки, содержащей index.html
	basePath := "../frontend/"
	envBasePath := os.Getenv("FRONTEND_PATH")
	if envBasePath != "" {
		basePath = envBasePath
	}
	log.Printf("frontendPath: %s", basePath)

	// Создаем обработчик для статических файлов
	fs := http.FileServer(http.Dir(basePath))

	// Определяем путь для обработчика файлов
	mux.Handle("/", http.StripPrefix("/", fs))

	// Запуск сервера
	port := ":8080"
	log.Fatal(http.ListenAndServe(port, mux))
}
