package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sea_battle/game"
	"sync"
)

var currentGame *game.Game
var gameMutex = &sync.Mutex{}
var saveExists bool

const saveFilename = "savegame.json"

func main() {
	_, err := os.Stat(saveFilename)
	saveExists = !os.IsNotExist(err)

	if saveExists {
		fmt.Println("Найден файл сохранения. Сервер готов к загрузке по запросу")
	} else {
		fmt.Println("Файл сохранения не найден")
	}

	currentGame = game.NewGame()

	router := newRouter()

	port := ":8080"
	fmt.Printf("Сервер запущен на порту %s\n", port)
	log.Fatal(http.ListenAndServe(port, corsMiddleware(router)))
}
