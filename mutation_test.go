package rel

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Testentity struct {
	Field1 string `db:",primary"`
	Field2 bool
	Field3 *string
	Field4 int
	Field5 int
}

func TestApply(t *testing.T) {
	var (
		entity   = Testentity{}
		doc      = NewDocument(&entity)
		mutators = []Mutator{
			Set("field1", "string"),
			Set("field2", true),
			Set("field3", "string pointer"),
			IncBy("field4", 2),
			DecBy("field5", 2),
			SetFragment("field6=?", true),
		}
		mutation = Mutation{
			Cascade: true,
			Mutates: map[string]Mutate{
				"field1":   Set("field1", "string"),
				"field2":   Set("field2", true),
				"field3":   Set("field3", "string pointer"),
				"field4":   IncBy("field4", 2),
				"field5":   DecBy("field5", 2),
				"field6=?": SetFragment("field6=?", true),
			},
			Reload: true,
		}
	)

	assert.Equal(t, mutation, Apply(doc, mutators...))
	assert.Equal(t, "string", entity.Field1)
	assert.Equal(t, true, entity.Field2)
	assert.Equal(t, "string pointer", *entity.Field3)

	// non set op won't update the struct
	assert.Equal(t, 0, entity.Field4)
	assert.Equal(t, 0, entity.Field5)
}

func TestApply_Options(t *testing.T) {
	var (
		entity   = Testentity{}
		doc      = NewDocument(&entity)
		mutators = []Mutator{
			Unscoped(true),
			Reload(true),
			Cascade(true),
			OnConflictIgnore(),
		}
		mutation = Mutation{
			Unscoped: true,
			Cascade:  true,
			Mutates: map[string]Mutate{
				"field2": Set("field2", false),
				"field3": Set("field3", nil),
				"field4": Set("field4", 0),
				"field5": Set("field5", 0),
			},
			Reload:     true,
			OnConflict: OnConflict{Keys: []string{"field1"}, Ignore: true},
		}
	)

	assert.Equal(t, mutation, Apply(doc, mutators...))
}

func TestApplyMutation_setValueError(t *testing.T) {
	var (
		entity = Testentity{}
		doc    = NewDocument(&entity)
	)

	assert.Panics(t, func() {
		Apply(doc, Set("field1", 1))
	})
	assert.Equal(t, "", entity.Field1)
}

func TestApplyMutation_incValueError(t *testing.T) {
	var (
		entity = Testentity{}
		doc    = NewDocument(&entity)
	)

	assert.Panics(t, func() {
		Apply(doc, Inc("field1"))
	})
	assert.Equal(t, "", entity.Field1)
}

func TestApplyMutation_unknownFieldValueError(t *testing.T) {
	var (
		entity = Testentity{}
		doc    = NewDocument(&entity)
	)

	assert.Panics(t, func() {
		Apply(doc, Dec("field0"))
	})
	assert.Equal(t, "", entity.Field1)
}

func TestApplyMutation_Reload(t *testing.T) {
	var (
		entity   = Testentity{}
		doc      = NewDocument(&entity)
		mutators = []Mutator{
			Set("field1", "string"),
			Reload(true),
		}
		mutation = Mutation{
			Mutates: map[string]Mutate{
				"field1": Set("field1", "string"),
			},
			Reload:  true,
			Cascade: true,
		}
	)

	assert.Equal(t, mutation, Apply(doc, mutators...))
	assert.Equal(t, "string", entity.Field1)
}

func TestApplyMutation_Cascade(t *testing.T) {
	var (
		entity   = Testentity{}
		doc      = NewDocument(&entity)
		mutators = []Mutator{
			Set("field1", "string"),
			Cascade(false),
		}
		mutation = Mutation{
			Mutates: map[string]Mutate{
				"field1": Set("field1", "string"),
			},
			Cascade: false,
		}
	)

	assert.Equal(t, mutation, Apply(doc, mutators...))
	assert.Equal(t, "string", entity.Field1)
}

func TestMutator_String(t *testing.T) {
	assert.Equal(t, "rel.Set(\"field\", 1)", fmt.Sprint(Set("field", 1)))
	assert.Equal(t, "rel.Set(\"field\", true)", fmt.Sprint(Set("field", true)))
	assert.Equal(t, "rel.Set(\"field\", \"value\")", fmt.Sprint(Set("field", "value")))
	assert.Equal(t, "rel.IncBy(\"count\", -1)", fmt.Sprint(Dec("count")))
	assert.Equal(t, "rel.IncBy(\"count\", 1)", fmt.Sprint(Inc("count")))
	assert.Equal(t, "rel.SetFragment(\"field = (?, ?, ?)\", 1, true, \"value\")", fmt.Sprint(SetFragment("field = (?, ?, ?)", 1, true, "value")))
	assert.Equal(t, "rel.Cascade(true)", fmt.Sprint(Cascade(true)))
}
