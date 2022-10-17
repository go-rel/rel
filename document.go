package rel

import (
	"database/sql"
	"reflect"
	"strings"
)

// Document provides an abstraction over reflect to easily works with struct for database purpose.
type Document struct {
	v    any
	rv   reflect.Value
	rt   reflect.Type
	meta DocumentMeta
}

// ReflectValue of referenced document.
func (d Document) ReflectValue() reflect.Value {
	return d.rv
}

// Table returns name of the table.
func (d Document) Table() string {
	// TODO: handle anonymous struct
	return d.meta.Table()
}

// PrimaryFields column name of this document.
func (d Document) PrimaryFields() []string {
	return d.meta.PrimaryFields()
}

// PrimaryField column name of this document.
// panic if document uses composite key.
func (d Document) PrimaryField() string {
	return d.meta.PrimaryField()
}

// PrimaryValues of this document.
func (d Document) PrimaryValues() []any {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryValues()
	}

	if len(d.meta.primaryIndex) == 0 {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	var (
		pValues = make([]any, len(d.meta.primaryIndex))
	)

	for i := range pValues {
		pValues[i] = reflectValueFieldByIndex(d.rv, d.meta.primaryIndex[i], false).Interface()
	}

	return pValues
}

// PrimaryValue of this document.
// panic if document uses composite key.
func (d Document) PrimaryValue() any {
	if values := d.PrimaryValues(); len(values) == 1 {
		return values[0]
	}

	panic("rel: composite primary key is not supported")
}

// Persisted returns true if document primary key is not zero.
func (d Document) Persisted() bool {
	var (
		pValues = d.PrimaryValues()
	)

	for i := range pValues {
		if !isZero(pValues[i]) {
			return true
		}
	}

	return false
}

// Index returns map of column name and it's struct index.
func (d Document) Index() map[string][]int {
	return d.meta.Index()
}

// Fields returns list of fields available on this document.
func (d Document) Fields() []string {
	return d.meta.Fields()
}

// Type returns reflect.Type of given field. if field does not exist, second returns value will be false.
func (d Document) Type(field string) (reflect.Type, bool) {
	return d.meta.Type(field)
}

// Value returns value of given field. if field does not exist, second returns value will be false.
func (d Document) Value(field string) (any, bool) {
	if i, ok := d.meta.index[field]; ok {

		var (
			value any
			fv    = reflectValueFieldByIndex(d.rv, i, false)
			ft    = fv.Type()
		)

		if ft.Kind() == reflect.Ptr {
			if !fv.IsNil() {
				value = fv.Elem().Interface()
			}
		} else {
			value = fv.Interface()
		}

		return value, true
	}

	return nil, false
}

// SetValue of the field, it returns false if field does not exist, or it's not assignable.
func (d Document) SetValue(field string, value any) bool {
	if i, ok := d.meta.index[field]; ok {
		var (
			rv reflect.Value
			rt reflect.Type
			fv = reflectValueFieldByIndex(d.rv, i, true)
			ft = fv.Type()
		)

		switch v := value.(type) {
		case nil:
			rv = reflect.Zero(ft)
		case reflect.Value:
			rv = reflect.Indirect(v)
		default:
			rv = reflect.Indirect(reflect.ValueOf(value))
		}

		rt = rv.Type()

		if fv.Type() == rt || rt.AssignableTo(ft) {
			fv.Set(rv)
			return true
		}

		if rt.ConvertibleTo(ft) {
			return setConvertValue(ft, fv, rt, rv)
		}

		if ft.Kind() == reflect.Ptr {
			return setPointerValue(ft, fv, rt, rv)
		}
	}

	return false
}

