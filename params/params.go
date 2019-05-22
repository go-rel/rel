// Package params defines different types of params used for changeset input.
package params

import (
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// Params is interface used by changeset when casting parameters to changeset.
type Params interface {
	Exists(name string) bool
	Get(name string) interface{}
	GetWithType(name string, typ reflect.Type) (interface{}, bool)
	GetParams(name string) (Params, bool)
	GetParamsSlice(name string) ([]Params, bool)
}
