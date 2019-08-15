package grimoire

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCollection_Slice(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			users = []User{}
			col   = newCollection(&users)
		)

		assert.Equal(t, 0, col.Len())

		doc := col.Add()
		assert.Len(t, users, 1)
		assert.Equal(t, 1, col.Len())
		assert.Equal(t, newDocument(&users[0]), doc)
		assert.Equal(t, newDocument(&users[0]), col.Get(0))

		col.Reset()
		assert.Len(t, users, 0)
		assert.Equal(t, 0, col.Len())
	})
}

func TestCollection(t *testing.T) {
	tests := []struct {
		entity interface{}
		panics bool
	}{
		{
			entity: &[]User{},
		},
		{
			entity: newCollection(&[]User{}),
		},
		{
			entity: reflect.ValueOf(&[]User{}),
		},
		{
			entity: reflect.ValueOf([]User{}),
			panics: true,
		},
		{
			entity: reflect.ValueOf(&User{}),
			panics: true,
		},
		{
			entity: reflect.TypeOf(&[]User{}),
			panics: true,
		},
		{
			entity: nil,
			panics: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test.entity), func(t *testing.T) {
			if test.panics {
				assert.Panics(t, func() {
					newCollection(test.entity)
				})
			} else {
				assert.NotPanics(t, func() {
					newCollection(test.entity)
				})
			}
		})
	}
}

func TestCollection_notPtr(t *testing.T) {
	assert.Panics(t, func() {
		newCollection([]User{}).Table()
	})
}

func TestCollection_notPtrOfSlice(t *testing.T) {
	assert.Panics(t, func() {
		newCollection(&User{}).Table()
	})
}
