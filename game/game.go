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

func (g *Game) handleComputerTurn() (AttackResult, error) {
	result, err := g.CurrentPlayer.EnemyBoard.AttackBot(g.Player2)
	return result, err
}

func (g *Game) switchPlayer() {
	if g.CurrentPlayer == g.Player1 {
		g.CurrentPlayer = g.Player2
	} else {
		g.CurrentPlayer = g.Player1
	}
}

func (g *Game) displayTurnInfo() {
	fmt.Printf("Ход игрока: %s\n", g.CurrentPlayer.Name)
	if g.CurrentPlayer.Name == "Player" {
		fmt.Println("Ваше поле:")
		g.Player1.MyBoard.printField(false)

		fmt.Println("Поле противника:")
		g.Player1.EnemyBoard.printField(false) // потом на фолс поменть
	}
}

func (g *Game) StartGame() {
CurrentGameLoop:
	for {
		g.displayTurnInfo()

		var result AttackResult
		var err error

		if g.CurrentPlayer == g.Player1 {
			result, err = g.handleHumanTurn()
		} else {
			result, err = g.handleComputerTurn()
		}

		if err != nil {
			fmt.Println("Ошибка во время хода.", err)
		}

		if g.CurrentPlayer.EnemyBoard.AllShipSunk() {
			fmt.Printf("Победил игрок: %s\n", g.CurrentPlayer.Name)
			break CurrentGameLoop
		}

		if result == ResultMiss {
			g.switchPlayer()
		}
	}
}
