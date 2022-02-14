package rel

import (
	"reflect"
	"sync"

	"github.com/serenize/snaker"
)

// AssociationType defines the type of association in database.
type AssociationType uint8

const (
	// BelongsTo association.
	BelongsTo = iota
	// HasOne association.
	HasOne
	// HasMany association.
	HasMany
)

type associationKey struct {
	rt reflect.Type
	// string repr of index, because []int is not hashable
	index string
}

type associationData struct {
	typ            AssociationType
	targetIndex    []int
	referenceField string
	referenceIndex []int
	foreignField   string
	foreignIndex   []int
	through        string
	autoload       bool
	autosave       bool
}

var associationCache sync.Map

// Association provides abstraction to work with association of document or collection.
type Association struct {
	data associationData
	rv   reflect.Value
}

// Type of association.
func (a Association) Type() AssociationType {
	return a.data.typ
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
		rv = reflectValueFieldByIndex(a.rv, a.data.targetIndex, !lazy)
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
		rv     = reflectValueFieldByIndex(a.rv, a.data.targetIndex, true)
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
		rv = reflectValueFieldByIndex(a.rv, a.data.targetIndex, false)
	)

	return isDeepZero(reflect.Indirect(rv), 1)
}

// ReferenceField of the association.
func (a Association) ReferenceField() string {
	return a.data.referenceField
}

// ReferenceValue of the association.
func (a Association) ReferenceValue() interface{} {
	return indirectInterface(reflectValueFieldByIndex(a.rv, a.data.referenceIndex, false))
}

// ForeignField of the association.
func (a Association) ForeignField() string {
	return a.data.foreignField
}

// ForeignValue of the association.
// It'll panic if association type is has many.
func (a Association) ForeignValue() interface{} {
	if a.Type() == HasMany {
		panic("rel: cannot infer foreign value for has many or many to many association")
	}

	var (
		rv = reflectValueFieldByIndex(a.rv, a.data.targetIndex, false)
	)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return indirectInterface(reflectValueFieldByIndex(rv, a.data.foreignIndex, false))
}

// Through return intermediary association.
func (a Association) Through() string {
	return a.data.through
}

// Autoload assoc setting when parent is loaded.
func (a Association) Autoload() bool {
	return a.data.autoload
}

// Autosave setting when parent is created/updated/deleted.
func (a Association) Autosave() bool {
	return a.data.autosave
}

func newAssociation(rv reflect.Value, index []int) Association {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return Association{
		data: extractAssociationData(rv.Type(), index),
		rv:   rv,
	}
}

func extractAssociationData(rt reflect.Type, index []int) associationData {
	var (
		key = associationKey{
			rt:    rt,
			index: encodeIndices(index),
		}
	)

	if val, cached := associationCache.Load(key); cached {
		return val.(associationData)
	}

	var (
		sf        = rt.FieldByIndex(index)
		ft        = sf.Type
		ref       = sf.Tag.Get("ref")
		fk        = sf.Tag.Get("fk")
		fName, _  = fieldName(sf)
		assocData = associationData{
			targetIndex: index,
			through:     sf.Tag.Get("through"),
			autoload:    sf.Tag.Get("auto") == "true" || sf.Tag.Get("autoload") == "true",
			autosave:    sf.Tag.Get("auto") == "true" || sf.Tag.Get("autosave") == "true",
		}
	)

	if assocData.autosave && assocData.through != "" {
		panic("rel: autosave is not supported for has one/has many through association")
	}

	for ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Slice {
		ft = ft.Elem()
	}

	var (
		refDocData = extractDocumentData(rt, true)
		fkDocData  = extractDocumentData(ft, true)
	)

	// Try to guess ref and fk if not defined.
	if ref == "" || fk == "" {
		// TODO: replace "id" with inferred primary field
		if assocData.through != "" {
			ref = "id"
			fk = "id"
		} else if _, isBelongsTo := refDocData.index[fName+"_id"]; isBelongsTo {
			ref = fName + "_id"
			fk = "id"
		} else {
			ref = "id"
			fk = snaker.CamelToSnake(rt.Name()) + "_id"
		}
	}

	if id, exist := refDocData.index[ref]; !exist {
		panic("rel: references (" + ref + ") field not found ")
	} else {
		assocData.referenceIndex = id
		assocData.referenceField = ref
	}

	if id, exist := fkDocData.index[fk]; !exist {
		panic("rel: foreign_key (" + fk + ") field not found")
	} else {
		assocData.foreignIndex = id
		assocData.foreignField = fk
	}

	// guess assoc type
	if sf.Type.Kind() == reflect.Slice || (sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Slice) {
		assocData.typ = HasMany
	} else {
		if len(assocData.referenceField) > len(assocData.foreignField) {
			assocData.typ = BelongsTo
		} else {
			assocData.typ = HasOne
		}
	}

	associationCache.Store(key, assocData)

	return assocData
}
