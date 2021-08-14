package reltest

import (
	"context"
	"io"
	"reflect"

	"github.com/go-rel/rel"
)

type iterate []*MockIterate

func (i *iterate) register(ctxData ctxData, query rel.Query, options ...rel.IteratorOption) *MockIterate {
	mi := &MockIterate{ctxData: ctxData, argQuery: query, argOptions: options}
	*i = append(*i, mi)
	return mi
}

func (i iterate) execute(ctx context.Context, query rel.Query, options ...rel.IteratorOption) rel.Iterator {
	for _, mi := range i {
		if fetchContext(ctx) == mi.ctxData &&
			reflect.DeepEqual(mi.argOptions, options) &&
			matchQuery(mi.argQuery, query) {
			return mi
		}
	}

	panic("TODO: Query doesn't match")
}

type data interface {
	Len() int
	Get(index int) *rel.Document
}

// MockIterate asserts and simulate Delete function for test.
type MockIterate struct {
	ctxData    ctxData
	result     data
	current    int
	err        error
	argQuery   rel.Query
	argOptions []rel.IteratorOption
}

// Result sets the result of preload.
func (mi *MockIterate) Result(result interface{}) {
	rt := reflect.TypeOf(result)
	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		mi.result = rel.NewCollection(result, true)
	} else {
		mi.result = rel.NewDocument(result, true)
	}
}

// Error sets error to be returned.
func (mi *MockIterate) Error(err error) {
	mi.err = err
}

// ConnectionClosed sets this error to be returned.
func (mi *MockIterate) ConnectionClosed() {
	mi.Error(ErrConnectionClosed)
}

func (mi MockIterate) Close() error {
	return nil
}

func (mi *MockIterate) Next(record interface{}) error {
	if mi.err != nil {
		return mi.err
	}

	if mi.result == nil || mi.current == mi.result.Len() {
		return io.EOF
	}

	var (
		doc = mi.result.Get(mi.current)
	)

	reflect.ValueOf(record).Elem().Set(doc.ReflectValue())

	mi.current++
	return nil
}
