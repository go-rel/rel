package main

import (
	"context"

	"github.com/Fs02/rel"
)

// Process is blank function used for transaction example.
func Process(ctx context.Context, trx Transaction) {
}

// Transactions docs example.
func Transactions(ctx context.Context, repo rel.Repository) error {
	var book Book
	var transaction Transaction
	/// [transactions]
	err := repo.Transaction(ctx, func(ctx context.Context) error {
		repo.Update(ctx, &book, rel.Dec("stock"))

		// Any database calls inside other function will be using the same transaction as long as it share the same context.
		Process(ctx, transaction)

		return repo.Update(ctx, &transaction, rel.Set("status", "paid"))
	})
	/// [transactions]

	return err
}
