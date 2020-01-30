package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreload(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Rel for dummies", AuthorID: 1}
		author = Author{ID: 1, Name: "Kia"}
	)

	repo.ExpectPreload("author").Result(author)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "author"))
	assert.Equal(t, author, result.Author)
	repo.AssertExpectations(t)

	repo.ExpectPreload("author").Result(author)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "author")
	})
	assert.Equal(t, author, result.Author)
	repo.AssertExpectations(t)
}

func TestPreload_nested(t *testing.T) {
	var (
		repo   = New()
		result = Rating{
			Book: &Book{ID: 2, Title: "Rel for dummies", AuthorID: 1},
		}
		author = Author{ID: 1, Name: "Kia"}
	)

	repo.ExpectPreload("book.author").Result(author)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "book.author"))
	assert.Equal(t, author, result.Book.Author)
	repo.AssertExpectations(t)

	repo.ExpectPreload("book.author").Result(author)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "book.author")
	})
	assert.Equal(t, author, result.Book.Author)
	repo.AssertExpectations(t)
}

func TestPreload_slice(t *testing.T) {
	var (
		repo   = New()
		result = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
		ratings = []Rating{
			{ID: 1, BookID: 1, Score: 10},
			{ID: 2, BookID: 1, Score: 8},
			{ID: 3, BookID: 2, Score: 9},
		}
	)

	repo.ExpectPreload("ratings").Result(ratings)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "ratings"))
	assert.Len(t, result[0].Ratings, 2)
	assert.Equal(t, ratings[:2], result[0].Ratings)
	assert.Len(t, result[1].Ratings, 1)
	assert.Equal(t, ratings[2:], result[1].Ratings)
	repo.AssertExpectations(t)

	repo.ExpectPreload("ratings").Result(ratings)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "ratings")
	})
	assert.Len(t, result[0].Ratings, 2)
	assert.Equal(t, ratings[:2], result[0].Ratings)
	assert.Len(t, result[1].Ratings, 1)
	assert.Equal(t, ratings[2:], result[1].Ratings)
	repo.AssertExpectations(t)
}

func TestPreload_sliceNested(t *testing.T) {
	var (
		repo   = New()
		result = []Author{
			{
				Books: []Book{
					{ID: 1, Title: "Golang for dummies"},
					{ID: 2, Title: "Rel for dummies"},
				},
			},
			{
				Books: nil,
			},
		}
		ratings = []Rating{
			{ID: 1, BookID: 1, Score: 10},
			{ID: 2, BookID: 1, Score: 8},
			{ID: 3, BookID: 2, Score: 9},
		}
	)

	repo.ExpectPreload("books.ratings").Result(ratings)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "books.ratings"))
	assert.Len(t, result[0].Books[0].Ratings, 2)
	assert.Equal(t, ratings[:2], result[0].Books[0].Ratings)
	assert.Len(t, result[0].Books[1].Ratings, 1)
	assert.Equal(t, ratings[2:], result[0].Books[1].Ratings)
	repo.AssertExpectations(t)

	repo.ExpectPreload("books.ratings").Result(ratings)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "books.ratings")
	})
	assert.Len(t, result[0].Books[0].Ratings, 2)
	assert.Equal(t, ratings[:2], result[0].Books[0].Ratings)
	assert.Len(t, result[0].Books[1].Ratings, 1)
	assert.Equal(t, ratings[2:], result[0].Books[1].Ratings)
	repo.AssertExpectations(t)
}

func TestPreload_nilReferenceValue(t *testing.T) {
	var (
		repo   = New()
		result = struct {
			ID       int
			Author   *Author
			AuthorID *int
		}{
			ID: 1,
		}
		author = Author{ID: 1, Name: "Kia"}
	)

	repo.ExpectPreload("author").Result(author)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "author"))
	assert.Nil(t, result.Author)
	repo.AssertExpectations(t)

	repo.ExpectPreload("author").Result(author)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "author")
	})
	assert.Nil(t, result.Author)
	repo.AssertExpectations(t)
}

func TestPreload_For(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Rel for dummies", AuthorID: 1}
	)

	repo.ExpectPreload("author").For(&result)
	assert.Nil(t, repo.Preload(context.TODO(), &result, "author"))
	repo.AssertExpectations(t)

	repo.ExpectPreload("author").For(&result)
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "author")
	})
	repo.AssertExpectations(t)
}

func TestPreload_ForType(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Rel for dummies", AuthorID: 1}
	)

	repo.ExpectPreload("author").ForType("reltest.Book")
	assert.Nil(t, repo.Preload(context.TODO(), &result, "author"))
	repo.AssertExpectations(t)

	repo.ExpectPreload("author").ForType("reltest.Book")
	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &result, "author")
	})
	repo.AssertExpectations(t)
}

func TestPreload_error(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Rel for dummies", AuthorID: 1}
	)

	repo.ExpectPreload("author").ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.Preload(context.TODO(), &result, "author"))
	repo.AssertExpectations(t)

	repo.ExpectPreload("author").ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustPreload(context.TODO(), &result, "author")
	})
	repo.AssertExpectations(t)
}
