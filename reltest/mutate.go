package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type mutate []*MockMutate

func (m *mutate) register(name string, ctxData ctxData, mutators ...rel.Mutator) *MockMutate {
	mm := &MockMutate{
		assert:      &Assert{ctxData: ctxData},
		name:        name,
		argMutators: mutators,
	}
	*m = append(*m, mm)
	return mm
}

func (m mutate) execute(ctx context.Context, record interface{}, mutators ...rel.Mutator) error {
	for _, mm := range m {
		if (mm.argRecord == nil || reflect.DeepEqual(mm.argRecord, record)) &&
			(mm.argRecordType == "" || mm.argRecordType == reflect.TypeOf(record).String()) &&
			(mm.argRecordTable == "" || mm.argRecordTable == rel.NewDocument(record, true).Table()) &&
			(mm.argRecordContains == nil || matchContains(mm.argRecordContains, record)) &&
			matchMutators(mm.argMutators, mutators) &&
			mm.assert.call(ctx) {
			return mm.retError
		}
	}

	panic(failExecuteMessage(MockMutate{argRecord: record, argMutators: mutators}, m))
}

func (m *mutate) assert(t T) bool {
	for _, mm := range *m {
		if !mm.assert.assert(t, mm) {
			return false
		}
	}

	*m = nil
	return true
}

// MockMutate asserts and simulate Insert function for test.
type MockMutate struct {
	assert            *Assert
	name              string
	argRecord         interface{}
	argRecordType     string
	argRecordTable    string
	argRecordContains interface{}
	argMutators       []rel.Mutator
	retError          error
}

// For assert calls for given record.
func (mm *MockMutate) For(record interface{}) *MockMutate {
	mm.argRecord = record
	return mm
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (mm *MockMutate) ForType(typ string) *MockMutate {
	mm.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return mm
}

// ForTable assert calls for given table.
func (mm *MockMutate) ForTable(typ string) *MockMutate {
	mm.argRecordTable = typ
	return mm
}

// ForContains assert calls to contains some value of given struct.
func (mm *MockMutate) ForContains(contains interface{}) *MockMutate {
	mm.argRecordContains = contains
	return mm
}

// Error sets error to be returned.
func (mm *MockMutate) Error(err error) *Assert {
	mm.retError = err
	return mm.assert
}

// Success sets no error to be returned.
func (mm *MockMutate) Success() *Assert {
	return mm.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (mm *MockMutate) ConnectionClosed() *Assert {
	return mm.Error(ErrConnectionClosed)
}

// NotUnique sets not unique error to be returned.
func (mm *MockMutate) NotUnique(key string) *Assert {
	return mm.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

// String representation of mocked call.
func (mm MockMutate) String() string {
	argRecord := "<Any>"
	if mm.argRecord != nil {
		argRecord = csprint(mm.argRecord, true)
	} else if mm.argRecordContains != nil {
		argRecord = fmt.Sprintf("<Contains: %s>", csprint(mm.argRecordContains, true))
	} else if mm.argRecordType != "" {
		argRecord = fmt.Sprintf("<Type: %s>", mm.argRecordType)
	} else if mm.argRecordTable != "" {
		argRecord = fmt.Sprintf("<Table: %s>", mm.argRecordTable)
	}

	argMutators := ""
	for i := range mm.argMutators {
		argMutators += fmt.Sprintf(", %v", mm.argMutators[i])
	}

	return fmt.Sprintf("%s(ctx, %s%s)", mm.name, argRecord, argMutators)
}

// ExpectString representation of mocked call.
func (mm MockMutate) ExpectString() string {
	argMutators := ""
	for i := range mm.argMutators {
		if i > 0 {
			argMutators += fmt.Sprintf(", %v", mm.argMutators[i])
		} else {
			argMutators += fmt.Sprintf("%v", mm.argMutators[i])
		}
	}

	return fmt.Sprintf("Expect%s(%s).ForType(\"%T\")", mm.name, argMutators, mm.argRecord)
}
