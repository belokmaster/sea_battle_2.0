package game

import (
	"fmt"
	"math/rand"
)

func (a *ArtilleryStrike) Apply(g *Game) string {
	enemyBoard := g.CurrentPlayer.EnemyBoard
	enemyLivesCells := []Point{}

	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			cellStatus := enemyBoard.Grid[i][j]
			if cellStatus == ShipCell {
				enemyLivesCells = append(enemyLivesCells, Point{X: i, Y: j})
			}
		}
	}

	if len(enemyLivesCells) == 0 {
		return "Нет целей для артиллерийского удара"
	}

	randomPointInd := rand.Intn(len(enemyLivesCells))
	randomPoint := enemyLivesCells[randomPointInd]
	enemyBoard.Attack(&randomPoint, g.CurrentPlayer)

	s := fmt.Sprintf("Артиллерийский удар нанесен по клетке (%d, %d)", randomPoint.X, randomPoint.Y)
	return s
}

func (a *ArtilleryStrike) Name() string {
	return "Артиллерийский удар"
}

func (s *Scanner) Apply(g *Game) string {
	return "Для использования сканера необходимо указать координаты"
}

func (s *Scanner) ApplyWithTarget(g *Game, target Point) string {
	enemyBoard := g.CurrentPlayer.EnemyBoard
	countShips := 0

	scanPoint := target
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			checkX, checkY := scanPoint.X+dx, scanPoint.Y+dy
			if checkX >= 0 && checkX < 10 && checkY >= 0 && checkY < 10 {
				currentPointStatus := enemyBoard.Grid[checkX][checkY]
				if currentPointStatus == ShipCell {
					countShips++
				}
			}
		}
	}

	str := fmt.Sprintf("Сканирование области 3x3 в точке (%d, %d)... Обнаружено %d сегментов кораблей.", scanPoint.X, scanPoint.Y, countShips)
	return str
}

func (s *Scanner) Name() string {
	return "Сканнер."
}

func (d *DoubleDamage) Apply(g *Game) string {
	g.CurrentPlayer.HasDoubleDamage = true
	s := "Следующая атака нанесет двойной урон!"
	return s
}

func (s *DoubleDamage) Name() string {
	return "Двойной урон."
}

func (p *Player) AddRandomAbility() {
	abilities := []Ability{&ArtilleryStrike{}, &Scanner{}, &DoubleDamage{}}
	rd := rand.Intn(len(abilities))
	ability := abilities[rd]
	p.Abilities = append(p.Abilities, ability)
}
