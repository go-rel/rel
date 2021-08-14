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
