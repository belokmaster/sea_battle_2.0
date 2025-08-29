package game

import (
	"fmt"
	"math/rand"
)

func (g *Game) HandleHumanTurn(x, y int) (AttackResult, string, error) {
	attackPoint := Point{X: x, Y: y}
	result, err := g.CurrentPlayer.EnemyBoard.Attack(&attackPoint, g.CurrentPlayer)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ResultMiss, "", err
	}

	var msg string
	switch result {
	case ResultHit:
		msg = "Попадание! Вы ходите еще раз"
	case ResultSunk:
		msg = "Корабль потоплен! Вы ходите еще раз и вам добавлена способность!"
		g.CurrentPlayer.AddRandomAbility()
	case ResultMiss:
		msg = "Промах! Ход переходит"
	}

	return result, msg, nil
}

func contains(points []Point, p Point) bool {
	for _, pt := range points {
		if pt.X == p.X && pt.Y == p.Y {
			return true
		}
	}
	return false
}

func (g *Game) findNextTarget() (Point, bool) {
	computer := g.CurrentPlayer
	enemyBoard := computer.EnemyBoard

	for _, hitPoint := range computer.TargetHits {
		directions := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}

		for _, dir := range directions {
			candidate := Point{X: hitPoint.X + dir.X, Y: hitPoint.Y + dir.Y}

			if candidate.X >= 0 && candidate.X < 10 && candidate.Y >= 0 && candidate.Y < 10 {
				cell := enemyBoard.Grid[candidate.X][candidate.Y]
				if cell != MissCell && cell != HitCell {
					return candidate, true
				}
			}
		}
	}

	var targetPoint Point
	for {
		x, y := rand.Intn(10), rand.Intn(10)
		targetPoint = Point{X: x, Y: y}
		if contains(computer.VerifiedPoints, targetPoint) || contains(computer.AllHits, targetPoint) {
			continue
		}
		break
	}
	return targetPoint, false
}

func (g *Game) HandleComputerTurn() (AttackResult, error) {
	computer := g.CurrentPlayer
	enemyBoard := computer.EnemyBoard
	var targetPoint Point
	var flagTarget bool = true

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
		targetPoint, flagTarget = g.findNextTarget()
		if flagTarget {
			fmt.Printf("Режим добивания. Бот атакует клетку (%d, %d)\n", targetPoint.X, targetPoint.Y)
		} else {
			fmt.Printf("Режим поиска. Бот атакует клетку (%d, %d)\n", targetPoint.X, targetPoint.Y)
		}
	}

	result, err := enemyBoard.Attack(&targetPoint, computer)
	if err != nil {
		return ResultMiss, err
	}

	switch result {
	case ResultHit:
		computer.AllHits = append(computer.AllHits, targetPoint)
		computer.State = FinishingOff

		if computer.State == FinishingOff && !flagTarget {
			computer.TargetHits = []Point{targetPoint} // новое добивание
		} else {
			computer.TargetHits = append(computer.TargetHits, targetPoint)
		}
	case ResultSunk:
		computer.AllHits = append(computer.AllHits, targetPoint)
		computer.TargetHits = []Point{}
		computer.State = Searching
	case ResultMiss:
		computer.VerifiedPoints = append(computer.VerifiedPoints, targetPoint)

		if computer.State == FinishingOff && !flagTarget {
			fmt.Println("Случайный выстрел в режиме добивания был промахом. Сброс цели")
			computer.TargetHits = []Point{} // забытие старой цели
			computer.State = Searching
		}
	}
	return result, nil
}

func (p *Player) ShipSunkBot(makedPoints []Point) {
	if p.Name == "Computer" {
		p.VerifiedPoints = append(p.VerifiedPoints, makedPoints...)
		fmt.Printf("Бот обновил данные после потопления корабля")
	}
}
