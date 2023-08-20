package rel

import (
	"reflect"
)

type slice interface {
	table
	Reset()
	NewDocument() *Document
	Append(doc *Document)
	Get(index int) *Document
	Len() int
	Meta() DocumentMeta
}

// Collection provides an abstraction over reflect to easily works with slice for database purpose.
type Collection struct {
	v       any
	rv      reflect.Value
	rt      reflect.Type
	meta    DocumentMeta
	swapper func(i, j int)
}

// ReflectValue of referenced document.
func (c Collection) ReflectValue() reflect.Value {
	return c.rv
}

// Table returns name of the table.
func (c *Collection) Table() string {
	return c.meta.Table()
}

// PrimaryFields column name of this collection.
func (c Collection) PrimaryFields() []string {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryFields()
	}

	if len(c.meta.primaryField) == 0 {
		panic("rel: failed to infer primary key for type " + c.rt.String())
	}

	return c.meta.primaryField
}

// PrimaryField column name of this document.
// panic if document uses composite key.
func (c Collection) PrimaryField() string {
	if fields := c.PrimaryFields(); len(fields) == 1 {
		return fields[0]
	}

	panic("rel: composite primary key is not supported")
}

// PrimaryValues of collection.
// Returned value will be interface of slice interface.
func (c Collection) PrimaryValues() []any {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryValues()
	}

	var (
		index   = c.meta.primaryIndex
		pValues = make([]any, len(c.PrimaryFields()))
	)

	if index != nil {
		for i := range index {
			var (
				idxLen = c.rv.Len()
				values = make([]any, 0, idxLen)
			)

			for j := 0; j < idxLen; j++ {
				if item := c.rvIndex(j); item.IsValid() {
					values = append(values, reflectValueFieldByIndex(item, index[i], false).Interface())
				}
			}

			pValues[i] = values
		}
	} else {
		// using interface.
		var (
			tmp = make([][]any, len(pValues))
		)

		for i := 0; i < c.rv.Len(); i++ {
			item := c.rvIndex(i)
			if !item.IsValid() {
				continue
			}
			for j, id := range item.Interface().(primary).PrimaryValues() {
				tmp[j] = append(tmp[j], id)
			}
		}

		for i := range tmp {
			pValues[i] = tmp[i]
		}
	}

	return pValues
}

// PrimaryValue of this document.
// panic if document uses composite key.
func (c Collection) PrimaryValue() any {
	if values := c.PrimaryValues(); len(values) == 1 {
		return values[0]
	}

	panic("rel: composite primary key is not supported")
}

func (c Collection) rvIndex(index int) reflect.Value {
	return reflect.Indirect(c.rv.Index(index))
}

// Get an element from the underlying slice as a document.
func (c Collection) Get(index int) *Document {
	return NewDocument(c.rvIndex(index).Addr())
}

// Len of the underlying slice.
func (c Collection) Len() int {
	return c.rv.Len()
}

// Meta returns document meta.
func (c Collection) Meta() DocumentMeta {
	return c.meta
}

// Reset underlying slice to be zero length.
func (c Collection) Reset() {
	c.rv.Set(reflect.MakeSlice(c.rt, 0, 0))
}

// Add new document into collection.
func (c Collection) Add() *Document {
	c.Append(c.NewDocument())
	return c.Get(c.Len() - 1)
}

// NewDocument returns new document with zero values.
func (c Collection) NewDocument() *Document {
	return newZeroDocument(c.rt.Elem())
}

// Append new document into collection.
func (c Collection) Append(doc *Document) {
	if c.rt.Elem().Kind() == reflect.Ptr {
		c.rv.Set(reflect.Append(c.rv, doc.rv.Addr()))
	} else {
		c.rv.Set(reflect.Append(c.rv, doc.rv))
	}
}

// Truncate collection.
func (c Collection) Truncate(i, j int) {
	c.rv.Set(c.rv.Slice(i, j))
}

// Slice returns a new collection that is a slice of the original collection.s
func (c Collection) Slice(i, j int) *Collection {
	return NewCollection(c.rv.Slice(i, j), true)
}

// Swap element in the collection.
func (c Collection) Swap(i, j int) {
	if c.swapper == nil {
		c.swapper = reflect.Swapper(c.rv.Interface())
	}

	c.swapper(i, j)
}

// NewCollection used to create abstraction to work with slice.
// COllection can be created using interface or reflect.Value.
func NewCollection(entities any, readonly ...bool) *Collection {
	switch v := entities.(type) {
	case *Collection:
		return v
	case reflect.Value:
		return newCollection(v.Interface(), v, len(readonly) > 0 && readonly[0])
	case reflect.Type:
		panic("rel: cannot use reflect.Type")
	case nil:
		panic("rel: cannot be nil")
	default:
		return newCollection(v, reflect.ValueOf(v), len(readonly) > 0 && readonly[0])
	}
}

func newCollection(v any, rv reflect.Value, readonly bool) *Collection {
	var (
		rt = rv.Type()
	)

	if rt.Kind() != reflect.Ptr {
		if !readonly {
			panic("rel: must be a pointer to slice")
		}
	} else {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Slice {
		panic("rel: must be a slice or pointer to a slice")
	}

	return &Collection{
		v:    v,
		rv:   rv,
		rt:   rt,
		meta: getDocumentMeta(indirectReflectType(rt.Elem()), false),
	}
}
