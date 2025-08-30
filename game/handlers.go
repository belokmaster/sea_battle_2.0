package game

import (
	"fmt"
	"math/rand"
)

func (g *Game) HandleHumanTurn(x, y int) (AttackResult, []Point, string, error) {
	attackPoint := Point{X: x, Y: y}
	result, markedPoints, err := g.CurrentPlayer.EnemyBoard.Attack(&attackPoint, g.CurrentPlayer)
	if err != nil {
		fmt.Println("Ошибка:", err)
		return ResultMiss, nil, "", err
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

	return result, markedPoints, msg, nil
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

	availableTargets := g.findAvailableTargets(computer)
	if len(availableTargets) > 0 {
		randInd := rand.Intn(len(availableTargets))
		return availableTargets[randInd], true
	}

	targetPoint := g.searchingNewTarget()
	return targetPoint, false
}

func (g *Game) findAvailableTargets(computer *Player) []Point {
	var availableTargets []Point

	attackedSet := make(map[Point]bool)
	for _, p := range computer.AllHits {
		attackedSet[p] = true
	}
	for _, p := range computer.VerifiedPoints {
		attackedSet[p] = true
	}

	for _, hitPoint := range computer.TargetHits {
		directions := []Point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for _, dir := range directions {
			candidate := Point{X: hitPoint.X + dir.X, Y: hitPoint.Y + dir.Y}
			if candidate.IsValidPoint() && !attackedSet[candidate] {
				availableTargets = append(availableTargets, candidate)
			}
		}
	}

	return availableTargets
}

func (g *Game) searchingNewTarget() Point {
	var targetPoint Point
	computer := g.CurrentPlayer
	for {
		x, y := rand.Intn(10), rand.Intn(10)
		targetPoint = Point{X: x, Y: y}

		if !contains(computer.AllHits, targetPoint) && !contains(computer.VerifiedPoints, targetPoint) {
			fmt.Printf("Режим поиска. Бот атакует клетку (%d, %d)\n", targetPoint.X, targetPoint.Y)
			break
		}
	}
	return targetPoint
}

func (g *Game) HandleComputerTurn() (Point, AttackResult, []Point, error) {
	computer := g.CurrentPlayer

	var targetPoint Point
	if computer.State == FinishingOff && len(computer.TargetHits) > 0 {
		targetPoint, _ = g.findNextTarget()
	} else {
		targetPoint = g.searchingNewTarget()
		computer.State = Searching
	}

	result, newlyMarkedPoints, err := computer.EnemyBoard.Attack(&targetPoint, computer)
	if err != nil {
		return targetPoint, ResultMiss, nil, err
	}

	switch result {
	case ResultHit:
		computer.AllHits = append(computer.AllHits, targetPoint)
		computer.TargetHits = append(computer.TargetHits, targetPoint)
		computer.State = FinishingOff

	case ResultSunk:
		computer.AllHits = append(computer.AllHits, targetPoint)
		g.CurrentPlayer.shipSunkBot(newlyMarkedPoints)

		computer.TargetHits = []Point{}
		computer.State = Searching

	case ResultMiss:
		computer.VerifiedPoints = append(computer.VerifiedPoints, targetPoint)

		if computer.State == FinishingOff && len(g.findAvailableTargets(computer)) == 0 {
			computer.TargetHits = []Point{}
			computer.State = Searching
		}
	}

	return targetPoint, result, newlyMarkedPoints, nil
}

func (p *Player) shipSunkBot(makedPoints []Point) {
	if p.Name == "Computer" {
		for _, mp := range makedPoints {
			if !contains(p.VerifiedPoints, mp) {
				p.VerifiedPoints = append(p.VerifiedPoints, mp)
			}
		}
	}
}

func (p Point) IsValidPoint() bool {
	return p.X >= 0 && p.X < 10 && p.Y >= 0 && p.Y < 10
}