// Scanners returns slice of sql.Scanner for given fields.
func (d Document) Scanners(fields []string) []any {
	var (
		result    = make([]any, len(fields))
		assocRefs map[string]struct {
			fields  []string
			indexes []int
		}
	)

	for index, field := range fields {
		if structIndex, ok := d.meta.index[field]; ok {
			var (
				fv = reflectValueFieldByIndex(d.rv, structIndex, true)
				ft = fv.Type()
			)

			if ft.Kind() == reflect.Ptr {
				result[index] = fv.Addr().Interface()
			} else {
				result[index] = Nullable(fv.Addr().Interface())
			}
		} else if split := strings.SplitN(field, ".", 2); len(split) == 2 {
			if assocRefs == nil {
				assocRefs = make(map[string]struct {
					fields  []string
					indexes []int
				})
			}

			refs := assocRefs[split[0]]
			refs.fields = append(refs.fields, split[1])
			refs.indexes = append(refs.indexes, index)
			assocRefs[split[0]] = refs
		} else {
			result[index] = &sql.RawBytes{}
		}
	}

	// get scanners from associations
	for assocName, refs := range assocRefs {
		if assoc, ok := d.association(assocName); ok && assoc.Type() == BelongsTo || assoc.Type() == HasOne {
			var (
				assocDoc, _   = assoc.Document()
				assocScanners = assocDoc.Scanners(refs.fields)
			)

			for i, index := range refs.indexes {
				result[index] = assocScanners[i]
			}
		} else {
			for _, index := range refs.indexes {
				result[index] = &sql.RawBytes{}
			}
		}
	}

	return result
}

// BelongsTo fields of this document.
func (d Document) BelongsTo() []string {
	return d.meta.BelongsTo()
}

// HasOne fields of this document.
func (d Document) HasOne() []string {
	return d.meta.HasOne()
}

// HasMany fields of this document.
func (d Document) HasMany() []string {
	return d.meta.HasMany()
}

// Preload fields of this document.
func (d Document) Preload() []string {
	return d.meta.Preload()
}

// Association of this document with given name.
func (d Document) Association(name string) Association {
	if assoc, ok := d.association(name); ok {
		return assoc
	}

	panic("rel: no field named (" + name + ") in type " + d.rt.String() + " found ")
}

func (d Document) association(name string) (Association, bool) {
	index, ok := d.meta.index[name]
	if !ok {
		return Association{}, false
	}

	return newAssociation(d.rv, index), true
}

// Reset this document, this is a noop for compatibility with collection.
func (d Document) Reset() {
}

// Add returns this document.
func (d *Document) Add() *Document {
	// if d.rv is a null pointer, set it to a new struct.
	if d.rv.Kind() == reflect.Ptr && d.rv.IsNil() {
		d.rv.Set(reflect.New(d.rv.Type().Elem()))
		d.rv = d.rv.Elem()
	}

	return d
}

// Get always returns this document, this is a noop for compatibility with collection.
func (d *Document) Get(index int) *Document {
	return d
}

// Len always returns 1 for document, this is a noop for compatibility with collection.
func (d *Document) Len() int {
	return 1
}

// Meta returns document meta.
func (d Document) Meta() DocumentMeta {
	return d.meta
}

// Flag returns true if struct contains specified flag.
func (d Document) Flag(flag DocumentFlag) bool {
	return d.meta.Flag(flag)
}

// NewDocument used to create abstraction to work with struct.
// Document can be created using interface or reflect.Value.
func NewDocument(entity any, readonly ...bool) *Document {
	switch v := entity.(type) {
	case *Document:
		return v
	case reflect.Value:
		return newDocument(v.Interface(), v, len(readonly) > 0 && readonly[0])
	case reflect.Type:
		panic("rel: cannot use reflect.Type")
	case nil:
		panic("rel: cannot be nil")
	default:
		return newDocument(v, reflect.ValueOf(v), len(readonly) > 0 && readonly[0])
	}
}

func newDocument(v any, rv reflect.Value, readonly bool) *Document {
	var (
		rt = rv.Type()
	)

	if rt.Kind() != reflect.Ptr {
		if !readonly {
			panic("rel: must be a pointer to struct")
		}
	} else {
		if !rv.IsNil() {
			rv = rv.Elem()
		}
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("rel: must be a struct or pointer to a struct")
	}

	return &Document{
		v:    v,
		rv:   rv,
		rt:   rt,
		meta: getDocumentMeta(rt, false),
	}
}
