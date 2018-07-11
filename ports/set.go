package ports

// Set represents a set of port numbers.
type Set map[int]bool

// Append appends port to set.
func (set Set) Append(port int) {
	set[port] = true
}

// Has returns true if set contains port.
func (set Set) Has(port int) bool {
	_, ok := set[port]
	return ok
}
