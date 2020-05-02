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

// DocumentFlag stores information about document as a flag.
type DocumentFlag int8

// Is returns true if it's defined.
func (df DocumentFlag) Is(flag DocumentFlag) bool {
	return (df & flag) == flag
}

const (
	// Invalid flag.
	Invalid DocumentFlag = 1 << iota
	// HasCreatedAt flag.
	HasCreatedAt
	// HasUpdatedAt flag.
	HasUpdatedAt
	// HasDeletedAt flag.
	HasDeletedAt
)

var (
	tablesCache       sync.Map
	primariesCache    sync.Map
	fieldsCache       sync.Map
	typesCache        sync.Map
	documentDataCache sync.Map
	rtTime            = reflect.TypeOf(time.Time{})
	rtTable           = reflect.TypeOf((*table)(nil)).Elem()
	rtPrimary         = reflect.TypeOf((*primary)(nil)).Elem()
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
	index        map[string]int
	fields       []string
	belongsTo    []string
	hasOne       []string
	hasMany      []string
	primaryField string
	primaryIndex int
	flag         DocumentFlag
}

// Document provides an abstraction over reflect to easily works with struct for database purpose.
type Document struct {
	v    interface{}
	rv   reflect.Value
	rt   reflect.Type
	data documentData
}

// ReflectValue of referenced document.
func (d Document) ReflectValue() reflect.Value {
	return d.rv
}

// Table returns name of the table.
func (d Document) Table() string {
	if tn, ok := d.v.(table); ok {
		return tn.Table()
	}

	// TODO: handle anonymous struct
	return tableName(d.rt)
}

// PrimaryField column name of this document.
func (d Document) PrimaryField() string {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryField()
	}

	if d.data.primaryField == "" {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	return d.data.primaryField
}

// PrimaryValue of this document.
func (d Document) PrimaryValue() interface{} {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryValue()
	}

	if d.data.primaryIndex < 0 {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	return d.rv.Field(d.data.primaryIndex).Interface()
}

// Index returns map of column name and it's struct index.
func (d Document) Index() map[string]int {
	return d.data.index
}

// Fields returns list of fields available on this document.
func (d Document) Fields() []string {
	return d.data.fields
}

// Type returns reflect.Type of given field. if field does not exist, second returns value will be false.
func (d Document) Type(field string) (reflect.Type, bool) {
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

// Value returns value of given field. if field does not exist, second returns value will be false.
func (d Document) Value(field string) (interface{}, bool) {
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

// SetValue of the field, it returns false if field does not exist, or it's not assignable.
func (d Document) SetValue(field string, value interface{}) bool {
	if i, ok := d.data.index[field]; ok {
		var (
			rv reflect.Value
			rt reflect.Type
			fv = d.rv.Field(i)
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

func setPointerValue(ft reflect.Type, fv reflect.Value, rt reflect.Type, rv reflect.Value) bool {
	if ft.Elem() != rt && !rt.AssignableTo(ft.Elem()) {
		return false
	}

	if fv.IsNil() {
		fv.Set(reflect.New(ft.Elem()))
	}
	fv.Elem().Set(rv)

	return true
}

func setConvertValue(ft reflect.Type, fv reflect.Value, rt reflect.Type, rv reflect.Value) bool {
	var (
		rk = rt.Kind()
		fk = ft.Kind()
	)

	// prevents unintentional convertion
	if (rk >= reflect.Int || rk <= reflect.Uint64) && fk == reflect.String {
		return false
	}

	fv.Set(rv.Convert(ft))
	return true
}

// Scanners returns slice of sql.Scanner for given fields.
func (d Document) Scanners(fields []string) []interface{} {
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

// BelongsTo fields of this document.
func (d Document) BelongsTo() []string {
	return d.data.belongsTo
}

// HasOne fields of this document.
func (d Document) HasOne() []string {
	return d.data.hasOne
}

// HasMany fields of this document.
func (d Document) HasMany() []string {
	return d.data.hasMany
}

// Association of this document with given name.
func (d Document) Association(name string) Association {
	index, ok := d.data.index[name]
	if !ok {
		panic("rel: no field named (" + name + ") in type " + d.rt.String() + " found ")
	}

	return newAssociation(d.rv, index)
}

// Reset this document, this is a noop for compatibility with collection.
func (d Document) Reset() {
}

// Add returns this document, this is a noop for compatibility with collection.
func (d *Document) Add() *Document {
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

// Flag returns true if struct contains specified flag.
func (d Document) Flag(flag DocumentFlag) bool {
	return d.data.flag.Is(flag)
}

// NewDocument used to create abstraction to work with struct.
// Document can be created using interface or reflect.Value.
func NewDocument(record interface{}, readonly ...bool) *Document {
	switch v := record.(type) {
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

func newDocument(v interface{}, rv reflect.Value, readonly bool) *Document {
	var (
		rt = rv.Type()
	)

	if rt.Kind() != reflect.Ptr {
		if !readonly {
			panic("rel: must be a pointer to struct")
		}
	} else {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("rel: must be a struct or pointer to a struct")
	}

	return &Document{
		v:    v,
		rv:   rv,
		rt:   rt,
		data: extractDocumentData(rt, false),
	}
}

func extractDocumentData(rt reflect.Type, skipAssoc bool) documentData {
	if data, cached := documentDataCache.Load(rt); cached {
		return data.(documentData)
	}

	var (
		data = documentData{
			index: make(map[string]int, rt.NumField()),
		}
	)

	// TODO probably better to use slice index instead.
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

		if typ.Kind() != reflect.Struct {
			data.fields = append(data.fields, name)
			continue
		}

		if flag := extractFlag(typ, name); flag != Invalid {
			data.fields = append(data.fields, name)
			data.flag |= flag
			continue
		}

		// struct without primary key is a field
		// TODO: test by scanner/valuer instead?
		if pk, _ := searchPrimary(typ); pk == "" {
			data.fields = append(data.fields, name)
			continue
		}

		if !skipAssoc {
			switch extractAssociationData(rt, i).typ {
			case BelongsTo:
				data.belongsTo = append(data.belongsTo, name)
			case HasOne:
				data.hasOne = append(data.hasOne, name)
			case HasMany:
				data.hasMany = append(data.hasMany, name)
			}
		}
	}

	data.primaryField, data.primaryIndex = searchPrimary(rt)

	if !skipAssoc {
		documentDataCache.Store(rt, data)
	}

	return data
}

func extractFlag(rt reflect.Type, name string) DocumentFlag {
	flag := Invalid
	if rt != rtTime {
		return flag
	}

	switch name {
	case "created_at", "inserted_at":
		flag = HasCreatedAt
	case "updated_at":
		flag = HasUpdatedAt
	case "deleted_at":
		flag = HasDeletedAt
	}

	return flag
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

	if rt.Implements(rtPrimary) {
		var (
			v = reflect.Zero(rt).Interface().(primary)
		)

		field = v.PrimaryField()
		index = -2 // special index to mark interface usage
	} else {
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)

			if tag := sf.Tag.Get("db"); strings.HasSuffix(tag, ",primary") {
				index = i
				field = fieldName(sf)
				continue
			}

			// check fallback for id field
			if strings.EqualFold("id", sf.Name) {
				index = i
				field = "id"
			}
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
