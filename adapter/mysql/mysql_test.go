package mysql

import (
	"testing"

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
