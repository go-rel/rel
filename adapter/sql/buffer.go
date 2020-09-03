package sql

import (
	"strings"
)

// Buffer used to strings buffer and argument of the query.
type Buffer struct {
	strings.Builder
	Arguments []interface{}
}

// Append argumetns.
func (b *Buffer) Append(args ...interface{}) {
	b.Arguments = append(b.Arguments, args...)
}

// Reset buffer.
func (b *Buffer) Reset() {
	b.Builder.Reset()
	b.Arguments = nil
}
