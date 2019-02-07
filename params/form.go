package params

import (
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Form is param type alias for url.Values
type Form map[string][]interface{}

var _ Params = (*Form)(nil)

// Exists returns true if key exists.
func (form Form) Exists(name string) bool {
	_, exists := form[name]
	return exists
}

// Get returns value as interface.
// returns nil if value doens't exists.
// the value returned is slice of interface{}
func (form Form) Get(name string) interface{} {
	if val, exists := form[name]; exists {
		return val
	}

	return nil
}

// GetWithType returns value given from given name and type.
// second return value will only be false if the type of parameter is not convertible to requested type.
// If value is not convertible to type, it'll return nil, false
// If value is not exists, it will return nil, true
func (form Form) GetWithType(name string, typ reflect.Type) (interface{}, bool) {
	value, exist := form[name]
	if !exist {
		return nil, true
	}

	result, valid := interface{}(nil), false

	if typ.Kind() == reflect.Slice {
		rv := reflect.MakeSlice(typ, len(value), len(value))
		elmType := typ.Elem()

		for i, elm := range value {
			if str, ok := elm.(string); ok {
				elmValue, valid := form.convert(str, elmType)
				if !valid {
					return nil, false
				}

				rv.Index(i).Set(reflect.ValueOf(elmValue))
			} else {
				return nil, false
			}
		}

		result, valid = rv.Interface(), true
	} else if len(value) > 0 {
		if str, ok := value[0].(string); ok {
			result, valid = form.convert(str, typ)
		}
	}

	return result, valid
}

func (form Form) convert(str string, typ reflect.Type) (interface{}, bool) {
	result := interface{}(nil)
	valid := false

	switch typ.Kind() {
	case reflect.Bool:
		if parsed, err := strconv.ParseBool(str); err == nil {
			result, valid = parsed, true
		}
	case reflect.Int:
		if parsed, err := strconv.ParseInt(str, 10, 0); err == nil {
			result, valid = int(parsed), true
		}
	case reflect.Int8:
		if parsed, err := strconv.ParseInt(str, 10, 8); err == nil {
			result, valid = int8(parsed), true
		}
	case reflect.Int16:
		if parsed, err := strconv.ParseInt(str, 10, 16); err == nil {
			result, valid = int16(parsed), true
		}
	case reflect.Int32:
		if parsed, err := strconv.ParseInt(str, 10, 32); err == nil {
			result, valid = int32(parsed), true
		}
	case reflect.Int64:
		if parsed, err := strconv.ParseInt(str, 10, 64); err == nil {
			result, valid = int64(parsed), true
		}
	case reflect.Uint:
		if parsed, err := strconv.ParseUint(str, 10, 0); err == nil {
			result, valid = uint(parsed), true
		}
	case reflect.Uint8:
		if parsed, err := strconv.ParseUint(str, 10, 8); err == nil {
			result, valid = uint8(parsed), true
		}
	case reflect.Uint16:
		if parsed, err := strconv.ParseUint(str, 10, 16); err == nil {
			result, valid = uint16(parsed), true
		}
	case reflect.Uint32:
		if parsed, err := strconv.ParseUint(str, 10, 32); err == nil {
			result, valid = uint32(parsed), true
		}
	case reflect.Uint64:
		if parsed, err := strconv.ParseUint(str, 10, 64); err == nil {
			result, valid = uint64(parsed), true
		}
	case reflect.Uintptr:
		if parsed, err := strconv.ParseUint(str, 10, 8); err == nil {
			result, valid = uintptr(parsed), true
		}
	case reflect.Float32:
		if parsed, err := strconv.ParseFloat(str, 32); err == nil {
			result, valid = float32(parsed), true
		}
	case reflect.Float64:
		if parsed, err := strconv.ParseFloat(str, 64); err == nil {
			result, valid = float64(parsed), true
		}
	case reflect.String:
		result, valid = str, true
	case reflect.Struct:
		if typ == timeType {
			if parsed, err := time.Parse(time.RFC3339, str); err == nil {
				result, valid = parsed, true
			}
		}
	}

	// handle type alias
	if valid && typ.PkgPath() != "" && typ.Kind() != reflect.Struct && typ.Kind() != reflect.Slice && typ.Kind() != reflect.Array {
		rv := reflect.ValueOf(result)
		return rv.Convert(typ).Interface(), true
	}

	return result, valid
}

// GetParams returns nested param
func (form Form) GetParams(name string) (Params, bool) {
	if val, exist := form[name]; exist && len(val) == 1 {
		if par, ok := val[0].(Params); ok {
			return par, ok
		}
	}

	return nil, false
}

// GetParamsSlice returns slice of nested param
func (form Form) GetParamsSlice(name string) ([]Params, bool) {
	if val, exist := form[name]; exist {
		pars := make([]Params, len(val))
		ok := true

		for i := range val {
			pars[i], ok = val[i].(Form)

			if !ok {
				return nil, false
			}
		}

		return pars, true
	}

	return nil, false
}

func (form Form) assigns(pfield string, cfield string, index int, values interface{}) {
	if _, exist := form[pfield]; !exist && pfield != "" && cfield != "" {
		form[pfield] = []interface{}{Form{}}
	}

	if cfield != "" {
		if pfield != "" {
			form = form[pfield][0].(Form)
		}

		if stringValues, ok := values.([]string); ok {
			for _, value := range stringValues {
				form[cfield] = append(form[cfield], value)
			}
		}
	} else {
		// expand
		if index >= len(form[pfield]) {
			form[pfield] = append(form[pfield], make([]interface{}, index-len(form[pfield])+1)...)
		}

		if form[pfield][index] == nil {
			if values == nil {
				form[pfield][index] = Form{}
			} else {
				form[pfield][index] = values
			}
		}
	}
}

// ParseForm form from url values.
func ParseForm(raw url.Values) Form {
	result := make(Form, len(raw))

	for k, v := range raw {
		if len(v) == 0 {
			continue
		}

		fields := strings.FieldsFunc(k, fieldsExtractor)

		pfield, cfield := "", ""
		form := result
		for i := range fields {
			pfield = cfield
			cfield = fields[i]

			if index, err := strconv.Atoi(cfield); err == nil {
				if i == len(fields)-1 {
					form.assigns(pfield, "", index, v[0])
				} else {
					form.assigns(pfield, "", index, nil)
					form = form[pfield][index].(Form)
					cfield = "" // set cfield empty, so unnecesary nesting wont be created in the next loop
				}
			} else {
				if i == len(fields)-1 {
					form.assigns(pfield, cfield, -1, v)
				} else {
					form.assigns(pfield, cfield, -1, nil)
					if pfield != "" {
						index := len(form[pfield]) - 1
						form = form[pfield][index].(Form)
					}
				}
			}
		}
	}

	return result
}

func fieldsExtractor(c rune) bool {
	return c == '[' || c == ']' || c == '.'
}
