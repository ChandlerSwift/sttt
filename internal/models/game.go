package models

import (
	"crypto/rand"
	"math/big"
	"time"

	"gorm.io/gorm"
	"sttt.chandlerswift.com/sttt/pkg/sttt"
)

// State describes what state the game is in
type State int

const (
	// InProgress describes a game that has been created, but not won or tied.
	InProgress State = iota

	// Won represents a game that either one player won, or the other player
	// resigned.
	Won

	// Abandoned references a game that is not won, but is more than a month
	// old.
	Abandoned

	// Tied references a game that has been played into such a state that
	// neither player can win. For example, if the top level grid is
	//   X | O |
	//  ---+---+---
	//   O |   | X
	//  ---+---+---
	//     | X | O
	// then no matter who takes the last three squares, nobody can attain three
	// top-level squares in a row.
	Tied
)

// Game wraps the sttt Game class, providing some extra information for our use,
// like what userids are participating. It does _not_ save the full game state;
// that's just saved in memory. We only hit the database when saving a completed
// game.
type Game struct {
	gorm.Model
	sttt.Game
	JoinID    string `gorm:"unique"` // The joining key that players use to pair manually
	Player1ID int
	Player1   User // "X"
	Player2ID int
	Player2   User // "O"
	WinnerID  int
	Winner    User
	State     State
	StartTime time.Time
	Duration  time.Duration
	Password  string
}

// GetLink returns the URL to view or participate in a game.
func (g Game) GetLink() string {
	return "/join/" + g.JoinID
}

// NewGame creates and initializes a new game.
func NewGame() *Game {
	// TODO: check for collisions, or increment, or...
	alphabet := []rune("ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789")
	gameJoinIDLen := 4
	joinID := make([]rune, gameJoinIDLen)
	for i := range joinID {
		randInt, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			panic("Failed to rand.Int()")
		}
		joinID[i] = alphabet[randInt.Int64()]
	}
	return &Game{
		JoinID:    string(joinID),
		StartTime: time.Now(),
	}
}
