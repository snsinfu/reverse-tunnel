package main

var (
	authorities = map[string]AuthScope
)

// AuthScope
type AuthScope struct {
	ports []NetPort
}
