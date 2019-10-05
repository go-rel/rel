package rel

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

// type Association interface {
// 	Type() AssociationType
// 	Document() (Document, bool)
// 	Collection() (Collection, bool)
// 	IsZero() bool
// 	ReferenceField() string
// 	ReferenceValue() interface{}
// 	ForeignField() string
// 	ForeignValue() interface{}
// }

type associationKey struct {
	rt    reflect.Type
	index int
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

type Association struct {
	data associationData
	rv   reflect.Value
}

func (a Association) Type() AssociationType {
	return a.data.typ
}

// // Replace Target with Document and Collection
// func (a Association) Target() (*Collection, bool) {
// 	var (
// 		rv = a.rv.FieldByIndex(a.data.targetIndex)
// 	)

// 	switch rv.Kind() {
// 	case reflect.Slice:
// 		var (
// 			loaded = !rv.IsNil()
// 		)

// 		if !loaded {
// 			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
// 		}

// 		return newCollection(rv.Addr()), loaded
// 	case reflect.Ptr:
// 		var (
// 			loaded = !rv.IsNil()
// 		)

// 		if !loaded {
// 			rv.Set(reflect.New(rv.Type().Elem()))
// 		}

// 		if rv.Elem().Kind() == reflect.Slice {
// 			rv.Elem().Set(reflect.MakeSlice(rv.Elem().Type(), 0, 0))

// 			return newCollection(rv), loaded
// 		}

// 		var (
// 			doc = newDocument(rv)
// 			id  = doc.PrimaryValue()
// 		)

// 		return doc, loaded && !isZero(id)
// 	default:
// 		var (
// 			doc = newDocument(rv.Addr())
// 			id  = doc.PrimaryValue()
// 		)

// 		return doc, !isZero(id)
// 	}
// }

func (a Association) Document() (*Document, bool) {
	var (
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	switch rv.Kind() {
	case reflect.Ptr:
		if rv.IsNil() {
			rv.Set(reflect.New(rv.Type().Elem()))
			return newDocument(rv), false
		}

		var (
			doc = newDocument(rv)
			id  = doc.PrimaryValue()
		)

		return doc, !isZero(id)
	default:
		var (
			doc = newDocument(rv.Addr())
			id  = doc.PrimaryValue()
		)

		return doc, !isZero(id)
	}
}

func (a Association) Collection() (*Collection, bool) {
	var (
		rv     = a.rv.FieldByIndex(a.data.targetIndex)
		loaded = !rv.IsNil()
	)

	if rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Slice {
		if !loaded {
			rv.Set(reflect.New(rv.Type().Elem()))
			rv.Elem().Set(reflect.MakeSlice(rv.Elem().Type(), 0, 0))
		}

		return newCollection(rv), loaded
	}

	if !loaded {
		rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
	}

	return newCollection(rv.Addr()), loaded
}

func (a Association) IsZero() bool {
	var (
		rv = a.rv.FieldByIndex(a.data.targetIndex)
	)

	return isDeepZero(rv, 1)
}

func (a Association) ReferenceField() string {
	return a.data.referenceColumn
}

func (a Association) ReferenceValue() interface{} {
	return indirect(a.rv.FieldByIndex(a.data.referenceIndex))
}

func (a Association) ForeignField() string {
	return a.data.foreignField
}

func (a Association) ForeignValue() interface{} {
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

	return indirect(rv.FieldByIndex(a.data.foreignIndex))
}

func newAssociation(rv reflect.Value, index int) Association {
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}

	var (
		rt  = rv.Type()
		key = associationKey{
			rt:    rt,
			index: index,
		}
	)

	if val, cached := associationCache.Load(key); cached {
		return Association{
			data: val.(associationData),
			rv:   rv,
		}
	}

	// TODO: maybe use column name instead of field name for ref and fk key
	var (
		st   = rt.Field(index)
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
		panic("rel: references (" + ref + ") field not found ")
	} else {
		data.referenceIndex = reft.Index
		data.referenceColumn = fieldName(reft)
	}

	if fkt, exist := ft.FieldByName(fk); !exist {
		panic("rel: foreign_key (" + fk + ") field not found " + fk)
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

	return Association{
		data: data,
		rv:   rv,
	}
}
