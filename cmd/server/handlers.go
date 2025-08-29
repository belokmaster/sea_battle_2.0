package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sea_battle/game"
	"strconv"
)

func gameStatusHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	response := struct {
		Game       *game.Game `json:"game"`
		SaveExists bool       `json:"save_exists"`
	}{
		Game:       currentGame,
		SaveExists: saveExists,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Ошибка кодирования состояния игры", http.StatusInternalServerError)
	}
}

func loadGameHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if !saveExists {
		http.Error(w, "Файл сохранения не найден.", http.StatusNotFound)
		return
	}

	loadedGame, err := game.LoadGame(saveFilename)
	if err != nil {
		http.Error(w, "Не удалось загрузить игру: "+err.Error(), http.StatusInternalServerError)
		return
	}

	currentGame = loadedGame

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Игра успешно загружена.",
	})
}

func newGameHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	currentGame = game.NewGame()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Новая игра успешно создана")
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
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if currentGame.CurrentPlayer != currentGame.Player1 {
		http.Error(w, "Сейчас не ваш ход", http.StatusForbidden)
		return
	}

	resultMessage, err := applyAbility(r.URL.Query())
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": resultMessage,
	})
}

func HandlerCoords(w http.ResponseWriter, r *http.Request) (int, int, error) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return 0, 0, fmt.Errorf("метод не разрешен")
	}

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
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	if currentGame.CurrentPlayer != currentGame.Player1 {
		http.Error(w, "Сейчас не ваш ход", http.StatusForbidden)
		return
	}

	query := r.URL.Query()
	var abilityResultMessage string

	if query.Has("ability_name") {
		var err error
		abilityResultMessage, err = applyAbility(query)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	x, y, err := HandlerCoords(w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var msg string
	result, msg, err := currentGame.HandleHumanTurn(x, y)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if currentGame.CurrentPlayer.MyBoard.AllShipSunk() {
		handleGameOver(w, result, abilityResultMessage, currentGame)
		currentGame = game.NewGame()
		return
	}

	if result == game.ResultMiss {
		currentGame.SwitchPlayer()
		for {
			result, err := currentGame.HandleComputerTurn()
			if err != nil {
				http.Error(w, "Ошибка в ходе бота: "+err.Error(), http.StatusInternalServerError)
				return
			}

			if currentGame.CurrentPlayer.EnemyBoard.AllShipSunk() {
				handleGameOver(w, result, abilityResultMessage, currentGame)
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
		"attack_result":  fmt.Sprintf("%v", result),
		"message":        msg,
		"ability_result": abilityResultMessage,
		"game_over":      false,
		"winner":         "",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleGameOver(w http.ResponseWriter, result game.AttackResult, abilityResultMessage string, currentGame *game.Game) {
	msg := fmt.Sprintf("Победа. Все корабли противника потоплены. Победитель: %s", currentGame.CurrentPlayer.Name)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"attack_result":  fmt.Sprintf("%v", result),
		"message":        msg,
		"ability_result": abilityResultMessage,
		"game_over":      true,
		"winner":         currentGame.CurrentPlayer.Name,
	})
}

func saveGameHandler(w http.ResponseWriter, r *http.Request) {
	gameMutex.Lock()
	defer gameMutex.Unlock()

	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	err := currentGame.SaveGame(saveFilename)
	if err != nil {
		http.Error(w, "Не удалось сохранить игру: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Игра успешно сохранена",
	})
}
