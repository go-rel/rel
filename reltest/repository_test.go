package reltest

import (
	"database/sql"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
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
}

func TestRepository_Adapter(t *testing.T) {
	assert.Nil(t, (&Repository{}).Adapter())
}

func TestRepository_Aggregate(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	sum, err := repo.Aggregate(rel.From("books"), "sum", "id")
	assert.Nil(t, err)
	assert.Equal(t, 3, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").Result(3)
	assert.NotPanics(t, func() {
		sum := repo.MustAggregate(rel.From("books"), "sum", "id")
		assert.Equal(t, 3, sum)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Aggregate_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	sum, err := repo.Aggregate(rel.From("books"), "sum", "id")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, sum)
	repo.AssertExpectations(t)

	repo.ExpectAggregate(rel.From("books"), "sum", "id").ConnectionClosed()
	assert.Panics(t, func() {
		sum := repo.MustAggregate(rel.From("books"), "sum", "id")
		assert.Equal(t, 0, sum)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Count(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectCount("books").Result(2)
	count, err := repo.Count("books")
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	repo.AssertExpectations(t)

	repo.ExpectCount("books").Result(2)
	assert.NotPanics(t, func() {
		count := repo.MustCount("books")
		assert.Equal(t, 2, count)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Count_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectCount("books").ConnectionClosed()
	count, err := repo.Count("books")
	assert.Equal(t, sql.ErrConnDone, err)
	assert.Equal(t, 0, count)
	repo.AssertExpectations(t)

	repo.ExpectCount("books").ConnectionClosed()
	assert.Panics(t, func() {
		count := repo.MustCount("books")
		assert.Equal(t, 0, count)
	})
	repo.AssertExpectations(t)
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
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).Result(book)
	assert.NotPanics(t, func() {
		repo.MustFind(&result, where.Eq("id", 2))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
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
	repo.AssertExpectations(t)

	repo.ExpectFind(where.Eq("id", 2)).NoResult()
	assert.Panics(t, func() {
		repo.MustFind(&result, where.Eq("id", 2))
		assert.NotEqual(t, book, result)
	})
	repo.AssertExpectations(t)
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
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).Result(books)
	assert.NotPanics(t, func() {
		repo.MustFindAll(&result, where.Like("title", "%dummies%"))
		assert.Equal(t, books, result)
	})
	repo.AssertExpectations(t)
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
	repo.AssertExpectations(t)

	repo.ExpectFindAll(where.Like("title", "%dummies%")).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustFindAll(&result, where.Like("title", "%dummies%"))
		assert.NotEqual(t, books, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert(t *testing.T) {
	var (
		repo   Repository
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectInsert()
	assert.Nil(t, repo.Insert(&result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert()
	assert.NotPanics(t, func() {
		repo.MustInsert(&result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert_nested(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			Title:  "Rel for dummies",
			Author: Author{Name: "Kia"},
			Ratings: []Rating{
				{Score: 9},
				{Score: 10},
			},
			Poster: Poster{Image: "http://image.url"},
		}
		book = Book{
			ID:       1,
			Title:    "Rel for dummies",
			Author:   Author{ID: 1, Name: "Kia"},
			AuthorID: 1,
			Ratings: []Rating{
				{ID: 1, Score: 9, BookID: 1},
				{ID: 1, Score: 10, BookID: 1},
			},
			Poster: Poster{ID: 1, BookID: 1, Image: "http://image.url"},
		}
	)

	repo.ExpectInsert()
	assert.Nil(t, repo.Insert(&result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert()
	assert.NotPanics(t, func() {
		repo.MustInsert(&result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert_record(t *testing.T) {
	var (
		repo   Repository
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectInsert().ForType("reltest.Book")
	assert.Nil(t, repo.Insert(&result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert().ForType("reltest.Book")
	assert.NotPanics(t, func() {
		repo.MustInsert(&result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert_set(t *testing.T) {
	var (
		repo   Repository
		result Book
		book   = Book{ID: 1, Title: "Rel for dummies"}
	)

	repo.ExpectInsert(rel.Set("title", "Rel for dummies"))
	assert.Nil(t, repo.Insert(&result, rel.Set("title", "Rel for dummies")))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert(rel.Set("title", "Rel for dummies"))
	assert.NotPanics(t, func() {
		repo.MustInsert(&result, rel.Set("title", "Rel for dummies"))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert_map(t *testing.T) {
	var (
		repo   Repository
		result Book
		book   = Book{
			ID:       1,
			Title:    "Rel for dummies",
			Author:   Author{ID: 2, Name: "Kia"},
			AuthorID: 2,
			Ratings: []Rating{
				{ID: 1, Score: 9, BookID: 1},
				{ID: 1, Score: 10, BookID: 1},
			},
			Poster: Poster{ID: 1, BookID: 1, Image: "http://image.url"},
		}
		ch = rel.Map{
			"title": "Rel for dummies",
			"author": rel.Map{
				"id":   2,
				"name": "Kia",
			},
			"ratings": []rel.Map{
				{"score": 9},
				{"score": 10},
			},
			"poster": rel.Map{
				"image": "http://image.url",
			},
		}
	)

	repo.ExpectInsert(ch)
	assert.Nil(t, repo.Insert(&result, ch))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert(ch)
	assert.NotPanics(t, func() {
		repo.MustInsert(&result, ch)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert_unknownField(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectInsert(rel.Set("titles", "Rel for dummies"))
	assert.Panics(t, func() {
		_ = repo.Insert(&result, rel.Set("titles", "Rel for dummies"))
	})
	repo.AssertExpectations(t)
}

func TestRepository_Insert_notUnique(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectInsert(rel.Set("title", "Rel for dummies")).NotUnique("title")
	assert.Equal(t,
		rel.ConstraintError{Key: "title", Type: rel.UniqueConstraint},
		repo.Insert(&result, rel.Set("title", "Rel for dummies")),
	)
	repo.AssertExpectations(t)
}

func TestRepository_InsertAll(t *testing.T) {
	var (
		repo    Repository
		results = []Book{
			{Title: "Golang for dummies"},
			{Title: "Rel for dummies"},
		}
		books = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 1, Title: "Rel for dummies"},
		}
	)

	repo.ExpectInsertAll()
	assert.Nil(t, repo.InsertAll(&results))
	assert.Equal(t, books, results)
	repo.AssertExpectations(t)

	repo.ExpectInsertAll()
	assert.NotPanics(t, func() {
		repo.MustInsertAll(&results)
		assert.Equal(t, books, results)
	})
	repo.AssertExpectations(t)
}

func TestRepository_InsertAll_map(t *testing.T) {
	var (
		repo    Repository
		results []Book
		books   = []Book{
			{ID: 1, Title: "Golang for dummies"},
			{ID: 1, Title: "Rel for dummies"},
		}
	)

	repo.ExpectInsertAll(
		rel.BuildChanges(rel.Map{"title": "Golang for dummies"}),
		rel.BuildChanges(rel.Map{"title": "Rel for dummies"}),
	)
	assert.Nil(t, repo.InsertAll(&results,
		rel.BuildChanges(rel.Map{"title": "Golang for dummies"}),
		rel.BuildChanges(rel.Map{"title": "Rel for dummies"}),
	))
	assert.Equal(t, books, results)
	repo.AssertExpectations(t)
}

func TestRepository_Update(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate()
	assert.Nil(t, repo.Update(&result))
	repo.AssertExpectations(t)

	repo.ExpectUpdate()
	assert.NotPanics(t, func() {
		repo.MustUpdate(&result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_nested(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:       2,
			Title:    "Rel for dummies",
			Author:   Author{ID: 2, Name: "Kia"},
			AuthorID: 2,
			Ratings: []Rating{
				{ID: 1, BookID: 2, Score: 9},
				{ID: 1, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
		book = Book{
			ID:       2,
			Title:    "Rel for dummies",
			Author:   Author{ID: 2, Name: "Kia"},
			AuthorID: 2,
			Ratings: []Rating{
				{ID: 1, BookID: 2, Score: 9},
				{ID: 1, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
	)

	repo.ExpectUpdate()
	assert.Nil(t, repo.Update(&result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate()
	assert.NotPanics(t, func() {
		repo.MustUpdate(&result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_nestedInsert(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:     2,
			Title:  "Rel for dummies",
			Author: Author{Name: "Kia"},
			Ratings: []Rating{
				{Score: 9},
				{Score: 10},
			},
			Poster: Poster{Image: "http://image.url"},
		}
		book = Book{
			ID:       2,
			Title:    "Rel for dummies",
			Author:   Author{ID: 1, Name: "Kia"},
			AuthorID: 1,
			Ratings: []Rating{
				{ID: 1, BookID: 2, Score: 9},
				{ID: 1, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
	)

	repo.ExpectUpdate()
	assert.Nil(t, repo.Update(&result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate()
	assert.NotPanics(t, func() {
		repo.MustUpdate(&result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_record(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate().For(&result)
	assert.Nil(t, repo.Update(&result))
	repo.AssertExpectations(t)

	repo.ExpectUpdate().For(&result)
	assert.NotPanics(t, func() {
		repo.MustUpdate(&result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_withoutPrimaryValue(t *testing.T) {
	var (
		repo   Repository
		result = Book{Title: "Golang for dummies"}
	)

	repo.ExpectUpdate().For(&result)
	assert.Panics(t, func() {
		_ = repo.Update(&result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_set(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectUpdate(rel.Set("title", "Rel for dummies"))
	assert.Nil(t, repo.Update(&result, rel.Set("title", "Rel for dummies")))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate(rel.Set("title", "Rel for dummies"))
	assert.NotPanics(t, func() {
		repo.MustUpdate(&result, rel.Set("title", "Rel for dummies"))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_setNil(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
		book   = Book{ID: 2, Title: ""}
	)

	repo.ExpectUpdate(rel.Set("title", nil))
	assert.Nil(t, repo.Update(&result, rel.Set("title", nil)))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestRepository_Update_map(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:    2,
			Title: "Golang for dummies",
			Ratings: []Rating{
				{ID: 4, BookID: 2, Score: 15},
				{ID: 2, BookID: 2, Score: 5},
				{ID: 3, BookID: 2, Score: 6},
			},
		}
		book = Book{
			ID:       2,
			Title:    "Rel for dummies",
			Author:   Author{ID: 2, Name: "Kia"},
			AuthorID: 2,
			Ratings: []Rating{
				{ID: 2, BookID: 2, Score: 9},
				{ID: 1, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
		ch = rel.Map{
			"title": "Rel for dummies",
			"author": rel.Map{
				"id":   2,
				"name": "Kia",
			},
			"ratings": []rel.Map{
				{"id": 2, "score": 9},
				{"score": 10},
			},
			"poster": rel.Map{
				"image": "http://image.url",
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Nil(t, repo.Update(&result, ch))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate(ch)
	assert.NotPanics(t, func() {
		repo.MustUpdate(&result, ch)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_belongsToInconsistentPk(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:       2,
			Title:    "Golang for dummies",
			AuthorID: 2,
			Author:   Author{ID: 2, Name: "Kia"},
		}
		ch = rel.Map{
			"author": rel.Map{
				"id":   1,
				"name": "Koa",
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_belongsToInconsistentFk(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:       2,
			Title:    "Golang for dummies",
			AuthorID: 1,
			Author:   Author{ID: 2, Name: "Kia"},
		}
		ch = rel.Map{
			"author": rel.Map{
				"id":   2,
				"name": "Koa",
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_hasOneInconsistentPk(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:     2,
			Title:  "Golang for dummies",
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
		ch = rel.Map{
			"poster": rel.Map{
				"id":    2,
				"image": "http://image.url/other",
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_hasOneInconsistentFk(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:     2,
			Title:  "Golang for dummies",
			Poster: Poster{ID: 1, BookID: 1, Image: "http://image.url"},
		}
		ch = rel.Map{
			"poster": rel.Map{
				"id":    1,
				"image": "http://image.url/other",
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_hasManyNotLoaded(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:    2,
			Title: "Golang for dummies",
		}
		ch = rel.Map{
			"ratings": []rel.Map{
				{"id": 2, "score": 9},
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_hasManyInconsistentPk(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:    2,
			Title: "Golang for dummies",
			Ratings: []Rating{
				{ID: 2, BookID: 2, Score: 5},
			},
		}
		ch = rel.Map{
			"ratings": []rel.Map{
				{"id": 1, "score": 9},
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_hasManyInconsistentFk(t *testing.T) {
	var (
		repo   Repository
		result = Book{
			ID:    2,
			Title: "Golang for dummies",
			Ratings: []Rating{
				{ID: 2, BookID: 1, Score: 5},
			},
		}
		ch = rel.Map{
			"ratings": []rel.Map{
				{"id": 2, "score": 9},
			},
		}
	)

	repo.ExpectUpdate(ch)
	assert.Panics(t, func() {
		_ = repo.Update(&result, ch)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_unknownField(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate(rel.Set("titles", "Rel for dummies"))
	assert.Panics(t, func() {
		_ = repo.Update(&result, rel.Set("titles", "Rel for dummies"))
	})
	repo.AssertExpectations(t)
}

func TestRepository_Update_notUnique(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate(rel.Set("title", "Rel for dummies")).NotUnique("title")
	assert.Equal(t,
		rel.ConstraintError{Key: "title", Type: rel.UniqueConstraint},
		repo.Update(&result, rel.Set("title", "Rel for dummies")),
	)
	repo.AssertExpectations(t)
}

func TestRepository_Save(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectSave()
	assert.Nil(t, repo.Save(&result))
	repo.AssertExpectations(t)

	repo.ExpectSave()
	assert.NotPanics(t, func() {
		repo.MustSave(&result)
	})
	repo.AssertExpectations(t)
}

func TestRepository_Delete(t *testing.T) {
	var (
		repo Repository
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

func TestRepository_Delete_forType(t *testing.T) {
	var (
		repo Repository
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

func TestRepository_Delete_error(t *testing.T) {
	var (
		repo Repository
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

func TestRepository_DeleteAll(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.Nil(t, repo.DeleteAll(rel.From("books").Where(where.Eq("id", 1))))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestRepository_DeleteAll_error(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Equal(t, sql.ErrConnDone, repo.DeleteAll(rel.From("books").Where(where.Eq("id", 1))))
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1))).ConnectionClosed()
	assert.Panics(t, func() {
		repo.MustDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
	})
	repo.AssertExpectations(t)
}

func TestRepository_DeleteAll_noTable(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll()
	assert.Panics(t, func() {
		repo.MustDeleteAll()
	})
	repo.AssertExpectations(t)
}

func TestRepository_DeleteAll_unsafe(t *testing.T) {
	var (
		repo Repository
	)

	repo.ExpectDeleteAll(rel.From("books"))
	assert.Panics(t, func() {
		repo.MustDeleteAll(rel.From("books"))
	})
	repo.AssertExpectations(t)

	repo.ExpectDeleteAll(rel.From("books")).Unsafe()
	assert.NotPanics(t, func() {
		repo.MustDeleteAll(rel.From("books"))
	})
	repo.AssertExpectations(t)
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
