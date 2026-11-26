package rel

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type iterationErrorCursor struct {
	*sql.Rows
}

func (c *iterationErrorCursor) Fields() ([]string, error) {
	return c.Columns()
}

func (c *iterationErrorCursor) NopScanner() any {
	return &sql.RawBytes{}
}

type iterationErrorConnector struct {
	rowCount int
	err      error
}

func (c iterationErrorConnector) Connect(context.Context) (driver.Conn, error) {
	return iterationErrorConn{c}, nil
}

func (c iterationErrorConnector) Driver() driver.Driver {
	return iterationErrorDriver{}
}

type iterationErrorDriver struct{}

func (iterationErrorDriver) Open(string) (driver.Conn, error) {
	return nil, errors.New("use the connector")
}

type iterationErrorConn struct {
	iterationErrorConnector
}

func (iterationErrorConn) Prepare(string) (driver.Stmt, error) {
	return nil, errors.New("prepare is not supported")
}

func (iterationErrorConn) Close() error { return nil }

func (iterationErrorConn) Begin() (driver.Tx, error) {
	return nil, errors.New("transactions are not supported")
}

func (c iterationErrorConn) QueryContext(context.Context, string, []driver.NamedValue) (driver.Rows, error) {
	return &iterationErrorRows{remaining: c.rowCount, err: c.err}, nil
}

type iterationErrorRows struct {
	remaining int
	err       error
}

func (*iterationErrorRows) Columns() []string { return []string{"id"} }
func (*iterationErrorRows) Close() error      { return nil }

func (r *iterationErrorRows) Next(values []driver.Value) error {
	if r.remaining == 0 {
		if r.err != nil {
			return r.err
		}
		return io.EOF
	}
	r.remaining--
	values[0] = int64(10)
	return nil
}

func newIterationErrorCursor(t *testing.T, rowCount int, err error) *iterationErrorCursor {
	t.Helper()
	db := sql.OpenDB(iterationErrorConnector{rowCount: rowCount, err: err})
	t.Cleanup(func() { assert.NoError(t, db.Close()) })
	rows, queryErr := db.QueryContext(context.Background(), "SELECT id FROM users")
	require.NoError(t, queryErr)
	t.Cleanup(func() { assert.NoError(t, rows.Close()) })
	return &iterationErrorCursor{Rows: rows}
}

func TestRepository_CursorIterationError(t *testing.T) {
	for _, iterationErr := range []error{
		context.DeadlineExceeded,
		errors.New("Error 1205: Lock wait timeout exceeded"),
		errors.New("Error 1213: Deadlock found when trying to get lock"),
	} {
		for _, method := range []string{"Find", "FindAll", "FindAndCountAll"} {
			for _, rowCount := range []int{0, 1} {
				if method == "Find" && rowCount > 0 {
					continue
				}
				t.Run(fmt.Sprintf("%s/rows=%d/%s", method, rowCount, iterationErr), func(t *testing.T) {
					ctx := context.Background()
					cur := newIterationErrorCursor(t, rowCount, iterationErr)
					adapter := &testAdapter{}
					repo := New(adapter)
					query := From("users")
					expectedQuery := query
					if method == "Find" {
						expectedQuery = query.Limit(1)
					}
					adapter.On("Query", expectedQuery).Return(cur, nil).Once()
					adapter.On("Aggregate", query, "count", "*").Return(2, nil).Maybe()

					var (
						users []User
						err   error
						count int
					)
					switch method {
					case "Find":
						var user User
						err = repo.Find(ctx, &user, query)
					case "FindAll":
						err = repo.FindAll(ctx, &users, query)
					case "FindAndCountAll":
						count, err = repo.FindAndCountAll(ctx, &users, query)
					}

					t.Logf("returned error=%v, rows=%d, count=%d, rows.Err()=%v", err, len(users), count, cur.Err())
					assert.ErrorIs(t, cur.Err(), iterationErr)
					assert.ErrorIs(t, err, iterationErr)
					if method == "FindAndCountAll" {
						assert.Zero(t, count)
						adapter.AssertNotCalled(t, "Aggregate", mock.Anything, mock.Anything, mock.Anything)
					}
					adapter.AssertExpectations(t)
				})
			}
		}
	}
}

func TestScanMulti_CursorIterationError(t *testing.T) {
	for _, rowCount := range []int{0, 1} {
		t.Run(fmt.Sprintf("rows=%d", rowCount), func(t *testing.T) {
			cur := newIterationErrorCursor(t, rowCount, context.DeadlineExceeded)
			var users []User
			cols := map[any][]slice{10: {NewCollection(&users)}}

			err := scanMulti(cur, "id", reflect.TypeOf(0), cols)

			assert.ErrorIs(t, cur.Err(), context.DeadlineExceeded)
			assert.ErrorIs(t, err, context.DeadlineExceeded)
		})
	}
}

func TestRepository_CursorNormalEOF(t *testing.T) {
	for _, method := range []string{"Find", "FindAll", "FindAndCountAll"} {
		for _, rowCount := range []int{0, 1} {
			t.Run(fmt.Sprintf("%s/rows=%d", method, rowCount), func(t *testing.T) {
				ctx := context.Background()
				cur := newIterationErrorCursor(t, rowCount, nil)
				adapter := &testAdapter{}
				repo := New(adapter)
				query := From("users")
				expectedQuery := query
				if method == "Find" {
					expectedQuery = query.Limit(1)
				}
				adapter.On("Query", expectedQuery).Return(cur, nil).Once()

				var users []User
				var err error
				switch method {
				case "Find":
					var user User
					err = repo.Find(ctx, &user, query)
					if rowCount == 0 {
						assert.ErrorAs(t, err, new(NotFoundError))
						adapter.AssertExpectations(t)
						return
					}
					assert.Equal(t, 10, user.ID)
				case "FindAll":
					err = repo.FindAll(ctx, &users, query)
					assert.Len(t, users, rowCount)
				case "FindAndCountAll":
					adapter.On("Aggregate", query, "count", "*").Return(rowCount, nil).Once()
					var count int
					count, err = repo.FindAndCountAll(ctx, &users, query)
					assert.Equal(t, rowCount, count)
					assert.Len(t, users, rowCount)
				}
				assert.NoError(t, err)
				assert.NoError(t, cur.Err())
				adapter.AssertExpectations(t)
			})
		}
	}
}
