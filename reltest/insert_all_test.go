package reltest

import (
	"context"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestInsertAll(t *testing.T) {
	var (
		repo    = New()
		results = []Book{
			{Title: "Golang for dummies"},
			{Title: "Rel for dummies"},
		}
		books = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectInsertAll().Success()
	assert.Nil(t, repo.InsertAll(context.TODO(), &results))
	assert.Equal(t, books, results)
	repo.AssertExpectations(t)

	repo.ExpectInsertAll().Success()
	assert.NotPanics(t, func() {
		repo.MustInsertAll(context.TODO(), &results)
		assert.Equal(t, books, results)
	})
	repo.AssertExpectations(t)
}

func TestInsertAll_error(t *testing.T) {
	var (
		repo    = New()
		results []Book
	)

	repo.ExpectInsertAll().ConnectionClosed()
	assert.Equal(t, ErrConnectionClosed, repo.InsertAll(context.TODO(), &results))
	repo.AssertExpectations(t)

	repo.ExpectInsertAll().NotUnique("title")
	assert.Equal(t, rel.ConstraintError{
		Key:  "title",
		Type: rel.UniqueConstraint,
	}, repo.InsertAll(context.TODO(), &results))
	repo.AssertExpectations(t)
}

func TestInsertAll_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectInsertAll().For(&[]Book{{Title: "Golang"}})

	assert.Panics(t, func() {
		repo.InsertAll(context.TODO(), &[]Book{{Title: "Go"}})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tInsertAll(ctx, &[]reltest.Book{reltest.Book{Title: Golang}})", nt.lastLog)
}

func TestInsertAll_assertForTable(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectInsertAll().ForTable("users")

	assert.Panics(t, func() {
		repo.InsertAll(context.TODO(), &[]Book{})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tInsertAll(ctx, <Table: users>)", nt.lastLog)
}

func TestInsertAll_assertForType(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectInsertAll().ForType("[]User")

	assert.Panics(t, func() {
		repo.InsertAll(context.TODO(), &[]Book{})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: Mock defined but not called:\n\tInsertAll(ctx, <Type: *[]User>)", nt.lastLog)
}

func TestInsertAll_String(t *testing.T) {
	var (
		mockInsertAll = MockInsertAll{assert: &Assert{}, argRecord: &[]Book{}}
	)

	assert.Equal(t, "InsertAll(ctx, &[]reltest.Book{})", mockInsertAll.String())
	assert.Equal(t, "InsertAll().ForType(\"*[]reltest.Book\")", mockInsertAll.ExpectString())
}
