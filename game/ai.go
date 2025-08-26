package game

import (
	"fmt"
	"math/rand"
)

func contains(points []Point, p Point) bool {
	for _, pt := range points {
		if pt.X == p.X && pt.Y == p.Y {
			return true
		}
	}
	return false
}

func (g *Game) handleComputerTurn() (AttackResult, error) {
	computer := g.CurrentPlayer
	enemyBoard := computer.EnemyBoard
	var targetPoint Point

	switch computer.state {
	case Searching:
		for {
			x, y := rand.Intn(10), rand.Intn(10)
			targetPoint = Point{X: x, Y: y}

			if contains(computer.allHits, targetPoint) || contains(computer.verifiedPoints, targetPoint) {
				continue
			}

			fmt.Printf("Режим поиска. Бот атакует клетку (%d, %d)\n", targetPoint.X, targetPoint.Y)
			break
		}

	case FinishingOff:
		targetFound := false
		for _, hitPoint := range computer.targetHits {
			directions := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

			for _, d := range directions {
				candidate := Point{X: hitPoint.X + d.X, Y: hitPoint.Y + d.Y}

				if candidate.X < 0 || candidate.X >= 10 || candidate.Y < 0 || candidate.Y >= 10 {
					continue
				}

				if contains(computer.allHits, candidate) || contains(computer.verifiedPoints, candidate) {
					continue
				}

				targetPoint = candidate
				targetFound = true
				fmt.Printf("Режим добивания. Бот атакует клетку (%d, %d)\n", targetPoint.X, targetPoint.Y)
				break
			}

			if targetFound {
				break
			}
		}

		if !targetFound {
			fmt.Println("Бот не нашел клетку корабля. Возврат в режим поиска.")
			computer.state = Searching
			for {
				x, y := rand.Intn(10), rand.Intn(10)
				targetPoint = Point{X: x, Y: y}
				if contains(computer.verifiedPoints, targetPoint) || contains(computer.allHits, targetPoint) {
					continue
				}
				break
			}
		}
	}

	result, err := enemyBoard.Attack(&targetPoint, computer)
	if err != nil {
		return ResultMiss, err
	}

	switch result {
	case ResultHit:
		computer.allHits = append(computer.allHits, targetPoint)
		computer.targetHits = append(computer.targetHits, targetPoint)
		computer.state = FinishingOff
	case ResultSunk:
		computer.allHits = append(computer.allHits, targetPoint)
		computer.targetHits = []Point{}
		computer.state = Searching
	case ResultMiss:
		computer.verifiedPoints = append(computer.verifiedPoints, targetPoint)
	}
	return result, nil
}

func (p *Player) ShipSunkBot(makedPoints []Point) {
	if p.Name == "Computer" {
		p.verifiedPoints = append(p.verifiedPoints, makedPoints...)
		fmt.Printf("Бот обновил данные после потопления корабля")
	}
}
