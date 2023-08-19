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

func (i Item) PrimaryFields() []string {
	return []string{"_uuid"}
}

func (i Item) PrimaryValues() []any {
	return []any{i.UUID}
}

func TestDocument_ReflectValue(t *testing.T) {
	var (
		entity = User{}
		doc    = NewDocument(&entity)
	)

	assert.Equal(t, doc.rv, doc.ReflectValue())
}

func TestDocument_Table(t *testing.T) {
	var (
		entity = User{}
		doc    = NewDocument(&entity)
	)

	// infer table name
	assert.Equal(t, "users", doc.Table())
}

func TestDocument_Table_usingInterface(t *testing.T) {
	var (
		entity = Item{}
		doc    = NewDocument(&entity)
	)

	// infer table name
	assert.Equal(t, "_items", doc.Table())
}

func TestDocument_Primary(t *testing.T) {
	var (
		entity = User{ID: 1}
		doc    = NewDocument(&entity)
	)

	// infer primary key
	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 1, doc.PrimaryValue())

	entity.ID = 2

	// infer primary key using cache
	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 2, doc.PrimaryValue())
}

func TestDocument_PrimaryEmbedded(t *testing.T) {
	type UserWrapper struct {
		User
		ActivityScore float32
	}

	var (
		entity = UserWrapper{User: User{ID: 1}}
		doc    = NewDocument(&entity)
	)

	assert.Equal(t, "id", doc.PrimaryField())
	assert.Equal(t, 1, doc.PrimaryValue())
}

func TestDocument_Primary_usingInterface(t *testing.T) {
	var (
		entity = Item{
			UUID: "abc123",
		}
		doc = NewDocument(&entity)
	)

	// infer primary key
	assert.Equal(t, "_uuid", doc.PrimaryField())
	assert.Equal(t, "abc123", doc.PrimaryValue())
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
		doc = NewDocument(&entity)
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
		doc = NewDocument(&entity)
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
		doc = NewDocument(&entity)
	)

	assert.Panics(t, func() {
		doc.PrimaryField()
	})

	assert.Panics(t, func() {
		doc.PrimaryValue()
	})
}

func TestDocument_Primary_composite(t *testing.T) {
	var (
		userRole = UserRole{UserID: 1, RoleID: 2}
		doc      = NewDocument(&userRole)
	)

	assert.Panics(t, func() {
		doc.PrimaryField()
	})

	assert.Panics(t, func() {
		doc.PrimaryValue()
	})

	assert.Equal(t, []string{"user_id", "role_id"}, doc.PrimaryFields())
	assert.Equal(t, []any{1, 2}, doc.PrimaryValues())
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
		doc    = NewDocument(&entity)
		fields = []string{"a", "b", "c", "D"}
	)

	assert.Equal(t, fields, doc.Fields())
}

func TestDocument_Index(t *testing.T) {
	var (
		entity = struct {
			A string
			B *int
			C []byte     `db:",primary"`
			D bool       `db:"D"`
			E []*float64 `db:"-"`
		}{}
		doc   = NewDocument(&entity)
		index = map[string][]int{
			"a": {0},
			"b": {1},
			"c": {2},
			"D": {3},
		}
	)

	assert.Equal(t, index, doc.Index())
}

func TestDocument_IndexEmbedded(t *testing.T) {
	type FirstEmbedded struct {
		A int
		B int
	}
	type SecondEmbedded struct {
		D float32
	}
	var (
		entity = struct {
			FirstEmbedded `db:"first_"`
			C             string
			*SecondEmbedded
		}{}
		doc   = NewDocument(&entity)
		index = map[string][]int{
			"first_a": {0, 0},
			"first_b": {0, 1},
			"c":       {1},
			"d":       {2, 0},
		}
	)

	assert.Equal(t, index, doc.Index())
}

func TestDocument_IndexFieldEmbedded(t *testing.T) {
	type FirstEmbedded struct {
		A int
		B int
	}
	type SecondEmbedded struct {
		D float32
	}
	var (
		entity = struct {
			First  FirstEmbedded `db:"first_,embedded"`
			C      string
			Second SecondEmbedded `db:",embedded"`
			E      int            `db:"embedded"` // this field is not embedded, but only called so
		}{}
		doc   = NewDocument(&entity)
		index = map[string][]int{
			"first_a":  {0, 0},
			"first_b":  {0, 1},
			"c":        {1},
			"d":        {2, 0},
			"embedded": {3},
		}
	)

	assert.Equal(t, index, doc.Index())
}

