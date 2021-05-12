package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAll(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAll(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).Result(1)
	updatedCount, err := repo.UpdateAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	assert.Nil(t, err)
	assert.Equal(t, 1, updatedCount)
	repo.AssertExpectations(t)

	repo.ExpectUpdateAll(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)).Result(1)
	assert.NotPanics(t, func() {
		updatedCount = repo.MustUpdateAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
		assert.Equal(t, 1, updatedCount)
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).Result(1)
	deletedCount, err := repo.DeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).Result(1)
	assert.NotPanics(t, func() {
		deletedCount = repo.MustDeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
		assert.Equal(t, 1, deletedCount)
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	_, err := repo.DeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	assert.Equal(t, sql.ErrConnDone, err)
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_noTable(t *testing.T) {
	var (
		repo  = New()
		query = rel.Where(where.Eq("id", 1))
	)

	repo.ExpectDeleteAll(query)
	assert.Panics(t, func() {
		repo.MustDeleteAll(context.TODO(), query)
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_unsafe(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll(rel.From("books"))
	assert.Panics(t, func() {
		repo.MustDeleteAll(context.TODO(), rel.From("books"))
	})
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books")).Unsafe()
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), rel.From("books"))
	})
	repo.AssertExpectations(t)
}
