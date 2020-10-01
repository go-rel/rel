package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestAggregate(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	sum, err := repo.Aggregate(context.TODO(), rel.From("books"), "sum", "id")
	assert.Nil(t, err)
	assert.Equal(t, 3, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	assert.NotPanics(t, func() {
		sum := repo.MustAggregate(context.TODO(), rel.From("books"), "sum", "id")
		assert.Equal(t, 3, sum)
	})
	repo.AssertExpectations(t)
}

func TestAggregate_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	sum, err := repo.Aggregate(context.TODO(), rel.From("books"), "sum", "id")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	assert.Panics(t, func() {
		sum := repo.MustAggregate(context.TODO(), rel.From("books"), "sum", "id")
		assert.Equal(t, 0, sum)
	})
	repo.AssertExpectations(t)
}

func TestAggregate_Count(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectCount("books").Result(2)
	count, err := repo.Count(context.TODO(), "books")
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	repo.AssertExpectations(t)

	repo.ExpectCount("books").Result(2)
	assert.NotPanics(t, func() {
		count := repo.MustCount(context.TODO(), "books")
		assert.Equal(t, 2, count)
	})
	repo.AssertExpectations(t)
}

func TestAggregate_Count_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectCount("books").ConnectionClosed()
	count, err := repo.Count(context.TODO(), "books")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, count)
	repo.AssertExpectations(t)

	repo.ExpectCount("books").ConnectionClosed()
	assert.Panics(t, func() {
		count := repo.MustCount(context.TODO(), "books")
		assert.Equal(t, 0, count)
	})
	repo.AssertExpectations(t)
}
