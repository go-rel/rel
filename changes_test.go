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
			Fields: map[string]int{
				"field1":   0,
				"field2":   1,
				"field3":   2,
				"field4":   3,
				"field5":   4,
				"field6=?": 5,
			},
			Changes: []Change{
				Set("field1", 10),
				Inc("field2"),
				IncBy("field3", 2),
				Dec("field4"),
				DecBy("field5", 2),
				ChangeFragment("field6=?", true),
			},
			Assoc: map[string]int{},
		}
	)

	assert.Equal(t, changes, BuildChanges(changers...))
}
