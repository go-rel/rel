package rel

import (
	"fmt"
	"reflect"
)

// Modifier is interface for a record modifier.
type Modifier interface {
	Apply(doc *Document, modification *Modification)
}

// Apply using given modifiers.
func Apply(doc *Document, modifiers ...Modifier) Modification {
	modification := Modification{
		Modifies: make(map[string]Modify),
		Assoc:    make(map[string]AssocModification),
	}

	// FIXME: supports db default
	for i := range modifiers {
		modifiers[i].Apply(doc, &modification)
	}

	return modification
}

// AssocModification represents modification for association.
type AssocModification struct {
	Modifications []Modification
	DeletedIDs    []interface{}
}

// Modification represents value to be inserted or updated to database.
// It's not safe to be used multiple time. some operation my alter modification data.
type Modification struct {
	Modifies map[string]Modify
	Assoc    map[string]AssocModification
	Reload   bool
}

// Add a modify.
func (m *Modification) Add(mod Modify) {
	m.Modifies[mod.Field] = mod
}

// SetAssoc modification.
func (m *Modification) SetAssoc(field string, mods ...Modification) {
	assoc := m.Assoc[field]
	assoc.Modifications = mods
	m.Assoc[field] = assoc
}

// SetDeletedIDs modification.
// nil slice will clear association.
func (m *Modification) SetDeletedIDs(field string, ids []interface{}) {
	assoc := m.Assoc[field]
	assoc.DeletedIDs = ids
	m.Assoc[field] = assoc
}

// ChangeOp represents type of modify operation.
type ChangeOp int

const (
	// ChangeInvalidOp operation.
	ChangeInvalidOp ChangeOp = iota
	// ChangeSetOp operation.
	ChangeSetOp
	// ChangeIncOp operation.
	ChangeIncOp
	// ChangeDecOp operation.
	ChangeDecOp
	// ChangeFragmentOp operation.
	ChangeFragmentOp
)

// Modify stores modification instruction.
type Modify struct {
	Type  ChangeOp
	Field string
	Value interface{}
}

// Apply modification.
func (m Modify) Apply(doc *Document, modification *Modification) {
	invalid := false

	switch m.Type {
	case ChangeSetOp:
		if !doc.SetValue(m.Field, m.Value) {
			invalid = true
		}
	case ChangeFragmentOp:
		modification.Reload = true
	default:
		if typ, ok := doc.Type(m.Field); ok {
			kind := typ.Kind()
			invalid = (m.Type == ChangeIncOp || m.Type == ChangeDecOp) &&
				(kind < reflect.Int || kind > reflect.Uint64)
		} else {
			invalid = true
		}

		modification.Reload = true
	}

	if invalid {
		panic(fmt.Sprint("rel: cannot assign ", m.Value, " as ", m.Field, " into ", doc.Table()))
	}

	modification.Add(m)
}

// Set create a modify using set operation.
func Set(field string, value interface{}) Modify {
	return Modify{
		Type:  ChangeSetOp,
		Field: field,
		Value: value,
	}
}

// Inc create a modify using increment operation.
func Inc(field string) Modify {
	return IncBy(field, 1)
}

// IncBy create a modify using increment operation with custom increment value.
func IncBy(field string, n int) Modify {
	return Modify{
		Type:  ChangeIncOp,
		Field: field,
		Value: n,
	}
}

// Dec create a modify using deccrement operation.
func Dec(field string) Modify {
	return DecBy(field, 1)
}

// DecBy create a modify using decrement operation with custom decrement value.
func DecBy(field string, n int) Modify {
	return Modify{
		Type:  ChangeDecOp,
		Field: field,
		Value: n,
	}
}

// SetFragment create a modify operation using randoc fragment operation.
// Only available for Update.
func SetFragment(raw string, args ...interface{}) Modify {
	return Modify{
		Type:  ChangeFragmentOp,
		Field: raw,
		Value: args,
	}
}

// Setf is an alias for SetFragment
var Setf = SetFragment
