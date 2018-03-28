package sqlutil

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

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
			case *uint:
				*(dest[i].(*uint)) = r.values[i].(uint)
			case *string:
				*(dest[i].(*string)) = r.values[i].(string)
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
	assert.Equal(t, User{uint(10), "string", "string", "string", "", Custom{}}, user)
}

func TestScanColumnError(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(errors.New("error"))

	user := User{}
	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScanScanError(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(errors.New("error"))

	user := User{}
	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScanPanicWhenNotPointer(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)

	user := User{}
	assert.Panics(t, func() {
		Scan(user, rows)
	})
}

func TestScanSlice(t *testing.T) {
	rows := createRows()
	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	users := []User{}
	count, err := Scan(&users, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, User{uint(10), "string", "string", "string", "", Custom{}}, users[0])
}

func TestScanScanner(t *testing.T) {
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
	intefaces := fieldPtr(rv, index, columns)

	reflect.ValueOf(intefaces[0]).Elem().SetUint(10)
	reflect.ValueOf(intefaces[1]).Elem().SetString("string")
	reflect.ValueOf(intefaces[3]).Elem().SetString("string")
	reflect.ValueOf(intefaces[4]).Elem().SetString("string")

	assert.Equal(t, User{uint(10), "string", "string", "string", "", Custom{}}, user)
}

func TestFieldIndex(t *testing.T) {
	index := fieldIndex(reflect.TypeOf(User{}))
	assert.Equal(t, map[string]int{
		"id":         0,
		"name":       1,
		"other_info": 2,
		"real_name":  3,
		"custom":     5,
	}, index)
}
