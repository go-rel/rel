package params

import (
	"reflect"
)

// Map is param type alias for map[string]interface{}
type Map map[string]interface{}

var _ Params = (*Map)(nil)

// Exists returns true if key exists.
func (m Map) Exists(name string) bool {
	_, exists := m[name]
	return exists
}

// Get returns value as interface.
// returns nil if value doens't exists.
func (m Map) Get(name string) interface{} {
	return m[name]
}

// GetWithType returns value given from given name and type.
// second return value will only be false if the type of parameter is not convertible to requested type.
// If value is not convertible to type, it'll return nil, false
// If value is not exists, it will return nil, true
func (m Map) GetWithType(name string, typ reflect.Type) (interface{}, bool) {
	value := m[name]

	if value == nil {
		return nil, true
	}

	rv := reflect.ValueOf(value)
	rt := rv.Type()

	if rt.Kind() == reflect.Ptr {
		rv = rv.Elem()
		rt = rt.Elem()
	}

	if !rv.IsValid() {
		return nil, true
	}

	if typ.Kind() == reflect.Slice && (rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array) {
		result := reflect.MakeSlice(typ, rv.Len(), rv.Len())
		elemTyp := typ.Elem()

		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i)
			if elem.Kind() == reflect.Interface {
				elem = elem.Elem()
			}

			if elem.Type().ConvertibleTo(elemTyp) {
				result.Index(i).Set(elem.Convert(elemTyp))
			} else {
				return nil, false
			}
		}

		return result.Interface(), true
	}

	if !rt.ConvertibleTo(typ) {
		return nil, false
	}

	return rv.Convert(typ).Interface(), true
}

// GetParams returns nested param
func (m Map) GetParams(name string) (Params, bool) {
	if val, exist := m[name]; exist {
		if par, ok := val.(Params); ok {
			return par, ok
		}

		if par, ok := val.(map[string]interface{}); ok {
			return Map(par), ok
		}
	}

	return nil, false
}

// GetParamsSlice returns slice of nested param
func (m Map) GetParamsSlice(name string) ([]Params, bool) {
	if val, exist := m[name]; exist {
		if pars, ok := val.([]Params); ok {
			return pars, ok
		}

		if pars, ok := val.([]Map); ok {
			mpar := make([]Params, len(pars))
			for i, par := range pars {
				mpar[i] = Map(par)
			}
			return mpar, true
		}

		if pars, ok := val.([]map[string]interface{}); ok {
			mpar := make([]Params, len(pars))
			for i, par := range pars {
				mpar[i] = Map(par)
			}
			return mpar, true
		}
	}

	return nil, false
}
