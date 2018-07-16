package ports

// Set represents a set of port numbers.
type Set map[int]bool

// Add adds a port number to set.
func (set Set) Add(port int) {
	set[port] = true
}

// Has returns true if set contains given port number.
func (set Set) Has(port int) bool {
	_, ok := set[port]
	return ok
}
