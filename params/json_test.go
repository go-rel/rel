package params_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

func TestJSON_Exists(t *testing.T) {
	p := params.ParseJSON(`{"exists": true}`)
	assert.True(t, p.Exists("exists"))
	assert.False(t, p.Exists("not-exists"))
}

func TestJSON_Get(t *testing.T) {
	p := params.ParseJSON(`{"exists": true}`)
	assert.Equal(t, true, p.Get("exists"))
	assert.Equal(t, nil, p.Get("not-exists"))
}

func TestJSON_GetWithType(t *testing.T) {
	p := params.ParseJSON(`{
		"nil": null,
		"incorrect type": "some string",
		"boolean": true,
		"boolean slice": [true, false],
		"number": 1,
		"number slice": [1, 2],
		"float": 1.5,
		"float slice": [1.5, 2.5],
		"string": "string",
		"string slice": ["string1", "string2"],
		"time": "2016-11-28T23:00:00+07:00",
		"time slice": ["2016-11-28T23:00:00+07:00", "2016-11-28T23:30:00+07:00"]
	}`)

	t1, _ := time.Parse(time.RFC3339, "2016-11-28T23:00:00+07:00")
	t2, _ := time.Parse(time.RFC3339, "2016-11-28T23:30:00+07:00")

	tests := []struct {
		name  string
		field string
		typ   reflect.Type
		value interface{}
		valid bool
	}{
		{
			name:  "not exist",
			field: "not exist",
			typ:   reflect.TypeOf(true),
			value: nil,
			valid: true,
		},
		{
			name:  "not exist alias",
			field: "not exist alias",
			typ:   reflect.TypeOf(Number(0)),
			value: nil,
			valid: true,
		},
		{
			name:  "not exist struct",
			field: "not exist struct",
			typ:   reflect.TypeOf(time.Time{}),
			value: nil,
			valid: true,
		},
		{
			name:  "nil",
			field: "nil",
			typ:   reflect.TypeOf(true),
			value: nil,
			valid: true,
		},
		{
			name:  "incorrect type",
			field: "incorrect type",
			typ:   reflect.TypeOf(true),
			value: nil,
			valid: false,
		},
		{
			name:  "incorrect type Number",
			field: "incorrect type",
			typ:   reflect.TypeOf(Number(0)),
			value: nil,
			valid: false,
		},
		{
			name:  "boolean",
			field: "boolean",
			typ:   reflect.TypeOf(true),
			value: true,
			valid: true,
		},
		{
			name:  "boolean slice",
			field: "boolean slice",
			typ:   reflect.TypeOf([]bool{}),
			value: []bool{true, false},
			valid: true,
		},
		{
			name:  "int",
			field: "number",
			typ:   reflect.TypeOf(int(0)),
			value: int(1),
			valid: true,
		},
		{
			name:  "int slice",
			field: "number slice",
			typ:   reflect.TypeOf([]int{}),
			value: []int{1, 2},
			valid: true,
		},
		{
			name:  "int8",
			field: "number",
			typ:   reflect.TypeOf(int8(0)),
			value: int8(1),
			valid: true,
		},
		{
			name:  "int8 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]int8{}),
			value: []int8{1, 2},
			valid: true,
		},
		{
			name:  "int16",
			field: "number",
			typ:   reflect.TypeOf(int16(0)),
			value: int16(1),
			valid: true,
		},
		{
			name:  "int16 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]int16{}),
			value: []int16{1, 2},
			valid: true,
		},
		{
			name:  "int32",
			field: "number",
			typ:   reflect.TypeOf(int32(0)),
			value: int32(1),
			valid: true,
		},
		{
			name:  "int32 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]int32{}),
			value: []int32{1, 2},
			valid: true,
		},
		{
			name:  "int64",
			field: "number",
			typ:   reflect.TypeOf(int64(0)),
			value: int64(1),
			valid: true,
		},
		{
			name:  "int64 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]int64{}),
			value: []int64{1, 2},
			valid: true,
		},
		{
			name:  "uint",
			field: "number",
			typ:   reflect.TypeOf(uint(0)),
			value: uint(1),
			valid: true,
		},
		{
			name:  "uint slice",
			field: "number slice",
			typ:   reflect.TypeOf([]uint{}),
			value: []uint{1, 2},
			valid: true,
		},
		{
			name:  "uint8",
			field: "number",
			typ:   reflect.TypeOf(uint8(0)),
			value: uint8(1),
			valid: true,
		},
		{
			name:  "uint8 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]uint8{}),
			value: []uint8{1, 2},
			valid: true,
		},
		{
			name:  "uint16",
			field: "number",
			typ:   reflect.TypeOf(uint16(0)),
			value: uint16(1),
			valid: true,
		},
		{
			name:  "uint16 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]uint16{}),
			value: []uint16{1, 2},
			valid: true,
		},
		{
			name:  "uint32",
			field: "number",
			typ:   reflect.TypeOf(uint32(0)),
			value: uint32(1),
			valid: true,
		},
		{
			name:  "uint32 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]uint32{}),
			value: []uint32{1, 2},
			valid: true,
		},
		{
			name:  "uint64",
			field: "number",
			typ:   reflect.TypeOf(uint64(0)),
			value: uint64(1),
			valid: true,
		},
		{
			name:  "uint64 slice",
			field: "number slice",
			typ:   reflect.TypeOf([]uint64{}),
			value: []uint64{1, 2},
			valid: true,
		},
		{
			name:  "uintptr",
			field: "number",
			typ:   reflect.TypeOf(uintptr(0)),
			value: uintptr(1),
			valid: true,
		},
		{
			name:  "uintptr slice",
			field: "number slice",
			typ:   reflect.TypeOf([]uintptr{}),
			value: []uintptr{1, 2},
			valid: true,
		},
		{
			name:  "type Number",
			field: "number",
			typ:   reflect.TypeOf(Number(0)),
			value: Number(1),
			valid: true,
		},
		{
			name:  "type Number slice",
			field: "number slice",
			typ:   reflect.TypeOf([]Number{}),
			value: []Number{1, 2},
			valid: true,
		},
		{
			name:  "float32",
			field: "float",
			typ:   reflect.TypeOf(float32(0)),
			value: float32(1.5),
			valid: true,
		},
		{
			name:  "float32 slice",
			field: "float slice",
			typ:   reflect.TypeOf([]float32{}),
			value: []float32{1.5, 2.5},
			valid: true,
		},
		{
			name:  "float64",
			field: "float",
			typ:   reflect.TypeOf(float64(0)),
			value: float64(1.5),
			valid: true,
		},
		{
			name:  "float64 slice",
			field: "float slice",
			typ:   reflect.TypeOf([]float64{}),
			value: []float64{1.5, 2.5},
			valid: true,
		},
		{
			name:  "string",
			field: "string",
			typ:   reflect.TypeOf(""),
			value: "string",
			valid: true,
		},
		{
			name:  "string slice",
			field: "string slice",
			typ:   reflect.TypeOf([]string{}),
			value: []string{"string1", "string2"},
			valid: true,
		},
		{
			name:  "time",
			field: "time",
			typ:   reflect.TypeOf(time.Time{}),
			value: t1,
			valid: true,
		},
		{
			name:  "time slice",
			field: "time slice",
			typ:   reflect.TypeOf([]time.Time{}),
			value: []time.Time{t1, t2},
			valid: true,
		},
		{
			name:  "time invalid",
			field: "string",
			typ:   reflect.TypeOf(time.Time{}),
			value: nil,
			valid: false,
		},
		{
			name:  "time slice invalid",
			field: "string slice",
			typ:   reflect.TypeOf([]time.Time{}),
			value: nil,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			value, valid := p.GetWithType(tt.field, tt.typ)
			assert.Equal(t, tt.value, value)
			assert.Equal(t, tt.valid, valid)
		})
	}
}

