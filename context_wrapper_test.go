package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestContextWrapper(t *testing.T) {
	var (
		cw      contextWrapper
		adapter = &testAdapter{}
		ctx     = context.TODO()
	)

	t.Run("fetch empty context", func(t *testing.T) {
		cw = fetchContext(ctx, adapter)
		assert.Equal(t, ctx, cw.ctx)
		assert.Equal(t, adapter, cw.adapter)
	})

	t.Run("wrap context", func(t *testing.T) {
		adapter = &testAdapter{result: 1}
		cw = wrapContext(ctx, adapter)
		ctx = cw.ctx

		assert.Equal(t, ctx, cw.ctx)
		assert.Equal(t, adapter, cw.adapter)
	})

	t.Run("fetch wrapped context", func(t *testing.T) {
		cw = fetchContext(ctx, adapter)
		assert.Equal(t, ctx, cw.ctx)
		assert.Equal(t, adapter, cw.adapter)
	})
}
