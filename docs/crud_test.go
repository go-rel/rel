package main

import (
	"context"
	"errors"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/stretchr/testify/assert"
)

func TestInsert(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert]
	repo.ExpectInsert()
	/// [insert]

	assert.Nil(t, Insert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsert_forType(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-for-type]
	repo.ExpectInsert().ForType("main.Book")
	/// [insert-for-type]

	assert.Nil(t, Insert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsert_specific(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-specific]
	repo.ExpectInsert().For(&Book{
		Title:    "Rel for dummies",
		Category: "education",
	})
	/// [insert-specific]

	assert.Nil(t, Insert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsert_error(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-error]
	repo.ExpectInsert().ForType("main.Book").Error(errors.New("oops"))
	/// [insert-error]

	assert.Equal(t, errors.New("oops"), Insert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsertMap(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-map]
	repo.ExpectInsert(rel.Map{
		"title":    "Rel for dummies",
		"category": "education",
	}).ForType("main.Book")
	/// [insert-map]

	assert.Nil(t, InsertMap(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsertSet(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-set]
	repo.ExpectInsert(
		rel.Set("title", "Rel for dummies"),
		rel.Set("category", "education"),
	).ForType("main.Book")
	/// [insert-set]

	assert.Nil(t, InsertSet(ctx, repo))
	repo.AssertExpectations(t)
}

func TestInsertAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-all]
	repo.ExpectInsertAll()
	/// [insert-all]

	assert.Nil(t, InsertAll(ctx, repo))
	repo.AssertExpectations(t)
}
