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

var (
	tableRt   = reflect.TypeOf((*table)(nil)).Elem()
	primaryRt = reflect.TypeOf((*primary)(nil)).Elem()
)

type Collection struct {
	v  interface{}
	rv reflect.Value
	rt reflect.Type
}

func (c *Collection) reflect() {
	if c.rv.IsValid() {
		return
	}

	c.rv = reflect.ValueOf(c.v)
	if c.rv.Kind() != reflect.Ptr {
		panic("rel: must be a pointer")
	}

	c.rv = c.rv.Elem()
	c.rt = c.rv.Type()

	if c.rt.Kind() != reflect.Slice {
		panic("rel: must be a pointer to a slice")
	}
}

func (c *Collection) Table() string {
	if tn, ok := c.v.(table); ok {
		return tn.Table()
	}

	return c.tableName()
}

func (c *Collection) tableName() string {
	c.reflect()

	var (
		rt = c.rt.Elem()
	)

	// check for cache
	if name, cached := tablesCache.Load(rt); cached {
		return name.(string)
	}

	if rt.Implements(tableRt) {
		var (
			v = reflect.Zero(rt).Interface().(table)
		)

		tablesCache.Store(rt, v.Table())
		return v.Table()
	}

	return tableName(rt)
}

func (c *Collection) PrimaryField() string {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryField()
	}

	var (
		field, _ = c.searchPrimary()
	)

	return field
}

func (c *Collection) PrimaryValue() interface{} {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryValue()
	}

	var (
		_, index = c.searchPrimary()
		ids      = make([]interface{}, c.rv.Len())
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

func (c *Collection) searchPrimary() (string, int) {
	c.reflect()

	var (
		rt = c.rt.Elem()
	)

	if result, cached := primariesCache.Load(rt); cached {
		p := result.(primaryData)
		return p.field, p.index
	}

	if rt.Implements(primaryRt) {
		var (
			v     = reflect.Zero(rt).Interface().(primary)
			field = v.PrimaryField()
			index = -2 // special index to mark interface usage
		)

		primariesCache.Store(rt, primaryData{
			field: field,
			index: index,
		})

		return field, index
	}

	return searchPrimary(rt)
}

func (c *Collection) Get(index int) *Document {
	c.reflect()

	return NewDocument(c.rv.Index(index).Addr().Interface())
}

func (c *Collection) Len() int {
	c.reflect()

	return c.rv.Len()
}

func (c *Collection) Reset() {
	c.reflect()

	c.rv.Set(reflect.Zero(c.rt))
}

func (c *Collection) Add() *Document {
	c.reflect()

	var (
		index = c.Len()
		typ   = c.rt.Elem()
		drv   = reflect.Zero(typ)
	)

	c.rv.Set(reflect.Append(c.rv, drv))

	return NewDocument(c.rv.Index(index).Addr().Interface())
}

func NewCollection(records interface{}) *Collection {
	switch v := records.(type) {
	case *Collection:
		return v
	case reflect.Value:
		if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
			panic("rel: must be a pointer to a slice")
		}

		return &Collection{
			v:  v.Interface(),
			rv: v.Elem(),
			rt: v.Elem().Type(),
		}
	case reflect.Type:
		panic("rel: cannot use reflect.Type")
	case nil:
		panic("rel: cannot be nil")
	default:
		return &Collection{v: v}
	}
}
