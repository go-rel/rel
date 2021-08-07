package reltest

import (
	"context"
	"database/sql"
	"testing"

	"github.com/go-rel/rel"
	"github.com/go-rel/rel/where"
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
	ID         int
	Title      string
	Author     Author
	AuthorID   *int
	Ratings    []Rating `auto:"true"`
	Poster     Poster   `autosave:"true"`
	AbstractID int
	Abstract   Abstract `autosave:"true"`
	Views      int
}

type Abstract struct {
	ID      int
	Content string
}

func TestRepository_Adapter(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = New()
	)

	assert.NotNil(t, repo.Adapter(ctx))
}

func TestRepository_Instrumentation(t *testing.T) {
	assert.NotPanics(t, func() {
		New().Instrumentation(rel.DefaultLogger)
	})
}

func TestRepository_Ping(t *testing.T) {
	assert.Nil(t, New().Ping(context.TODO()))
}

func TestRepository_Transaction(t *testing.T) {
	var (
		repo   = New()
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies", Ratings: []Rating{}}
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectInsert()

		repo.ExpectTransaction(func(repo *Repository) {
			repo.ExpectFind(where.Eq("id", 1))

			repo.ExpectTransaction(func(repo *Repository) {
				repo.ExpectDelete()
			})
		})
	})

	assert.Nil(t, repo.Transaction(context.TODO(), func(ctx context.Context) error {
		repo.MustInsert(ctx, &result)

		return repo.Transaction(ctx, func(ctx context.Context) error {
			repo.MustFind(ctx, &result, where.Eq("id", 1))

			return repo.Transaction(ctx, func(ctx context.Context) error {
				return repo.Delete(ctx, &result)
			})
		})
	}))

	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestRepository_Transaction_error(t *testing.T) {
	var (
		repo   = New()
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectTransaction(func(repo *Repository) {
		repo.ExpectInsert().ConnectionClosed()
	})

	assert.Equal(t, sql.ErrConnDone, repo.Transaction(context.TODO(), func(ctx context.Context) error {
		repo.MustInsert(ctx, &result)
		return nil
	}))

	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestRepository_Transaction_panic(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectTransaction(func(repo *Repository) {
	})

	assert.Panics(t, func() {
		_ = repo.Transaction(context.TODO(), func(ctx context.Context) error {
			panic("error")
		})
	})

	repo.AssertExpectations(t)
}

func TestRepository_Transaction_runtimerError(t *testing.T) {
	var (
		book *Book
		repo = New()
	)

	repo.ExpectTransaction(func(repo *Repository) {
	})

	assert.Panics(t, func() {
		_ = repo.Transaction(context.TODO(), func(ctx context.Context) error {
			_ = book.ID
			return nil
		})
	})

	repo.AssertExpectations(t)
}
