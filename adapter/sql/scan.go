package sql

import (
	"reflect"
	"unicode"
)

var tag = "db"

func fieldIndex(rt reflect.Type) map[string]int {
	if rt.Kind() != reflect.Struct {
		panic("must be a struct")
	}

	fields := make(map[string]int)
	for i := 0; i < rt.NumField(); i++ {
		f := rt.Field(i)
		if tag := f.Tag.Get(tag); tag != "" {
			if tag == "-" {
				continue
			}

			fields[tag] = i
		} else {
			fields[toSnake(f.Name)] = i
		}
	}

	return fields
}

// convert string to snake case
// https://gist.github.com/elwinar/14e1e897fdbe4d3432e1
func toSnake(in string) string {
	runes := []rune(in)
	length := len(runes)

	var out []rune
	for i := 0; i < length; i++ {
		if i > 0 && unicode.IsUpper(runes[i]) && ((i+1 < length && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}

	return string(out)
}
