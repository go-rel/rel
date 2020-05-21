package rel

import (
	"fmt"
	"reflect"
)

// Mutator is interface for a record mutator.
type Mutator interface {
	Apply(doc *Document, mutation *Mutation)
}

// Apply using given mutators.
func Apply(doc *Document, mutators ...Mutator) Mutation {
	var (
		mutation     Mutation
		optionsCount int
	)

	for i := range mutators {
		switch mut := mutators[i].(type) {
		case Unscoped, Reload:
			optionsCount++
			mut.Apply(doc, &mutation)
		default:
			mut.Apply(doc, &mutation)
		}
	}

	// fallback to structset.
	if optionsCount == len(mutators) {
		newStructset(doc, false).Apply(doc, &mutation)
	}

	return mutation
}

// AssocMutation represents mutation for association.
type AssocMutation struct {
	Mutations  []Mutation
	DeletedIDs []interface{}
}

// Mutation represents value to be inserted or updated to database.
// It's not safe to be used multiple time. some operation my alter mutation data.
type Mutation struct {
	Mutates  map[string]Mutate
	Assoc    map[string]AssocMutation
	Unscoped Unscoped
	Reload   Reload
}

func (m *Mutation) initMutates() {
	if m.Mutates == nil {
		m.Mutates = make(map[string]Mutate)
	}
}

func (m *Mutation) initAssoc() {
	if m.Assoc == nil {
		m.Assoc = make(map[string]AssocMutation)
	}
}

// IsEmpty returns true if no mutates operation and assoc's mutation is defined.
func (m *Mutation) IsEmpty() bool {
	return m.IsMutatesEmpty() && m.IsAssocEmpty()
}

// IsMutatesEmpty returns true if no mutates operation is defined.
func (m *Mutation) IsMutatesEmpty() bool {
	return len(m.Mutates) == 0
}

// IsAssocEmpty returns true if no assoc's mutation is defined.
func (m *Mutation) IsAssocEmpty() bool {
	return len(m.Assoc) == 0
}

// Add a mutate.
func (m *Mutation) Add(mut Mutate) {
	m.initMutates()

	m.Mutates[mut.Field] = mut
}

// SetAssoc mutation.
func (m *Mutation) SetAssoc(field string, mods ...Mutation) {
	m.initAssoc()

	assoc := m.Assoc[field]
	assoc.Mutations = mods
	m.Assoc[field] = assoc
}

// SetDeletedIDs mutation.
// nil slice will clear association.
func (m *Mutation) SetDeletedIDs(field string, ids []interface{}) {
	m.initAssoc()

	assoc := m.Assoc[field]
	assoc.DeletedIDs = ids
	m.Assoc[field] = assoc
}

// ChangeOp represents type of mutate operation.
type ChangeOp int

const (
	// ChangeInvalidOp operation.
	ChangeInvalidOp ChangeOp = iota
	// ChangeSetOp operation.
	ChangeSetOp
	// ChangeIncOp operation.
	ChangeIncOp
	// ChangeFragmentOp operation.
	ChangeFragmentOp
)

// Mutate stores mutation instruction.
type Mutate struct {
	Type  ChangeOp
	Field string
	Value interface{}
}

// Apply mutation.
func (m Mutate) Apply(doc *Document, mutation *Mutation) {
	invalid := false

	switch m.Type {
	case ChangeSetOp:
		if !doc.SetValue(m.Field, m.Value) {
			invalid = true
		}
	case ChangeFragmentOp:
		mutation.Reload = true
	default:
		if typ, ok := doc.Type(m.Field); ok {
			kind := typ.Kind()
			invalid = m.Type == ChangeIncOp && (kind < reflect.Int || kind > reflect.Uint64)
		} else {
			invalid = true
		}

		mutation.Reload = true
	}

	if invalid {
		panic(fmt.Sprint("rel: cannot assign ", m.Value, " as ", m.Field, " into ", doc.Table()))
	}

	mutation.Add(m)
}

// Set create a mutate using set operation.
func Set(field string, value interface{}) Mutate {
	return Mutate{
		Type:  ChangeSetOp,
		Field: field,
		Value: value,
	}
}

// Inc create a mutate using increment operation.
func Inc(field string) Mutate {
	return IncBy(field, 1)
}

// IncBy create a mutate using increment operation with custom increment value.
func IncBy(field string, n int) Mutate {
	return Mutate{
		Type:  ChangeIncOp,
		Field: field,
		Value: n,
	}
}

// Dec create a mutate using deccrement operation.
func Dec(field string) Mutate {
	return DecBy(field, 1)
}

// DecBy create a mutate using decrement operation with custom decrement value.
func DecBy(field string, n int) Mutate {
	return Mutate{
		Type:  ChangeIncOp,
		Field: field,
		Value: -n,
	}
}

// SetFragment create a mutate operation using randoc fragment operation.
// Only available for Update.
func SetFragment(raw string, args ...interface{}) Mutate {
	return Mutate{
		Type:  ChangeFragmentOp,
		Field: raw,
		Value: args,
	}
}

// Setf is an alias for SetFragment
var Setf = SetFragment

// Reload force reload after insert/update.
type Reload bool

// Apply mutation.
func (r Reload) Apply(doc *Document, mutation *Mutation) {
	mutation.Reload = r
}
