package sql

import (
	"bytes"
)

// Buffer used to strings buffer and argument of the query.
type Buffer struct {
	bytes.Buffer
	Arguments []interface{}
}

// Append argumetns.
func (b *Buffer) Append(args ...interface{}) {
	b.Arguments = append(b.Arguments, args...)
}

// Reset buffer.
func (b *Buffer) Reset() {
	b.Buffer.Reset()
	b.Arguments = nil
}
