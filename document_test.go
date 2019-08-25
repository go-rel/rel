package grimoire

import (
	"database/sql"
	"fmt"
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

func (i Item) BelongsTo() []string {
	return nil
}

func (i Item) HasOne() []string {
	return nil
}

func (i Item) HasMany() []string {
	return nil
}

func (i Item) Association(field string) Association {
	return nil
}

func TestDocument_Table(t *testing.T) {
	var (
		record = User{}
		rt     = reflect.TypeOf(record)
		doc    = newDocument(&record)
	)

	// infer table name
	assert.Equal(t, "users", doc.Table())

	// cached
	_, cached := tablesCache.Load(rt)
	assert.True(t, cached)
}

func TestDocument_Table_usingInterface(t *testing.T) {
	var (
		record = Item{}
		rt     = reflect.TypeOf(record)
		doc    = newDocument(&record)
	)

	// infer table name
	assert.Equal(t, "_items", doc.Table())

	// never cache
	_, cached := tablesCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Primary(t *testing.T) {
	var (
		record = User{ID: 1}
		rt     = reflect.TypeOf(record)
		doc    = newDocument(&record)
	)

	// infer primary key
	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 1, doc.PrimaryValue())

	// cached
	_, cached := primariesCache.Load(rt)
	assert.True(t, cached)

	record.ID = 2

	// infer primary key using cache
	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 2, doc.PrimaryValue())

	primariesCache.Delete(rt)
}

func TestDocument_Primary_usingInterface(t *testing.T) {
	var (
		record = Item{
			UUID: "abc123",
		}
		rt  = reflect.TypeOf(record)
		doc = newDocument(&record)
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
		record = struct {
			ID         uint
			ExternalID int `db:",primary"`
			Name       string
		}{
			ExternalID: 12345,
		}
		doc = newDocument(&record)
	)

	// infer primary key
	assert.Equal(t, "external_id", doc.PrimaryField())
	assert.Equal(t, 12345, doc.PrimaryValue())
}

func TestDocument_Primary_usingTagAmdCustomName(t *testing.T) {
	var (
		record = struct {
			ID         uint
			ExternalID int `db:"partner_id,primary"`
			Name       string
		}{
			ExternalID: 1111,
		}
		doc = newDocument(&record)
	)

	// infer primary key
	assert.Equal(t, "partner_id", doc.PrimaryField())
	assert.Equal(t, 1111, doc.PrimaryValue())
}

