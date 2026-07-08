package game

import "errors"

var (
	ErrGameNotFound  = errors.New("game not found")
	ErrGameFull      = errors.New("game is already full")
	ErrAlreadyJoined = errors.New("cannot join a game you created")
)
