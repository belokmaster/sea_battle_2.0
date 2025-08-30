package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sea_battle/game"
	"strconv"
	"time"
)

func sendJSON(w http.ResponseWriter, data interface{}, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func sendJSONError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"Message": message})
}

func gameStatusHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()
	sendJSON(w, map[string]interface{}{"game": currentGame, "save_exists": saveExists}, http.StatusOK)
}

func loadGameHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		sendJSONError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if !saveExists {
		sendJSONError(w, "Файл сохранения не найден", http.StatusNotFound)
		return
	}

	loadedGame, err := game.LoadGame(saveFilename)
	if err != nil {
		sendJSONError(w, "Не удалось загрузить игру: "+err.Error(), http.StatusInternalServerError)
		return
	}

	currentGame = loadedGame
	sendJSON(w, map[string]string{"message": "Игра успешно загружена"}, http.StatusOK)
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()
	if r.Method != http.MethodPost {
		sendJSONError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}
	currentGame = game.NewGame()
	sendJSON(w, map[string]string{"message": "Новая игра успешно создана"}, http.StatusOK)
}

func applyAbility(query url.Values) (string, error) {
	abilityName := query.Get("ability_name")
	if abilityName == "" {
		return "", fmt.Errorf("параметр 'ability_name' обязателен")
	}

	var selectedAbility game.Ability
	abilityIndex := -1
	player := currentGame.Player1

	for i, ab := range player.Abilities {
		var name string
		switch ab.(type) {
		case *game.ArtilleryStrike:
			name = "artillery"
		case *game.Scanner:
			name = "scanner"
		case *game.DoubleDamage:
			name = "doubledamage"
		}
		if name == abilityName {
			selectedAbility = ab
			abilityIndex = i
			break
		}
	}

	if selectedAbility == nil {
		return "", fmt.Errorf("у вас нет такой способности или она не существует")
	}

	var resultMessage string
	if scanner, ok := selectedAbility.(*game.Scanner); ok {
		xStr, yStr := query.Get("ability_x"), query.Get("ability_y")
		if xStr == "" || yStr == "" {
			return "", fmt.Errorf("для сканера нужны координаты ability_x и ability_y")
		}
		x, _ := strconv.Atoi(xStr)
		y, _ := strconv.Atoi(yStr)
		resultMessage = scanner.ApplyWithTarget(currentGame, game.Point{X: x, Y: y})
	} else {
		resultMessage = selectedAbility.Apply(currentGame)
	}

	player.Abilities = append(player.Abilities[:abilityIndex], player.Abilities[abilityIndex+1:]...)
	return resultMessage, nil
}

func abilityHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		sendJSONError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if currentGame.CurrentPlayer.Name != "Player" {
		sendJSONError(w, "Сейчас не ваш ход", http.StatusForbidden)
		return
	}

	resultMessage, err := applyAbility(r.URL.Query())
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// тут юз sendJson
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": resultMessage,
	})
}

func HandlerCoords(w http.ResponseWriter, r *http.Request) (int, int, error) {
	query := r.URL.Query()
	xStr := query.Get("x")
	yStr := query.Get("y")

	if xStr == "" || yStr == "" {
		return 0, 0, fmt.Errorf("не заданы координаты x и y")
	}

	x, errX := strconv.Atoi(xStr)
	if errX != nil {
		return 0, 0, fmt.Errorf("координата x должна быть числом")
	}

	y, errY := strconv.Atoi(yStr)
	if errY != nil {
		return 0, 0, fmt.Errorf("координата y должна быть числом")
	}

	return x, y, nil
}

func attackHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		sendJSONError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if currentGame.CurrentPlayer.Name != "Player" {
		sendJSONError(w, "Сейчас не ваш ход", http.StatusForbidden)
		return
	}

	query := r.URL.Query()
	var abilityResultMessage string

	if query.Has("ability_name") {
		var err error
		abilityResultMessage, err = applyAbility(query)
		if err != nil {
			sendJSONError(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	x, y, err := HandlerCoords(w, r)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	var msg string
	result, markedPoints, msg, err := currentGame.HandleHumanTurn(x, y)
	if err != nil {
		sendJSONError(w, err.Error(), http.StatusBadRequest)
		return
	}

	if currentGame.Player2.MyBoard.AllShipSunk() {
		response := handlerGameOver(abilityResultMessage, currentGame.Player1)
		sendJSON(w, response, http.StatusOK)
		currentGame = game.NewGame()
		return
	}

	var computerMoves []map[string]interface{}

	if result == game.ResultMiss {
		currentGame.SwitchPlayer()
		for {
			time.Sleep(300 * time.Millisecond)

			compTarget, result, newlyMarked, err := currentGame.HandleComputerTurn()
			if err != nil {
				sendJSONError(w, "Ошибка в ходе бота: "+err.Error(), http.StatusInternalServerError)
				return
			}

			computerMoves = append(computerMoves, map[string]interface{}{
				"x":             compTarget.X,
				"y":             compTarget.Y,
				"result":        result,
				"marked_points": newlyMarked,
			})
			log.Printf("Ход компьютера: %+v, Результат: %v", compTarget, result)

			if currentGame.Player1.MyBoard.AllShipSunk() {
				response := handlerGameOver(abilityResultMessage, currentGame.Player2)
				sendJSON(w, response, http.StatusOK)
				currentGame = game.NewGame()
				return
			}

			if result == game.ResultMiss {
				msg = "Бот промахнулся. Теперь ваш ход"
				currentGame.SwitchPlayer()
				break
			}

			log.Printf("Ход компьютера: %v", result)
		}

	}

	response := map[string]interface{}{
		"message":        msg,
		"ability_result": abilityResultMessage,
		"game_over":      false,
		"winner":         "",
		"computer_moves": computerMoves,
		"human_move_result": map[string]interface{}{
			"x":             x,
			"y":             y,
			"result":        result,
			"marked_points": markedPoints,
		},
	}
	sendJSON(w, response, http.StatusOK)
}

func handlerGameOver(abilityResultMessage string, winner *game.Player) map[string]interface{} {
	msg := fmt.Sprintf("Игра окончена! Победитель: %s", winner.Name)
	return map[string]interface{}{
		"message":        msg,
		"ability_result": abilityResultMessage,
		"game_over":      true,
		"winner":         winner.Name,
	}
}

func saveGameHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		sendJSONError(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	err := currentGame.SaveGame(saveFilename)
	if err != nil {
		sendJSONError(w, "Не удалось сохранить игру: "+err.Error(), http.StatusInternalServerError)
		return
	}

	sendJSON(w, map[string]string{"message": "Игра успешно сохранена"}, http.StatusOK)
}
