package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDelete(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().For(&Book{ID: 1})
	assert.Nil(t, repo.Delete(context.TODO(), &Book{ID: 1}))
	repo.AssertExpectations(t)

	repo.ExpectDelete().For(&Book{ID: 1})
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
		repo.MustDelete(context.TODO(), &Book{ID: 1})
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

func TestDelete_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectDelete().ForTable("users")

	assert.Panics(t, func() {
		repo.Delete(context.TODO(), &Book{})
	})
	assert.False(t, repo.delete.assert(nt))
	assert.Equal(t, "FAIL: The code you are testing needs to call:\n\tDelete(ctx, <Table: users>)", nt.lastLog)
}

func TestDelete_String(t *testing.T) {
	var (
		mockDelete = MockDelete{assert: &Assert{}, argRecord: &Book{}}
	)

	assert.Equal(t, "Delete(ctx, &reltest.Book{ID:0, Title:\"\", Author:reltest.Author{ID:0, Name:\"\", Books:[]reltest.Book(nil)}, AuthorID:(*int)(nil), Ratings:[]reltest.Rating(nil), Poster:reltest.Poster{ID:0, Image:\"\", BookID:0}, AbstractID:0, Abstract:reltest.Abstract{ID:0, Content:\"\"}, Views:0})", mockDelete.String())
	assert.Equal(t, "ExpectDelete().ForType(\"*reltest.Book\")", mockDelete.ExpectString())
}
