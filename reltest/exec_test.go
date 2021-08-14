package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExec(t *testing.T) {
	var (
		repo      = New()
		statement = "UPDATE users SET something = ? WHERE something2 = ?;"
		args      = []interface{}{1, "2"}
	)

	repo.ExpectExec(statement, args).Result(1, 2)
	lastInsertedId, rowsAffected, err := repo.Exec(context.TODO(), statement, args...)
	assert.Nil(t, err)
	assert.Equal(t, 1, lastInsertedId)
	assert.Equal(t, 2, rowsAffected)
	repo.AssertExpectations(t)

	repo.ExpectExec(statement, args).Result(1, 2)
	assert.NotPanics(t, func() {
		lastInsertedId, rowsAffected := repo.MustExec(context.TODO(), statement, args...)
		assert.Equal(t, 1, lastInsertedId)
		assert.Equal(t, 2, rowsAffected)
	})
	repo.AssertExpectations(t)
}

func TestExec_error(t *testing.T) {
	var (
		repo      = New()
		statement = "UPDATE users SET something = ? WHERE something2 = ?;"
		args      = []interface{}{1, "2"}
	)

	repo.ExpectExec(statement, args).ConnectionClosed()
	lastInsertedId, rowsAffected, err := repo.Exec(context.TODO(), statement, args...)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, lastInsertedId)
	assert.Equal(t, 0, rowsAffected)
	repo.AssertExpectations(t)

	repo.ExpectExec(statement, args).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustExec(context.TODO(), statement, args...)
	})
	repo.AssertExpectations(t)
}

func TestExec_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectExec("UPDATE users SET something = ? WHERE something2 = ?;", 1, 2)

	assert.Panics(t, func() {
		repo.Exec(context.TODO(), "UPDATE users SET something = ? WHERE something2 = ?;", 1, 2)
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tExec(ctx, \"UPDATE users SET something = ? WHERE something2 = ?;\", 1, 2)", nt.lastLog)
}

func TestExec_String(t *testing.T) {
	var (
		mockExec = MockExec{assert: &Assert{}, argStatement: "UPDATE users SET something = ? WHERE something2 = ?;", argArgs: []interface{}{1, 2}}
	)

	assert.Equal(t, "Exec(ctx, \"UPDATE users SET something = ? WHERE something2 = ?;\", 1, 2)", mockExec.String())
	assert.Equal(t, "ExpectString(\"UPDATE users SET something = ? WHERE something2 = ?;\", 1, 2)", mockExec.ExpectString())
}
