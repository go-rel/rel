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
		fields: make(map[string]int),
		assoc:  make(map[string]int),
	}

	// FIXME: supports db default
	for i := range modifiers {
		modifiers[i].Apply(doc, &modification)
	}

	return modification
}

// Modification represents value to be inserted or updated to database.
// It's not safe to be used multiple time. some operation my alter modification data.
type Modification struct {
	fields            map[string]int // TODO: not copy friendly
	modification      []Modify
	assoc             map[string]int
	assocModification [][]Modification
	reload            bool
}

// Empty returns true if modification is empty.
func (m Modification) Empty() bool {
	return len(m.modification) == 0
}

// Count returns count of modification.
func (m Modification) Count() int {
	return len(m.modification)
}

// AssocCount returns count of associations being modification.
func (m Modification) AssocCount() int {
	return len(m.assocModification)
}

// All return array of modify.
func (m Modification) All() []Modify {
	return m.modification
}

// Get a modify by field name.
func (m Modification) Get(field string) (Modify, bool) {
	if index, ok := m.fields[field]; ok {
		return m.modification[index], true
	}

	return Modify{}, false
}

// Set a modify op directly, will existing value replace if it's already exists.
func (m *Modification) Set(mod Modify) {
	if index, exist := m.fields[mod.Field]; exist {
		m.modification[index] = mod
	} else {
		m.fields[mod.Field] = len(m.modification)
		m.modification = append(m.modification, mod)
	}
}

// GetValue of modify by field name.
func (m Modification) GetValue(field string) (interface{}, bool) {
	var (
		mod, ok = m.Get(field)
	)

	return mod.Value, ok
}

// SetValue using field name and changed value.
func (m *Modification) SetValue(field string, value interface{}) {
	m.Set(Set(field, value))
}

// GetAssoc by field name.
func (m Modification) GetAssoc(field string) ([]Modification, bool) {
	if index, ok := m.assoc[field]; ok {
		return m.assocModification[index], true
	}

	return nil, false
}

// SetAssoc by field name.
func (m *Modification) SetAssoc(field string, mods ...Modification) {
	if index, exist := m.assoc[field]; exist {
		m.assocModification[index] = mods
	} else {
		m.appendAssoc(field, mods)
	}
}

func (m *Modification) appendAssoc(field string, ac []Modification) {
	m.assoc[field] = len(m.assocModification)
	m.assocModification = append(m.assocModification, ac)
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
		modification.reload = true
	default:
		if typ, ok := doc.Type(m.Field); ok {
			kind := typ.Kind()
			invalid = (m.Type == ChangeIncOp || m.Type == ChangeDecOp) &&
				(kind < reflect.Int || kind > reflect.Uint64)
		} else {
			invalid = true
		}

		modification.reload = true
	}

	if invalid {
		panic(fmt.Sprint("rel: cannot assign ", m.Value, " as ", m.Field, " into ", doc.Table()))
	}

	modification.Set(m)
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
