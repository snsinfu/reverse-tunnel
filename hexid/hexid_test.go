package hexid

import (
	"encoding/hex"
	"testing"
)

func Test_New_ReturnsHexString(t *testing.T) {
	sizes := []int{1, 2, 3, 4, 5}

	for _, n := range sizes {
		id := New(n)
		if len(id) != 2*n {
			t.Fatalf("unexpected string length: got %d, want %d", len(id), 2*n)
		}

		data, err := hex.DecodeString(id)
		if err != nil {
			t.Fatalf("unexpected error: %s", err)
		}

		if len(data) != n {
			t.Fatalf("unexpected data length: got %d, want %d", len(data), n)
		}
	}
}

func Test_New_ReturnsUniqueStringForEachCall(t *testing.T) {
	s1 := New(4)
	s2 := New(4)

	if s1 == s2 {
		t.Errorf("generated strings should be different: '%s' and '%s'", s1, s2)
	}
}
