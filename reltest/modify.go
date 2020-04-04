package reltest

import (
	"strings"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// Modify asserts and simulate insert or update function for test.
type Modify struct {
	*Expect
}

// For match expect calls for given record.
func (m *Modify) For(record interface{}) *Modify {
	m.Arguments[0] = record
	return m
}

// ForType match expect calls for given type.
// Type must include package name, example: `model.User`.
func (m *Modify) ForType(typ string) *Modify {
	return m.For(mock.AnythingOfType("*" + strings.TrimPrefix(typ, "*")))
}

// NotUnique sets not unique error to be returned.
func (m *Modify) NotUnique(key string) {
	m.Error(rel.ConstraintError{
		Key:  key,
		Type: rel.UniqueConstraint,
	})
}

func expectModify(r *Repository, methodName string, modifiers []rel.Modifier) *Modify {
	em := &Modify{
		Expect: newExpect(r, methodName,
			[]interface{}{mock.Anything, modifiers},
			[]interface{}{nil},
		),
	}

	return em
}

// ExpectInsert to be called with given field and queries.
func ExpectInsert(r *Repository, modifiers []rel.Modifier) *Modify {
	return expectModify(r, "Insert", modifiers)
}

// ExpectUpdate to be called with given field and queries.
func ExpectUpdate(r *Repository, modifiers []rel.Modifier) *Modify {
	return expectModify(r, "Update", modifiers)
}

// ExpectInsertAll to be called.
func ExpectInsertAll(r *Repository) *Modify {
	em := &Modify{
		Expect: newExpect(r, "InsertAll",
			[]interface{}{mock.Anything},
			[]interface{}{nil},
		),
	}

	return em
}
