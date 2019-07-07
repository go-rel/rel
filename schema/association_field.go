package schema

import (
	"reflect"
	"sync"
)

var associationFieldCache sync.Map

type associationFieldKey struct {
	rt    reflect.Type
	field string
}

type AssociationField struct {
	ReferenceIndex  []int
	ForeignIndex    []int
	ForeignColumn   string
	ReferenceColumn string
}

// InferAssociationField from a struct type.
func InferAssociationField(rt reflect.Type, field string) AssociationField {
	key := associationFieldKey{
		rt:    rt,
		field: field,
	}

	if val, cached := associationFieldCache.Load(key); cached {
		return val.(AssociationField)
	}

	st, exist := rt.FieldByName(field)
	if !exist {
		panic("grimoire: field named (" + field + ") not found ")
	}

	var (
		ft   = st.Type
		ref  = st.Tag.Get("references")
		fk   = st.Tag.Get("foreign_key")
		data = AssociationField{}
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
		data.ReferenceIndex = reft.Index
		data.ReferenceColumn = inferFieldColumn(reft)
	}

	if fkt, exist := ft.FieldByName(fk); !exist {
		panic("grimoire: foreign_key (" + fk + ") field not found " + fk)
	} else {
		data.ForeignIndex = fkt.Index
		data.ForeignColumn = inferFieldColumn(fkt)
	}

	associationFieldCache.Store(key, data)
	return data
}
