package main

import (
	"context"
	"errors"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestCrudInsert(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert]
	repo.ExpectInsert()
	/// [insert]

	assert.Nil(t, CrudInsert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudInsert_forType(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-for-type]
	repo.ExpectInsert().ForType("main.Book")
	/// [insert-for-type]

	assert.Nil(t, CrudInsert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudInsert_specific(t *testing.T) {
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

	assert.Nil(t, CrudInsert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudInsert_error(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-error]
	repo.ExpectInsert().ForType("main.Book").Error(errors.New("oops"))
	/// [insert-error]

	assert.Equal(t, errors.New("oops"), CrudInsert(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudInsertMap(t *testing.T) {
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

	assert.Nil(t, CrudInsertMap(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudInsertSet(t *testing.T) {
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

	assert.Nil(t, CrudInsertSet(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudInsertAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-all]
	repo.ExpectInsertAll()
	/// [insert-all]

	assert.Nil(t, CrudInsertAll(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudFind(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find]
	book := Book{
		Title:    "Rel for dummies",
		Category: "education",
	}

	repo.ExpectFind(rel.Eq("id", 1)).Result(book)
	/// [find]

	assert.Nil(t, CrudFind(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudFindAlias_error(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find-alias-error]
	repo.ExpectFind(where.Eq("id", 1)).NotFound()
	/// [find-alias-error]

	assert.Equal(t, rel.ErrNotFound, CrudFindAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudFindAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find-all]
	books := []Book{
		{
			Title:    "Rel for dummies",
			Category: "education",
		},
	}

	repo.ExpectFindAll(
		where.Like("title", "%dummies%").AndEq("category", "education"),
		rel.Limit(10),
	).Result(books)
	/// [find-all]

	assert.Nil(t, CrudFindAll(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudFindAllChained(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find-all-chained]
	books := []Book{
		{
			Title:    "Rel for dummies",
			Category: "education",
		},
	}

	query := rel.Select("title", "category").Where(where.Eq("category", "education")).SortAsc("title")
	repo.ExpectFindAll(query).Result(books)
	/// [find-all-chained]

	assert.Nil(t, CrudFindAllChained(ctx, repo))
	repo.AssertExpectations(t)
}
