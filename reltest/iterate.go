package reltest

import (
	"io"
	"reflect"

	"github.com/go-rel/rel"
)

type data interface {
	Len() int
	Get(index int) *rel.Document
}

type iterator struct {
	current int
	data    data
	err     error
}

func (i iterator) Close() error {
	return nil
}

func (i *iterator) Next(record interface{}) error {
	if i.err != nil {
		return i.err
	}

	if i.data == nil || i.current == i.data.Len() {
		return io.EOF
	}

	var (
		doc = i.data.Get(i.current)
	)

	reflect.ValueOf(record).Elem().Set(doc.ReflectValue())

	i.current++
	return nil
}

// Iterate asserts and simulate iterate function for test.
type Iterate iterator

// Error sets error to be returned.
func (i *Iterate) Error(err error) {
	i.err = err
}

// ConnectionClosed sets this error to be returned.
func (i *Iterate) ConnectionClosed() {
	i.Error(ErrConnectionClosed)
}

// Result sets the result of this query.
func (i *Iterate) Result(records interface{}) {
	rt := reflect.TypeOf(records)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		i.data = rel.NewCollection(records, true)
	} else {
		i.data = rel.NewDocument(records, true)
	}
}

// ExpectIterate to be called.
func ExpectIterate(r *Repository, query rel.Query, options []rel.IteratorOption) *Iterate {
	iterate := &Iterate{}
	r.mock.On("Iterate", r.ctxData, query, options).Return(iterate).Once()
	return iterate
}
