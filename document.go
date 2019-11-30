package rel

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/azer/snakecase"
	"github.com/jinzhu/inflection"
)

var (
	tablesCache       sync.Map
	primariesCache    sync.Map
	fieldsCache       sync.Map
	typesCache        sync.Map
	documentDataCache sync.Map
	rtTime            = reflect.TypeOf(time.Time{})
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

type documentData struct {
	index     map[string]int
	fields    []string
	belongsTo []string
	hasOne    []string
	hasMany   []string
}

type Document struct {
	v    interface{}
	rv   reflect.Value
	rt   reflect.Type
	data documentData
}

func (d *Document) Table() string {
	if tn, ok := d.v.(table); ok {
		return tn.Table()
	}

	// TODO: handle anonymous struct
	return tableName(d.rt)
}

func (d *Document) PrimaryField() string {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryField()
	}

	var (
		field, _ = searchPrimary(d.rt)
	)

	if field == "" {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	return field
}

func (d *Document) PrimaryValue() interface{} {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryValue()
	}

	var (
		_, index = searchPrimary(d.rt)
	)

	if index < 0 {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	return d.rv.Field(index).Interface()
}

func (d *Document) Index() map[string]int {
	return d.data.index
}

func (d *Document) Fields() []string {
	return d.data.fields
}

func (d *Document) Type(field string) (reflect.Type, bool) {
	if i, ok := d.data.index[field]; ok {
		var (
			ft = d.rt.Field(i).Type
		)

		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		} else if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
			ft = reflect.SliceOf(ft.Elem().Elem())
		}

		return ft, true
	}

	return nil, false
}

func (d *Document) Value(field string) (interface{}, bool) {
	if i, ok := d.data.index[field]; ok {
		var (
			value interface{}
			fv    = d.rv.Field(i)
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

func (d *Document) Scanners(fields []string) []interface{} {
	var (
		result = make([]interface{}, len(fields))
	)

	for index, field := range fields {
		if structIndex, ok := d.data.index[field]; ok {
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
			result[index] = &sql.RawBytes{}
		}
	}

	return result
}

func (d *Document) BelongsTo() []string {
	return d.data.belongsTo
}

func (d *Document) HasOne() []string {
	return d.data.hasOne
}

func (d *Document) HasMany() []string {
	return d.data.hasMany
}

func (d *Document) Association(name string) Association {
	index, ok := d.data.index[name]
	if !ok {
		panic("rel: no field named (" + name + ") in type " + d.rt.String() + " found ")
	}

	return newAssociation(d.rv, index)
}

func (d *Document) Reset() {
}

func (d *Document) Add() *Document {
	return d
}

func (d *Document) Get(index int) *Document {
	return d
}

func (d *Document) Len() int {
	return 1
}

func NewDocument(record interface{}) *Document {
	switch v := record.(type) {
	case *Document:
		return v
	case reflect.Value:
		if v.Kind() != reflect.Ptr || v.Elem().Kind() == reflect.Slice {
			panic("rel: must be a pointer to a struct")
		}

		var (
			rv = v.Elem()
			rt = rv.Type()
		)

		return &Document{
			v:    v.Interface(),
			rv:   rv,
			rt:   rt,
			data: extractdocumentData(rv, rt),
		}
	case reflect.Type:
		panic("rel: cannot use reflect.Type")
	case nil:
		panic("rel: cannot be nil")
	default:
		var (
			rv = reflect.ValueOf(v)
			rt = rv.Type()
		)

		if rt.Kind() != reflect.Ptr && rt.Elem().Kind() != reflect.Struct {
			panic("rel: must be a pointer to struct")
		}

		rv = rv.Elem()
		rt = rt.Elem()

		return &Document{
			v:    v,
			rv:   rv,
			rt:   rt,
			data: extractdocumentData(rv, rt),
		}
	}
}

func extractdocumentData(rv reflect.Value, rt reflect.Type) documentData {
	if data, cached := documentDataCache.Load(rv.Type()); cached {
		return data.(documentData)
	}

	var (
		data = documentData{
			index: make(map[string]int, rt.NumField()),
		}
	)

	for i := 0; i < rt.NumField(); i++ {
		var (
			sf   = rt.Field(i)
			typ  = sf.Type
			name = fieldName(sf)
		)

		if name == "" {
			continue
		}

		data.index[name] = i

		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface || typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}

		if typ.Kind() != reflect.Struct || typ == rtTime {
			data.fields = append(data.fields, name)
			continue
		}

		// struct without primary key is a field
		// TODO: test by scanner/valuer instead?
		if pk, _ := searchPrimary(typ); pk == "" {
			data.fields = append(data.fields, name)
			continue
		}

		switch newAssociation(rv, i).Type() {
		case BelongsTo:
			data.belongsTo = append(data.belongsTo, name)
		case HasOne:
			data.hasOne = append(data.hasOne, name)
		case HasMany:
			data.hasMany = append(data.hasMany, name)
		}
	}

	documentDataCache.Store(rt, data)

	return data
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
