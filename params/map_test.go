package params_test

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestMap_Exists(t *testing.T) {
	p := params.Map{"exists": true}
	assert.True(t, p.Exists("exists"))
	assert.False(t, p.Exists("not-exists"))
}

func TestMap_Get(t *testing.T) {
	p := params.Map{"exists": true}
	assert.Equal(t, true, p.Get("exists"))
	assert.Equal(t, nil, p.Get("not-exists"))
}

func TestMap_GetWithType(t *testing.T) {
	p := params.Map{
		"nil":                      (*bool)(nil),
		"incorrect type":           "some string",
		"correct type":             true,
		"slice":                    []bool{true, false},
		"slice of interface":       []interface{}{true, false},
		"slice of interface mixed": []interface{}{true, 0},
	}

	tests := []struct {
		name  string
		typ   reflect.Type
		value interface{}
		valid bool
	}{
		{
			name:  "not exist",
			typ:   reflect.TypeOf(true),
			value: nil,
			valid: true,
		},
		{
			name:  "nil",
			typ:   reflect.TypeOf(true),
			value: nil,
			valid: true,
		},
		{
			name:  "incorrect type",
			typ:   reflect.TypeOf(true),
			value: nil,
			valid: false,
		},
		{
			name:  "correct type",
			typ:   reflect.TypeOf(true),
			value: true,
			valid: true,
		},
		{
			name:  "slice",
			typ:   reflect.TypeOf([]bool{}),
			value: []bool{true, false},
			valid: true,
		},
		{
			name:  "slice of interface",
			typ:   reflect.TypeOf([]bool{}),
			value: []bool{true, false},
			valid: true,
		},
		{
			name:  "slice of interface mixed",
			typ:   reflect.TypeOf([]bool{}),
			value: nil,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, valid := p.GetWithType(tt.name, tt.typ)
			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestMap_GetParams(t *testing.T) {
	p := params.Map{
		"params.Map":       params.Map{},
		"params.Map slice": []params.Map{},
		"map":              map[string]interface{}{},
		"map slice":        []map[string]interface{}{},
		"invalid":          true,
	}

	tests := []struct {
		name  string
		param params.Params
		valid bool
	}{
		{
			name:  "params.Map",
			param: params.Map{},
			valid: true,
		},
		{
			name:  "params.Map slice",
			param: nil,
			valid: false,
		},
		{
			name:  "map",
			param: params.Map{},
			valid: true,
		},
		{
			name:  "map slice",
			param: nil,
			valid: false,
		},
		{
			name:  "invalid",
			param: nil,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, valid := p.GetParams(tt.name)
			assert.Equal(t, tt.param, param)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestMap_GetParamsSlice(t *testing.T) {
	p := params.Map{
		"params.Params slice": []params.Params{params.Map{}},
		"params.Map":          params.Map{},
		"params.Map slice":    []params.Map{params.Map{}},
		"map":                 map[string]interface{}{},
		"map slice":           []map[string]interface{}{map[string]interface{}{}},
		"invalid":             true,
	}

	tests := []struct {
		name   string
		params []params.Params
		valid  bool
	}{
		{
			name:   "params.Params slice",
			params: []params.Params{params.Map{}},
			valid:  true,
		},
		{
			name:   "params.Map",
			params: nil,
			valid:  false,
		},
		{
			name:   "params.Map slice",
			params: []params.Params{params.Map{}},
			valid:  true,
		},
		{
			name:   "map",
			params: nil,
			valid:  false,
		},
		{
			name:   "map slice",
			params: []params.Params{params.Map{}},
			valid:  true,
		},
		{
			name:   "invalid",
			params: nil,
			valid:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, valid := p.GetParamsSlice(tt.name)
			assert.Equal(t, tt.params, params)
			assert.Equal(t, tt.valid, valid)
		})
	}
}
