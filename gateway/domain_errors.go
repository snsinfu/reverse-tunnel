package main

import (
	"errors"
)

var (
	ErrMissingKey       = errors.New("missing key")
	ErrInvalidKey       = errors.New("invalid key")
	ErrInvalidPort      = errors.New("invalid port")
	ErrInvalidSessionID = errors.New("invalid session ID")
)
