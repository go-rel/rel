package main

import (
	"context"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
)

// CrudInsert docs example.
func CrudInsert(ctx context.Context, repo rel.Repository) error {
	/// [insert]
	book := Book{
		Title:    "Rel for dummies",
		Category: "education",
		Author: Author{
			Name: "CZ2I28 Delta",
		},
	}

	// Insert directly using struct.
	err := repo.Insert(ctx, &book)
	/// [insert]

	return err
}

// CrudInsertAll docs example.
func CrudInsertAll(ctx context.Context, repo rel.Repository) error {
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

// CrudFind docs example.
func CrudFind(ctx context.Context, repo rel.Repository) error {
	/// [find]
	var book Book
	err := repo.Find(ctx, &book, rel.Eq("id", 1))
	/// [find]

	return err
}

// CrudFindAlias docs example.
func CrudFindAlias(ctx context.Context, repo rel.Repository) error {
	/// [find-alias]
	var book Book
	err := repo.Find(ctx, &book, where.Eq("id", 1))
	/// [find-alias]

	return err
}

// CrudFindAll docs example.
func CrudFindAll(ctx context.Context, repo rel.Repository) error {
	/// [find-all]
	var books []Book
	err := repo.FindAll(ctx, &books,
		where.Like("title", "%dummies%").AndEq("category", "education"),
		rel.Limit(10),
	)
	/// [find-all]

	return err
}

// CrudFindAllChained docs example.
func CrudFindAllChained(ctx context.Context, repo rel.Repository) error {
	/// [find-all-chained]
	var books []Book
	query := rel.Select("title", "category").Where(where.Eq("category", "education")).SortAsc("title")
	err := repo.FindAll(ctx, &books, query)
	/// [find-all-chained]

	return err
}

// CrudUpdate docs example.
func CrudUpdate(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [update]
	book.Title = "REL for dummies"
	err := repo.Update(ctx, &book)
	/// [update]

	return err
}

// CrudUpdateAll docs example.
func CrudUpdateAll(ctx context.Context, repo rel.Repository) error {
	/// [update-all]
	err := repo.UpdateAll(ctx, rel.From("books").Where(where.Lt("stock", 100)), rel.Set("discount", true))
	/// [update-all]

	return err
}

// CrudDelete docs example.
func CrudDelete(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [delete]
	err := repo.Delete(ctx, &book)
	/// [delete]

	return err
}

// CrudDeleteAll docs example.
func CrudDeleteAll(ctx context.Context, repo rel.Repository) error {
	/// [delete-all]
	err := repo.DeleteAll(ctx, rel.From("books").Where(where.Eq("id", 1)))
	/// [delete-all]

	return err
}
