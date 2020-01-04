package reltest

import (
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().For(&Book{ID: 1})
	assert.Nil(t, repo.Delete(&Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().For(&Book{ID: 1})
	assert.NotPanics(t, func() {
		repo.MustDelete(&Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDelete_forType(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForType("reltest.Book")
	assert.Nil(t, repo.Delete(&Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().ForType("reltest.Book")
	assert.NotPanics(t, func() {
		repo.MustDelete(&Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDelete_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.Delete(&Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDelete(&Book{ID: 1})
	})
	repo.AssertExpectations(t)
}
