package rel

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

func TestDocument_Table(t *testing.T) {
	var (
		record = User{}
		rt     = reflect.TypeOf(record)
		doc    = NewDocument(&record)
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
		doc    = NewDocument(&record)
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
		doc    = NewDocument(&record)
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
		doc = NewDocument(&record)
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
		doc = NewDocument(&record)
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
		doc = NewDocument(&record)
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
		doc = NewDocument(&record)
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
		doc    = NewDocument(&record)
		fields = []string{"a", "b", "c", "D"}
	)

	assert.Equal(t, fields, doc.Fields())
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
		doc   = NewDocument(&record)
		types = map[string]reflect.Type{
			"a": reflect.TypeOf(""),
			"b": reflect.TypeOf(0),
			"c": reflect.TypeOf([]byte{}),
			"d": reflect.TypeOf(false),
			"e": reflect.TypeOf([]float64{}),
			"f": reflect.TypeOf(userDefined(0)),
			"g": reflect.TypeOf(time.Time{}),
		}
	)

	for field, etyp := range types {
		typ, ok := doc.Type(field)
		assert.True(t, ok)
		assert.Equal(t, etyp, typ)
	}
}

func TestDocument_Value(t *testing.T) {
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
		doc    = NewDocument(&record)
		values = map[string]interface{}{
			"id":     1,
			"name":   "name",
			"number": 10.5,
			"data":   []byte("data"),
		}
	)

	for field, evalue := range values {
		value, ok := doc.Value(field)
		assert.True(t, ok)
		assert.Equal(t, evalue, value)
	}
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
		doc      = NewDocument(&record)
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

func TestDocument_Slice(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			doc = NewDocument(&Item{})
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
				doc = NewDocument(test.record)
			)

			assert.Equal(t, test.belongsTo, doc.BelongsTo())
			assert.Equal(t, test.hasOne, doc.HasOne())
			assert.Equal(t, test.hasMany, doc.HasMany())
		})
	}
}

// func TestDocument_Association_interface(t *testing.T) {
// 	var (
// 		doc = NewDocument(&Item{})
// 	)

// 	assert.NotPanics(t, func() {
// 		assert.Nil(t, doc.Association("empty"))
// 	})
// }

func TestDocument(t *testing.T) {
	tests := []struct {
		record interface{}
		panics bool
	}{
		{
			record: &User{},
		},
		{
			record: NewDocument(&User{}),
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
					NewDocument(test.record)
				})
			} else {
				assert.NotPanics(t, func() {
					NewDocument(test.record)
				})
			}
		})
	}
}

func TestDocument_notPtr(t *testing.T) {
	assert.Panics(t, func() {
		NewDocument(User{}).Table()
	})
}

func TestDocument_notPtrOfStruct(t *testing.T) {
	assert.Panics(t, func() {
		i := 1
		NewDocument(&i).Table()
	})
}
