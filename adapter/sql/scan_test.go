package sql

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type Custom struct{}

var _ sql.Scanner = (*Custom)(nil)

func (c *Custom) Scan(interface{}) error {
	return nil
}

type User struct {
	ID        uint
	Name      string
	OtherInfo string
	OtherName string `db:"real_name"`
	Ignore    string `db:"-"`
	PtrString *string
	Custom    Custom
}

// testRows is a mock version of sql.Rows which can only scan uint and strings
type testRows struct {
	mock.Mock
	columns []string
	values  []interface{}
	count   int
}

func (r *testRows) Scan(dest ...interface{}) error {
	if len(dest) == len(r.values) {
		for i := range r.values {
			v := reflect.ValueOf(dest[i])
			if v.Kind() != reflect.Ptr {
				panic("Not a pointer!")
			}

			switch dest[i].(type) {
			case **uint:
				**(dest[i].(**uint)) = r.values[i].(uint)
			case **string:
				if *(dest[i].(**string)) != nil {
					**(dest[i].(**string)) = r.values[i].(string)
				}
			default:
				// Do nothing.
			}
		}
	}

	args := r.Called()
	return args.Error(0)
}

func (r *testRows) Columns() ([]string, error) {
	args := r.Called()
	return r.columns, args.Error(0)
}

func (r *testRows) Next() bool {
	r.count++
	return r.count == 1
}

func (r *testRows) addValue(c string, v interface{}) {
	r.columns = append(r.columns, c)
	r.values = append(r.values, v)
}

func createRows() *testRows {
	rows := new(testRows)
	rows.addValue("id", uint(10))
	rows.addValue("name", "string")
	rows.addValue("other_info", "string")
	rows.addValue("real_name", "string")
	rows.addValue("ignore", "string")
	rows.addValue("ptr_string", "string")

	return rows
}

func TestScan(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	user := User{}
	count, err := Scan(&user, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, User{uint(10), "string", "string", "string", "", nil, Custom{}}, user)
}

func TestScan_columnError(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(errors.New("error"))

	user := User{}
	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScan_scanError(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(errors.New("error"))

	user := User{}
	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScan_panicWhenNotPointer(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)

	user := User{}
	assert.Panics(t, func() {
		Scan(user, rows)
	})
}

func TestScan_slice(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	users := []User{}
	count, err := Scan(&users, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, User{uint(10), "string", "string", "string", "", nil, Custom{}}, users[0])
}

func TestScan_scanner(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	custom := Custom{}
	count, err := Scan(&custom, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}

func TestFieldPtr(t *testing.T) {
	user := User{ID: 5}
	rv := reflect.ValueOf(&user).Elem()
	index := fieldIndex(rv.Type())
	columns := []string{"id", "name", "fake1", "other_info", "real_name", "fake2"}
	ptr, reset := fieldPtr(rv, index, columns)

	reflect.ValueOf(ptr[0]).Elem().Elem().SetUint(10)
	reflect.ValueOf(ptr[1]).Elem().Elem().SetString("string")
	reflect.ValueOf(ptr[3]).Elem().Elem().SetString("string")
	reflect.ValueOf(ptr[4]).Elem().Elem().SetString("string")

	assert.Equal(t, User{uint(10), "string", "string", "string", "", nil, Custom{}}, user)
	assert.Equal(t, 4, len(reset))
}

func TestFieldIndex(t *testing.T) {
	var obj struct {
		ID                 int
		Name               string
		Other              bool
		SkippedStructSlice []User
		OtherID            int64
		SkippedIntSlice    []int
		Score              float64
		SkippedStruct      User
		Custom             Custom
		SkippedStringSlice []string
		CreatedAt          time.Time
		DeletedAt          *time.Time
		CustomPtr          *Custom
		SkippedStructPtr   *User
		SkippedIntSlicePtr *[]int
		ArrayInt           [1]int
		ArrayIntPtr        *[2]int
	}

	index := fieldIndex(reflect.TypeOf(obj))
	assert.Equal(t, map[string]int{
		"id":         0,
		"name":       1,
		"other":      2,
		"other_id":   4,
		"score":      6,
		"custom":     8,
		"created_at": 10,
		"deleted_at": 11,
		"custom_ptr": 12,
	}, index)
}