func TestDocument_EmbeddedNameConfict(t *testing.T) {
	type Embedded struct {
		Name string
	}
	entity := struct {
		Embedded
		Name string
	}{}

	assert.Panics(t, func() {
		NewDocument(&entity)
	})
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
		doc   = NewDocument(&entity)
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
		doc    = NewDocument(&entity)
		values = map[string]any{
			"id":      1,
			"name":    "name",
			"number":  10.5,
			"address": address,
			"data":    []byte("data"),
		}
	)

	t.Run("ok", func(t *testing.T) {
		for field, evalue := range values {
			value, ok := doc.Value(field)
			assert.True(t, ok)
			assert.Equal(t, evalue, value)
		}
	})

	t.Run("field not exists", func(t *testing.T) {
		_, ok := doc.Value("not_exists")
		assert.False(t, ok)
	})
}

func TestDocument_ValueEmbedded(t *testing.T) {
	type Embedded struct {
		ID int
	}
	var (
		entity = struct {
			*Embedded
			Name string
		}{}
		doc = NewDocument(&entity)
	)

	value, ok := doc.Value("id")
	assert.True(t, ok)
	assert.Nil(t, value)

	value = doc.PrimaryValue()
	assert.Nil(t, value)

	doc.SetValue("id", 1)

	value, ok = doc.Value("id")
	assert.True(t, ok)
	assert.Equal(t, 1, value)

	value = doc.PrimaryValue()
	assert.Equal(t, 1, value)

	value, ok = doc.Value("name")
	assert.True(t, ok)
	assert.Equal(t, "", value)
}

func TestDocument_SetValue(t *testing.T) {
	var (
		entity struct {
			ID      int
			Name    string
			Skip    bool `db:"-"`
			Number  float64
			Address *string
			Data    []byte
		}
		doc = NewDocument(&entity)
	)

	t.Run("ok", func(t *testing.T) {
		assert.True(t, doc.SetValue("id", 1))
		assert.True(t, doc.SetValue("name", "name"))
		assert.True(t, doc.SetValue("number", 10.5))
		assert.True(t, doc.SetValue("data", []byte("data")))
		assert.True(t, doc.SetValue("address", "address"))

		assert.Equal(t, 1, entity.ID)
		assert.Equal(t, "name", entity.Name)
		assert.Equal(t, false, entity.Skip)
		assert.Equal(t, 10.5, entity.Number)
		assert.Equal(t, "address", *entity.Address)
		assert.Equal(t, []byte("data"), entity.Data)
	})

	t.Run("zero", func(t *testing.T) {
		assert.True(t, doc.SetValue("id", nil))
		assert.True(t, doc.SetValue("name", nil))
		assert.True(t, doc.SetValue("number", nil))
		assert.True(t, doc.SetValue("data", nil))
		assert.True(t, doc.SetValue("address", nil))

		assert.Equal(t, 0, entity.ID)
		assert.Equal(t, "", entity.Name)
		assert.Equal(t, float64(0), entity.Number)
		assert.Equal(t, (*string)(nil), entity.Address)
		assert.Equal(t, []byte(nil), entity.Data)
	})

	t.Run("convert", func(t *testing.T) {
		assert.True(t, doc.SetValue("id", uint(2)))
		assert.True(t, doc.SetValue("number", 10))
		assert.Equal(t, 2, entity.ID)
		assert.Equal(t, float64(10), entity.Number)
	})

	t.Run("reflect", func(t *testing.T) {
		assert.True(t, doc.SetValue("id", reflect.ValueOf(21)))
		assert.True(t, doc.SetValue("address", reflect.ValueOf("continassa")))
		assert.Equal(t, 21, entity.ID)
		assert.Equal(t, "continassa", *entity.Address)
	})

	t.Run("field not exists", func(t *testing.T) {
		assert.False(t, doc.SetValue("id", "a"))
		assert.False(t, doc.SetValue("skip", true))
		assert.False(t, doc.SetValue("address", []byte("a")))
	})
}

