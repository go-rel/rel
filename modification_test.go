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

func TestApplyModification(t *testing.T) {
	var (
		record    = TestRecord{}
		doc       = NewDocument(&record)
		modifiers = []Modifier{
			Set("field1", "string"),
			Set("field2", true),
			Set("field3", "string pointer"),
			IncBy("field4", 2),
			DecBy("field5", 2),
			SetFragment("field6=?", true),
		}
		modification = Modification{
			Modifies: map[string]Modify{
				"field1":   Set("field1", "string"),
				"field2":   Set("field2", true),
				"field3":   Set("field3", "string pointer"),
				"field4":   IncBy("field4", 2),
				"field5":   DecBy("field5", 2),
				"field6=?": SetFragment("field6=?", true),
			},
			Assoc:  map[string]AssocModification{},
			Reload: true,
		}
	)

	assert.Equal(t, modification, Apply(doc, modifiers...))
	assert.Equal(t, "string", record.Field1)
	assert.Equal(t, true, record.Field2)
	assert.Equal(t, "string pointer", *record.Field3)

	// non set op won't update the struct
	assert.Equal(t, 0, record.Field4)
	assert.Equal(t, 0, record.Field5)
}

func TestApplyModification_Reload(t *testing.T) {
	var (
		record    = TestRecord{}
		doc       = NewDocument(&record)
		modifiers = []Modifier{
			Set("field1", "string"),
			Reload(true),
		}
		modification = Modification{
			Modifies: map[string]Modify{
				"field1": Set("field1", "string"),
			},
			Assoc:  map[string]AssocModification{},
			Reload: true,
		}
	)

	assert.Equal(t, modification, Apply(doc, modifiers...))
	assert.Equal(t, "string", record.Field1)
}

func TestApplyModification_setValueError(t *testing.T) {
	var (
		record = TestRecord{}
		doc    = NewDocument(&record)
	)

	assert.Panics(t, func() {
		Apply(doc, Set("field1", 1))
	})
	assert.Equal(t, "", record.Field1)
}

func TestApplyModification_incValueError(t *testing.T) {
	var (
		record = TestRecord{}
		doc    = NewDocument(&record)
	)

	assert.Panics(t, func() {
		Apply(doc, Inc("field1"))
	})
	assert.Equal(t, "", record.Field1)
}

func TestApplyModification_unknownFieldValueError(t *testing.T) {
	var (
		record = TestRecord{}
		doc    = NewDocument(&record)
	)

	assert.Panics(t, func() {
		Apply(doc, Dec("field0"))
	})
	assert.Equal(t, "", record.Field1)
}
