package sql

import (
	"bytes"
)

type Buffer struct {
	bytes.Buffer
	Arguments []interface{}
}

func (b *Buffer) Append(args ...interface{}) {
	b.Arguments = append(b.Arguments, args...)
}

func (b *Buffer) Reset() {
	b.Buffer.Reset()
	b.Arguments = nil
}
