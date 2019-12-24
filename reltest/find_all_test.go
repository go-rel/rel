package reltest

import (
	"database/sql"
	"testing"

	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	var (
		repo   Repository
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.Nil(t, repo.FindAll(&result, where.Like("title", "%dummies%")))
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(&result, where.Like("title", "%dummies%"))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestFindAll_error(t *testing.T) {
	var (
		repo   Repository
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.FindAll(&result, where.Like("title", "%dummies%")))
	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAll(&result, where.Like("title", "%dummies%"))
		assert.NotEqual(t, books, result)
	})
	repo.AssertExpectations(t)
}
