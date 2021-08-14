package reltest

import (
	"context"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type mutate []*MockMutate

func (m *mutate) register(ctxData ctxData, mutators ...rel.Mutator) *MockMutate {
	mm := &MockMutate{ctxData: ctxData, argMutators: mutators}
	*m = append(*m, mm)
	return mm
}

func (m mutate) execute(ctx context.Context, record interface{}, mutators ...rel.Mutator) error {
	for _, mm := range m {
		if fetchContext(ctx) == mm.ctxData &&
			(mm.argRecord == nil || reflect.DeepEqual(mm.argRecord, record)) &&
			(mm.argRecordType == "" || mm.argRecordType == reflect.TypeOf(record).String()) &&
			(mm.argRecordTable == "" || mm.argRecordTable == rel.NewDocument(record, true).Table()) &&
			(mm.argRecordContains == nil || matchContains(mm.argRecordContains, record)) &&
			matchMutators(mm.argMutators, mutators) {
			return mm.retError
		}
	}

	panic("TODO: Query doesn't match")
}

// MockMutate asserts and simulate Insert function for test.
type MockMutate struct {
	ctxData           ctxData
	argRecord         interface{}
	argRecordType     string
	argRecordTable    string
	argRecordContains interface{}
	argMutators       []rel.Mutator
	retError          error
}

// For expect calls for given record.
func (mm *MockMutate) For(record interface{}) *MockMutate {
	mm.argRecord = record
	return mm
}

// ForType expect calls for given type.
// Type must include package name, example: `model.User`.
func (mm *MockMutate) ForType(typ string) *MockMutate {
	mm.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return mm
}

// ForTable expect calls for given table.
func (mm *MockMutate) ForTable(typ string) *MockMutate {
	mm.argRecordTable = typ
	return mm
}

// ForContains expect calls to contains some value of given struct.
func (mm *MockMutate) ForContains(contains interface{}) *MockMutate {
	mm.argRecordContains = contains
	return mm
}

// Error sets error to be returned.
func (mm *MockMutate) Error(err error) {
	mm.retError = err
}

// ConnectionClosed sets this error to be returned.
func (mm *MockMutate) ConnectionClosed() {
	mm.Error(ErrConnectionClosed)
}

// NotUnique sets not unique error to be returned.
func (mm *MockMutate) NotUnique(key string) {
	mm.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}
