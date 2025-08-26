package game

import (
	"errors"
	"fmt"
	"math/rand"
)

type CellState int

const (
	EmptyCell CellState = iota
	ShipCell
	MissCell
	HitCell
)

type AttackResult int

const (
	ResultMiss AttackResult = iota
	ResultHit
	ResultSunk
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

type Board struct {
	Grid  [10][10]CellState
	Ships []Ship
}

func NewBoard() *Board {
	return &Board{
		Grid:  [10][10]CellState{},
		Ships: []Ship{},
	}
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
}

func (b *Board) AllShipSunk() bool {
	for i := range b.Ships {
		if !b.Ships[i].IsSunk {
			return false
		}
	}
	return true
}

func (b *Board) markSunkShip(ship *Ship, pl *Player) {
	ship.IsSunk = true

	points := ship.Position
	for _, p := range points {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				checkX, checkY := p.X+dx, p.Y+dy
				if checkX >= 0 && checkX < 10 && checkY >= 0 && checkY < 10 {
					if pl.Name == "Computer" {
						pl.verifiedPoints = append(pl.verifiedPoints, Point{X: checkX, Y: checkY})
					}

					if b.Grid[checkX][checkY] == EmptyCell {
						b.Grid[checkX][checkY] = MissCell
					}
				}
			}
		}
	}
}

func (b *Board) Attack(p *Point, attacker *Player) (AttackResult, error) {
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
					ship.Hits++

					if attacker.HasDoubleDamage {
						if ship.Hits < ship.Size {
							ship.Hits++
						}
						attacker.HasDoubleDamage = false
					}

					if ship.Hits >= ship.Size {
						b.markSunkShip(ship, attacker)
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
