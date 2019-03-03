package params

import (
	"reflect"
	"sync"
	"time"

	"github.com/tidwall/gjson"
)

// JSON is param type for json document.
type JSON struct {
	gjson.Result
	results sync.Map
}

var _ Params = (*JSON)(nil)

// Exists returns true if key exists.
func (json *JSON) Exists(name string) bool {
	return json.fetch(name).Exists()
}

// Get returns value as interface.
// returns nil if value doens't exists.
func (json *JSON) Get(name string) interface{} {
	return json.Result.Get(name).Value()
}

// GetWithType returns value given from given name and type.
// second return value will only be false if the type of parameter is not convertible to requested type.
// If value is not convertible to type, it'll return nil, false
// If value is not exists, it will return nil, true
func (json *JSON) GetWithType(name string, typ reflect.Type) (interface{}, bool) {
	value := json.fetch(name)
	if value.IsArray() && typ.Kind() == reflect.Slice {
		array := value.Array()
		result := reflect.MakeSlice(typ, len(array), len(array))
		elmType := typ.Elem()

		for i, elm := range array {
			elmValue, valid := json.convert(elm, elmType)
			if !valid {
				return nil, false
			}

			result.Index(i).Set(reflect.ValueOf(elmValue))
		}
		return result.Interface(), true
	}

	return json.convert(value, typ)
}

// GetParams returns nested param
func (json *JSON) GetParams(name string) (Params, bool) {
	if value := json.fetch(name); value.IsObject() {
		return &JSON{Result: value}, true
	}

	return nil, false
}

// GetParamsSlice returns slice of nested param
func (json *JSON) GetParamsSlice(name string) ([]Params, bool) {
	if value := json.fetch(name); value.IsArray() {
		pars := value.Array()
		mpar := make([]Params, len(pars))

		for i, par := range pars {
			if !par.IsObject() {
				return nil, false
			}

			mpar[i] = &JSON{Result: par}
		}
		return mpar, true
	}

	return nil, false
}

func (json *JSON) convert(value gjson.Result, typ reflect.Type) (interface{}, bool) {
	if value.Type == gjson.Null {
		return nil, true
	}

	// handle type alias
	if typ.PkgPath() != "" && typ.Kind() != reflect.Struct && typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		rv := reflect.ValueOf(value.Value())
		if !rv.Type().ConvertibleTo(typ) {
			return nil, false
		}

		return rv.Convert(typ).Interface(), true
	}

	switch value.Type {
	case gjson.False, gjson.True:
		if typ.Kind() == reflect.Bool {
			return value.Bool(), true
		}
	case gjson.Number:
		switch typ.Kind() {
		case reflect.Int:
			return int(value.Int()), true
		case reflect.Int8:
			return int8(value.Int()), true
		case reflect.Int16:
			return int16(value.Int()), true
		case reflect.Int32:
			return int32(value.Int()), true
		case reflect.Int64:
			return value.Int(), true
		case reflect.Uint:
			return uint(value.Uint()), true
		case reflect.Uint8:
			return uint8(value.Uint()), true
		case reflect.Uint16:
			return uint16(value.Uint()), true
		case reflect.Uint32:
			return uint32(value.Uint()), true
		case reflect.Uint64:
			return value.Uint(), true
		case reflect.Uintptr:
			return uintptr(value.Uint()), true
		case reflect.Float32:
			return float32(value.Float()), true
		case reflect.Float64:
			return value.Float(), true
		}
	case gjson.String:
		if typ.Kind() == reflect.String {
			return value.String(), true
		} else if typ == timeType {
			if res, err := time.Parse(time.RFC3339, value.String()); err == nil {
				return res, true
			}
		}
	}

	return nil, false
}

func (json *JSON) fetch(name string) gjson.Result {
	if result, ok := json.results.Load(name); ok {
		return result.(gjson.Result)
	}

	result := json.Result.Get(name)
	json.results.Store(name, result)
	return result
}

// ParseJSON as params
func ParseJSON(json string) Params {
	return &JSON{Result: gjson.Parse(json)}
}
