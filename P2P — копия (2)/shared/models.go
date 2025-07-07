package shared

type Card struct {
	ID    int
	Value string
	Image string
}

type GameLevel int

const (
	Easy GameLevel = iota
	Medium
	Hard
	Expert
)

type GameRequest struct {
	PlayerName string
	Level      GameLevel
}

type GameResponse struct {
	ID      string
	Cards   []Card
	Level   GameLevel
	Rows    int
	Columns int
}

type Move struct {
	GameID   string
	PlayerID string
	CardID   int
}

type Score struct {
	PlayerName string
	Score      int
	Level      GameLevel
	Time       int // in seconds
	Moves      int
}

type Achievement struct {
	Name        string
	Description string
	Unlocked    bool
}
