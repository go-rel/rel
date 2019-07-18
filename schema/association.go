package schema

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
	TargetAddr() (interface{}, bool)
	ReferenceColumn() string
	ReferenceValue() interface{}
	ForeignColumn() string
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
	foreignColumn   string
	foreignIndex    []int
}

var associationCache sync.Map

type association struct {
	data  associationData
	value reflect.Value
}

func (a association) Type() AssociationType {
	return a.data.typ
}

func (a association) TargetAddr() (interface{}, bool) {
	var (
		rv = a.value.FieldByIndex(a.data.targetIndex)
	)

	switch rv.Kind() {
	case reflect.Slice:
		if rv.IsNil() {
			rv.Set(reflect.MakeSlice(rv.Type(), 0, 0))
		}

		return rv.Addr().Interface(), rv.Len() != 0
	case reflect.Ptr:
		var (
			loaded = !rv.IsNil()
		)

		if !loaded {
			rv.Set(reflect.New(rv.Type().Elem()))
		}

		return rv.Interface(), loaded
	default:
		var (
			target = rv.Addr().Interface()
			_, pv  = InferPrimaryKey(target, true)
		)

		return target, !isZero(pv[0])
	}
}

func (a association) ReferenceColumn() string {
	return a.data.referenceColumn
}

func (a association) ReferenceValue() interface{} {
	return a.value.FieldByIndex(a.data.referenceIndex).Interface()
}

func (a association) ForeignColumn() string {
	return a.data.foreignColumn
}

func (a association) ForeignValue() interface{} {
	if a.Type() == HasMany {
		panic("cannot infer foreign value for has many association")
	}

	var (
		rv = a.value.FieldByIndex(a.data.targetIndex)
	)

	if rv.Kind() == reflect.Ptr {
		if rv.IsNil() {
			return nil
		}

		rv = rv.Elem()
	}

	return rv.FieldByIndex(a.data.foreignIndex).Interface()
}

func InferAssociation(rv reflect.Value, name string) Association {
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
			data:  val.(associationData),
			value: rv,
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
		data.referenceColumn = inferFieldColumn(reft)
	}

	if fkt, exist := ft.FieldByName(fk); !exist {
		panic("grimoire: foreign_key (" + fk + ") field not found " + fk)
	} else {
		data.foreignIndex = fkt.Index
		data.foreignColumn = inferFieldColumn(fkt)
	}

	// guess assoc type
	if st.Type.Kind() == reflect.Slice || st.Type.Kind() == reflect.Array {
		data.typ = HasMany
	} else {
		if typ == "belongs_to" || len(data.referenceColumn) > len(data.foreignColumn) {
			data.typ = BelongsTo
		} else {
			data.typ = HasOne
		}
	}

	associationCache.Store(key, data)

	return association{
		data:  data,
		value: rv,
	}
}

type Associations interface {
	BelongsTo() []string
	HasOne() []string
	HasMany() []string
	Association(field string) Association
}

var associationsFieldCache sync.Map

type associationsField struct {
	belongsTo []string
	hasOne    []string
	hasMany   []string
}

type associations struct {
	rv     reflect.Value
	fields associationsField
}

func (a associations) BelongsTo() []string {
	return a.fields.belongsTo
}

func (a associations) HasOne() []string {
	return a.fields.hasOne
}

func (a associations) HasMany() []string {
	return a.fields.hasMany
}

func (a associations) Association(name string) Association {
	return InferAssociation(a.rv, name)
}

func InferAssociations(record interface{}) Associations {
	if s, ok := record.(Associations); ok {
		return s
	}

	var (
		rv = reflectValuePtr(record)
		rt = rv.Type()
	)

	if fields, cached := associationsFieldCache.Load(rt); cached {
		return associations{
			rv:     rv,
			fields: fields.(associationsField),
		}
	}

	var (
		fields associationsField
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = sf.Name
			typ  = sf.Type
		)

		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface || typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}

		if typ.Kind() != reflect.Struct {
			continue
		}

		// must have a primary key
		if pk, _ := searchPrimaryKey(typ); pk == "" {
			continue
		}

		switch InferAssociation(rv, name).Type() {
		case BelongsTo:
			fields.belongsTo = append(fields.belongsTo, name)
		case HasOne:
			fields.hasOne = append(fields.hasOne, name)
		case HasMany:
			fields.hasMany = append(fields.hasMany, name)
		}
	}

	associationsFieldCache.Store(rt, fields)

	return associations{
		rv:     rv,
		fields: fields,
	}
}
