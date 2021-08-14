package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteAll(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll().For(&Book{ID: 1})
	assert.Nil(t, repo.DeleteAll(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll().For(&Book{ID: 1})
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_ForType(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll().ForType("reltest.Book")
	assert.Nil(t, repo.DeleteAll(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll().ForType("reltest.Book")
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_ForTable(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll().ForTable("books")
	assert.Nil(t, repo.DeleteAll(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll().ForTable("books")
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDeleteAll_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDeleteAll().ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.DeleteAll(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll().ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAll(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}
