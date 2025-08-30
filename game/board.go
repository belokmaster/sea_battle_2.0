package game

import (
	"errors"
	"math/rand"
)

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

func (b *Board) markSunkShip(ship *Ship) []Point {
	points := ship.Position
	markedCells := []Point{}
	for _, p := range points {
		for dx := -1; dx <= 1; dx++ {
			for dy := -1; dy <= 1; dy++ {
				checkX, checkY := p.X+dx, p.Y+dy
				if checkX >= 0 && checkX < 10 && checkY >= 0 && checkY < 10 {
					markedCells = append(markedCells, Point{X: checkX, Y: checkY})

					if b.Grid[checkX][checkY] == EmptyCell {
						b.Grid[checkX][checkY] = MissCell
					}
				}
			}
		}
	}
	return markedCells
}

func (b *Board) Attack(p *Point, attacker *Player) (AttackResult, []Point, error) {
	if p.X < 0 || p.X >= 10 || p.Y < 0 || p.Y >= 10 {
		return ResultMiss, nil, errors.New("атака вне поля")
	}

	currentSquare := b.Grid[p.X][p.Y]
	if currentSquare == MissCell || currentSquare == HitCell {
		return ResultMiss, nil, errors.New("по этой клетке уже стреляли")
	}

	if currentSquare == EmptyCell {
		b.Grid[p.X][p.Y] = MissCell
		return ResultMiss, nil, nil
	}

	b.Grid[p.X][p.Y] = HitCell
	for i := range b.Ships {
		ship := &b.Ships[i]

		isTargetShip := false
		for j := 0; j < len(ship.Position); j++ {
			if p.X == ship.Position[j].X && p.Y == ship.Position[j].Y {
				isTargetShip = true
				break
			}
		}

		if isTargetShip {
			ship.Hits++

			if attacker.HasDoubleDamage {
				// if ship.Hits < ship.Size {
				// 	ship.Hits++
				// } // изменить логику
				attacker.HasDoubleDamage = false
			}

			if ship.Hits >= ship.Size {
				ship.IsSunk = true
				markedCells := b.markSunkShip(ship)

				if observer, ok := interface{}(attacker).(AttackObserver); ok {
					observer.ShipSunk(markedCells)
				}

				return ResultSunk, markedCells, nil
			}
			return ResultHit, nil, nil
		}
	}

	return ResultMiss, nil, errors.New("ошибка состояния: клетка корабля есть, а самого корабля нет")
}
