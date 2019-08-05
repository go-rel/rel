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
	tableNamesCache   sync.Map
	primaryKeysCache  sync.Map
	fieldsCache       sync.Map
	fieldMappingCache sync.Map
	typesCache        sync.Map
	associationsCache sync.Map
)

type tableName interface {
	TableName() string
}

type primaryKey interface {
	PrimaryKey() (string, interface{})
}

type primaryKeyData struct {
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

type Model interface {
	tableName
	primaryKey
	fields
	types
	scanners
	values
	associations
}

type model struct {
	v         interface{}
	rv        reflect.Value
	rt        reflect.Type
	belongsTo []string
	hasOne    []string
	hasMany   []string
}

func (m *model) reflect() {
	if m.rv.IsValid() {
		return
	}

	m.rv = reflect.ValueOf(m.v)
	if m.rv.Kind() != reflect.Ptr {
		panic("grimoire: must be a pointer")
	}

	m.rv = m.rv.Elem()
	m.rt = m.rv.Type()

	if m.rt.Kind() != reflect.Struct {
		panic("grimoire: must be a pointer to struct")
	}
}

func (m *model) TableName() string {
	if tn, ok := m.v.(tableName); ok {
		return tn.TableName()
	}

	m.reflect()

	// check for cache
	if name, cached := tableNamesCache.Load(m.rt); cached {
		return name.(string)
	}

	name := inflection.Plural(m.rt.Name())
	name = snakecase.SnakeCase(name)

	tableNamesCache.Store(m.rt, name)

	return name
}

func (m *model) PrimaryKey() (string, interface{}) {
	if pk, ok := m.v.(primaryKey); ok {
		key, value := pk.PrimaryKey()
		return key, []interface{}{value}
	}

	m.reflect()

	var (
		field string
		index int
	)

	if result, cached := primaryKeysCache.Load(m.rt); cached {
		pk := result.(primaryKeyData)
		field = pk.field
		index = pk.index
	} else {
		field, index = m.searchPrimaryKey()
		if field == "" {
			panic("grimoire: failed to infer primary key for type " + m.rt.String())
		}

		primaryKeysCache.Store(m.rt, primaryKeyData{
			field: field,
			index: index,
		})
	}

	return field, m.rv.Field(index).Interface()
}

func (m *model) searchPrimaryKey() (string, int) {
	var (
		field = ""
		index = 0
	)

	for i := 0; i < m.rt.NumField(); i++ {
		sf := m.rt.Field(i)

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

	return field, index
}

func (m *model) Fields() map[string]int {
	if s, ok := m.v.(fields); ok {
		return s.Fields()
	}

	m.reflect()

	// check for cache
	if v, cached := fieldsCache.Load((m.rt)); cached {
		return v.(map[string]int)
	}

	var (
		index  = 0
		fields = make(map[string]int, m.rt.NumField())
	)

	for i := 0; i < m.rt.NumField(); i++ {
		var (
			sf   = m.rt.Field(i)
			name = fieldName(sf)
		)

		if name != "" {
			fields[name] = index
			index++
		}
	}

	fieldsCache.Store(m.rt, fields)

	return fields
}

func (m *model) Types() []reflect.Type {
	if v, ok := m.v.(types); ok {
		return v.Types()
	}

	m.reflect()

	// check for cache
	if v, cached := typesCache.Load(m.rt); cached {
		return v.([]reflect.Type)
	}

	var (
		fields  = m.Fields()
		mapping = fieldMapping(m.rt)
		types   = make([]reflect.Type, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			ft          = m.rt.Field(structIndex).Type
		)

		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		} else if ft.Kind() == reflect.Slice && ft.Elem().Kind() == reflect.Ptr {
			ft = reflect.SliceOf(ft.Elem().Elem())
		}

		types[index] = ft
	}

	typesCache.Store(m.rt, types)

	return types
}

func (m *model) Scanners(fields []string) []interface{} {
	if v, ok := m.v.(scanners); ok {
		return v.Scanners(fields)
	}

	if s, ok := m.v.(sql.Scanner); ok {
		return []interface{}{s}
	}

	m.reflect()

	var (
		mapping   = fieldMapping(m.rt)
		result    = make([]interface{}, len(fields))
		tempValue = sql.RawBytes{}
	)

	for index, field := range fields {
		if structIndex, ok := mapping[field]; ok {
			var (
				fv = m.rv.Field(structIndex)
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

func (m *model) Values() []interface{} {
	if v, ok := m.v.(values); ok {
		return v.Values()
	}

	m.reflect()

	var (
		fields  = m.Fields()
		mapping = fieldMapping(m.rt)
		values  = make([]interface{}, len(fields))
	)

	for name, index := range fields {
		var (
			structIndex = mapping[name]
			fv          = m.rv.Field(structIndex)
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

func (m *model) inferAssociations() {
	// if one of assocs fields is not a nil array
	// doesn't neet to check all, because it'll be initialized here.
	if m.belongsTo != nil {
		return
	}

	if s, ok := m.v.(associations); ok {
		m.belongsTo = s.BelongsTo()
		m.hasOne = s.HasOne()
		m.hasMany = s.HasMany()
		return
	}

	m.reflect()

	if result, cached := associationsCache.Load(m.rt); cached {
		fields := result.(associationsData)
		m.belongsTo = fields.belongsTo
		m.hasOne = fields.hasOne
		m.hasMany = fields.hasMany
		return
	}

	for i := 0; i < m.rt.NumField(); i++ {
		var (
			sf   = m.rt.Field(i)
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

		switch InferAssociation(m.rv, name).Type() {
		case BelongsTo:
			m.belongsTo = append(m.belongsTo, name)
		case HasOne:
			m.hasOne = append(m.hasOne, name)
		case HasMany:
			m.hasMany = append(m.hasMany, name)
		}
	}

	associationsCache.Store(m.rt, associationsData{
		belongsTo: m.belongsTo,
		hasOne:    m.hasOne,
		hasMany:   m.hasMany,
	})
}

func (m *model) BelongsTo() []string {
	m.inferAssociations()

	return m.belongsTo
}

func (m *model) HasOne() []string {
	m.inferAssociations()

	return m.hasOne
}

func (m *model) HasMany() []string {
	m.inferAssociations()

	return m.hasMany
}

func (m *model) Association(name string) Association {
	if s, ok := m.v.(associations); ok {
		return s.Association(name)
	}

	m.reflect()

	return InferAssociation(m.rv, name)
}

func InferModel(record interface{}) Model {
	return &model{v: record}
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

func searchPrimaryKey(rt reflect.Type) (string, int) {
	var (
		field = ""
		index = 0
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

	return field, index
}
