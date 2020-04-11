package main

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestIteration(t *testing.T) {
	var (
		ctx   = context.TODO()
		repo  = reltest.New()
		users = make([]User, 5)
	)

	/// [batch-iteration]
	repo.ExpectIterate(rel.From("users"), rel.BatchSize(500)).Result(users)
	/// [batch-iteration]

	assert.Nil(t, Iteration(ctx, repo))
	repo.AssertExpectations(t)
}

func TestIteration_error(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [batch-iteration-connection-error]
	repo.ExpectIterate(rel.From("users"), rel.BatchSize(500)).ConnectionClosed()
	/// [batch-iteration-connection-error]

	assert.Equal(t, reltest.ErrConnectionClosed, Iteration(ctx, repo))
	repo.AssertExpectations(t)
}
