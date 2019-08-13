package grimoire

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Item struct {
	UUID  string
	Price int
}

func (i Item) Table() string {
	return "_items"
}

func (i Item) PrimaryField() string {
	return "_uuid"
}

func (i Item) PrimaryValue() interface{} {
	return i.UUID
}

func TestDocument_Table(t *testing.T) {
	type User struct{}

	var (
		user = User{}
		rt   = reflect.TypeOf(user)
		doc  = newDocument(&user)
	)

	// infer table name
	assert.Equal(t, "users", doc.Table())

	// cached
	_, cached := tablesCache.Load(rt)
	assert.True(t, cached)
}

func TestDocument_Table_usingInterface(t *testing.T) {
	var (
		item = Item{}
		rt   = reflect.TypeOf(item)
		doc  = newDocument(&item)
	)

	// infer table name
	assert.Equal(t, "_items", doc.Table())

	// never cache
	_, cached := tablesCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Primary(t *testing.T) {
	var (
		user = User{ID: 1}
		rt   = reflect.TypeOf(user)
		doc  = newDocument(&user)
	)

	// infer primary key
	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 1, doc.PrimaryValue())

	// cached
	_, cached := primariesCache.Load(rt)
	assert.True(t, cached)

	user.ID = 2

	// infer primary key using cache
	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 2, doc.PrimaryValue())
}

func TestDocument_Primary_usingInterface(t *testing.T) {
	var (
		item = Item{
			UUID: "abc123",
		}
		rt  = reflect.TypeOf(item)
		doc = newDocument(&item)
	)

	// should not be cached yet
	_, cached := primariesCache.Load(rt)
	assert.False(t, cached)

	// infer primary key
	assert.Equal(t, "_uuid", doc.PrimaryField())
	assert.Equal(t, "abc123", doc.PrimaryValue())

	// never cache
	_, cached = primariesCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Primary_usingTag(t *testing.T) {
	var (
		entity = struct {
			ID         uint
			ExternalID int `db:",primary"`
			Name       string
		}{
			ExternalID: 12345,
		}
		doc = newDocument(&entity)
	)

	// infer primary key
	assert.Equal(t, "external_id", doc.PrimaryField())
	assert.Equal(t, 12345, doc.PrimaryValue())
}

func TestDocument_Primary_usingTagAmdCustomName(t *testing.T) {
	var (
		entity = struct {
			ID         uint
			ExternalID int `db:"partner_id,primary"`
			Name       string
		}{
			ExternalID: 1111,
		}
		doc = newDocument(&entity)
	)

	// infer primary key
	assert.Equal(t, "partner_id", doc.PrimaryField())
	assert.Equal(t, 1111, doc.PrimaryValue())
}

func TestDocument_Primary_notFound(t *testing.T) {
	var (
		entity = struct {
			ExternalID int
			Name       string
		}{}
		doc = newDocument(&entity)
	)

	assert.Panics(t, func() {
		doc.PrimaryField()
	})

	assert.Panics(t, func() {
		doc.PrimaryValue()
	})
}
