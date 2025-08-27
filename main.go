package main

import (
	"bufio"
	"fmt"
	"os"
	"sea_battle/game"
)

func main() {
	const saveFile = "savegame.json"
	reader := bufio.NewReader(os.Stdin)

	_, err := os.Stat(saveFile)
	if os.IsNotExist(err) {
		gameInstance := game.NewGame()
		gameInstance.StartGame()
	} else if err != nil {
		fmt.Println("Ошибка при получении информации о файле:", err)
	} else {
		fmt.Println("Существует сохранение. Хотите ли вы загрузить файл? (y/n)")
		var ans string
		fmt.Scan(&ans)

		if ans == "y" || ans == "Y" {
			reader.ReadString('\n')
			gameInstance := game.NewGameFromFile(saveFile)
			if gameInstance != nil {
				gameInstance.StartGame()
			}
			fmt.Println("Критическая ошибка. Игра не инициализируется")
		} else {
			reader.ReadString('\n')
			gameInstance := game.NewGame()
			gameInstance.StartGame()
		}
	}
}
