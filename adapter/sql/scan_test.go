package sql

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

type User struct {
	ID        uint
	Name      string
	OtherInfo string
	OtherName string `db:"real_name"`
}

func TestFieldIndex(t *testing.T) {
	index := fieldIndex(reflect.TypeOf(User{}))
	assert.Equal(t, map[string]int{
		"id":         0,
		"name":       1,
		"other_info": 2,
		"real_name":  3,
	}, index)
}
