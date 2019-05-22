package sql

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
	PtrString *string
	Custom    Custom
}

// testRows is a mock version of sql.Rows which can only scan uint and strings
type testRows struct {
	mock.Mock
	columns []string
	values  []interface{}
	count   int
	total   int
}

func (r *testRows) Scan(dest ...interface{}) error {
	if len(dest) == len(r.values) {
		for i := range r.values {
			val := r.values[i]
			if v, ok := val.(uint); ok {
				val = v + uint(r.count)
			}

			if s, ok := dest[i].(sql.Scanner); ok {
				s.Scan(val)
				continue
			}

			rt := reflect.TypeOf(dest[i])
			if rt.Kind() != reflect.Ptr {
				panic("Not a pointer!")
			}

			switch dest[i].(type) {
			case **uint:
				**(dest[i].(**uint)) = val.(uint)
			case **string:
				if *(dest[i].(**string)) != nil {
					**(dest[i].(**string)) = val.(string)
				}
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
	return r.count <= r.total
}

func (r *testRows) addValue(c string, v interface{}) {
	r.columns = append(r.columns, c)
	r.values = append(r.values, v)
}

func createRows(total int) *testRows {
	rows := &testRows{total: total}
	rows.addValue("id", uint(0))
	rows.addValue("name", "name")
	rows.addValue("other_info", "other info")
	rows.addValue("real_name", "real name")
	rows.addValue("ignore", "ignore")
	rows.addValue("ptr_string", "ptr string")

	return rows
}

func TestScan(t *testing.T) {
	var (
		rows     = createRows(1)
		user     = User{}
		expected = User{
			ID:        uint(1),
			Name:      "name",
			OtherInfo: "other info",
			OtherName: "real name",
			Ignore:    "",
			PtrString: nil,
			Custom:    Custom{},
		}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	count, err := Scan(&user, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	assert.Equal(t, expected, user)
}

func TestScan_emptyRow(t *testing.T) {
	var (
		rows     = createRows(0)
		user     = User{}
		expected = User{}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	count, err := Scan(&user, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, expected, user)
}

func TestScan_columnError(t *testing.T) {
	var (
		rows = createRows(1)
		user = User{}
	)

	rows.On("Columns").Return(errors.New("error"))

	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScan_scanError(t *testing.T) {
	var (
		rows = createRows(1)
		user = User{}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(errors.New("error"))

	count, err := Scan(&user, rows)
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
}

func TestScan_panicWhenNotPointer(t *testing.T) {
	var (
		rows = createRows(1)
		user = User{}
	)

	rows.On("Columns").Return(nil)

	assert.Panics(t, func() {
		Scan(user, rows)
	})
}

func TestScan_slice(t *testing.T) {
	var (
		rows     = createRows(2)
		users    = []User{}
		expected = []User{
			{
				ID:        1,
				Name:      "name",
				OtherInfo: "other info",
				OtherName: "real name",
			},
			{
				ID:        2,
				Name:      "name",
				OtherInfo: "other info",
				OtherName: "real name",
			},
		}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	count, err := Scan(&users, rows)

	assert.Nil(t, err)
	assert.Equal(t, int64(2), count)
	assert.Equal(t, expected, users)
}

func TestScan_sliceEmptyRow(t *testing.T) {
	var (
		rows     = createRows(0)
		users    = []User{}
		expected = []User{}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	count, err := Scan(&users, rows)

	assert.Nil(t, err)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, expected, users)
}

func TestScan_sliceError(t *testing.T) {
	var (
		rows     = createRows(2)
		users    = []User{}
		expected = []User{}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(errors.New("error"))

	count, err := Scan(&users, rows)

	assert.NotNil(t, err)
	assert.Equal(t, int64(0), count)
	assert.Equal(t, expected, users)
}

func TestScan_scanner(t *testing.T) {
	var (
		rows   = createRows(1)
		custom = Custom{}
	)

	rows.On("Columns").Return(nil)
	rows.On("Scan").Return(nil)

	count, err := Scan(&custom, rows)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
}
