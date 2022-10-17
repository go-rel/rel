//go:build go1.18
// +build go1.18

package rel

import (
	"io"
)

// EntityIterator allows iterating through all entity in database in batch.
type EntityIterator[T any] interface {
	io.Closer
	Next() (T, error)
}

type entityIterator[T any] struct {
	iterator Iterator
}

func (ei *entityIterator[T]) Close() error {
	return ei.iterator.Close()
}

func (ei *entityIterator[T]) Next() (T, error) {
	var entity T
	return entity, ei.iterator.Next(&entity)
}

func newEntityIterator[T any](iterator Iterator) EntityIterator[T] {
	return &entityIterator[T]{
		iterator: iterator,
	}
}
