// Package sttt defines the rules for super tic-tac-toe. It's based on the logic
// at  https://github.com/smit8397/super-tic-tac-toe/blob/master/game.py, and
// functions fundamentally the same, but has some extra error checking and other
// functionality built in. It is intended to conform to the description at
// https://en.wikipedia.org/wiki/Ultimate_tic-tac-toe
package sttt

import (
	"errors"
	"fmt"
)

// A Token is placed by a Move, represents a player, and is placed on a Subboard.
type Token int

// winningCombos is an enumeration of all the possible win conditions (including
// 3 horizontal, 3 vertical, and two diagonal).
var winningCombos = [][]int{{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, {0, 3, 6}, {1, 4, 7}, {2, 5, 8}, {0, 4, 8}, {2, 4, 6}}

const (
	Empty   Token = iota // Empty
	Player1              // "X"
	Player2              // "O"
)

// A Subboard represents the state of a super tic-tac-toe sub-board; that is,
// a "normal tic-tac-toe grid".
type Subboard struct {
	Finished         bool
	Winner           Token
	State            [9]Token
	placedTokenCount int
}

// PlaceToken places a token `token` at location `location` (0-indexed, so 0-8)
func (s Subboard) PlaceToken(token Token, location int) (err error) {
	if s.Finished {
		return errors.New("Board is already won")
	}
	if location < 0 || location > 8 {
		return fmt.Errorf("Invalid location %v", location)
	}
	if s.State[location] != Empty {
		return fmt.Errorf("Cannot play in location %v, %v is already there", location, s.State[location])
	}
	s.State[location] = token
	s.placedTokenCount++
	s.Finished = s.checkIfWon()
	return nil
}

// Checks if the subboard is won, and as a side effect updates the s.Winner.
func (s Subboard) checkIfWon() bool {
	// Check for winner
	for _, wc := range winningCombos {
		if s.State[wc[0]] != Empty &&
			s.State[wc[0]] == s.State[wc[1]] &&
			s.State[wc[0]] == s.State[wc[2]] {
			s.Winner = s.State[wc[0]]
			return true
		}
	}

	if s.placedTokenCount == 9 {
		// No need to set a winner, as nobody won.
		return true
	}
	return false
}

// A Move is a placed token. A list of these is kept to record a history of game state.
type Move struct {
	Token    Token
	Subboard int
	Location int
}

// A Game is a game of super tic-tac-toe.
type Game struct {
	Finished  bool       // true when game is won, false otherwise.
	Subboards []Subboard `gorm:"-"`
	Winner    Token
	Moves     []Move `gorm:"-"`
}

func (g Game) move(subboard int, location int, token Token) error {
	if g.Finished {
		return errors.New("The game is over")
	}
	if len(g.Moves) > 0 && g.Moves[len(g.Moves)-1].Token == token {
		return errors.New("It's not your turn")
	}

	if token == Empty {
		return errors.New("Empty is not a valid player")
	}
	// Make sure we're moving on a valid subboard
	if len(g.Moves) > 0 { // There's a previous move to check against
		prevMove := g.Moves[len(g.Moves)-1]
		destBoard := g.Subboards[prevMove.Location]
		if subboard != prevMove.Location && destBoard.Finished == false {
			return fmt.Errorf("Invalid subboard. You should be moving in subboard %v", prevMove.Location)
		}
	}

	g.Subboards[subboard].PlaceToken(token, location)
	g.Moves = append(g.Moves, Move{
		Token:    token,
		Subboard: subboard,
		Location: location,
	})

	g.Finished = g.checkIfWon()
	return nil
}

func (g Game) checkIfWon() bool {
	// Check for winner
	for _, wc := range winningCombos {
		if g.Subboards[wc[0]].Winner != Empty &&
			g.Subboards[wc[0]].Winner == g.Subboards[wc[1]].Winner &&
			g.Subboards[wc[0]].Winner == g.Subboards[wc[2]].Winner {
			g.Winner = g.Subboards[wc[0]].Winner
			return true
		}
	}

	completeSubboards := 0
	for _, subboard := range g.Subboards {
		if subboard.Finished {
			completeSubboards++
		}
	}
	if completeSubboards == 9 {
		// No need to set a winner, as nobody won.
		return true
	}
	return false
}
