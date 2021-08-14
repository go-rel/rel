package reltest

import (
	"context"
	"testing"

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

	repo.ExpectInsertAll()
	assert.Nil(t, repo.InsertAll(context.TODO(), &results))
	assert.Equal(t, books, results)
	repo.AssertExpectations(t)

	repo.ExpectInsertAll()
	assert.NotPanics(t, func() {
		repo.MustInsertAll(context.TODO(), &results)
		assert.Equal(t, books, results)
	})
	repo.AssertExpectations(t)
}

func TestInsertAll_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectInsertAll().ForTable("users")

	assert.Panics(t, func() {
		repo.InsertAll(context.TODO(), &[]Book{})
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: The code you are testing needs to call:\n\tInsertAll(ctx, <Table: users>)", nt.lastLog)
}

func TestInsertAll_String(t *testing.T) {
	var (
		mockInsertAll = MockInsertAll{assert: &Assert{}, argRecord: &[]Book{}}
	)

	assert.Equal(t, "InsertAll(ctx, &[]reltest.Book{})", mockInsertAll.String())
	assert.Equal(t, "InsertAll().ForType(\"*[]reltest.Book\")", mockInsertAll.ExpectString())
}
