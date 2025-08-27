package game

import (
	"encoding/json"
	"fmt"
	"os"
)

type RawPlayer struct {
	Name            string
	MyBoard         *Board
	EnemyBoard      *Board
	AbilityNames    []string `json:"abilities"`
	HasDoubleDamage bool

	State          AIState `json:"state"`           // поведение ИИ
	AllHits        []Point `json:"all_hits"`        // все попадания
	TargetHits     []Point `json:"target_hits"`     // добиваемый корабль
	VerifiedPoints []Point `json:"verified_points"` // промахи
}

func (p *Player) MarshalJSON() ([]byte, error) {
	abilities := make([]string, len(p.Abilities))
	for i, ability := range p.Abilities {
		abilities[i] = ability.Name()
	}

	raw := RawPlayer{
		Name:            p.Name,
		MyBoard:         p.MyBoard,
		EnemyBoard:      p.EnemyBoard,
		AbilityNames:    abilities,
		HasDoubleDamage: p.HasDoubleDamage,

		State:          p.State,
		AllHits:        p.AllHits,
		TargetHits:     p.TargetHits,
		VerifiedPoints: p.VerifiedPoints,
	}

	return json.Marshal(raw)
}

func (p *Player) UnmarshalJSON(data []byte) error {
	var raw RawPlayer
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	p.Name = raw.Name
	p.MyBoard = raw.MyBoard
	p.EnemyBoard = raw.EnemyBoard
	p.HasDoubleDamage = raw.HasDoubleDamage

	p.State = raw.State
	p.AllHits = raw.AllHits
	p.TargetHits = raw.TargetHits
	p.VerifiedPoints = raw.VerifiedPoints

	for _, abilityName := range raw.AbilityNames {
		switch abilityName {
		case "Артиллерийский удар":
			p.Abilities = append(p.Abilities, &ArtilleryStrike{})
		case "Сканнер.":
			p.Abilities = append(p.Abilities, &Scanner{})
		case "Двойной урон.":
			p.Abilities = append(p.Abilities, &DoubleDamage{})
		}
	}

	return nil
}

func (g *Game) SaveGame(filename string) error {
	data, err := json.MarshalIndent(g, "", " ")
	if err != nil {
		fmt.Println("Ошибка при попытке сохранить игру:", err)
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Ошибка при записи файла:", err)
		return err
	}

	fmt.Println("Игра успешно сохранена в", filename)
	return nil
}

func LoadGame(filename string) (*Game, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Ошибка при попытке загрузить игру:", err)
		return nil, err
	}

	var game Game
	err = json.Unmarshal(data, &game)
	if err != nil {
		fmt.Println("Ошибка при попытке чтения файла:", err)
		return nil, err
	}

	game.Player1.EnemyBoard = game.Player2.MyBoard
	game.Player2.EnemyBoard = game.Player1.MyBoard

	if game.CurrentPlayer.Name == game.Player1.Name {
		game.CurrentPlayer = game.Player1
	} else {
		game.CurrentPlayer = game.Player2
	}

	fmt.Println("Игра успешно загружена")
	return &game, nil
}
