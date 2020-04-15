package main

import (
	"context"

	"github.com/Fs02/rel"
)

// Transactions docs example.
func Transactions(ctx context.Context, repo rel.Repository) error {
	var book Book
	var transaction Transaction
	/// [transactions]
	err := repo.Transaction(ctx, func(repo rel.Repository) error {
		repo.Update(ctx, &book, rel.Dec("stock"))
		return repo.Update(ctx, &transaction, rel.Set("status", "paid"))
	})
	/// [transactions]

	return err
}
