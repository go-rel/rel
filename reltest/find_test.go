package reltest

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	var (
		repo   Repository
		result Book
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.Nil(t, repo.Find(&result, where.Eq("id", 2)))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.NotPanics(t, func() {
		repo.MustFind(&result, where.Eq("id", 2))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestFind_noResult(t *testing.T) {
	var (
		repo   Repository
		result Book
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).NoResult()

	assert.Equal(t, rel.NoResultError{}, repo.Find(&result, where.Eq("id", 2)))
	assert.NotEqual(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).NoResult()
	assert.Panics(t, func() {
		repo.MustFind(&result, where.Eq("id", 2))
		assert.NotEqual(t, book, result)
	})
	repo.AssertExpectations(t)
}
