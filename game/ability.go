package game

import (
	"errors"
	"fmt"
	"math/rand"
)

func (a *ArtilleryStrike) Apply(g *Game, target *Point) (*AbilityResult, error) {
	enemyBoard := g.CurrentPlayer.EnemyBoard
	availableTargets := []Point{}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			cellStatus := enemyBoard.Grid[i][j]
			if cellStatus == ShipCell || cellStatus == EmptyCell {
				availableTargets = append(availableTargets, Point{X: i, Y: j})
			}
		}
	}

	if len(availableTargets) == 0 {
		return &AbilityResult{Message: "Нет целей для артиллерийского удара"}, nil
	}

	randomPointInd := rand.Intn(len(availableTargets))
	randomPoint := availableTargets[randomPointInd]

	result, markedPoints, err := enemyBoard.Attack(&randomPoint, g.CurrentPlayer)
	if err != nil {
		return nil, fmt.Errorf("ошибка при использовании артиллерийского удара: %w", err)
	}

	if result == ResultSunk {
		g.CurrentPlayer.AddRandomAbility()
	}

	msg := fmt.Sprintf("Артиллерийский удар нанесен по (%d, %d)", randomPoint.X, randomPoint.Y)
	return &AbilityResult{
		Message: msg,
		AttackResult: &AttackResultData{
			Target:       randomPoint,
			Result:       result,
			MarkedPoints: markedPoints,
		},
	}, nil
}

func (a *ArtilleryStrike) Name() string {
	return "Артиллерийский удар"
}

func (a *ArtilleryStrike) RequiresTarget() bool {
	return false
}

func (s *Scanner) Apply(g *Game, target *Point) (*AbilityResult, error) {
	if target == nil {
		return nil, errors.New("для сканера необходимо указать координаты")
	}

	enemyBoard := g.CurrentPlayer.EnemyBoard
	countShips := 0
	var affectedPoints []Point

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			checkX, checkY := target.X+dx, target.Y+dy
			candidate := Point{X: checkX, Y: checkY}
			if candidate.IsValidPoint() {
				affectedPoints = append(affectedPoints, candidate)
				if enemyBoard.Grid[checkX][checkY] == ShipCell {
					countShips++
				}
			}
		}
	}

	msg := fmt.Sprintf("Сканирование области 3x3 в точке (%d, %d). Обнаружено %d сегментов кораблей", target.X, target.Y, countShips)
	return &AbilityResult{
		Message:        msg,
		AffectedPoints: affectedPoints,
	}, nil
}

func (s *Scanner) Name() string {
	return "Сканнер"
}

func (s *Scanner) RequiresTarget() bool {
	return true
}

func (d *DoubleDamage) Apply(g *Game, target *Point) (*AbilityResult, error) {
	g.CurrentPlayer.HasDoubleDamage = true
	return &AbilityResult{Message: "Следующая атака нанесет двойной урон!"}, nil
}

func (s *DoubleDamage) Name() string {
	return "Двойной урон"
}

func (d *DoubleDamage) RequiresTarget() bool {
	return false
}

func (p *Player) AddRandomAbility() {
	abilities := []Ability{&ArtilleryStrike{}, &Scanner{}, &DoubleDamage{}}
	rd := rand.Intn(len(abilities))
	ability := abilities[rd]
	p.Abilities = append(p.Abilities, ability)
}
