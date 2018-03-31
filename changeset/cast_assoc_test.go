package changeset

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Inner struct {
	Field4 int
	Field5 string
}

func TestCastAssocOne(t *testing.T) {
	var inner Inner
	var entity struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": map[string]interface{}{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, params["field3"].(map[string]interface{})),
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

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssocOnePointer(t *testing.T) {
	var inner Inner
	var entity struct {
		Field1 int
		Field2 string
		Field3 *Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": map[string]interface{}{
			"field4": 4,
			"field5": "5",
		},
	}

	expectedChanges := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": changeInner(inner, params["field3"].(map[string]interface{})),
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

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.Nil(t, ch.Errors())
	assert.Equal(t, expectedTypes, ch.types)
	assert.Equal(t, expectedValues, ch.values)
	assert.Equal(t, expectedChanges, ch.Changes())
}

func TestCastAssocOneErrorParamsNotAMAp(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": "3",
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field3 is invalid", ch.Error().Error())
}

func TestCastAssocOneInnerChangesetError(t *testing.T) {
	var entity struct {
		Field1 int
		Field2 string
		Field3 Inner
	}

	changeInner := func(entity interface{}, params map[string]interface{}) *Changeset {
		ch := Cast(entity, params, []string{"field4", "field5"})
		return ch
	}

	params := map[string]interface{}{
		"field1": 1,
		"field2": "2",
		"field3": map[string]interface{}{
			"field4": "4",
		},
	}

	ch := Cast(entity, params, []string{"field1", "field2"})
	CastAssoc(ch, "field3", changeInner)

	assert.NotNil(t, ch.Errors())
	assert.Equal(t, "field4 is invalid", ch.Error().Error())
}
