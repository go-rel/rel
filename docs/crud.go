package main

import (
	"context"

	"github.com/Fs02/rel"
)

// Insert docs example.
func Insert(ctx context.Context, repo rel.Repository) error {
	/// [insert]
	book := Book{
		Title:    "Rel for dummies",
		Category: "education",
	}

	// Insert directly using struct.
	err := repo.Insert(ctx, &book)
	/// [insert]

	return err
}

// InsertMap docs example.
func InsertMap(ctx context.Context, repo rel.Repository) error {
	/// [insert-map]
	var book Book
	data := rel.Map{
		"title":    "Rel for dummies",
		"category": "education",
	}

	// Insert using map.
	err := repo.Insert(ctx, &book, data)
	/// [insert-map]

	return err
}

// InsertSet docs example.
func InsertSet(ctx context.Context, repo rel.Repository) error {
	/// [insert-set]
	var book Book
	err := repo.Insert(ctx, &book,
		rel.Set("title", "Rel for dummies"),
		rel.Set("category", "education"),
	)
	/// [insert-set]

	return err
}

// InsertAll docs example.
func InsertAll(ctx context.Context, repo rel.Repository) error {
	/// [insert-all]
	books := []Book{
		{
			Title:    "Golang for dummies",
			Category: "education",
		},
		{
			Title:    "Rel for dummies",
			Category: "education",
		},
	}

	err := repo.InsertAll(ctx, &books)
	/// [insert-all]

	return err
}