package game

import (
	"fmt"
)

func NewGameFromFile(filename string) *Game {
	game, err := LoadGame(filename)
	if err != nil {
		fmt.Println("Ошибка при попытке загрузить файл для начала игры:", err)
		return nil
	}
	return game
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
		Abilities:       []Ability{},
		HasDoubleDamage: false,
	}

	p2 := Player{
		Name:            "Computer",
		MyBoard:         board2,
		EnemyBoard:      board1,
		Abilities:       []Ability{},
		HasDoubleDamage: false,
		State:           Searching,
		AllHits:         []Point{},
		TargetHits:      []Point{},
		VerifiedPoints:  []Point{},
	}

	game := Game{
		Player1:       &p1,
		Player2:       &p2,
		CurrentPlayer: &p1,
	}

	return &game
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
