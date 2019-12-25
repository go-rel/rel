package reltest

import (
	"database/sql"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestDeleteAll(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, repo.DeleteAll(rel.From("books").Where(where.Eq("id", 1))))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.DeleteAll(rel.From("books").Where(where.Eq("id", 1))))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_noTable(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll()
	assert.Panics(t, func() {
		repo.MustDeleteAll()
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_unsafe(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books"))
	assert.Panics(t, func() {
		repo.MustDeleteAll(rel.From("books"))
	})
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books")).Unsafe()
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(rel.From("books"))
	})
	repo.AssertExpectations(t)
}
