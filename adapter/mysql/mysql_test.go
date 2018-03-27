package mysql

import (
	"testing"

	"github.com/Fs02/grimoire/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestTransactionCommit(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	assert.Nil(t, adapter.Begin())
	assert.NotNil(t, adapter.tx)
	assert.Nil(t, adapter.Commit())
}

func TestTransactionRollback(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	assert.Nil(t, adapter.Begin())
	assert.NotNil(t, adapter.tx)
	assert.Nil(t, adapter.Rollback())
}

func TestTransactionCommitError(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	assert.NotNil(t, adapter.Commit())
}

func TestTransactionRollbackError(t *testing.T) {
	adapter := new(Adapter)
	adapter.Open(dsn() + "?charset=utf8&parseTime=True&loc=Local")
	defer adapter.Close()

	assert.NotNil(t, adapter.Rollback())
}

func TestError(t *testing.T) {
	adapter := new(Adapter)

	// error nil
	assert.Nil(t, adapter.Error(nil))

	// 1062 error
	rawerr := &mysql.MySQLError{Message: "duplicate", Number: 1062}
	duperr := errors.DuplicateError(rawerr.Message, "")
	assert.Equal(t, duperr, adapter.Error(rawerr))

	// other errors
	err := errors.UnexpectedError("error")
	assert.Equal(t, err, adapter.Error(err))
}