func TestJSON_GetParams(t *testing.T) {
	p := params.ParseJSON(`{
		"object": {},
		"array": [],
		"value": true
	}`)

	tests := []struct {
		name  string
		valid bool
	}{
		{
			name:  "object",
			valid: true,
		},
		{
			name:  "array",
			valid: false,
		},
		{
			name:  "value",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			param, valid := p.GetParams(tt.name)
			assert.Equal(t, tt.valid, valid)
			assert.Equal(t, tt.valid, param != nil)
		})
	}
}

func TestJSON_GetParamsSlice(t *testing.T) {
	p := params.ParseJSON(`{
		"array of object": [{}, {}],
		"array of array": [[], []],
		"array of value": [true, false],
		"array of mixed": [{}, true],
		"array": [],
		"object": {},
		"value": true
	}`)

	tests := []struct {
		name  string
		valid bool
	}{
		{
			name:  "array of object",
			valid: true,
		},
		{
			name:  "array of array",
			valid: false,
		},
		{
			name:  "array of value",
			valid: false,
		},
		{
			name:  "array of mixed",
			valid: false,
		},
		{
			name:  "array",
			valid: true,
		},
		{
			name:  "object",
			valid: false,
		},
		{
			name:  "value",
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			params, valid := p.GetParamsSlice(tt.name)
			assert.Equal(t, tt.valid, valid)
			assert.Equal(t, tt.valid, params != nil)
		})
	}
}
