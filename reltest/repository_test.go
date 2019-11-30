package reltest

import (
	"database/sql"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

type Author struct {
	ID   int
	Name string
}

type Book struct {
	ID     int
	Title  string
	Author Author
}

func TestRepository_Aggregate(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	sum, err := repo.Aggregate(rel.From("books"), "sum", "id")
	assert.Nil(t, err)
	assert.Equal(t, 3, sum)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	assert.NotPanics(t, func() {
		sum := repo.MustAggregate(rel.From("books"), "sum", "id")
		assert.Equal(t, 3, sum)
	})
}

func TestRepository_Aggregate_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	sum, err := repo.Aggregate(rel.From("books"), "sum", "id")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, sum)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	assert.Panics(t, func() {
		sum := repo.MustAggregate(rel.From("books"), "sum", "id")
		assert.Equal(t, 0, sum)
	})
}

func TestRepository_Count(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectCount("books").Result(2)
	count, err := repo.Count("books")
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	repo.ExpectCount("books").Result(2)
	assert.NotPanics(t, func() {
		count := repo.MustCount("books")
		assert.Equal(t, 2, count)
	})
}

func TestRepository_Count_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectCount("books").ConnectionClosed()
	count, err := repo.Count("books")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, count)

	repo.ExpectCount("books").ConnectionClosed()
	assert.Panics(t, func() {
		count := repo.MustCount("books")
		assert.Equal(t, 0, count)
	})
}

func TestRepository_Find(t *testing.T) {
	var (
		repo   Repository
		result Book
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.Nil(t, repo.Find(&result, where.Eq("id", 2)))
	assert.Equal(t, book, result)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.NotPanics(t, func() {
		repo.MustFind(&result, where.Eq("id", 2))
		assert.Equal(t, book, result)
	})
}

func TestRepository_Find_noResult(t *testing.T) {
	var (
		repo   Repository
		result Book
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectFind(where.Eq("id", 2)).NoResult()

	assert.Equal(t, rel.NoResultError{}, repo.Find(&result, where.Eq("id", 2)))
	assert.NotEqual(t, book, result)

	repo.ExpectFind(where.Eq("id", 2)).NoResult()
	assert.Panics(t, func() {
		repo.MustFind(&result, where.Eq("id", 2))
		assert.NotEqual(t, book, result)
	})
}

func TestRepository_FindAll(t *testing.T) {
	var (
		repo   Repository
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.Nil(t, repo.FindAll(&result, where.Like("title", "%dummies%")))
	assert.Equal(t, books, result)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(&result, where.Like("title", "%dummies%"))
		assert.Equal(t, books, result)
	})
}

func TestRepository_FindAll_error(t *testing.T) {
	var (
		repo   Repository
		result []Book
		books  = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 2, Title: "Rel for dummies"},
		}
	)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.FindAll(&result, where.Like("title", "%dummies%")))
	assert.NotEqual(t, books, result)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAll(&result, where.Like("title", "%dummies%"))
		assert.NotEqual(t, books, result)
	})
}

func TestRepository_Delete(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDelete().Record(&Book{ID: 1})
	assert.Nil(t, repo.Delete(&Book{ID: 1}))

	repo.ExpectDelete().Record(&Book{ID: 1})
	assert.NotPanics(t, func() {
		repo.MustDelete(&Book{ID: 1})
	})
}

func TestRepository_Delete_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDelete().ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.Delete(&Book{ID: 1}))

	repo.ExpectDelete().ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDelete(&Book{ID: 1})
	})
}

func TestRepository_DeleteAll(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, repo.DeleteAll(rel.From("books").Where(where.Eq("id", 1))))

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	})
}

func TestRepository_DeleteAll_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.DeleteAll(rel.From("books").Where(where.Eq("id", 1))))

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	})
}

func TestRepository_DeleteAll_noTable(t *testing.T) {
	var (
		repo Repository
	)

	assert.Panics(t, func() {
		repo.ExpectDeleteAll()
	})
}

func TestRepository_DeleteAll_unsafe(t *testing.T) {
	var (
		repo Repository
	)

	assert.Panics(t, func() {
		repo.ExpectDeleteAll(rel.From("books"))
		repo.MustDeleteAll(rel.From("books"))
	})

	assert.NotPanics(t, func() {
		repo.ExpectDeleteAll(rel.From("books")).Unsafe()
		repo.MustDeleteAll(rel.From("books"))
	})
}

func TestRepository_Preload(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Rel for dummies"}
		book   = Book{ID: 2, Title: "Rel for dummies", Author: Author{ID: 1, Name: "Kia"}}
	)

	repo.ExpectPreload("author").Result(book)
	assert.Nil(t, repo.Preload(&result, "author"))
	assert.Equal(t, book, result)

	repo.ExpectPreload("author").Result(book)
	assert.NotPanics(t, func() {
		repo.MustPreload(&result, "author")
		assert.Equal(t, book, result)
	})
}

func TestRepository_Preload_error(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Rel for dummies"}
		book   = Book{ID: 2, Title: "Rel for dummies", Author: Author{ID: 1, Name: "Kia"}}
	)

	repo.ExpectPreload("author").ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.Preload(&result, "author"))
	assert.NotEqual(t, book, result)

	repo.ExpectPreload("author").ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustPreload(&result, "author")
		assert.NotEqual(t, book, result)
	})
}
