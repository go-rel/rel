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

func TestCrudInsert_for(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-for]
	repo.ExpectInsert().For(&Book{
		Title:    "Rel for dummies",
		Category: "education",
		Author: Author{
			Name: "CZ2I28 Delta",
		},
	})
	/// [insert-for]

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

func TestCrudInsertAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [insert-all]
	repo.ExpectInsertAll().ForType("[]main.Book")
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

func TestCrudUpdate(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [update]
	repo.ExpectUpdate().ForType("main.Book")
	/// [update]

	assert.Nil(t, CrudUpdate(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudUpdateAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [update-all]
	repo.ExpectUpdateAll(rel.From("books").Where(where.Lt("stock", 100)), rel.Set("discount", true))
	/// [update-all]

	assert.Nil(t, CrudUpdateAll(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudDelete(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
		book Book
	)

	/// [delete]
	repo.ExpectDelete().For(&book)
	/// [delete]

	assert.Nil(t, CrudDelete(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCrudDeleteAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [delete-all]
	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	/// [delete-all]

	assert.Nil(t, CrudDeleteAll(ctx, repo))
	repo.AssertExpectations(t)
}
