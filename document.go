package rel

import (
	"database/sql"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
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
	// HasDeleted flag.
	HasDeleted
)

var (
	tablesCache       sync.Map
	primariesCache    sync.Map
	documentDataCache sync.Map
	rtTime            = reflect.TypeOf(time.Time{})
	rtBool            = reflect.TypeOf(false)
	rtTable           = reflect.TypeOf((*table)(nil)).Elem()
	rtPrimary         = reflect.TypeOf((*primary)(nil)).Elem()
)

type table interface {
	Table() string
}

type primary interface {
	PrimaryFields() []string
	PrimaryValues() []interface{}
}

type primaryData struct {
	field []string
	index [][]int
}

type documentData struct {
	index        map[string][]int
	fields       []string
	belongsTo    []string
	hasOne       []string
	hasMany      []string
	primaryField []string
	primaryIndex [][]int
	preload      []string
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
	// TODO: handle anonymous struct
	return tableName(d.rt)
}

// PrimaryFields column name of this document.
func (d Document) PrimaryFields() []string {
	if len(d.data.primaryField) == 0 {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	return d.data.primaryField
}

// PrimaryField column name of this document.
// panic if document uses composite key.
func (d Document) PrimaryField() string {
	if fields := d.PrimaryFields(); len(fields) == 1 {
		return fields[0]
	}

	panic("rel: composite primary key is not supported")
}

// PrimaryValues of this document.
func (d Document) PrimaryValues() []interface{} {
	if p, ok := d.v.(primary); ok {
		return p.PrimaryValues()
	}

	if len(d.data.primaryIndex) == 0 {
		panic("rel: failed to infer primary key for type " + d.rt.String())
	}

	var (
		pValues = make([]interface{}, len(d.data.primaryIndex))
	)

	for i := range pValues {
		pValues[i] = d.rv.FieldByIndex(d.data.primaryIndex[i]).Interface()
	}

	return pValues
}

// PrimaryValue of this document.
// panic if document uses composite key.
func (d Document) PrimaryValue() interface{} {
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
			ft = d.rt.FieldByIndex(i).Type
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
			fv    = d.rv.FieldByIndex(i)
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
			fv = d.rv.FieldByIndex(i)
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

// Init pointers to embedded structs for index path
func initPointersForIndices(rv reflect.Value, indices []int) {
	for depth := 0; depth < len(indices)-1; depth += 1 {
		field := rv.Field(indices[depth])
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				field.Set(reflect.New(field.Type().Elem()))
			}
			rv = field.Elem()
		} else {
			rv = field
		}
	}
}

