package uid

import (
	"crypto/rand"
	"encoding/hex"
)

// Generates random hexadecimal string. Panics if random source cannot be read.
func New() string {
	data := [8]byte{}

	_, err := rand.Read(data[:])
	if err != nil {
		panic("could not read from random source")
	}

	return hex.EncodeToString(data[:])
}
