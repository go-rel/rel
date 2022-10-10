package rel

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

func (it *Items) PrimaryFields() []string {
	return []string{"_uuid"}
}

func (it *Items) PrimaryValues() []any {
	var (
		ids = make([]any, len(*it))
	)

	for i := range *it {
		ids[i] = (*it)[i].UUID
	}

	return []any{ids}
}

func TestCollection_ReflectValue(t *testing.T) {
	var (
		record = []User{}
		doc    = NewCollection(&record)
	)

	assert.Equal(t, doc.rv, doc.ReflectValue())
}

func TestCollection_Table(t *testing.T) {
	var (
		records = []User{}
		col     = NewCollection(&records)
	)

	// infer table name
	assert.Equal(t, "users", col.Table())
}

func TestCollection_Table_usingInterface(t *testing.T) {
	var (
		records = Items{}
		col     = NewCollection(&records)
	)

	// infer table name
	assert.Equal(t, "_items", col.Table())
}

func TestCollection_Table_usingElemInterface(t *testing.T) {
	var (
		records = []Item{}
		col     = NewCollection(&records)
	)

	// infer table name
	assert.Equal(t, "_items", col.Table())
}

func TestCollection_Primary(t *testing.T) {
	var (
		records = []User{
			{ID: 1},
			{ID: 2},
		}
		rt  = reflect.TypeOf(records).Elem()
		col = NewCollection(&records)
	)

	// infer primary key
	assert.Equal(t, "id", col.PrimaryField())
	assert.Equal(t, []any{1, 2}, col.PrimaryValue())

	// cached
	_, cached := primariesCache.Load(rt)
	assert.True(t, cached)

	records[1].ID = 4

	// infer primary key using cache
	assert.Equal(t, "id", col.PrimaryField())
	assert.Equal(t, []any{1, 4}, col.PrimaryValue())

	primariesCache.Delete(rt)
}

func TestCollection_Primary_usingInterface(t *testing.T) {
	var (
		records = Items{
			{UUID: "abc123"},
			{UUID: "def456"},
		}
		col = NewCollection(&records)
	)

	// infer primary key
	assert.Equal(t, "_uuid", col.PrimaryField())
	assert.Equal(t, []any{"abc123", "def456"}, col.PrimaryValue())
}

func TestCollection_Primary_usingElemInterface(t *testing.T) {
	var (
		records = []Item{
			{UUID: "abc123"},
			{UUID: "def456"},
		}
		rt  = reflect.TypeOf(records).Elem()
		col = NewCollection(&records)
	)

	// infer primary key
	assert.Equal(t, "_uuid", col.PrimaryField())
	assert.Equal(t, []any{"abc123", "def456"}, col.PrimaryValue())

	primariesCache.Delete(rt)
}

func TestCollection_Primary_usingElemInterface_ptrElem(t *testing.T) {
	var (
		records = []*Item{
			{UUID: "abc123"},
			{UUID: "def456"},
			nil,
		}
		rt  = reflect.TypeOf(records).Elem()
		col = NewCollection(&records)
	)

	// infer primary key
	assert.Equal(t, "_uuid", col.PrimaryField())
	assert.Equal(t, []any{"abc123", "def456"}, col.PrimaryValue())

	primariesCache.Delete(rt)
}

func TestCollection_Primary_usingTag(t *testing.T) {
	var (
		records = []struct {
			ID         uint
			ExternalID int `db:",primary"`
			Name       string
		}{
			{ExternalID: 1},
			{ExternalID: 2},
		}
		col = NewCollection(&records)
	)

	// infer primary key
	assert.Equal(t, "external_id", col.PrimaryField())
	assert.Equal(t, []any{1, 2}, col.PrimaryValue())
}

func TestCollection_Primary_composite(t *testing.T) {
	var (
		userRole = []UserRole{
			{UserID: 1, RoleID: 2},
			{UserID: 3, RoleID: 4},
			{UserID: 5, RoleID: 6},
		}
		col = NewCollection(&userRole)
	)

	assert.Panics(t, func() {
		col.PrimaryField()
	})

	assert.Panics(t, func() {
		col.PrimaryValue()
	})

	assert.Equal(t, []string{"user_id", "role_id"}, col.PrimaryFields())
	assert.Equal(t, []any{
		[]any{1, 3, 5},
		[]any{2, 4, 6},
	}, col.PrimaryValues())
}

func TestCollection_Primary_notFound(t *testing.T) {
	var (
		records = []struct {
			ExternalID int
			Name       string
		}{}
		col = NewCollection(&records)
	)

	assert.Panics(t, func() {
		col.PrimaryField()
	})
}

func TestCollection_Truncate(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			users = []User{}
			col   = NewCollection(&users)
		)

		col.Add()
		assert.Equal(t, 1, col.Len())

		col.Truncate(0, 0)
		assert.Equal(t, 0, col.Len())
	})
}

func TestCollection_Swap(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			users = []User{
				{ID: 1},
				{ID: 2},
			}
			col = NewCollection(&users)
		)

		col.Swap(0, 1)

		assert.Equal(t, 2, users[0].ID)
		assert.Equal(t, 1, users[1].ID)
	})
}

func TestCollection_Slice(t *testing.T) {
	assert.NotPanics(t, func() {
		var (
			users = []User{}
			col   = NewCollection(&users)
		)

		assert.Equal(t, 0, col.Len())

		doc := col.Add()
		assert.Len(t, users, 1)
		assert.Equal(t, 1, col.Len())
		assert.Equal(t, NewDocument(&users[0]), doc)
		assert.Equal(t, NewDocument(&users[0]), col.Get(0))

		col.Reset()
		assert.Len(t, users, 0)
		assert.Equal(t, 0, col.Len())
		assert.Equal(t, &[]User{}, col.v)
	})
}

func TestCollection(t *testing.T) {
	tests := []struct {
		record any
		panics bool
	}{
		{
			record: &[]User{},
		},
		{
			record: NewCollection(&[]User{}),
		},
		{
			record: reflect.ValueOf(&[]User{}),
		},
		{
			record: reflect.ValueOf([]User{}),
			panics: true,
		},
		{
			record: reflect.ValueOf(&User{}),
			panics: true,
		},
		{
			record: reflect.TypeOf(&[]User{}),
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
					NewCollection(test.record)
				})
			} else {
				assert.NotPanics(t, func() {
					NewCollection(test.record)
				})
			}
		})
	}
}

func TestCollection_notPtr(t *testing.T) {
	assert.Panics(t, func() {
		NewCollection([]User{}).Table()
	})
}

func TestCollection_notPtrOfSlice(t *testing.T) {
	assert.Panics(t, func() {
		NewCollection(&User{}).Table()
	})
}
