package main

import (
	"context"
	"testing"

	"github.com/go-rel/rel/reltest"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestInstrumentation(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	assert.NotPanics(t, func() {
		Instrumentation(ctx, repo)

		repo.ExpectFind(where.Eq("id", 1)).Result(Book{})
		repo.MustFind(ctx, &Book{}, where.Eq("id", 1))
	})
}
