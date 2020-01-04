package reltest

import (
	"database/sql"
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

type Author struct {
	ID    int
	Name  string
	Books []Book
}

type Rating struct {
	ID     int
	Score  int
	BookID int
	Book   *Book
}

type Poster struct {
	ID     int
	Image  string
	BookID int
}

type Book struct {
	ID       int
	Title    string
	Author   Author
	AuthorID int
	Ratings  []Rating
	Poster   Poster
	Views    int
}

func TestRepository_Adapter(t *testing.T) {
	assert.Nil(t, (&Repository{}).Adapter())
}

func TestRepository_Transaction(t *testing.T) {
	var (
		repo   Repository
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectInsert()
	})

	assert.Nil(t, repo.Transaction(func(repo rel.Repository) error {
		return repo.Insert(&result)
	}))

	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestRepository_Transaction_error(t *testing.T) {
	var (
		repo   Repository
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectInsert().ConnectionClosed()
	})

	assert.Equal(t, sql.ErrConnDone, repo.Transaction(func(repo rel.Repository) error {
		repo.MustInsert(&result)
		return nil
	}))

	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestRepository_Transaction_panic(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectTransaction(func(repo *Repository) {
	})

	assert.Panics(t, func() {
		_ = repo.Transaction(func(repo rel.Repository) error {
			panic("error")
		})
	})

	repo.AssertExpectations(t)
}

func TestRepository_Transaction_runtimerError(t *testing.T) {
	var (
		book *Book
		repo Repository
	)

	repo.ExpectTransaction(func(repo *Repository) {
	})

	assert.Panics(t, func() {
		_ = repo.Transaction(func(repo rel.Repository) error {
			_ = book.ID
			return nil
		})
	})

	repo.AssertExpectations(t)
}
