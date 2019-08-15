package grimoire

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Items []Item

func (it *Items) Table() string {
	return "_items"
}

func (it *Items) PrimaryField() string {
	return "_uuid"
}

func (it *Items) PrimaryValue() interface{} {
	var (
		ids = make([]interface{}, len(*it))
	)

	for i := range *it {
		ids[i] = (*it)[i].UUID
	}

	return ids
}

func TestCollection_Table(t *testing.T) {
	var (
		entities = []User{}
		rt       = reflect.TypeOf(entities).Elem()
		col      = newCollection(&entities)
	)

	// infer table name
	assert.Equal(t, "users", col.Table())

	// cached
	_, cached := tablesCache.Load(rt)
	assert.True(t, cached)

	tablesCache.Delete(rt)
}

func TestCollection_Table_usingInterface(t *testing.T) {
	var (
		entities = Items{}
		rt       = reflect.TypeOf(entities).Elem()
		col      = newCollection(&entities)
	)

	// infer table name
	assert.Equal(t, "_items", col.Table())

	// never cache
	_, cached := tablesCache.Load(rt)
	assert.False(t, cached)
}

func TestCollection_Table_usingElemInterface(t *testing.T) {
	var (
		entities = []Item{}
		rt       = reflect.TypeOf(entities).Elem()
		col      = newCollection(&entities)
	)

	// infer table name
	assert.Equal(t, "_items", col.Table())

	// cache
	_, cached := tablesCache.Load(rt)
	assert.True(t, cached)

	tablesCache.Delete(rt)
}

func TestCollection_Primary(t *testing.T) {
	var (
		entities = []User{
			{ID: 1},
			{ID: 2},
		}
		rt  = reflect.TypeOf(entities).Elem()
		col = newCollection(&entities)
	)

	// infer primary key
	assert.Equal(t, "id", col.PrimaryField())
	assert.Equal(t, []interface{}{1, 2}, col.PrimaryValue())

	// cached
	_, cached := primariesCache.Load(rt)
	assert.True(t, cached)

	entities[1].ID = 4

	// infer primary key using cache
	assert.Equal(t, "id", col.PrimaryField())
	assert.Equal(t, []interface{}{1, 4}, col.PrimaryValue())

	primariesCache.Delete(rt)
}

func TestCollection_Primary_usingInterface(t *testing.T) {
	var (
		entities = Items{
			{UUID: "abc123"},
			{UUID: "def456"},
		}
		rt  = reflect.TypeOf(entities).Elem()
		col = newCollection(&entities)
	)

	// should not be cached yet
	_, cached := primariesCache.Load(rt)
	assert.False(t, cached)

	// infer primary key
	assert.Equal(t, "_uuid", col.PrimaryField())
	assert.Equal(t, []interface{}{"abc123", "def456"}, col.PrimaryValue())

	// never cache
	_, cached = primariesCache.Load(rt)
	assert.False(t, cached)
}

func TestCollection_Primary_usingElemInterface(t *testing.T) {
	var (
		entities = []Item{
			{UUID: "abc123"},
			{UUID: "def456"},
		}
		rt  = reflect.TypeOf(entities).Elem()
		col = newCollection(&entities)
	)

	// should not be cached yet
	_, cached := primariesCache.Load(rt)
	assert.False(t, cached)

	// infer primary key
	assert.Equal(t, "_uuid", col.PrimaryField())
	assert.Equal(t, []interface{}{"abc123", "def456"}, col.PrimaryValue())

	// cache
	_, cached = primariesCache.Load(rt)
	assert.True(t, cached)

	primariesCache.Delete(rt)
}

func TestCollection_Primary_usingTag(t *testing.T) {
	var (
		entities = []struct {
			ID         uint
			ExternalID int `db:",primary"`
			Name       string
		}{
			{ExternalID: 1},
			{ExternalID: 2},
		}
		doc = newCollection(&entities)
	)

	// infer primary key
	assert.Equal(t, "external_id", doc.PrimaryField())
	assert.Equal(t, []interface{}{1, 2}, doc.PrimaryValue())
}

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
