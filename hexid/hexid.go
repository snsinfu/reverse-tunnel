// Package hexid provides a dead simple hexadecimal random ID generator.
package hexid

import (
	"crypto/rand"
	"encoding/hex"
)

// New generates a lowercase hexadecimal string encoding n bytes of random data
// (hence the generated string is 2*n characters long). This function panics if
// the random source cannot be read.
func New(n int) string {
	data := make([]byte, n)

	_, err := rand.Read(data)
	if err != nil {
		panic("could not read from crypto/rand")
	}

	return hex.EncodeToString(data)
}
