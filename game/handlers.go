package game

import (
	"fmt"
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
