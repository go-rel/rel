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

func TestFind_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectFind(where.Eq("status", "paid"))

	assert.Panics(t, func() {
		repo.Find(context.TODO(), where.Eq("status", "pending"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tFind(ctx, <Any>, query todo)", nt.lastLog)
}

func TestFind_String(t *testing.T) {
	var (
		mockFind = MockFind{assert: &Assert{}, argQuery: rel.Where(where.Eq("status", "paid"))}
	)

	assert.Equal(t, "Find(ctx, <Any>, query todo)", mockFind.String())
	assert.Equal(t, "ExpectFind(query todo)", mockFind.ExpectString())
}
