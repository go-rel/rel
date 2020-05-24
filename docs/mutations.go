package main

import (
	"context"

	"github.com/Fs02/rel"
)

// MutationsBasicSet docs example.
func MutationsBasicSet(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [basic-set]
	err := repo.Update(ctx, &book,
		rel.Set("title", "REL for Dummies"),
		rel.Set("category", "technology"),
	)
	/// [basic-set]

	return err
}

// MutationsBasicDec docs example.
func MutationsBasicDec(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [basic-dec]
	err := repo.Update(ctx, &book, rel.DecBy("stock", 2))
	/// [basic-dec]

	return err
}

// MutationsBasicFragment docs example.
func MutationsBasicFragment(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [basic-fragment]
	err := repo.Update(ctx, &book, rel.SetFragment("title=?", "REL for dummies"))
	/// [basic-fragment]

	return err
}

// MutationsStructset docs example.
func MutationsStructset(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [structset]
	structset := rel.NewStructset(&book, false)
	err := repo.Insert(ctx, &book, structset)
	/// [structset]

	return err
}

// MutationsChangeset docs example.
func MutationsChangeset(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [changeset]
	changeset := rel.NewChangeset(&book)
	book.Price = 10
	if changeset.FieldChanged("price") {
		book.Discount = false
	}

	err := repo.Update(ctx, &book, changeset)
	/// [changeset]

	return err
}

// MutationsMap docs example.
func MutationsMap(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [map]
	data := rel.Map{
		"title":    "Rel for dummies",
		"category": "education",
		"author": rel.Map{
			"name": "CZ2I28 Delta",
		},
	}

	// Insert using map.
	err := repo.Insert(ctx, &book, data)
	/// [map]

	return err
}

// MutationsReload docs example.
func MutationsReload(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [reload]
	err := repo.Update(ctx, &book,
		rel.Set("title", "REL for Dummies"),
		rel.Reload(true),
	)
	/// [reload]

	return err
}

// MutationsCascade docs example.
func MutationsCascade(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [cascade]
	err := repo.Insert(ctx, &book, rel.Cascade(false))
	/// [cascade]

	return err
}

// MutationsDeleteCascade docs example.
func MutationsDeleteCascade(ctx context.Context, repo rel.Repository) error {
	var book Book

	/// [delete-cascade]
	err := repo.Delete(ctx, &book, rel.Cascade(true))
	/// [delete-cascade]

	return err
}
