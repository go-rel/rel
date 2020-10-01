package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	var (
		repo   = New()
		result Book
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.Nil(t, repo.Find(context.TODO(), &result, where.Eq("id", 2)))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.NotPanics(t, func() {
		repo.MustFind(context.TODO(), &result, where.Eq("id", 2))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestFind_noResult(t *testing.T) {
	var (
		result Book
		repo   = New()
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).NotFound()

	assert.Equal(t, rel.NotFoundError{}, repo.Find(context.TODO(), &result, where.Eq("id", 2)))
	assert.NotEqual(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).NotFound()
	assert.Panics(t, func() {
		repo.MustFind(context.TODO(), &result, where.Eq("id", 2))
		assert.NotEqual(t, book, result)
	})
	repo.AssertExpectations(t)
}
