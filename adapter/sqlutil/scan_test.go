package sqlutil

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type User struct {
	ID        uint
	Name      string
	OtherInfo string
	OtherName string `db:"real_name"`
}

// testRows is a mock version of sql.Rows which can only scan uint and strings
type testRows struct {
	columns []string
	values  []interface{}
	count   int
}

func (r *testRows) Scan(dest ...interface{}) error {
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
	return nil
}

func (r *testRows) Columns() ([]string, error) {
	return r.columns, nil
}

func (r *testRows) Next() bool {
	r.count++
	return r.count == 1
}

func (r *testRows) addValue(c string, v interface{}) {
	r.columns = append(r.columns, c)
	r.values = append(r.values, v)
}

func createRows() testRows {
	rows := testRows{}
	rows.addValue("id", uint(10))
	rows.addValue("name", "string")
	rows.addValue("other_info", "string")
	rows.addValue("real_name", "string")

	return rows
}

func TestScan(t *testing.T) {
	rows := createRows()

	user := User{}
	count, err := Scan(&user, &rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, User{uint(10), "string", "string", "string"}, user)
}

func TestScanSlice(t *testing.T) {
	rows := createRows()

	users := []User{}
	count, err := Scan(&users, &rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, 1, len(users))
	assert.Equal(t, User{uint(10), "string", "string", "string"}, users[0])
}

func TestScanScanner(t *testing.T) {
	t.Skip("PENDING")
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

	assert.Equal(t, User{uint(10), "string", "string", "string"}, user)
}

func TestFieldIndex(t *testing.T) {
	index := fieldIndex(reflect.TypeOf(User{}))
	assert.Equal(t, map[string]int{
		"id":         0,
		"name":       1,
		"other_info": 2,
		"real_name":  3,
	}, index)
}
