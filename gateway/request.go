package main

import (
	"errors"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var (
	ErrInvalidPort = errors.New("invalid port")
	ErrMissingKey  = errors.New("missing key")
	ErrInvalidKey  = errors.New("invalid key")
)

// ParsePort parses s as a port number.
func ParsePort(s string) (int, error) {
	const (
		portMin = 1
		portMax = 65535
	)

	port, err := strconv.Atoi(s)
	if err != nil {
		return 0, ErrInvalidPort
	}

	if port < portMin || portMax < port {
		return port, ErrInvalidPort
	}

	return port, nil
}

// ExtractAuthKey extracts Bearer token from request header.
func ExtractAuthKey(c echo.Context) (string, error) {
	const (
		authHeader = "Authorization"
		authScheme = "Bearer "
	)

	auth := c.Request().Header.Get(authHeader)
	if auth == "" {
		return "", ErrMissingKey
	}

	if !strings.HasPrefix(auth, authScheme) {
		return "", ErrInvalidKey
	}

	return strings.TrimPrefix(auth, authScheme), nil
}
