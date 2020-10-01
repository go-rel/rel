package reltest

import (
	"strings"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/mock"
)

// Mutate asserts and simulate insert or update function for test.
type Mutate struct {
	*Expect
}

// For match expect calls for given record.
func (m *Mutate) For(record interface{}) *Mutate {
	m.Arguments[1] = record
	return m
}

// ForType match expect calls for given type.
// Type must include package name, example: `model.User`.
func (m *Mutate) ForType(typ string) *Mutate {
	return m.For(mock.AnythingOfType("*" + strings.TrimPrefix(typ, "*")))
}

// NotUnique sets not unique error to be returned.
func (m *Mutate) NotUnique(key string) {
	m.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

func expectMutate(r *Repository, methodName string, mutators []rel.Mutator) *Mutate {
	mutatorsArgument := interface{}(mutators)
	if mutators == nil {
		mutatorsArgument = mock.Anything
	}

	em := &Mutate{
		Expect: newExpect(r, methodName,
			[]interface{}{r.ctxData, mock.Anything, mutatorsArgument},
			[]interface{}{nil},
		),
	}

	return em
}

// ExpectInsert to be called with given field and queries.
func ExpectInsert(r *Repository, mutators []rel.Mutator) *Mutate {
	return expectMutate(r, "Insert", mutators)
}

// ExpectUpdate to be called with given field and queries.
func ExpectUpdate(r *Repository, mutators []rel.Mutator) *Mutate {
	return expectMutate(r, "Update", mutators)
}

// ExpectInsertAll to be called.
func ExpectInsertAll(r *Repository) *Mutate {
	em := &Mutate{
		Expect: newExpect(r, "InsertAll",
			[]interface{}{r.ctxData, mock.Anything},
			[]interface{}{nil},
		),
	}

	return em
}