// Scanners returns slice of sql.Scanner for given fields.
func (d Document) Scanners(fields []string) []interface{} {
	var (
		result = make([]interface{}, len(fields))
	)

	for index, field := range fields {
		if structIndex, ok := d.data.index[field]; ok {
			if len(structIndex) > 1 {
				initPointersForIndices(d.rv, structIndex)
			}

			var (
				fv = d.rv.FieldByIndex(structIndex)
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

// Preload fields of this document.
func (d Document) Preload() []string {
	return d.data.preload
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

// Flag returns true if struct contains specified flag.
func (d Document) Flag(flag DocumentFlag) bool {
	return d.data.flag.Is(flag)
}

// Adds a prefix to field names
func prefixFieldNames(fieldNames []string, prefix string) []string {
	if prefix == "" {
		return fieldNames
	}
	// this copy is necessary as embedded structs can be reused
	newNames := make([]string, len(fieldNames))
	for i, name := range fieldNames {
		newNames[i] = prefix + name
	}
	return newNames
}

// Adds a field index and checks for conflicts
func (d *documentData) addFieldIndex(name string, index []int) {
	if _, ok := d.index[name]; ok {
		panic("rel: conflicting field (" + name + ") in struct")
	}
	d.index[name] = index
}

// Transfer values from other document data
func (d *documentData) mergeEmbedded(other documentData, indexPrefix int, namePrefix string) {
	for name, path := range other.index {
		d.addFieldIndex(namePrefix+name, append([]int{indexPrefix}, path...))
	}
	appendWithPrefix := func(slice, newNames []string) []string {
		return append(slice, prefixFieldNames(newNames, namePrefix)...)
	}
	d.fields = appendWithPrefix(d.fields, other.fields)
	d.belongsTo = appendWithPrefix(d.belongsTo, other.belongsTo)
	d.hasOne = appendWithPrefix(d.hasOne, other.hasOne)
	d.hasMany = appendWithPrefix(d.hasMany, other.hasMany)
	d.primaryField = appendWithPrefix(d.primaryField, other.primaryField)
	for index := range other.primaryIndex {
		d.primaryIndex = append(d.primaryIndex, append([]int{indexPrefix}, index))
	}
	d.preload = appendWithPrefix(d.preload, other.preload)
	d.flag |= other.flag
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
		data: extractDocumentData(rt, false),
	}
}

func extractDocumentData(rt reflect.Type, skipAssoc bool) documentData {
	if data, cached := documentDataCache.Load(rt); cached {
		return data.(documentData)
	}

	var (
		data = documentData{
			index: make(map[string][]int, rt.NumField()),
		}
	)

	// TODO probably better to use slice index instead.
	for i := 0; i < rt.NumField(); i++ {
		var (
			sf           = rt.Field(i)
			typ          = sf.Type
			name, tagged = fieldName(sf)
		)

		if c := sf.Name[0]; c < 'A' || c > 'Z' || name == "" {
			continue
		}

		for typ.Kind() == reflect.Ptr || typ.Kind() == reflect.Interface || typ.Kind() == reflect.Slice {
			typ = typ.Elem()
		}

		if typ.Kind() == reflect.Struct && sf.Anonymous {
			embedded := extractDocumentData(typ, skipAssoc)
			embeddedName := ""
			if tagged {
				embeddedName = name
			}
			data.mergeEmbedded(embedded, i, embeddedName)
			continue
		}

		data.addFieldIndex(name, sf.Index)

		if flag := extractFlag(typ, name); flag != Invalid {
			data.fields = append(data.fields, name)
			data.flag |= flag
			continue
		}

		if typ.Kind() != reflect.Struct {
			data.fields = append(data.fields, name)
			continue
		}

		// struct without primary key is a field
		// TODO: test by scanner/valuer instead?
		if pk, _ := searchPrimary(typ); len(pk) == 0 {
			data.fields = append(data.fields, name)
			continue
		}

		if !skipAssoc {
			var (
				assocData = extractAssociationData(rt, sf.Index)
			)

			switch assocData.typ {
			case BelongsTo:
				data.belongsTo = append(data.belongsTo, name)
			case HasOne:
				data.hasOne = append(data.hasOne, name)
			case HasMany:
				data.hasMany = append(data.hasMany, name)
			}

			if assocData.autoload {
				data.preload = append(data.preload, name)
			}
		}
	}

	primaryField, primaryIndex := searchPrimary(rt)
	data.primaryField = append(data.primaryField, primaryField...)
	data.primaryIndex = append(data.primaryIndex, primaryIndex...)

	if !skipAssoc {
		documentDataCache.Store(rt, data)
	}

	return data
}

func extractTimeFlag(name string) DocumentFlag {
	switch name {
	case "created_at", "inserted_at":
		return HasCreatedAt
	case "updated_at":
		return HasUpdatedAt
	case "deleted_at":
		return HasDeletedAt
	}
	return Invalid
}

func extractBoolFlag(name string) DocumentFlag {
	if name == "deleted" {
		return HasDeleted
	}
	return Invalid
}

func extractFlag(rt reflect.Type, name string) DocumentFlag {
	if rt == rtTime {
		return extractTimeFlag(name)
	}
	if rt == rtBool {
		return extractBoolFlag(name)
	}
	return Invalid
}

func fieldName(sf reflect.StructField) (string, bool) {
	if tag := sf.Tag.Get("db"); tag != "" {
		name := strings.Split(tag, ",")[0]

		if name == "-" {
			return "", true
		}

		if name != "" {
			return name, true
		}
	}

	return snaker.CamelToSnake(sf.Name), false
}

func searchPrimary(rt reflect.Type) ([]string, [][]int) {
	if result, cached := primariesCache.Load(rt); cached {
		p := result.(primaryData)
		return p.field, p.index
	}

	var (
		field         []string
		index         [][]int
		fallbackIndex = -1
	)

	if rt.Implements(rtPrimary) {
		var (
			v = reflect.Zero(rt).Interface().(primary)
		)

		field = v.PrimaryFields()
		// index kept nil to mark interface usage
	} else {
		for i := 0; i < rt.NumField(); i++ {
			sf := rt.Field(i)

			if tag := sf.Tag.Get("db"); strings.HasSuffix(tag, ",primary") {
				index = append(index, sf.Index)
				name, _ := fieldName(sf)
				field = append(field, name)
				continue
			}

			// check fallback for id field
			if strings.EqualFold("id", sf.Name) {
				fallbackIndex = i
			}
		}
	}

	if len(field) == 0 && fallbackIndex >= 0 {
		field = []string{"id"}
		index = [][]int{{fallbackIndex}}
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

	var name string
	if rt.Implements(rtTable) {
		name = reflect.Zero(rt).Interface().(table).Table()
	} else {
		name = inflection.Plural(rt.Name())
		name = snaker.CamelToSnake(name)
	}

	tablesCache.Store(rt, name)

	return name
}
