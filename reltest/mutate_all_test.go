package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestUpdateAll(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectUpdateAll(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	assert.Nil(t, repo.UpdateAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true)))
	repo.AssertExpectations(t)

	repo.ExpectUpdateAll(rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	assert.NotPanics(t, func() {
		repo.MustUpdateAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)), rel.Set("discount", true))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, repo.DeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1))))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.DeleteAll(context.TODO(), rel.From("books").Where(where.Eq("id", 1))))
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
