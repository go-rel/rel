package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCount(t *testing.T) {
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

func TestCount_error(t *testing.T) {
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
