package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.Nil(t, repo.FindAll(context.TODO(), &result, where.Like("title", "%dummies%")))
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(context.TODO(), &result, where.Like("title", "%dummies%"))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_any(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Eq("title", Any)).Result(books)
	assert.Nil(t, repo.FindAll(context.TODO(), &result, where.Eq("title", "Golang")))
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Eq("title", Any)).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(context.TODO(), &result, where.Eq("title", "Golang"))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_error(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.FindAll(context.TODO(), &result, where.Like("title", "%dummies%")))
	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAll(context.TODO(), &result, where.Like("title", "%dummies%"))
		assert.NotEqual(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_noMatch(t *testing.T) {
	var (
		repo   = New()
		result []Book
	)

	repo.ExpectFindAll(where.Eq("title", "a"))
	assert.PanicsWithError(t, "TODO: Query doesn't match", func() {
		repo.FindAll(context.TODO(), &result, where.Eq("title", "b"))
	})
	repo.AssertExpectations(t)
}
