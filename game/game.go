package game

import (
	"fmt"
)

type Player struct {
	Name            string
	MyBoard         *Board
	EnemyBoard      *Board
	Abilities       []Ability
	HasDoubleDamage bool
}

type Game struct {
	Player1       *Player
	Player2       *Player
	CurrentPlayer *Player
}

func NewGame() *Game {
	board1 := NewBoard()
	board2 := NewBoard()
	board1.PlaceBoard()
	board2.PlaceBoard()

	p1 := Player{
		Name:            "Player",
		MyBoard:         board1,
		EnemyBoard:      board2,
		Abilities:       []Ability{&ArtilleryStrike{}, &Scanner{}, &DoubleDamage{}},
		HasDoubleDamage: false,
	}

	p2 := Player{
		Name:            "Computer",
		MyBoard:         board2,
		EnemyBoard:      board1,
		Abilities:       []Ability{&ArtilleryStrike{}, &Scanner{}, &DoubleDamage{}},
		HasDoubleDamage: false,
	}

	game := Game{
		Player1:       &p1,
		Player2:       &p2,
		CurrentPlayer: &p1,
	}

	return &game
}

func (g *Game) StartGame() {
CurrentGameLoop:
	for {
		fmt.Printf("Ход игрока: %s\n", g.CurrentPlayer.Name)
		if g.CurrentPlayer.Name == "Player" {
			fmt.Println("Ваше поле:")
			g.Player1.MyBoard.printField(false)

			fmt.Println("Поле противника:")
			g.Player1.EnemyBoard.printField(true)
		}

		var result AttackResult
		var err error

		if g.CurrentPlayer == g.Player1 {
			abilities := g.CurrentPlayer.Abilities
			if len(abilities) > 0 {
				fmt.Println("Хотите использовать способность? (y/n)")
				var ans string
				fmt.Scan(&ans)

				if ans == "y" || ans == "Y" {
					fmt.Println("Выберите способность, которую хотите использовать.")
					fmt.Println("Ваши способности:")
					for i, ab := range abilities {
						fmt.Printf("%d: %s\n", i, ab.Name())
					}
					var n int
					for {
						fmt.Scan(&n)
						if n >= len(abilities) || n < 0 {
							fmt.Println("Неккоректный ввод. Повторите еще раз")
							continue
						} else {
							s := fmt.Sprintf("Вы выбрали способность: %s", abilities[n].Name())
							fmt.Println(s)
							abilityResult := abilities[n].Apply(g)
							fmt.Println(abilityResult)
							g.CurrentPlayer.Abilities = append(g.CurrentPlayer.Abilities[:n], g.CurrentPlayer.Abilities[n+1:]...)
							break
						}
					}
				}
			}

			fmt.Print("Введите координаты для атаки: ")
			var x, y int

			n, err := fmt.Scan(&x, &y)
			if n != 2 || err != nil {
				fmt.Println("Некорректный ввод. Пожалуйста, введите два числа через пробел.")
				continue
			}

			attackPoint := Point{X: x, Y: y}
			result, err = g.CurrentPlayer.EnemyBoard.Attack(&attackPoint, g.CurrentPlayer)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				switch result {
				case ResultHit:
					fmt.Println("Попадание! Вы ходите еще раз.")
				case ResultSunk:
					fmt.Println("Корабль потоплен! Вы ходите еще раз и вам добавлена способность!")
					g.CurrentPlayer.AddRandomAbility()
				case ResultMiss:
					fmt.Println("Промах! Ход переходит.")
				}
			}
		} else {
			result, err = g.CurrentPlayer.EnemyBoard.AttackBot(g.Player2)
			if err != nil {
				fmt.Println("Ошибка в ходе бота:", err)
			}
		}

		if g.CurrentPlayer.EnemyBoard.AllShipSunk() {
			fmt.Printf("Победил игрок: %s\n", g.CurrentPlayer.Name)
			break CurrentGameLoop
		}

		if result == ResultMiss {
			if g.CurrentPlayer == g.Player1 {
				g.CurrentPlayer = g.Player2
			} else {
				g.CurrentPlayer = g.Player1
			}
		}
	}
}
