package main

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestTransactions(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [transactions]
	repo.ExpectTransaction(func(repo *reltest.Repository) {
		repo.ExpectUpdate(rel.Dec("stock")).ForType("main.Book")
		repo.ExpectUpdate(rel.Set("status", "paid")).ForType("main.Transaction")
	})
	/// [transactions]

	assert.Nil(t, Transactions(ctx, repo))
	repo.AssertExpectations(t)
}
