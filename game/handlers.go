package game

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func (g *Game) handleHumanTurn() (AttackResult, error) {
	reader := bufio.NewReader(os.Stdin)

	if len(g.CurrentPlayer.Abilities) > 0 {
		fmt.Println("У вас есть способности. Хотите их использовать? (y/n)")
		inputBytes, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		ans := strings.TrimSpace(inputBytes)

		if ans == "y" || ans == "Y" {
			fmt.Println("Выберите способность, которую хотите использовать.")
			for i, ab := range g.CurrentPlayer.Abilities {
				fmt.Printf("%d: %s\n", i, ab.Name())
			}

			var n int
			for {
				fmt.Scan(&n)
				if n < 0 || n >= len(g.CurrentPlayer.Abilities) {
					fmt.Println("Некорректный ввод. Повторите еще раз")
					reader.ReadString('\n')
					continue
				}

				fmt.Printf("Вы выбрали способность: %s\n", g.CurrentPlayer.Abilities[n].Name())
				fmt.Println(g.CurrentPlayer.Abilities[n].Apply(g))
				g.CurrentPlayer.Abilities = append(
					g.CurrentPlayer.Abilities[:n],
					g.CurrentPlayer.Abilities[n+1:]...,
				)
				break
			}
		}
		reader.ReadString('\n')
	}

	var input string
	var x, y int
	for {
		fmt.Print("Введите координаты для атаки или команду 'save' для сохранения: ")
		inputBytes, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		input = strings.TrimSpace(inputBytes)

		if input == "save" {
			err := g.SaveGame("savegame.json")
			if err != nil {
				fmt.Println("Ошибка при попытке сохранить игру: ", err)
				fmt.Println("Повторите еще раз")
				continue
			}
			continue
		}

		n, err := fmt.Sscanf(input, "%d %d", &x, &y)
		if n != 2 || err != nil || x < 0 || x > 9 || y < 0 || y > 9 {
			fmt.Println("Некорректный ввод. Повторите еще раз")
			continue
		}

		break
	}

	attackPoint := Point{X: x, Y: y}
	result, err := g.CurrentPlayer.EnemyBoard.Attack(&attackPoint, g.CurrentPlayer)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ResultMiss, err
	}

	switch result {
	case ResultHit:
		fmt.Println("Попадание! Вы ходите еще раз.")
	case ResultSunk:
		fmt.Println("Корабль потоплен! Вы ходите еще раз и вам добавлена способность!")
		g.CurrentPlayer.AddRandomAbility()
	case ResultMiss:
		fmt.Println("Промах! Ход переходит.")
	}

	return result, nil
}
