package rel

import (
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/inflection"
	"github.com/serenize/snaker"
)

var (
	primariesCache    sync.Map
	documentMetaCache sync.Map
	rtTime            = reflect.TypeOf(time.Time{})
	rtBool            = reflect.TypeOf(false)
	rtInt             = reflect.TypeOf(int(0))
	rtTable           = reflect.TypeOf((*table)(nil)).Elem()
	rtPrimary         = reflect.TypeOf((*primary)(nil)).Elem()
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
	// Versioning
	HasVersioning
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

type cachedDocumentMeta struct {
	table        string
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

// Adds a prefix to field names
func appendWithPrefix(target, fieldNames []string, prefix string) []string {
	if prefix == "" {
		return append(target, fieldNames...)
	}
	for _, name := range fieldNames {
		target = append(target, prefix+name)
	}
	return target
}

// Adds a field index and checks for conflicts
func (cdm *cachedDocumentMeta) addFieldIndex(name string, index []int) {
	if _, ok := cdm.index[name]; ok {
		panic("rel: conflicting field (" + name + ") in struct")
	}
	cdm.index[name] = index
}

// Transfer values from other document data
func (cdm *cachedDocumentMeta) mergeEmbedded(other cachedDocumentMeta, indexPrefix int, namePrefix string) {
	for name, path := range other.index {
		cdm.addFieldIndex(namePrefix+name, append([]int{indexPrefix}, path...))
	}
	cdm.fields = appendWithPrefix(cdm.fields, other.fields, namePrefix)
	cdm.belongsTo = appendWithPrefix(cdm.belongsTo, other.belongsTo, namePrefix)
	cdm.hasOne = appendWithPrefix(cdm.hasOne, other.hasOne, namePrefix)
	cdm.hasMany = appendWithPrefix(cdm.hasMany, other.hasMany, namePrefix)
	cdm.primaryField = appendWithPrefix(cdm.primaryField, other.primaryField, namePrefix)
	for index := range other.primaryIndex {
		cdm.primaryIndex = append(cdm.primaryIndex, append([]int{indexPrefix}, index))
	}
	cdm.preload = appendWithPrefix(cdm.preload, other.preload, namePrefix)
	cdm.flag |= other.flag
}

type DocumentMeta struct {
	rt reflect.Type
	cachedDocumentMeta
}

// Table returns name of the table.
func (dm DocumentMeta) Table() string {
	return dm.table
}

// PrimaryFields column name of this document.
func (dm DocumentMeta) PrimaryFields() []string {
	if len(dm.primaryField) == 0 {
		panic("rel: failed to infer primary key for type " + dm.rt.String())
	}

	return dm.primaryField
}

// PrimaryField column name of this document.
// panic if document uses composite key.
func (dm DocumentMeta) PrimaryField() string {
	if fields := dm.PrimaryFields(); len(fields) == 1 {
		return fields[0]
	}

	panic("rel: composite primary key is not supported")
}

// Index returns map of column name and it's struct index.
func (dm DocumentMeta) Index() map[string][]int {
	return dm.index
}

// Fields returns list of fields available on this document.
func (dm DocumentMeta) Fields() []string {
	return dm.fields
}

// Type returns reflect.Type of given field. if field does not exist, second returns value will be false.
func (dm DocumentMeta) Type(field string) (reflect.Type, bool) {
	if i, ok := dm.index[field]; ok {
		var (
			ft = dm.rt.FieldByIndex(i).Type
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

// BelongsTo fields of this document.
func (dm DocumentMeta) BelongsTo() []string {
	return dm.belongsTo
}

// HasOne fields of this document.
func (dm DocumentMeta) HasOne() []string {
	return dm.hasOne
}

// HasMany fields of this document.
func (dm DocumentMeta) HasMany() []string {
	return dm.hasMany
}

// Preload fields of this document.
func (dm DocumentMeta) Preload() []string {
	return dm.preload
}

// Association of this document with given name.
func (dm DocumentMeta) Association(name string) AssociationMeta {
	if assoc, ok := dm.association(name); ok {
		return assoc
	}

	panic("rel: no field named (" + name + ") in type " + dm.rt.String() + " found ")
}

func (dm DocumentMeta) association(name string) (AssociationMeta, bool) {
	index, ok := dm.index[name]
	if !ok {
		return AssociationMeta{}, false
	}

	return getAssociationMeta(dm.rt, index), true
}

// Flag returns true if struct contains specified flag.
func (dm DocumentMeta) Flag(flag DocumentFlag) bool {
	return dm.flag.Is(flag)
}

func getDocumentMeta(rt reflect.Type, skipAssoc bool) DocumentMeta {
	if meta, cached := documentMetaCache.Load(rt); cached {
		return DocumentMeta{
			cachedDocumentMeta: meta.(cachedDocumentMeta),
			rt:                 rt,
		}
	}

	var (
		meta = cachedDocumentMeta{
			table: tableName(rt),
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
			embedded := getDocumentMeta(typ, skipAssoc)
			embeddedName := ""
			if tagged {
				embeddedName = name
			}
			meta.mergeEmbedded(embedded.cachedDocumentMeta, i, embeddedName)
			continue
		}

		meta.addFieldIndex(name, sf.Index)

		if flag := extractFlag(typ, name); flag != Invalid {
			meta.fields = append(meta.fields, name)
			meta.flag |= flag
			continue
		}

		if typ.Kind() != reflect.Struct {
			meta.fields = append(meta.fields, name)
			continue
		}

		// struct without primary key is a field
		// TODO: test by scanner/valuer instead?
		if pk, _ := searchPrimary(typ); len(pk) == 0 {
			meta.fields = append(meta.fields, name)
			continue
		}

		if !skipAssoc {
			var (
				assocMeta = getAssociationMeta(rt, sf.Index)
			)

			switch assocMeta.typ {
			case BelongsTo:
				meta.belongsTo = append(meta.belongsTo, name)
			case HasOne:
				meta.hasOne = append(meta.hasOne, name)
			case HasMany:
				meta.hasMany = append(meta.hasMany, name)
			}

			if assocMeta.autoload {
				meta.preload = append(meta.preload, name)
			}
		}
	}

	primaryField, primaryIndex := searchPrimary(rt)
	meta.primaryField = append(meta.primaryField, primaryField...)
	meta.primaryIndex = append(meta.primaryIndex, primaryIndex...)

	if !skipAssoc {
		documentMetaCache.Store(rt, meta)
	}

	return DocumentMeta{
		rt:                 rt,
		cachedDocumentMeta: meta,
	}
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

func extractIntFlag(name string) DocumentFlag {
	if name == "lock_version" {
		return HasVersioning
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
	if rt == rtInt {
		return extractIntFlag(name)
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
	var name string
	if rt.Implements(rtTable) {
		name = reflect.Zero(rt).Interface().(table).Table()
	} else {
		name = inflection.Plural(rt.Name())
		name = snaker.CamelToSnake(name)
	}

	return name
}
