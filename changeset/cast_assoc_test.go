package changeset

import (
	"reflect"
	"testing"

	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/params"
	"github.com/stretchr/testify/assert"
)

type Inner struct {
	Field4 int
	Field5 string
}

func TestCastAssoc_one(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": params.Map{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, input["field3"].(params.Map)),
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(inner),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": inner,
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssoc_oneDifferentSourceField(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"fieldX": params.Map{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, input["fieldX"].(params.Map)),
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(inner),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": inner,
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner, SourceField("fieldX"))

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssoc_onePointer(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 *Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": params.Map{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, input["field3"].(params.Map)),
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf(inner),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssoc_oneErrorParamsNotAMap(t *testing.T) {
	var data struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": "3",
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssoc_oneInnerChangesetError(t *testing.T) {
	var data struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": params.Map{
			"field4": "4",
		},
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field4 is invalid", ch.Error().Error())
	assert.Equal(t, "field3.field4", ch.Error().(errors.Error).Field)
}

func TestCastAssoc_many(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": []params.Map{
			{
				"field4": 14,
				"field5": "15",
			},
			{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	field3 := input["field3"].([]params.Map)
	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []*Changeset{
			changeInner(inner, field3[0]),
			changeInner(inner, field3[1]),
		},
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf([]Inner{}),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": []Inner(nil),
	}

	// with map assoc
	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssoc_manyPointer(t *testing.T) {
	var inner Inner
	var data struct {
		Field1 int
		Field2 string
		Field3 []*Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": []params.Map{
			{
				"field4": 14,
				"field5": "15",
			},
			{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	field3 := input["field3"].([]params.Map)
	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": []*Changeset{
			changeInner(inner, field3[0]),
			changeInner(inner, field3[1]),
		},
	}

	expectedTypes := map[string]reflect.Type{
		"field1": reflect.TypeOf(0),
		"field2": reflect.TypeOf(""),
		"field3": reflect.TypeOf([]Inner{}),
	}

	expectedValues := map[string]interface{}{
		"field1": 0,
		"field2": "",
		"field3": []*Inner(nil),
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssoc_manyErrorParamsNotASliceOfAMap(t *testing.T) {
	var data struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": "3",
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssoc_manyErrorMixed(t *testing.T) {
	var data struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": []interface{}{
			params.Map{
				"field4": 14,
				"field5": "15",
			},
			"3",
		},
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssoc_manyInnerChangesetError(t *testing.T) {
	var data struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	input := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": []params.Map{
			{
				"field4": "14",
			},
		},
	}

	ch := Cast(data, input, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field4 is invalid", ch.Error().Error())
	assert.Equal(t, "field3[0].field4", ch.Error().(errors.Error).Field)
}

func TestCastAssoc_optionRequired(t *testing.T) {
	var data struct {
		Field1 int
		Field2 string
		Field3 []Inner
	}

	changeInner := func(data interface{}, input params.Params) *Changeset {
		ch := Cast(data, input, []string{"field4", "field5"})
		return ch
	}

	invalidInput := params.Map{
		"field1": 1,
		"field2": "2",
	}

	validInput := params.Map{
		"field1": 1,
		"field2": "2",
		"field3": []params.Map{
			{
				"field4": 14,
				"field5": "15",
			},
			{
				"field4": 24,
				"field5": "25",
			},
		},
	}

	invalidCh := Cast(data, invalidInput, []string{"field1", "field2"})
	CastAssoc(invalidCh, "field3", changeInner, Required(true))

	validCh := Cast(data, validInput, []string{"field1", "field2"})
	CastAssoc(validCh, "field3", changeInner, Required(true))

	assert.NotNil(t, invalidCh.Errors())
	assert.Nil(t, validCh.Errors())
}
