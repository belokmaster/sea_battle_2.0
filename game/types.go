package game

type CellState int

const (
	EmptyCell CellState = iota
	ShipCell
	MissCell
	HitCell
)

type AttackResult int

const (
	ResultMiss AttackResult = iota
	ResultHit
	ResultSunk
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
	Grid  [10][10]CellState
	Ships []Ship
}

type AttackObserver interface {
	ShipSunk(newlyMarkedPoints []Point)
}

type Player struct {
	Name            string
	MyBoard         *Board
	EnemyBoard      *Board
	Abilities       []Ability
	HasDoubleDamage bool

	State          AIState `json:"state"`           // поведение ИИ
	AllHits        []Point `json:"all_hits"`        // все попадания
	TargetHits     []Point `json:"target_hits"`     // добиваемый корабль
	VerifiedPoints []Point `json:"verified_points"` // промахи
}

type Game struct {
	Player1       *Player
	Player2       *Player
	CurrentPlayer *Player
}

type AIState int

const (
	Searching    AIState = iota // поиск случайной цели
	FinishingOff                // добивание подбитого корабля
)

type Ability interface {
	Apply(g *Game) string
	Name() string
}

type ArtilleryStrike struct{}
type Scanner struct{}
type DoubleDamage struct{}
