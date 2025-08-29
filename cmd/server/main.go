package main

import (
	"fmt"
	"log"
	"net/http"
	"sea_battle/game"
	"sync"
)

var currentGame *game.Game
var gameMutex = &sync.Mutex{}

func main() {
	currentGame = game.NewGame()

	router := newRouter()

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(port, corsMiddleware(router)))
}
