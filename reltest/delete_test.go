package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().For(&Book{ID: 1}).Success()
	assert.Nil(t, repo.Delete(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().For(&Book{ID: 1}).Success()
	assert.NotPanics(t, func() {
		repo.MustDelete(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDelete_ForType(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForType("reltest.Book")
	assert.Nil(t, repo.Delete(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().ForType("reltest.Book")
	assert.NotPanics(t, func() {
		repo.MustDelete(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDelete_ForTable(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForTable("books")
	assert.Nil(t, repo.Delete(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().ForTable("books")
	assert.NotPanics(t, func() {
		repo.MustDelete(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDelete_ForContains(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForContains(Book{Title: "Golang"})
	assert.Nil(t, repo.Delete(context.TODO(), &Book{ID: 1, Title: "Golang"}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().ForContains(Book{Title: "Golang"})
	assert.NotPanics(t, func() {
		repo.MustDelete(context.TODO(), &Book{ID: 1, Title: "Golang"})
	})
	repo.AssertExpectations(t)
}

func TestDelete_error(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.Delete(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDelete(context.TODO(), &Book{ID: 1})
	})
	repo.AssertExpectations(t)
}

func TestDelete_noMatchCascade(t *testing.T) {
	var (
		repo = New()
	)

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{}, rel.Cascade(true))
	})
}

func TestDelete_assertFor(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().For(&Book{ID: 2})

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{ID: 1})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tDelete(ctx, &reltest.Book{ID: 2})", nt.lastLog)
}

func TestDelete_assertForTable(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForTable("users")

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tDelete(ctx, <Table: users>)", nt.lastLog)
}

func TestDelete_assertForType(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForType("User")

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tDelete(ctx, <Type: *User>)", nt.lastLog)
}

func TestDelete_assertForContains(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForContains(Book{ID: 3})

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{ID: 1})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tDelete(ctx, <Contains: reltest.Book{ID: 3}>)", nt.lastLog)
}

func TestDelete_assertCascade(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete(rel.Cascade(true))

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tDelete(ctx, <Any>, rel.Cascade(true))", nt.lastLog)
}

func TestDelete_String(t *testing.T) {
	var (
		mockDelete = MockDelete{assert: &Assert{}, argRecord: &Book{}}
	)

	assert.Equal(t, "Delete(ctx, &reltest.Book{})", mockDelete.String())
	assert.Equal(t, "ExpectDelete().ForType(\"*reltest.Book\")", mockDelete.ExpectString())
}
