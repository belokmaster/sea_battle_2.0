package main

import (
	"errors"
	"fmt"
	"math/rand"
)

type Point struct {
	X int
	Y int
}

type Ship struct {
	Size       int
	IsVertical bool
	Hits       int
	IsSunk     bool
	Position   []Point
}

type CellState int

const (
	EmptyCell CellState = iota
	ShipCell
	MissCell
	HitCell
)

type Board struct {
	Grid  [10][10]CellState
	Ships []Ship
}

type AttackResult int

const (
	ResultMiss AttackResult = iota
	ResultHit
	ResultSunk
)

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

func newBoard() *Board {
	return &Board{
		Grid:  [10][10]CellState{},
		Ships: []Ship{},
	}
}

func NewGame() *Game {
	board1 := newBoard()
	board2 := newBoard()
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
			fmt.Scan(&x, &y)
			attackPoint := Point{X: x, Y: y}
			result, err = g.CurrentPlayer.EnemyBoard.Attack(&attackPoint)
			if err != nil {
				fmt.Println("Ошибка:", err)
			}

			if g.CurrentPlayer.EnemyBoard.AllShipSunk() {
				fmt.Printf("Победил игрок: %s\n", g.CurrentPlayer.Name)
				break CurrentGameLoop
			}
		} else {
			_, _ = g.CurrentPlayer.EnemyBoard.AttackBot()
			if g.CurrentPlayer.EnemyBoard.AllShipSunk() {
				fmt.Print("Победил бот", g.CurrentPlayer.Name)
				break CurrentGameLoop
			}
		}

		if g.CurrentPlayer.EnemyBoard.AllShipSunk() {
			fmt.Printf("Победил игрок: %s\n", g.CurrentPlayer.Name)
			break CurrentGameLoop
		}

		if result == ResultMiss {
			fmt.Println("Промах! Ход переходит.")
			if g.CurrentPlayer == g.Player1 {
				g.CurrentPlayer = g.Player2
			} else {
				g.CurrentPlayer = g.Player1
			}
		} else {
			fmt.Println("Попадание! Вы ходите еще раз.")
		}
	}
}

func (b *Board) AttackBot() (AttackResult, error) {
	var currentPoint Point
	for {
		x := rand.Intn(10)
		y := rand.Intn(10)
		if b.Grid[x][y] != MissCell && b.Grid[x][y] != HitCell {
			currentPoint = Point{X: x, Y: y}
			fmt.Printf("Бот атакует клетку: %d %d\n", currentPoint.X, currentPoint.Y)
			break
		}
	}
	return b.Attack(&currentPoint)
}

func (b *Board) PlaceBoard() {
	possibleShips := [10]int{4, 3, 3, 2, 2, 2, 1, 1, 1, 1}
	for _, size := range possibleShips {
		for {
			ship := Ship{
				Size:       size,
				IsVertical: (rand.Intn(2) == 1),
			}

			x := rand.Intn(10)
			y := rand.Intn(10)
			startPoint := Point{X: x, Y: y}
			if b.placeShip(&ship, startPoint) == nil {
				break
			}
		}
	}
	b.printField(false)
}

func (b *Board) AllShipSunk() bool {
	for i := range b.Ships {
		if !b.Ships[i].IsSunk {
			return false
		}
	}
	return true
}

func (b *Board) Attack(p *Point) (AttackResult, error) {
	if p.X < 0 || p.X >= 10 || p.Y < 0 || p.Y >= 10 {
		return ResultMiss, errors.New("атака вне поля")
	}

	currentSquare := b.Grid[p.X][p.Y]
	if currentSquare == MissCell || currentSquare == HitCell {
		return ResultMiss, errors.New("по этой клетке уже стреляли")
	} else if currentSquare == EmptyCell {
		b.Grid[p.X][p.Y] = MissCell
		return ResultMiss, nil
	} else {
		b.Grid[p.X][p.Y] = HitCell
		for i := range b.Ships {
			ship := &b.Ships[i]
			for j := 0; j < ship.Size; j++ {
				x, y := ship.Position[j].X, ship.Position[j].Y
				if p.X == x && p.Y == y {
					fmt.Printf("Удар по клетке %d %d нанесен\n", p.X, p.Y)
					ship.Hits++

					if ship.Hits == ship.Size {
						fmt.Println("Корабль был разрушен")
						ship.IsSunk = true
						return ResultSunk, nil
					}
					return ResultHit, nil
				}
			}
		}
		return ResultMiss, errors.New("ошибка состояния: клетка корабля есть, а самого корабля нет")
	}
}

func (b *Board) printField(isEnemyView bool) {
	fmt.Println()
	fmt.Print("    ")
	for c := 0; c < 10; c++ {
		fmt.Printf("%c  ", 'A'+c)
	}
	fmt.Println()

	for i := 0; i < 10; i++ {
		fmt.Printf("%2d  ", i)
		for j := 0; j < 10; j++ {
			cell := b.Grid[i][j]
			char := ""
			switch cell {
			case HitCell:
				char = "X  "
			case MissCell:
				char = "o  "
			case EmptyCell:
				char = ".  "
			case ShipCell:
				if isEnemyView {
					char = ".  "
				} else {
					char = "■  "
				}
			default:
				char = "?  "
			}
			fmt.Print(char)
		}
		fmt.Println()
	}
	fmt.Println()
}

func (b *Board) placeShip(ship *Ship, startPoint Point) error {
	shipPoints := make([]Point, ship.Size)
	for i := 0; i < ship.Size; i++ {
		if ship.IsVertical {
			shipPoints[i] = Point{X: startPoint.X + i, Y: startPoint.Y}
		} else {
			shipPoints[i] = Point{X: startPoint.X, Y: startPoint.Y + i}
		}
	}

	for _, p := range shipPoints {
		if p.X < 0 || p.X >= 10 || p.Y < 0 || p.Y >= 10 {
			return errors.New("корабль выходит за пределы поля")
		}

		for dx := -1; dx <= 1; dx++ { // проверка 3x3 квадрата вокруг точки корабля
			for dy := -1; dy <= 1; dy++ {
				checkX, checkY := p.X+dx, p.Y+dy
				if checkX >= 0 && checkX < 10 && checkY >= 0 && checkY < 10 {
					if b.Grid[checkX][checkY] == ShipCell {
						return errors.New("корабль соприкасается или пересекается с другим")
					}
				}
			}
		}
	}

	for _, p := range shipPoints {
		b.Grid[p.X][p.Y] = ShipCell
	}

	ship.Position = shipPoints
	b.Ships = append(b.Ships, *ship)

	return nil
}

func main() {
	game := NewGame()
	game.StartGame()
}
