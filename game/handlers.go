package game

import "fmt"

func (g *Game) handleHumanTurn() (AttackResult, error) {
	if len(g.CurrentPlayer.Abilities) > 0 {
		fmt.Println("Хотите использовать способность? (y/n)")
		var ans string
		fmt.Scan(&ans)

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
	}

	var x, y int
	fmt.Print("Введите координаты для атаки: ")
	for {
		n, err := fmt.Scan(&x, &y)
		if n != 2 || err != nil || x < 0 || x > 9 || y < 0 || y > 9 {
			fmt.Println("Некорректный ввод. Повторите еще раз.")
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
