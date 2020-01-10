package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildChanges(t *testing.T) {
	var (
		changers = []Changer{
			Set("field1", 10),
			Inc("field2"),
			IncBy("field3", 2),
			Dec("field4"),
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
				Set("field1", 10),
				Inc("field2"),
				IncBy("field3", 2),
				Dec("field4"),
				DecBy("field5", 2),
				ChangeFragment("field6=?", true),
			},
			assoc: map[string]int{},
		}
		result, err = ApplyChanges(nil, changers...)
	)

	assert.Nil(t, err)
	assert.Equal(t, changes, result)
}
