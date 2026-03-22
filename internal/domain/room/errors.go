package room

import "errors"

var (
	ErrRoomNotFound           = errors.New("room not found")
	ErrNameRequired           = errors.New("room name is required")
	ErrCapacityMustBePositive = errors.New("room capacity must be positive")
)
