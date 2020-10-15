package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestNopAdapter_Apply(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = &nopAdapter{}
	)

	assert.Nil(t, adapter.Apply(ctx, rel.Table{}))
}

func TestNopAdapter_Exec(t *testing.T) {
	var (
		ctx     = context.TODO()
		adapter = &nopAdapter{}
	)

	rowsAffected, lastInsertedID, err := adapter.Exec(ctx, "", nil)
	assert.Zero(t, rowsAffected)
	assert.Zero(t, lastInsertedID)
	assert.Nil(t, err)
}
