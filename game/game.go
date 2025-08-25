package game

import "fmt"

type Player struct {
	Name       string
	MyBoard    *Board
	EnemyBoard *Board
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
		Name:       "Player",
		MyBoard:    board1,
		EnemyBoard: board2,
	}

	p2 := Player{
		Name:       "Computer",
		MyBoard:    board2,
		EnemyBoard: board1,
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
			fmt.Print("Введите координаты для атаки: ")
			var x, y int

			n, err := fmt.Scan(&x, &y)
			if n != 2 || err != nil {
				fmt.Println("Некорректный ввод. Пожалуйста, введите два числа через пробел.")
				continue
			}

			attackPoint := Point{X: x, Y: y}
			result, err = g.CurrentPlayer.EnemyBoard.Attack(&attackPoint)
			if err != nil {
				fmt.Println("Ошибка:", err)
			} else {
				switch result {
				case ResultHit:
					fmt.Println("Попадание! Вы ходите еще раз.")
				case ResultSunk:
					fmt.Println("Корабль потоплен! Вы ходите еще раз.")
				case ResultMiss:
					fmt.Println("Промах! Ход переходит.")
				}
			}
		} else {
			result, err = g.CurrentPlayer.EnemyBoard.AttackBot()
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
