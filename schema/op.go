package schema

// Op type.
type Op uint8

const (
	// Add operation.
	Add Op = iota
	// Alter operation.
	Alter
	// Rename operation.
	Rename
	// Drop operation.
	Drop
)
