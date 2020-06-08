package main

import (
	"context"
	"testing"

	"github.com/Fs02/rel/reltest"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestQuickExample(t *testing.T) {
	repo := reltest.New()

	/// [quick-example]

	// Mock insert a Book.
	repo.ExpectInsert().ForType("main.Book")

	// Mock find and return result
	repo.ExpectFind(where.Eq("id", 1)).Result(Book{
		ID:    1,
		Title: "REL for Dummies",
	})

	// Mock update
	repo.ExpectUpdate().ForType("main.Book")

	// Mock delete.
	repo.ExpectDelete().ForType("main.Book")

	/// [quick-example]

	// run and asserts
	assert.Nil(t, QuickExample(context.Background(), repo))
	repo.AssertExpectations(t)

}

func TestQuickExample_DeleteError(t *testing.T) {
	repo := reltest.New()

	// Mock insert a Book.
	repo.ExpectInsert().ForType("main.Book")

	// Mock Find and return result
	repo.ExpectFind(where.Eq("id", 1)).Result(Book{ID: 1})

	// Mock update
	repo.ExpectUpdate().ForType("main.Book")

	// Mock delete.
	repo.ExpectDelete().ForType("main.Book").ConnectionClosed()

	// run and asserts
	assert.Equal(t, reltest.ErrConnectionClosed, QuickExample(context.Background(), repo))
	repo.AssertExpectations(t)
}

func TestQuickExample_UpdateError(t *testing.T) {
	repo := reltest.New()

	// Mock insert a Book.
	repo.ExpectInsert().ForType("main.Book")

	// Mock Find and return result
	repo.ExpectFind(where.Eq("id", 1)).Result(Book{ID: 1})

	// Mock update
	repo.ExpectUpdate().ForType("main.Book").ConnectionClosed()

	// run and asserts
	assert.Equal(t, reltest.ErrConnectionClosed, QuickExample(context.Background(), repo))
	repo.AssertExpectations(t)
}

func TestQuickExample_FindError(t *testing.T) {
	repo := reltest.New()

	// Mock insert a Book.
	repo.ExpectInsert().ForType("main.Book")

	// Mock Find and return result
	repo.ExpectFind(where.Eq("id", 1)).ConnectionClosed()

	// run and asserts
	assert.Equal(t, reltest.ErrConnectionClosed, QuickExample(context.Background(), repo))
	repo.AssertExpectations(t)
}

func TestQuickExample_InsertError(t *testing.T) {
	repo := reltest.New()

	// Mock insert a Book.
	repo.ExpectInsert().ForType("main.Book").ConnectionClosed()

	// run and asserts
	assert.Equal(t, reltest.ErrConnectionClosed, QuickExample(context.Background(), repo))
	repo.AssertExpectations(t)
}
