package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAny(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).UpdatedCount(1)
	updatedCount, err := repo.UpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	assert.Nil(t, err)
	assert.Equal(t, 1, updatedCount)
	repo.AssertExpectations(t)
}

func TestUpdateAny_wildcard(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", Any)), rel.Set("discount", Any)).UpdatedCount(1)
	updatedCount, err := repo.UpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	assert.Nil(t, err)
	assert.Equal(t, 1, updatedCount)
	repo.AssertExpectations(t)
}

func TestUpdateAny_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", Any)).ConnectionClosed()
	_, err := repo.UpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", Any))
	assert.Equal(t, sql.ErrConnDone, err)
	repo.AssertExpectations(t)
}

func TestUpdateAny_noTable(t *testing.T) {
	var (
		repo  = New()
		query = rel.Where(where.Eq("id", 1))
	)

	repo.ExpectUpdateAny(query, rel.Set("discount", Any))
	assert.Panics(t, func() {
		repo.MustUpdateAny(context.TODO(), query, rel.Set("discount", Any))
	})
	repo.AssertExpectations(t)
}

func TestUpdateAny_unsafe(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books"), rel.Set("discount", Any)).Unsafe()
	assert.NotPanics(t, func() {
		repo.MustUpdateAny(context.TODO(), rel.From("books"), rel.Set("discount", Any))
	})
	repo.AssertExpectations(t)
}

func TestUpdateAny_unsafe_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books"), rel.Set("discount", Any))
	assert.Panics(t, func() {
		repo.MustUpdateAny(context.TODO(), rel.From("books"), rel.Set("discount", Any))
	})
	repo.AssertExpectations(t)
}

func TestMustUpdateAny(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).UpdatedCount(1)
	assert.NotPanics(t, func() {
		updatedCount := repo.MustUpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
		assert.Equal(t, 1, updatedCount)
	})
	repo.AssertExpectations(t)
}

func TestMustUpdateAny_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustUpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	})
	repo.AssertExpectations(t)
}

func TestUpdateAny_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAny(rel.From("users"))

	assert.Panics(t, func() {
		repo.UpdateAny(context.TODO(), rel.From("books"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tUpdateAny(ctx, query todo)", nt.lastLog)
}

func TestUpdateAny_String(t *testing.T) {
	var (
		mockUpdateAny = MockUpdateAny{assert: &Assert{}, argQuery: rel.From("users")}
	)

	assert.Equal(t, `UpdateAny(ctx, query todo)`, mockUpdateAny.String())
	assert.Equal(t, `ExpectUpdateAny(query todo)`, mockUpdateAny.ExpectString())
}
