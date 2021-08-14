package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAny(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).DeletedCount(1)
	deletedCount, err := repo.DeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)
	repo.AssertExpectations(t)
}

func TestDeleteAny_wildcard(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", Any))).DeletedCount(1)
	deletedCount, err := repo.DeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)
	repo.AssertExpectations(t)
}

func TestDeleteAny_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	_, err := repo.DeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	assert.Equal(t, sql.ErrConnDone, err)
	repo.AssertExpectations(t)
}

func TestDeleteAny_noTable(t *testing.T) {
	var (
		repo  = New()
		query = rel.Where(where.Eq("id", 1))
	)

	repo.ExpectDeleteAny(query)
	assert.Panics(t, func() {
		repo.MustDeleteAny(context.TODO(), query)
	})
	repo.AssertExpectations(t)
}

func TestDeleteAny_unsafe(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books")).Unsafe()
	assert.NotPanics(t, func() {
		repo.MustDeleteAny(context.TODO(), rel.From("books"))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAny_unsafe_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books"))
	assert.Panics(t, func() {
		repo.MustDeleteAny(context.TODO(), rel.From("books"))
	})
	repo.AssertExpectations(t)
}

func TestMustDeleteAny(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).DeletedCount(1)
	assert.NotPanics(t, func() {
		deletedCount := repo.MustDeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
		assert.Equal(t, 1, deletedCount)
	})
	repo.AssertExpectations(t)

}

func TestMustDeleteAny_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAny_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("users"))

	assert.Panics(t, func() {
		repo.DeleteAny(context.TODO(), rel.From("books"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tDeleteAny(ctx, query todo)", nt.lastLog)
}

func TestDeleteAny_String(t *testing.T) {
	var (
		mockDeleteAny = MockDeleteAny{assert: &Assert{}, argQuery: rel.From("users")}
	)

	assert.Equal(t, `DeleteAny(ctx, query todo)`, mockDeleteAny.String())
	assert.Equal(t, `ExpectDeleteAny(query todo)`, mockDeleteAny.ExpectString())
}
