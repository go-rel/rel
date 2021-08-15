package reltest

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-rel/rel"
)

type delete []*MockDelete

func (d *delete) register(ctxData ctxData, options ...rel.Cascade) *MockDelete {
	md := &MockDelete{
		assert:     &Assert{ctxData: ctxData},
		argOptions: options,
	}
	*d = append(*d, md)
	return md
}

func (d delete) execute(ctx context.Context, record interface{}, options ...rel.Cascade) error {
	for _, md := range d {
		if (md.argRecord == nil || reflect.DeepEqual(md.argRecord, record)) &&
			(md.argRecordType == "" || md.argRecordType == reflect.TypeOf(record).String()) &&
			(md.argRecordTable == "" || md.argRecordTable == rel.NewDocument(record, true).Table()) &&
			(md.argRecordContains == nil || matchContains(md.argRecordContains, record)) &&
			reflect.DeepEqual(md.argOptions, options) &&
			md.assert.call(ctx) {
			return md.retError
		}
	}

	panic(failExecuteMessage(MockDelete{argRecord: record, argOptions: options}, d))
}

func (d *delete) assert(t T) bool {
	for _, md := range *d {
		if !md.assert.assert(t, md) {
			return false
		}
	}

	*d = nil
	return true
}

// MockDelete asserts and simulate Delete function for test.
type MockDelete struct {
	assert            *Assert
	argRecord         interface{}
	argRecordType     string
	argRecordTable    string
	argRecordContains interface{}
	argOptions        []rel.Cascade
	retError          error
}

// For assert calls for given record.
func (md *MockDelete) For(record interface{}) *MockDelete {
	md.argRecord = record
	return md
}

// ForType assert calls for given type.
// Type must include package name, example: `model.User`.
func (md *MockDelete) ForType(typ string) *MockDelete {
	md.argRecordType = "*" + strings.TrimPrefix(typ, "*")
	return md
}

// ForTable assert calls for given table.
func (md *MockDelete) ForTable(typ string) *MockDelete {
	md.argRecordTable = typ
	return md
}

// ForContains assert calls to contains some value of given struct.
func (md *MockDelete) ForContains(contains interface{}) *MockDelete {
	md.argRecordContains = contains
	return md
}

// Error sets error to be returned.
func (md *MockDelete) Error(err error) *Assert {
	md.retError = err
	return md.assert
}

// Success sets no error to be returned.
func (md *MockDelete) Success() *Assert {
	return md.Error(nil)
}

// ConnectionClosed sets this error to be returned.
func (md *MockDelete) ConnectionClosed() *Assert {
	return md.Error(ErrConnectionClosed)
}

// String representation of mocked call.
func (md MockDelete) String() string {
	argRecord := "<Any>"
	if md.argRecord != nil {
		argRecord = csprint(md.argRecord, true)
	} else if md.argRecordContains != nil {
		argRecord = fmt.Sprintf("<Contains: %s>", csprint(md.argRecordContains, true))
	} else if md.argRecordType != "" {
		argRecord = fmt.Sprintf("<Type: %s>", md.argRecordType)
	} else if md.argRecordTable != "" {
		argRecord = fmt.Sprintf("<Table: %s>", md.argRecordTable)
	}

	argCascade := ""
	for i := range md.argOptions {
		argCascade += fmt.Sprintf(", %v", md.argOptions[i])
	}

	return fmt.Sprintf("Delete(ctx, %s%s)", argRecord, argCascade)
}

// ExpectString representation of mocked call.
func (md MockDelete) ExpectString() string {
	argOptions := ""
	for i := range md.argOptions {
		argOptions += fmt.Sprintf("%v", md.argOptions[i])
	}

	return fmt.Sprintf("ExpectDelete(%s).ForType(\"%T\")", argOptions, md.argRecord)
}
