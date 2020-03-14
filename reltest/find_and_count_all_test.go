package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFindAndCountAll(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAndCountAll(where.Like("title", "%dummies%")).Result(books, 2)

	count, err := repo.FindAndCountAll(context.TODO(), &result, where.Like("title", "%dummies%"))
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAndCountAll(where.Like("title", "%dummies%")).Result(books, 2)
	assert.NotPanics(t, func() {
		count := repo.MustFindAndCountAll(context.TODO(), &result, where.Like("title", "%dummies%"))
		assert.Equal(t, books, result)
		assert.Equal(t, 2, count)
	})
	repo.AssertExpectations(t)
}

func TestFindAndCountAll_error(t *testing.T) {
	var (
		repo   = New()
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAndCountAll(where.Like("title", "%dummies%")).ConnectionClosed()

	count, err := repo.FindAndCountAll(context.TODO(), &result, where.Like("title", "%dummies%"))
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, count)
	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAndCountAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAndCountAll(context.TODO(), &result, where.Like("title", "%dummies%"))
	})

	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)
}
