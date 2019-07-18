package schema

import (
	"reflect"
	"strings"
	"sync"

	"github.com/azer/snakecase"
)

type primaryKey interface {
	PrimaryKey() (string, interface{})
}

type primaryKeyData struct {
	field string
	index int
}

var primaryKeysCache sync.Map

// InferPrimaryKey from struct.
func InferPrimaryKey(record interface{}, returnValue bool) (string, []interface{}) {
	if pk, ok := record.(primaryKey); ok {
		key, value := pk.PrimaryKey()
		return key, []interface{}{value}
	}

	rt := reflect.TypeOf(record)
	for rt.Kind() == reflect.Ptr || rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}

	if rt.Kind() != reflect.Struct {
		panic("grimoire: must be a struct or slice of structs")
	}

	var (
		pkey   = inferPrimaryKeyData(rt)
		field  = pkey.field
		values = []interface{}(nil)
	)

	if returnValue {
		rv := reflect.ValueOf(record)
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}

		if rv.Kind() == reflect.Slice {
			values = make([]interface{}, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				values[i] = rv.Index(i).Field(pkey.index).Interface()
			}
		} else {
			values = []interface{}{rv.Field(pkey.index).Interface()}
		}
	}

	return field, values
}

func inferPrimaryKeyData(rt reflect.Type) primaryKeyData {
	if result, cached := primaryKeysCache.Load(rt); cached {
		return result.(primaryKeyData)
	}

	field, index := searchPrimaryKey(rt)
	if field == "" {
		panic("grimoire: failed to infer primary key for type " + rt.String())
	}

	result := primaryKeyData{
		field: field,
		index: index,
	}

	primaryKeysCache.Store(rt, result)

	return result
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
