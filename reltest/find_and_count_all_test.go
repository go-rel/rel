package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
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

func TestFindAndCountAll_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectFindAndCountAll(where.Eq("status", "paid"))

	assert.Panics(t, func() {
		repo.FindAndCountAll(context.TODO(), where.Eq("status", "pending"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: The code you are testing needs to call:\n\tFindAndCountAll(ctx, <Any>, query todo)", nt.lastLog)
}

func TestFindAndCountAll_String(t *testing.T) {
	var (
		mockFindAndCountAll = MockFindAndCountAll{assert: &Assert{}, argQuery: rel.Where(where.Eq("status", "paid"))}
	)

	assert.Equal(t, "FindAndCountAll(ctx, <Any>, query todo)", mockFindAndCountAll.String())
	assert.Equal(t, "ExpectFindAndCountAll(query todo)", mockFindAndCountAll.ExpectString())
}
