package schema

import (
	"reflect"
	"sync"
)

var associationCache sync.Map

type associationKey struct {
	rt    reflect.Type
	field string
}

type association struct {
	refIndex []int
	fkIndex  []int
	column   string
}

// InferAssociation from a field in a struct type.
func InferAssociation(rt reflect.Type, field string) ([]int, []int, string) {
	key := associationKey{
		rt:    rt,
		field: field,
	}

	if val, cached := associationCache.Load(key); cached {
		data := val.(association)
		return data.refIndex, data.fkIndex, data.column
	}

	st, exist := rt.FieldByName(field)
	if !exist {
		panic("grimoire: field named (" + field + ") not found ")
	}

	var (
		ft   = st.Type
		ref  = st.Tag.Get("references")
		fk   = st.Tag.Get("foreign_key")
		data = association{}
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
		data.refIndex = reft.Index
	}

	if fkt, exist := ft.FieldByName(fk); !exist {
		panic("grimoire: foreign_key (" + fk + ") field not found " + fk)
	} else {
		data.fkIndex = fkt.Index
		data.column = inferFieldName(fkt)
	}

	associationCache.Store(key, data)
	return data.refIndex, data.fkIndex, data.column
}
