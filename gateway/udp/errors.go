package udp

import (
	"errors"
)

var (
	ErrUnauthorizedKey   = errors.New("unauthorized key")
	ErrInsufficientScope = errors.New("insufficient scope")
	ErrNoSession         = errors.New("no such session")
)
