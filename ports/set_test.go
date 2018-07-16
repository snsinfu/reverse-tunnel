package ports

import (
	"reflect"
	"testing"
)

func TestSet_IsIterable(t *testing.T) {
	set := Set{}

	set.Add(100)
	set.Add(200)

	seen := map[int]int{}
	for port := range set {
		seen[port]++
	}

	expected := map[int]int{
		100: 1,
		200: 1,
	}

	if !reflect.DeepEqual(seen, expected) {
		t.Errorf("unexpected result: got %v, want %v", seen, expected)
	}
}

func TestSet_Has_ChecksForExistence(t *testing.T) {
	set := Set{}

	const (
		existing    = 100
		nonexisting = 200
	)
	set.Add(existing)

	if !set.Has(existing) {
		t.Errorf("Has reported false for existing member %d", existing)
	}

	if set.Has(nonexisting) {
		t.Errorf("Has reported true for non-existing member %d", nonexisting)
	}
}
