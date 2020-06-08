package main

import (
	"context"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
)

// QuickExample documentation.
func QuickExample(ctx context.Context, repo rel.Repository) error {
	/// [quick-example]
	book := Book{Title: "REL for Dummies"}

	// Insert a Book.
	if err := repo.Insert(ctx, &book); err != nil {
		return err
	}

	// Find a Book with id 1.
	if err := repo.Find(ctx, &book, where.Eq("id", 1)); err != nil {
		return err
	}

	// Update a Book.
	book.Title = "REL for Dummies 2nd Edition"
	if err := repo.Update(ctx, &book); err != nil {
		return err
	}

	// Delete a Book.
	if err := repo.Delete(ctx, &book); err != nil {
		return err
	}
	/// [quick-example]

	return nil
}
