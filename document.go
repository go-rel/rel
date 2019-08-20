package grimoire

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"

	"github.com/azer/snakecase"
	"github.com/jinzhu/inflection"
)

var (
	tablesCache       sync.Map
	primariesCache    sync.Map
	fieldsCache       sync.Map
	fieldMappingCache sync.Map
	typesCache        sync.Map
	associationsCache sync.Map
)

type table interface {
	Table() string
}

type primary interface {
	PrimaryField() string
	PrimaryValue() interface{}
}

type primaryData struct {
	field string
	index int
}

type fields interface {
	Fields() map[string]int
}

type types interface {
	Types() []reflect.Type
}

type scanners interface {
	Scanners([]string) []interface{}
}

type values interface {
	Values() []interface{}
}

type associations interface {
	BelongsTo() []string
	HasOne() []string
	HasMany() []string
	Association(field string) Association
}

type associationsData struct {
	belongsTo []string
	hasOne    []string
	hasMany   []string
}

type Document interface {
	table
	primary
	fields
	types
	scanners
	values
	associations
	slice
}

type document struct {
	v         interface{}
	rv        reflect.Value
	rt        reflect.Type
	belongsTo []string
	hasOne    []string
	hasMany   []string
}

func (d *document) reflect() {
	if d.rv.IsValid() {
		return
	}

	d.rv = reflect.ValueOf(d.v)
	if d.rv.Kind() != reflect.Ptr {
		panic("grimoire: must be a pointer")
	}

	d.rv = d.rv.Elem()
	d.rt = d.rv.Type()

	if d.rt.Kind() != reflect.Struct {
		panic("grimoire: must be a pointer to a struct")
	}
}

func (d *document) Table() string {
	if tn, ok := d.v.(table); ok {
		return tn.Table()
	}

	d.reflect()

	// TODO: handle anonymous struct
	return tableName(d.rt)
}

func (d *document) PrimaryField() string {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryField()
	}

	d.reflect()

	var (
		field, _ = searchPrimary(d.rt)
	)

	if field == "" {
		panic("grimoire: failed to infer primary key for type " + d.rt.String())
	}

	return field
}

func (d *document) PrimaryValue() interface{} {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryValue()
	}

	d.reflect()

	var (
		_, index = searchPrimary(d.rt)
	)

	if index < 0 {
		panic("grimoire: failed to infer primary key for type " + d.rt.String())
	}

	return d.rv.Field(index).Interface()
}

func (d *document) Fields() map[string]int {
	if s, ok := d.v.(fields); ok {
		return s.Fields()
	}

	d.reflect()

	// check for cache
	if v, cached := fieldsCache.Load((d.rt)); cached {
		return v.(map[string]int)
	}

	var (
		index  = 0
		fields = make(map[string]int, d.rt.NumField())
	)

	for i := 0; i < d.rt.NumField(); i++ {
		var (
			sf   = d.rt.Field(i)
			name = fieldName(sf)
		)

		if name != "" {
			fields[name] = index
			index++
		}
	}

	fieldsCache.Store(d.rt, fields)

	return fields
}

func (d *document) Types() []reflect.Type {
	if v, ok := d.v.(types); ok {
		return v.Types()
	}

	d.reflect()

	// check for cache
	if v, cached := typesCache.Load(d.rt); cached {
		return v.([]reflect.Type)
	}

	var (
		fields  = d.Fields()
		mapping = fieldMapping(d.rt)
		types   = make([]reflect.Type, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			ft          = d.rt.Field(structIndex).Type
		)

		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		} else if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
			ft = reflect.SliceOf(ft.Elem().Elem())
		}

		types[index] = ft
	}

	typesCache.Store(d.rt, types)

	return types
}

func (d *document) Scanners(fields []string) []interface{} {
	if v, ok := d.v.(scanners); ok {
		return v.Scanners(fields)
	}

	if s, ok := d.v.(sql.Scanner); ok {
		return []interface{}{s}
	}

	d.reflect()

	var (
		mapping   = fieldMapping(d.rt)
		result    = make([]interface{}, len(fields))
		tempValue = sql.RawBytes{}
	)

	for index, field := range fields {
		if structIndex, ok := mapping[field]; ok {
			var (
				fv = d.rv.Field(structIndex)
				ft = fv.Type()
			)

			if ft.Kind() == reflect.Ptr {
				result[index] = fv.Addr().Interface()
			} else {
				result[index] = Nullable(fv.Addr().Interface())
			}
		} else {
			result[index] = &tempValue
		}
	}

	return result
}

func (d *document) Values() []interface{} {
	if v, ok := d.v.(values); ok {
		return v.Values()
	}

	d.reflect()

	var (
		fields  = d.Fields()
		mapping = fieldMapping(d.rt)
		values  = make([]interface{}, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			fv          = d.rv.Field(structIndex)
			ft          = fv.Type()
		)

		if ft.Kind() == reflect.Ptr {
			if !fv.IsNil() {
				values[index] = fv.Elem().Interface()
			}
		} else {
			values[index] = fv.Interface()
		}
	}

	return values
}

