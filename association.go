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
	rt    reflect.Type
	index int
}

type associationData struct {
	typ             AssociationType
	targetIndex     []int
	referenceColumn string
	referenceIndex  int
	foreignField    string
	foreignIndex    int
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
	var (
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
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
		rv     = a.rv.FieldByIndex(a.data.targetIndex)
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
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	return isDeepZero(reflect.Indirect(rv), 1)
}

// ReferenceField of the association.
func (a Association) ReferenceField() string {
	return a.data.referenceColumn
}

// ReferenceValue of the association.
func (a Association) ReferenceValue() interface{} {
	return indirect(a.rv.Field(a.data.referenceIndex))
}

// ForeignField of the association.
func (a Association) ForeignField() string {
	return a.data.foreignField
}

// ForeignValue of the association.
// It'll panic if association type is has many.
func (a Association) ForeignValue() interface{} {
	if a.Type() == HasMany {
		panic("cannot infer foreign value for has many association")
	}

	var (
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return indirect(rv.Field(a.data.foreignIndex))
}

func newAssociation(rv reflect.Value, index int) Association {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	return Association{
		data: extractAssociationData(rv.Type(), index),
		rv:   rv,
	}
}

func extractAssociationData(rt reflect.Type, index int) associationData {
	var (
		key = associationKey{
			rt:    rt,
			index: index,
		}
	)

	if val, cached := associationCache.Load(key); cached {
		return val.(associationData)
	}

	var (
		sf        = rt.Field(index)
		ft        = sf.Type
		ref       = sf.Tag.Get("ref")
		fk        = sf.Tag.Get("fk")
		fName     = fieldName(sf)
		assocData = associationData{
			targetIndex: sf.Index,
		}
	)

	for ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Slice {
		ft = ft.Elem()
	}

	var (
		refDocData = extractDocumentData(rt, true)
		fkDocData  = extractDocumentData(ft, true)
	)

	// Try to guess ref and fk if not defined.
	if ref == "" || fk == "" {
		if _, isBelongsTo := refDocData.index[fName+"_id"]; isBelongsTo {
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
		assocData.referenceColumn = ref
	}

	if id, exist := fkDocData.index[fk]; !exist {
		panic("rel: foreign_key (" + fk + ") field not found")
	} else {
		assocData.foreignIndex = id
		assocData.foreignField = fk
	}

	// guess assoc type
	if sf.Type.Kind() == reflect.Slice ||
		(sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Slice) {
		assocData.typ = HasMany
	} else {
		if len(assocData.referenceColumn) > len(assocData.foreignField) {
			assocData.typ = BelongsTo
		} else {
			assocData.typ = HasOne
		}
	}

	associationCache.Store(key, assocData)

	return assocData
}
