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

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).Result(1)
	updatedCount, err := repo.UpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	assert.Nil(t, err)
	assert.Equal(t, 1, updatedCount)
	repo.AssertExpectations(t)

	repo.ExpectUpdateAny(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).Result(1)
	assert.NotPanics(t, func() {
		updatedCount = repo.MustUpdateAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
		assert.Equal(t, 1, updatedCount)
	})
	repo.AssertExpectations(t)
}

func TestDeleteAny(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).Result(1)
	deletedCount, err := repo.DeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)
	repo.AssertExpectations(t)

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).Result(1)
	assert.NotPanics(t, func() {
		deletedCount = repo.MustDeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
		assert.Equal(t, 1, deletedCount)
	})
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

	repo.ExpectDeleteAny(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAny(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	})
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

	repo.ExpectDeleteAny(rel.From("books"))
	assert.Panics(t, func() {
		repo.MustDeleteAny(context.TODO(), rel.From("books"))
	})
	repo.AssertExpectations(t)

	repo.ExpectDeleteAny(rel.From("books")).Unsafe()
	assert.NotPanics(t, func() {
		repo.MustDeleteAny(context.TODO(), rel.From("books"))
	})
	repo.AssertExpectations(t)
}
