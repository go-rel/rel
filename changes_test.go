package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestRecord struct {
	Field1 string
	Field2 bool
	Field3 *string
	Field4 int
	Field5 int
}

func TestApplyChanges(t *testing.T) {
	var (
		record   = TestRecord{}
		doc      = NewDocument(&record)
		changers = []Changer{
			Set("field1", "string"),
			Set("field2", true),
			Set("field3", "string pointer"),
			IncBy("field4", 2),
			DecBy("field5", 2),
			ChangeFragment("field6=?", true),
		}
		changes = Changes{
			fields: map[string]int{
				"field1":   0,
				"field2":   1,
				"field3":   2,
				"field4":   3,
				"field5":   4,
				"field6=?": 5,
			},
			changes: []Change{
				Set("field1", "string"),
				Set("field2", true),
				Set("field3", "string pointer"),
				IncBy("field4", 2),
				DecBy("field5", 2),
				ChangeFragment("field6=?", true),
			},
			assoc:  map[string]int{},
			reload: true,
		}
	)

	assert.Equal(t, changes, ApplyChanges(doc, changers...))
	assert.Equal(t, "string", record.Field1)
	assert.Equal(t, true, record.Field2)
	assert.Equal(t, "string pointer", *record.Field3)

	// non set op won't update the struct
	assert.Equal(t, 0, record.Field4)
	assert.Equal(t, 0, record.Field5)
}

func TestApplyChanges_setValueError(t *testing.T) {
	var (
		record = TestRecord{}
		doc    = NewDocument(&record)
	)

	assert.Panics(t, func() {
		ApplyChanges(doc, Set("field1", 1))
	})
	assert.Equal(t, "", record.Field1)
}

func TestApplyChanges_incValueError(t *testing.T) {
	var (
		record = TestRecord{}
		doc    = NewDocument(&record)
	)

	assert.Panics(t, func() {
		ApplyChanges(doc, IncBy("field1", 2))
	})
	assert.Equal(t, "", record.Field1)
}

func TestApplyChanges_unknownFieldValueError(t *testing.T) {
	var (
		record = TestRecord{}
		doc    = NewDocument(&record)
	)

	assert.Panics(t, func() {
		ApplyChanges(doc, DecBy("field0", 2))
	})
	assert.Equal(t, "", record.Field1)
}
