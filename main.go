package main

import (
	"errors"
	"fmt"
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
	Grid  [10][10]int
	Ships []Ship
}

func newBoard() *Board {
	return &Board{
		Grid:  [10][10]int{},
		Ships: []Ship{},
	}
}

func main() {
	gameBoard := newBoard()

	ship1 := Ship{Size: 3, IsVertical: false}
	if err := gameBoard.placeShip(&ship1, Point{X: 0, Y: 2}); err != nil {
		fmt.Println("Ошибка:", err)
	}
	gameBoard.printField()

	ship2 := Ship{Size: 2, IsVertical: true}
	fmt.Println("Ставим корабль (2, вертикальный) в точку B1 (x:1, y:1), вплотную к первому. Должна быть ошибка.")
	if err := gameBoard.placeShip(&ship2, Point{X: 0, Y: 1}); err != nil {
		fmt.Println("Поймали ожидаемую ошибку:", err)
	}

	gameBoard.printField()
}

func (b *Board) printField() {
	fmt.Println()
	fmt.Print("    ")
	for c := 0; c < 10; c++ {
		fmt.Printf("%c  ", 'A'+c)
	}
	fmt.Println()

	for i := 0; i < 10; i++ {
		fmt.Printf("%2d  ", i+1)
		for j := 0; j < 10; j++ {
			if b.Grid[i][j] == 1 {
				fmt.Print("■  ")
			} else {
				fmt.Print(".  ")
			}
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
					if b.Grid[checkX][checkY] == 1 {
						return errors.New("корабль соприкасается или пересекается с другим")
					}
				}
			}
		}
	}

	for _, p := range shipPoints {
		b.Grid[p.X][p.Y] = 1
	}

	ship.Position = shipPoints
	b.Ships = append(b.Ships, *ship)

	return nil
}
