package rel

import (
	"reflect"
)

type slice interface {
	table
	Reset()
	Add() *Document
	Get(index int) *Document
	Len() int
}

// Collection provides an abstraction over reflect to easily works with slice for database purpose.
type Collection struct {
	v       interface{}
	rv      reflect.Value
	rt      reflect.Type
	data    documentData
	index   map[interface{}]int
	swapper func(i, j int)
}

// ReflectValue of referenced document.
func (c Collection) ReflectValue() reflect.Value {
	return c.rv
}

// Table returns name of the table.
func (c *Collection) Table() string {
	if tn, ok := c.v.(table); ok {
		return tn.Table()
	}

	return c.tableName()
}

func (c Collection) tableName() string {
	var (
		rt = c.rt.Elem()
	)

	// check for cache
	if name, cached := tablesCache.Load(rt); cached {
		return name.(string)
	}

	if rt.Implements(rtTable) {
		var (
			v = reflect.Zero(rt).Interface().(table)
		)

		tablesCache.Store(rt, v.Table())
		return v.Table()
	}

	return tableName(rt)
}

// PrimaryField column name of this collection.
func (c Collection) PrimaryField() string {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryField()
	}

	if c.data.primaryField == "" {
		panic("rel: failed to infer primary key for type " + c.rt.String())
	}

	return c.data.primaryField
}

// PrimaryValue of collection.
// Returned value will be interface of slice interface.
func (c Collection) PrimaryValue() interface{} {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryValue()
	}

	var (
		index = c.data.primaryIndex
		ids   = make([]interface{}, c.rv.Len())
	)

	for i := 0; i < len(ids); i++ {
		var (
			fv = c.rv.Index(i)
		)

		if index == -2 {
			// using interface
			ids[i] = fv.Interface().(primary).PrimaryValue()
		} else {
			ids[i] = fv.Field(index).Interface()
		}
	}

	return ids
}

// Get an element from the underlying slice as a document.
func (c Collection) Get(index int) *Document {
	return NewDocument(c.rv.Index(index).Addr())
}

// Len of the underlying slice.
func (c Collection) Len() int {
	return c.rv.Len()
}

// Reset underlying slice to be zero length.
func (c Collection) Reset() {
	c.rv.Set(reflect.MakeSlice(c.rt, 0, 0))
}

// Add new document into collection.
func (c Collection) Add() *Document {
	var (
		index = c.Len()
		typ   = c.rt.Elem()
		drv   = reflect.Zero(typ)
	)

	c.rv.Set(reflect.Append(c.rv, drv))

	return NewDocument(c.rv.Index(index).Addr())
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
func NewCollection(records interface{}, readonly ...bool) *Collection {
	switch v := records.(type) {
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

func newCollection(v interface{}, rv reflect.Value, readonly bool) *Collection {
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
		data: extractDocumentData(rt.Elem(), false),
	}
}
