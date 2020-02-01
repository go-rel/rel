/// [example]
package main

import (
	"context"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestExample(t *testing.T) {
	// create a mocked repository.
	var (
		repo = reltest.New()
		book = Book{
			ID:       1,
			Title:    "Go for dummies",
			Category: "learning",
		}
	)

	// mock find and return other result
	repo.ExpectFind(where.Eq("id", 1)).Result(book)

	// mocks update
	book.Title = "rel for dummies"
	repo.ExpectUpdate().For(&book)

	// run and asserts
	assert.Nil(t, Example(context.Background(), repo))
	repo.AssertExpectations(t)
}

func TestExample_findNoResult(t *testing.T) {
	// create a mocked repository.
	var repo = reltest.New()

	// mock find and return other result
	repo.ExpectFind(where.Eq("id", 1)).NotFound()

	// run and asserts
	assert.Equal(t, rel.NotFoundError{}, Example(context.Background(), repo))
	repo.AssertExpectations(t)
}

func TestExample_updateError(t *testing.T) {
	// create a mocked repository.
	var repo = reltest.New()

	// mock find and return other result
	repo.ExpectFind(where.Eq("id", 1)).Result(Book{ID: 1})

	// mocks update
	repo.ExpectUpdate().ForType("main.Book").ConnectionClosed()

	// run and asserts
	assert.Equal(t, reltest.ErrConnectionClosed, Example(context.Background(), repo))
	repo.AssertExpectations(t)
}

/// [example]

func TestMain(t *testing.T) {
	assert.NotPanics(t, func() {
		main()
	})
}
