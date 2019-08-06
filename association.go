package grimoire

import (
	"reflect"
	"sync"
)

type AssociationType uint8

const (
	BelongsTo = iota
	HasOne
	HasMany
)

type Association interface {
	Type() AssociationType
	Target() (Collection, bool)
	ReferenceField() string
	ReferenceValue() interface{}
	ForeignField() string
	ForeignValue() interface{}
}

type associationKey struct {
	rt   reflect.Type
	name string
}

type associationData struct {
	typ             AssociationType
	targetIndex     []int
	referenceColumn string
	referenceIndex  []int
	foreignField    string
	foreignIndex    []int
}

var associationCache sync.Map

type association struct {
	data associationData
	rv   reflect.Value
}

func (a association) Type() AssociationType {
	return a.data.typ
}

func (a association) Target() (Collection, bool) {
	var (
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	switch rv.Kind() {
	case reflect.Slice:
		var (
			loaded = !rv.IsNil()
		)

		if !loaded {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
		}

		return newCollection(rv.Addr()), loaded
	case reflect.Ptr:
		var (
			loaded = !rv.IsNil()
		)

		if !loaded {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		if rv.Elem().Kind() == reflect.Slice {
			rv.Elem().Set(reflect.MakeSlice(rv.Elem().Type(), 0, 0))

			return newCollection(rv), loaded
		}

		return newDocument(rv), loaded
	default:
		var (
			doc = newDocument(rv.Addr())
			id  = doc.PrimaryValue()
		)

		return doc, !isZero(id)
	}
}

func (a association) ReferenceField() string {
	return a.data.referenceColumn
}

func (a association) ReferenceValue() interface{} {
	return a.rv.FieldByIndex(a.data.referenceIndex).Interface()
}

func (a association) ForeignField() string {
	return a.data.foreignField
}

func (a association) ForeignValue() interface{} {
	if a.Type() == HasMany {
		panic("cannot infer foreign value for has many association")
	}

	var (
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
	}

	return rv.FieldByIndex(a.data.foreignIndex).Interface()
}

func newAssociation(rv reflect.Value, name string) Association {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	var (
		rt  = rv.Type()
		key = associationKey{
			rt:   rt,
			name: name,
		}
	)

	if val, cached := associationCache.Load(key); cached {
		return association{
			data: val.(associationData),
			rv:   rv,
		}
	}

	st, exist := rt.FieldByName(name)
	if !exist {
		panic("grimoire: field named (" + name + ") not found ")
	}

	var (
		ft   = st.Type
		ref  = st.Tag.Get("references")
		fk   = st.Tag.Get("foreign_key")
		typ  = st.Tag.Get("association")
		data = associationData{
			targetIndex: st.Index,
		}
	)

	if ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Slice || ft.Kind() == reflect.Array {
		ft = ft.Elem()
	}

	// Try to guess ref and fk if not defined.
	if ref == "" || fk == "" {
		if _, isBelongsTo := rt.FieldByName(st.Name + "ID"); isBelongsTo {
			ref = st.Name + "ID"
			fk = "ID"
		} else {
			ref = "ID"
			fk = rt.Name() + "ID"
		}
	}

	if reft, exist := rt.FieldByName(ref); !exist {
		panic("grimoire: references (" + ref + ") field not found ")
	} else {
		data.referenceIndex = reft.Index
		data.referenceColumn = fieldName(reft)
	}

	if fkt, exist := ft.FieldByName(fk); !exist {
		panic("grimoire: foreign_key (" + fk + ") field not found " + fk)
	} else {
		data.foreignIndex = fkt.Index
		data.foreignField = fieldName(fkt)
	}

	// guess assoc type
	if st.Type.Kind() == reflect.Slice || st.Type.Kind() == reflect.Array {
		data.typ = HasMany
	} else {
		if typ == "belongs_to" || len(data.referenceColumn) > len(data.foreignField) {
			data.typ = BelongsTo
		} else {
			data.typ = HasOne
		}
	}

	associationCache.Store(key, data)

	return association{
		data: data,
		rv:   rv,
	}
}
