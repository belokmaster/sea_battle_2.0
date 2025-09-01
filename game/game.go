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
	playerBoard := NewBoard()
	computerBoard := NewBoard()
	playerBoard.PlaceBoard()
	computerBoard.PlaceBoard()

	p1 := Player{
		Name:            "Player",
		MyBoard:         playerBoard,
		EnemyBoard:      computerBoard,
		Abilities:       []Ability{},
		HasDoubleDamage: false,
	}

	p2 := Player{
		Name:            "Computer",
		MyBoard:         computerBoard,
		EnemyBoard:      playerBoard,
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

func NewGameManual(playerBoard *Board) *Game {
	computerBoard := NewBoard()
	computerBoard.PlaceBoard()

	p1 := Player{
		Name:            "Player",
		MyBoard:         playerBoard,
		EnemyBoard:      computerBoard,
		Abilities:       []Ability{},
		HasDoubleDamage: false,
	}

	p2 := Player{
		Name:            "Computer",
		MyBoard:         computerBoard,
		EnemyBoard:      playerBoard,
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

func (g *Game) SwitchPlayer() {
	if g.CurrentPlayer == g.Player1 {
		g.CurrentPlayer = g.Player2
	} else {
		g.CurrentPlayer = g.Player1
	}
}