func (d *document) initAssociations() {
	// if one of assocs fields is not a nil array
	// doesn't neet to check all, because it'll be initialized here.
	if d.belongsTo != nil {
		return
	}

	if s, ok := d.v.(associations); ok {
		d.belongsTo = s.BelongsTo()
		d.hasOne = s.HasOne()
		d.hasMany = s.HasMany()
		return
	}

	d.reflect()

	if result, cached := associationsCache.Load(d.rt); cached {
		fields := result.(associationsData)
		d.belongsTo = fields.belongsTo
		d.hasOne = fields.hasOne
		d.hasMany = fields.hasMany
		return
	}

	for name, index := range d.Fields() {
		var (
			sf  = d.rt.Field(index)
			typ = sf.Type
		)

		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface || typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}

		if typ.Kind() != reflect.Struct {
			continue
		}

		// must have a primary key
		if pk, _ := searchPrimary(typ); pk == "" {
			continue
		}

		switch newAssociation(d.rv, index).Type() {
		case BelongsTo:
			d.belongsTo = append(d.belongsTo, name)
		case HasOne:
			d.hasOne = append(d.hasOne, name)
		case HasMany:
			d.hasMany = append(d.hasMany, name)
		}
	}

	associationsCache.Store(d.rt, associationsData{
		belongsTo: d.belongsTo,
		hasOne:    d.hasOne,
		hasMany:   d.hasMany,
	})
}

func (d *document) BelongsTo() []string {
	d.initAssociations()

	return d.belongsTo
}

func (d *document) HasOne() []string {
	d.initAssociations()

	return d.hasOne
}

func (d *document) HasMany() []string {
	d.initAssociations()

	return d.hasMany
}

func (d *document) Association(name string) Association {
	if s, ok := d.v.(associations); ok {
		return s.Association(name)
	}

	d.reflect()

	index, ok := d.Fields()[name]
	if !ok {
		panic("grimoire: no field named (" + name + ") in type " + d.rt.String() + " found ")
	}

	return newAssociation(d.rv, index)
}

func (d *document) Reset() {
}

func (d *document) Add() Document {
	return d
}

func (d *document) Get(index int) Document {
	return d
}

func (d *document) Len() int {
	return 1
}

func newDocument(entity interface{}) Document {
	switch v := entity.(type) {
	case Document:
		return v
	case reflect.Value:
		if v.Kind() != reflect.Ptr || v.Elem().Kind() == reflect.Slice {
			panic("grimoire: must be a pointer to a struct")
		}

		return &document{
			v:  v.Interface(),
			rv: v.Elem(),
			rt: v.Elem().Type(),
		}
	case reflect.Type:
		panic("grimoire: cannot use reflect.Type")
	case nil:
		panic("grimoire: cannot be nil")
	default:
		return &document{v: v}
	}
}

func fieldName(sf reflect.StructField) string {
	if tag := sf.Tag.Get("db"); tag != "" {
		name := strings.Split(tag, ",")[0]

		if name == "-" {
			return ""
		}

		if name != "" {
			return name
		}
	}

	return snakecase.SnakeCase(sf.Name)
}

func fieldMapping(rt reflect.Type) map[string]int {
	// check for cache
	if v, cached := fieldMappingCache.Load((rt)); cached {
		return v.(map[string]int)
	}

	mapping := make(map[string]int, rt.NumField())

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			name = fieldName(sf)
		)

		if name != "" {
			mapping[name] = i
		}
	}

	fieldMappingCache.Store(rt, mapping)

	return mapping
}

func searchPrimary(rt reflect.Type) (string, int) {
	if result, cached := primariesCache.Load(rt); cached {
		p := result.(primaryData)
		return p.field, p.index
	}

	var (
		field = ""
		index = -1
	)

	for i := 0; i < rt.NumField(); i++ {
		sf := rt.Field(i)

		if tag := sf.Tag.Get("db"); strings.HasSuffix(tag, ",primary") {
			index = i
			if len(tag) > 8 { // has custom field name
				field = tag[:len(tag)-8]
			} else {
				field = snakecase.SnakeCase(sf.Name)
			}

			continue
		}

		// check fallback for id field
		if strings.EqualFold("id", sf.Name) {
			index = i
			field = "id"
		}
	}

	primariesCache.Store(rt, primaryData{
		field: field,
		index: index,
	})

	return field, index
}

func tableName(rt reflect.Type) string {
	// check for cache
	if name, cached := tablesCache.Load(rt); cached {
		return name.(string)
	}

	name := inflection.Plural(rt.Name())
	name = snakecase.SnakeCase(name)

	tablesCache.Store(rt, name)

	return name
}
