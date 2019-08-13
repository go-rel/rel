package grimoire

import (
	"database/sql"
	"reflect"
	"testing"
	"time"

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

func (i Item) Fields() map[string]int {
	return map[string]int{
		"_uuid":  0,
		"_price": 1,
	}
}

func (i Item) Types() []reflect.Type {
	return []reflect.Type{String, Int}
}

func (i Item) Values() []interface{} {
	return []interface{}{i.UUID, i.Price}
}

func (i *Item) Scanners(fields []string) []interface{} {
	var (
		scanners  = make([]interface{}, len(fields))
		tempValue = sql.RawBytes{}
	)

	for index, field := range fields {
		switch field {
		case "_uuid":
			scanners[index] = Nullable(&i.UUID)
		case "_price":
			scanners[index] = Nullable(&i.Price)
		default:
			scanners[index] = &tempValue
		}
	}

	return scanners
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
		entity = Item{}
		rt     = reflect.TypeOf(entity)
		doc    = newDocument(&entity)
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
		entity = Item{
			UUID: "abc123",
		}
		rt  = reflect.TypeOf(entity)
		doc = newDocument(&entity)
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

func TestDocument_Fields(t *testing.T) {
	var (
		entity = struct {
			A string
			B *int
			C []byte     `db:",primary"`
			D bool       `db:"D"`
			E []*float64 `db:"-"`
		}{}
		rt     = reflect.TypeOf(entity)
		doc    = newDocument(&entity)
		fields = map[string]int{
			"a": 0,
			"b": 1,
			"c": 2,
			"D": 3,
		}
	)

	assert.Equal(t, fields, doc.Fields())

	_, cached := fieldsCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, fields, doc.Fields())
}

func TestDocument_Fields_usingInterface(t *testing.T) {
	var (
		entity = Item{}
		rt     = reflect.TypeOf(entity)
		doc    = newDocument(&entity)
	)

	assert.Equal(t, entity.Fields(), doc.Fields())

	_, cached := fieldsCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Types(t *testing.T) {
	var (
		entity = struct {
			A string
			B *int
			C []byte
			D bool
			E []*float64
			F userDefined
			G time.Time
		}{}
		rt    = reflect.TypeOf(entity)
		doc   = newDocument(&entity)
		types = []reflect.Type{
			String,
			Int,
			Bytes,
			Bool,
			reflect.TypeOf([]float64{}),
			reflect.TypeOf(userDefined(0)),
			Time,
		}
	)
	assert.Equal(t, types, doc.Types())

	_, cached := typesCache.Load(rt)
	assert.True(t, cached)

	assert.Equal(t, types, doc.Types())
}

func TestDocument_Types_usingInterface(t *testing.T) {
	var (
		entity = Item{}
		rt     = reflect.TypeOf(entity)
		doc    = newDocument(&entity)
	)

	_, cached := typesCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, entity.Types(), doc.Types())

	_, cached = typesCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Scanners(t *testing.T) {
	var (
		address = "address"
		entity  = struct {
			ID      int
			Name    string
			Skip    bool `db:"-"`
			Number  float64
			Address *string
			Data    []byte
		}{
			ID:      1,
			Name:    "name",
			Number:  10.5,
			Address: &address,
			Data:    []byte("data"),
		}
		doc      = newDocument(&entity)
		fields   = []string{"name", "id", "skip", "data", "number", "address", "not_exist"}
		scanners = []interface{}{
			Nullable(&entity.Name),
			Nullable(&entity.ID),
			&sql.RawBytes{},
			Nullable(&entity.Data),
			Nullable(&entity.Number),
			&entity.Address,
			&sql.RawBytes{},
		}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Scanners_usingInterface(t *testing.T) {
	var (
		entity = Item{
			UUID:  "abc123",
			Price: 100,
		}
		doc      = newDocument(&entity)
		fields   = []string{"_uuid", "_price"}
		scanners = []interface{}{Nullable(&entity.UUID), Nullable(&entity.Price)}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Scanners_sqlScanner(t *testing.T) {
	var (
		entity   = sql.NullBool{}
		doc      = newDocument(&entity)
		fields   = []string{}
		scanners = []interface{}{&sql.NullBool{}}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Values(t *testing.T) {
	var (
		address = "address"
		entity  = struct {
			ID      int
			Name    string
			Skip    bool `db:"-"`
			Number  float64
			Address *string
			Data    []byte
		}{
			ID:      1,
			Name:    "name",
			Number:  10.5,
			Address: &address,
			Data:    []byte("data"),
		}
		doc    = newDocument(&entity)
		values = []interface{}{1, "name", 10.5, address, []byte("data")}
	)

	assert.Equal(t, values, doc.Values())
}

func TestDocument_Values_usingInterface(t *testing.T) {
	var (
		entity = Item{
			UUID:  "abc123",
			Price: 100,
		}
		doc    = newDocument(&entity)
		values = []interface{}{"abc123", 100}
	)

	assert.Equal(t, values, doc.Values())
}