func TestDocument_SetValueEmbedded(t *testing.T) {
	type Embedded struct {
		ID   int
		Name string
	}
	var (
		entity struct {
			Embedded
			Number float64
		}
		doc = NewDocument(&entity)
	)

	assert.True(t, doc.SetValue("id", 1))
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
		doc      = NewDocument(&entity)
		fields   = []string{"name", "id", "skip", "data", "number", "address", "not_exist"}
		scanners = []any{
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

func TestDocument_Scanners_withAssoc(t *testing.T) {
	var (
		entity = Transaction{
			ID:      1,
			BuyerID: 2,
			Status:  "SENT",
			Buyer: User{
				ID:   2,
				Name: "user",
				WorkAddress: &Address{
					Street: "Takeshita-dori",
				},
			},
		}
		doc      = NewDocument(&entity)
		fields   = []string{"id", "user_id", "buyer.id", "buyer.name", "buyer.work_address.street", "status", "invalid_assoc.id"}
		scanners = []any{
			Nullable(&entity.ID),
			Nullable(&entity.BuyerID),
			Nullable(&entity.Buyer.ID),
			Nullable(&entity.Buyer.Name),
			Nullable(&entity.Buyer.WorkAddress.Street),
			Nullable(&entity.Status),
			&sql.RawBytes{},
		}
	)

	assert.Equal(t, scanners, doc.Scanners(fields))
}

func TestDocument_Scanners_withUnitializedAssoc(t *testing.T) {
	var (
		entity   = Transaction{}
		doc      = NewDocument(&entity)
		fields   = []string{"id", "user_id", "buyer.id", "buyer.name", "status", "buyer.work_address.street"}
		result   = doc.Scanners(fields)
		expected = []any{
			Nullable(&entity.ID),
			Nullable(&entity.BuyerID),
			Nullable(&entity.Buyer.ID),
			Nullable(&entity.Buyer.Name),
			Nullable(&entity.Status),
			Nullable(&entity.Buyer.WorkAddress.Street),
		}
	)

	assert.Equal(t, expected, result)
}

func TestDocument_ScannersInitPointers(t *testing.T) {
	type Embedded1 struct {
		ID int
	}
	type Embedded2 struct {
		Embedded1
	}
	type Embedded3 struct {
		*Embedded2
	}
	var (
		entity = struct {
			*Embedded3
		}{}
		doc = NewDocument(&entity)
		_   = doc.Scanners([]string{"id"})
	)
	assert.NotNil(t, entity.Embedded2)
	assert.NotNil(t, entity.Embedded2.Embedded1)
}

func TestDocument_Slice(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			doc = NewDocument(&Item{})
		)

		doc.Reset()
		assert.Equal(t, 1, doc.Len())
		assert.Equal(t, doc, doc.Get(0))
	})
}

func TestDocument_Association(t *testing.T) {
	tests := []struct {
		name      string
		entity    any
		belongsTo []string
		hasOne    []string
		hasMany   []string
		preload   []string
	}{
		{
			name:    "User",
			entity:  &User{},
			hasOne:  []string{"address", "work_address"},
			hasMany: []string{"transactions", "user_roles", "emails", "roles", "follows", "followeds", "followings", "followers"},
		},
		{
			name:    "User Cached",
			entity:  &User{},
			hasOne:  []string{"address", "work_address"},
			hasMany: []string{"transactions", "user_roles", "emails", "roles", "follows", "followeds", "followings", "followers"},
		},
		{
			name:      "Transaction",
			entity:    &Transaction{},
			belongsTo: []string{"buyer", "address"},
			hasMany:   []string{"histories"},
			preload:   []string{"buyer"},
		},
		{
			name:      "Address",
			entity:    &Address{},
			belongsTo: []string{"user"},
		},
		{
			name:   "Item",
			entity: &Item{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				doc = NewDocument(test.entity)
			)

			assert.Equal(t, test.belongsTo, doc.BelongsTo())
			assert.Equal(t, test.hasOne, doc.HasOne())
			assert.Equal(t, test.hasMany, doc.HasMany())
			assert.Equal(t, test.preload, doc.Preload())
		})
	}
}

func TestDocument_AssociationEmbedded(t *testing.T) {
	var (
		testHasOne  = []string{"user_address", "user_work_address"}
		testHasMany = []string{"user_transactions", "user_user_roles", "user_emails", "user_roles", "user_follows", "user_followeds", "user_followings", "user_followers"}
		entity      = struct {
			User  `db:"user_"`
			Score float32
		}{}
		doc = NewDocument(&entity)
	)

	assert.Equal(t, testHasOne, doc.HasOne())
	assert.Equal(t, testHasMany, doc.HasMany())
}

func TestDocument_Association_notFOund(t *testing.T) {
	var (
		doc = NewDocument(&Item{})
	)

	assert.Panics(t, func() {
		doc.Association("empty")
	})
}

func TestDocument(t *testing.T) {
	tests := []struct {
		entity any
		panics bool
	}{
		{
			entity: &User{},
		},
		{
			entity: NewDocument(&User{}),
		},
		{
			entity: reflect.ValueOf(&User{}),
		},
		{
			entity: reflect.ValueOf(User{}),
			panics: true,
		},
		{
			entity: reflect.ValueOf(&[]User{}),
			panics: true,
		},
		{
			entity: reflect.TypeOf(&User{}),
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
					NewDocument(test.entity)
				})
			} else {
				assert.NotPanics(t, func() {
					NewDocument(test.entity)
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
