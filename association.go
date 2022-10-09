package rel

import (
	"reflect"
)

// Association provides abstraction to work with association of document or collection.
type Association struct {
	meta AssociationMeta
	rv   reflect.Value
}

// Type of association.
func (a Association) Type() AssociationType {
	return a.meta.Type()
}

// Document returns association target as document.
// If association is zero, second return value will be false.
func (a Association) Document() (*Document, bool) {
	return a.document(false)
}

// LazyDocument is a lazy version of Document.
// If rv is a null pointer, it returns a document that delays setting the value of rv
// until Document#Add() is called.
func (a Association) LazyDocument() (*Document, bool) {
	return a.document(true)
}

func (a Association) document(lazy bool) (*Document, bool) {
	var (
		rv = reflectValueFieldByIndex(a.rv, a.meta.targetIndex, !lazy)
	)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			if !lazy {
				rv.Set(reflect.New(rv.Type().Elem()))
			}

			return NewDocument(rv), false
		}

		var (
			doc = NewDocument(rv)
		)

		return doc, doc.Persisted()
	default:
		var (
			doc = NewDocument(rv.Addr())
		)

		return doc, doc.Persisted()
	}
}

// Collection returns association target as collection.
// If association is zero, second return value will be false.
func (a Association) Collection() (*Collection, bool) {
	var (
		rv     = reflectValueFieldByIndex(a.rv, a.meta.targetIndex, true)
		loaded = !rv.IsNil()
	)

	if rv.Kind() == reflect.Ptr {
		if !loaded {
			rv.Set(reflect.New(rv.Type().Elem()))
			rv.Elem().Set(reflect.MakeSlice(rv.Elem().Type(), 0, 0))
		}

		return NewCollection(rv), loaded
	}

	if !loaded {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}

	return NewCollection(rv.Addr()), loaded
}

// IsZero returns true if association is not loaded.
func (a Association) IsZero() bool {
	var (
		rv = reflectValueFieldByIndex(a.rv, a.meta.targetIndex, false)
	)

	return isDeepZero(reflect.Indirect(rv), 1)
}

// ReferenceField of the association.
func (a Association) ReferenceField() string {
	return a.meta.ReferenceField()
}

// ReferenceValue of the association.
func (a Association) ReferenceValue() any {
	return indirectInterface(reflectValueFieldByIndex(a.rv, a.meta.referenceIndex, false))
}

// ForeignField of the association.
func (a Association) ForeignField() string {
	return a.meta.ForeignField()
}

// ForeignValue of the association.
// It'll panic if association type is has many.
func (a Association) ForeignValue() any {
	if a.Type() == HasMany {
		panic("rel: cannot infer foreign value for has many or many to many association")
	}

	var (
		rv = reflectValueFieldByIndex(a.rv, a.meta.targetIndex, false)
	)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return indirectInterface(reflectValueFieldByIndex(rv, a.meta.foreignIndex, false))
}

// Through return intermediary association.
func (a Association) Through() string {
	return a.meta.Through()
}

// Autoload assoc setting when parent is loaded.
func (a Association) Autoload() bool {
	return a.meta.Autoload()
}

// Autosave setting when parent is created/updated/deleted.
func (a Association) Autosave() bool {
	return a.meta.Autosave()
}

func newAssociation(rv reflect.Value, index []int) Association {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return Association{
		meta: getAssociationMeta(rv.Type(), index),
		rv:   rv,
	}
}
