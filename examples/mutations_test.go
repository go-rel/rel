package main

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestMutationsBasicSet(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [basic-set]
	repo.ExpectUpdate(
		rel.Set("title", "REL for Dummies"),
		rel.Set("category", "technology"),
	).For(&book)
	/// [basic-set]

	assert.Nil(t, MutationsBasicSet(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsBasicDec(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [basic-dec]
	repo.ExpectUpdate(rel.DecBy("stock", 2)).For(&book)
	/// [basic-dec]

	assert.Nil(t, MutationsBasicDec(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsBasicFragment(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [basic-fragment]
	repo.ExpectUpdate(rel.SetFragment("title=?", "REL for dummies")).For(&book)
	/// [basic-fragment]

	assert.Nil(t, MutationsBasicFragment(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsStructset(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [structset]
	repo.ExpectInsert().For(&book)
	/// [structset]

	assert.Nil(t, MutationsStructset(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsChangeset(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [changeset]
	repo.ExpectUpdate().ForType("main.Book")
	/// [changeset]

	assert.Nil(t, MutationsChangeset(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsMap(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [map]
	data := rel.Map{
		"title":    "Rel for dummies",
		"category": "education",
		"author": rel.Map{
			"name": "CZ2I28 Delta",
		},
	}

	repo.ExpectInsert(data).ForType("main.Book")
	/// [map]

	assert.Nil(t, MutationsMap(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsReload(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [reload]
	repo.ExpectUpdate(
		rel.Set("title", "REL for Dummies"),
		rel.Reload(true),
	).For(&book)
	/// [reload]

	assert.Nil(t, MutationsReload(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsCascade(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [cascade]
	repo.ExpectInsert(rel.Cascade(false)).For(&book)
	/// [cascade]

	assert.Nil(t, MutationsCascade(ctx, repo))
	repo.AssertExpectations(t)
}

func TestMutationsDeleteCascade(t *testing.T) {
	var (
		book Book
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [delete-cascade]
	repo.ExpectDelete(rel.Cascade(true)).For(&book)
	/// [delete-cascade]

	assert.Nil(t, MutationsDeleteCascade(ctx, repo))
	repo.AssertExpectations(t)
}
