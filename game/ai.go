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

func (g *Game) HandleComputerTurn() (AttackResult, error) {
	computer := g.CurrentPlayer
	enemyBoard := computer.EnemyBoard
	var targetPoint Point

	switch computer.State {
	case Searching:
		for {
			x, y := rand.Intn(10), rand.Intn(10)
			targetPoint = Point{X: x, Y: y}

			if contains(computer.AllHits, targetPoint) || contains(computer.VerifiedPoints, targetPoint) {
				continue
			}

			fmt.Printf("Режим поиска. Бот атакует клетку (%d, %d)\n", targetPoint.X, targetPoint.Y)
			break
		}

	case FinishingOff:
		targetFound := false
		for _, hitPoint := range computer.TargetHits {
			directions := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

			for _, d := range directions {
				candidate := Point{X: hitPoint.X + d.X, Y: hitPoint.Y + d.Y}

				if candidate.X < 0 || candidate.X >= 10 || candidate.Y < 0 || candidate.Y >= 10 {
					continue
				}

				if contains(computer.AllHits, candidate) || contains(computer.VerifiedPoints, candidate) {
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
			computer.State = Searching
			for {
				x, y := rand.Intn(10), rand.Intn(10)
				targetPoint = Point{X: x, Y: y}
				if contains(computer.VerifiedPoints, targetPoint) || contains(computer.AllHits, targetPoint) {
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
		computer.AllHits = append(computer.AllHits, targetPoint)
		computer.TargetHits = append(computer.TargetHits, targetPoint)
		computer.State = FinishingOff
	case ResultSunk:
		computer.AllHits = append(computer.AllHits, targetPoint)
		computer.TargetHits = []Point{}
		computer.State = Searching
	case ResultMiss:
		computer.VerifiedPoints = append(computer.VerifiedPoints, targetPoint)
	}
	return result, nil
}

func (p *Player) ShipSunkBot(makedPoints []Point) {
	if p.Name == "Computer" {
		p.VerifiedPoints = append(p.VerifiedPoints, makedPoints...)
		fmt.Printf("Бот обновил данные после потопления корабля")
	}
}
