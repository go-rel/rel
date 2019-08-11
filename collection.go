package grimoire

import (
	"reflect"
)

type Collection interface {
	table
	primary
	Reset()
	Add() Document
	Get(index int) Document
	Len() int
}

type collection struct {
	v  interface{}
	rv reflect.Value
	rt reflect.Type
}

func (c *collection) reflect() {
	if c.rv.IsValid() {
		return
	}

	c.rv = reflect.ValueOf(c.v)
	if c.rv.Kind() != reflect.Ptr {
		panic("grimoire: must be a pointer")
	}

	c.rv = c.rv.Elem()
	c.rt = c.rv.Type()

	if c.rt.Kind() != reflect.Slice {
		panic("grimoire: must be a pointer to a slice")
	}
}

func (c *collection) Table() string {
	if tn, ok := c.v.(table); ok {
		return tn.Table()
	}

	c.reflect()

	return tableName(c.rt.Elem())
}

func (c *collection) PrimaryField() string {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryField()
	}

	c.reflect()

	var (
		field, _ = searchPrimary(c.rt.Elem())
	)

	return field
}

func (c *collection) PrimaryValue() interface{} {
	if p, ok := c.v.(primary); ok {
		return p.PrimaryField()
	}

	c.reflect()

	var (
		_, index = searchPrimary(c.rt.Elem())
		ids      = make([]interface{}, c.rv.Len())
	)

	for i := 0; i < len(ids); i++ {
		ids[i] = c.rv.Index(i).Field(index).Interface()
	}

	return ids
}

func (c *collection) Get(index int) Document {
	c.reflect()

	return newDocument(c.rv.Index(index).Interface())
}

func (c *collection) Len() int {
	c.reflect()

	return c.rv.Len()
}

func (c *collection) Reset() {
	// TODD
}

func (c *collection) Add() Document {
	c.reflect()

	var (
		index = c.Len()
		typ   = c.rt.Elem()
		drv   = reflect.Zero(typ)
	)

	c.rv.Set(reflect.Append(c.rv, drv))

	return newDocument(c.rv.Index(index).Addr().Interface())
}

func newCollection(entities interface{}) Collection {
	switch v := entities.(type) {
	case Collection:
		return v
	case reflect.Value:
		if v.Kind() != reflect.Ptr && v.Elem().Kind() != reflect.Slice {
			panic("grimoire: must be a pointer to a slice")
		}

		return &collection{
			v:  v.Interface(),
			rv: v.Elem(),
			rt: v.Elem().Type(),
		}
	default:
		return &collection{v: v}
	}
}
