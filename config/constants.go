package config

// BufferSize is the max size of a single message transferred via a tunnel. This
// should be smaller than MTU minus overhead to prevent fragmentation.
const BufferSize = 1400