func TestDocument_Primary_notFound(t *testing.T) {
	var (
		record = struct {
			ExternalID int
			Name       string
		}{}
		doc = newDocument(&record)
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
		record = struct {
			A string
			B *int
			C []byte     `db:",primary"`
			D bool       `db:"D"`
			E []*float64 `db:"-"`
		}{}
		rt     = reflect.TypeOf(record)
		doc    = newDocument(&record)
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

	fieldsCache.Delete(rt)
}

func TestDocument_Fields_usingInterface(t *testing.T) {
	var (
		record = Item{}
		rt     = reflect.TypeOf(record)
		doc    = newDocument(&record)
	)

	assert.Equal(t, record.Fields(), doc.Fields())

	_, cached := fieldsCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Types(t *testing.T) {
	var (
		record = struct {
			A string
			B *int
			C []byte
			D bool
			E []*float64
			F userDefined
			G time.Time
		}{}
		rt    = reflect.TypeOf(record)
		doc   = newDocument(&record)
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

	typesCache.Delete(rt)
}

func TestDocument_Types_usingInterface(t *testing.T) {
	var (
		record = Item{}
		rt     = reflect.TypeOf(record)
		doc    = newDocument(&record)
	)

	_, cached := typesCache.Load(rt)
	assert.False(t, cached)

	assert.Equal(t, record.Types(), doc.Types())

	_, cached = typesCache.Load(rt)
	assert.False(t, cached)
}

func TestDocument_Scanners(t *testing.T) {
	var (
		address = "address"
		record  = struct {
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
		doc      = newDocument(&record)
		fields   = []string{"name", "id", "skip", "data", "number", "address", "not_exist"}
		scanners = []interface{}{
			Nullable(&record.Name),
			Nullable(&record.ID),
			&sql.RawBytes{},
			Nullable(&record.Data),
			Nullable(&record.Number),
			&record.Address,
			&sql.RawBytes{},
		}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Scanners_usingInterface(t *testing.T) {
	var (
		record = Item{
			UUID:  "abc123",
			Price: 100,
		}
		doc      = newDocument(&record)
		fields   = []string{"_uuid", "_price"}
		scanners = []interface{}{Nullable(&record.UUID), Nullable(&record.Price)}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Scanners_sqlScanner(t *testing.T) {
	var (
		record   = sql.NullBool{}
		doc      = newDocument(&record)
		fields   = []string{}
		scanners = []interface{}{&sql.NullBool{}}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Values(t *testing.T) {
	var (
		address = "address"
		record  = struct {
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
		doc    = newDocument(&record)
		values = []interface{}{1, "name", 10.5, address, []byte("data")}
	)

	assert.Equal(t, values, doc.Values())
}

func TestDocument_Values_usingInterface(t *testing.T) {
	var (
		record = Item{
			UUID:  "abc123",
			Price: 100,
		}
		doc    = newDocument(&record)
		values = []interface{}{"abc123", 100}
	)

	assert.Equal(t, values, doc.Values())
}

func TestDocument_Slice(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			doc = newDocument(&Item{})
		)

		doc.Reset()
		assert.Equal(t, 1, doc.Len())
		assert.Equal(t, doc, doc.Get(0))
		assert.Equal(t, doc, doc.Add())
	})
}

func TestDocument_Association(t *testing.T) {
	tests := []struct {
		name      string
		record    interface{}
		belongsTo []string
		hasOne    []string
		hasMany   []string
	}{
		{
			name:    "User",
			record:  &User{},
			hasOne:  []string{"address"},
			hasMany: []string{"transactions"},
		},
		{
			name:    "User Cached",
			record:  &User{},
			hasOne:  []string{"address"},
			hasMany: []string{"transactions"},
		},
		{
			name:      "Transaction",
			record:    &Transaction{},
			belongsTo: []string{"buyer"},
		},
		{
			name:      "Address",
			record:    &Address{},
			belongsTo: []string{"user"},
		},
		{
			name:   "Item",
			record: &Item{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				doc = newDocument(test.record)
			)

			assert.Equal(t, test.belongsTo, doc.BelongsTo())
			assert.Equal(t, test.hasOne, doc.HasOne())
			assert.Equal(t, test.hasMany, doc.HasMany())
		})
	}
}

func TestDocument_Association_interface(t *testing.T) {
	var (
		doc = newDocument(&Item{})
	)

	assert.NotPanics(t, func() {
		assert.Nil(t, doc.Association("empty"))
	})
}

func TestDocument(t *testing.T) {
	tests := []struct {
		record interface{}
		panics bool
	}{
		{
			record: &User{},
		},
		{
			record: newDocument(&User{}),
		},
		{
			record: reflect.ValueOf(&User{}),
		},
		{
			record: reflect.ValueOf(User{}),
			panics: true,
		},
		{
			record: reflect.ValueOf(&[]User{}),
			panics: true,
		},
		{
			record: reflect.TypeOf(&User{}),
			panics: true,
		},
		{
			record: nil,
			panics: true,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%T", test.record), func(t *testing.T) {
			if test.panics {
				assert.Panics(t, func() {
					newDocument(test.record)
				})
			} else {
				assert.NotPanics(t, func() {
					newDocument(test.record)
				})
			}
		})
	}
}

func TestDocument_notPtr(t *testing.T) {
	assert.Panics(t, func() {
		newDocument(User{}).Table()
	})
}

func TestDocument_notPtrOfStruct(t *testing.T) {
	assert.Panics(t, func() {
		i := 1
		newDocument(&i).Table()
	})
}
