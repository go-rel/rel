package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Fs02/rel"
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
		query = rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10)
	)

	repo.ExpectFindAndCountAll(query).Result(books, 12)

	count, err := repo.FindAndCountAll(context.TODO(), &result, query)
	assert.Nil(t, err)
	assert.Equal(t, 12, count)
	assert.Equal(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAndCountAll(query).Result(books, 12)
	assert.NotPanics(t, func() {
		count := repo.MustFindAndCountAll(context.TODO(), &result, query)
		assert.Equal(t, books, result)
		assert.Equal(t, 12, count)
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
		query = rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10)
	)

	repo.ExpectFindAndCountAll(query).ConnectionClosed()

	count, err := repo.FindAndCountAll(context.TODO(), &result, query)
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, count)
	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)

	repo.ExpectFindAndCountAll(query).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAndCountAll(context.TODO(), &result, query)
	})

	assert.NotEqual(t, books, result)
	repo.AssertExpectations(t)
}
