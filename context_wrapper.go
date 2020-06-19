package rel

import (
	"context"
)

type contextKey int8

type contextData struct {
	adapter Adapter
}

type contextWrapper struct {
	ctx     context.Context
	adapter Adapter
}

var ctxKey contextKey

// fetchContext and use adapter passed by context if exists.
// it stores contextData values to struct for fast repeated access.
func fetchContext(ctx context.Context, adapter Adapter) contextWrapper {
	if adp, ok := ctx.Value(ctxKey).(Adapter); ok {
		adapter = adp
	}

	return contextWrapper{
		ctx:     ctx,
		adapter: adapter,
	}
}

// wrapContext wraps adapter inside context.
func wrapContext(ctx context.Context, adapter Adapter) contextWrapper {
	return contextWrapper{
		ctx:     context.WithValue(ctx, ctxKey, adapter),
		adapter: adapter,
	}
}
