package reltest

import (
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestModify_Insert(t *testing.T) {
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

func TestModify_Insert_nested(t *testing.T) {
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

func TestModify_Insert_record(t *testing.T) {
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

func TestModify_Insert_set(t *testing.T) {
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

func TestModify_Insert_inc(t *testing.T) {
	var (
		repo   Repository
		result Book
	)

	repo.ExpectInsert(rel.Inc("views"))
	assert.Panics(t, func() {
		assert.Nil(t, repo.Insert(&result, rel.Inc("views")))
	})

	repo.AssertExpectations(t)
}

func TestModify_Insert_dec(t *testing.T) {
	var (
		repo   Repository
		result Book
	)

	repo.ExpectInsert(rel.Dec("views"))
	assert.Panics(t, func() {
		assert.Nil(t, repo.Insert(&result, rel.Dec("views")))
	})

	repo.AssertExpectations(t)
}

func TestModify_Insert_map(t *testing.T) {
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

func TestModify_Insert_unknownField(t *testing.T) {
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

func TestModify_Insert_notUnique(t *testing.T) {
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

func TestModify_InsertAll(t *testing.T) {
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

func TestModify_InsertAll_map(t *testing.T) {
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

func TestModify_Update(t *testing.T) {
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

func TestModify_Update_nested(t *testing.T) {
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

func TestModify_Update_nestedInsert(t *testing.T) {
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

func TestModify_Update_record(t *testing.T) {
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

func TestModify_Update_withoutPrimaryValue(t *testing.T) {
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

func TestModify_Update_set(t *testing.T) {
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

func TestModify_Update_inc(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Views: 10}
		book   = Book{ID: 2, Views: 11}
	)

	repo.ExpectUpdate(rel.Inc("views"))
	assert.Nil(t, repo.Update(&result, rel.Inc("views")))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestModify_Update_dec(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Views: 10}
		book   = Book{ID: 2, Views: 9}
	)

	repo.ExpectUpdate(rel.Dec("views"))
	assert.Nil(t, repo.Update(&result, rel.Dec("views")))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestModify_Update_incOrDecFieldNotExists(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Views: 10}
	)

	repo.ExpectUpdate(rel.Inc("watistis"))
	assert.Panics(t, func() {
		assert.Nil(t, repo.Update(&result, rel.Inc("watistis")))
	})

	repo.AssertExpectations(t)
}

func TestModify_Update_incOrDecFieldInvalid(t *testing.T) {
	var (
		repo   Repository
		result = Book{ID: 2, Views: 10}
	)

	repo.ExpectUpdate(rel.Inc("title"))
	assert.Panics(t, func() {
		assert.Nil(t, repo.Update(&result, rel.Inc("title")))
	})

	repo.AssertExpectations(t)
}

func TestModify_Update_setNil(t *testing.T) {
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

func TestModify_Update_map(t *testing.T) {
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

func TestModify_Update_belongsToInconsistentPk(t *testing.T) {
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

func TestModify_Update_belongsToInconsistentFk(t *testing.T) {
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

func TestModify_Update_hasOneInconsistentPk(t *testing.T) {
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

func TestModify_Update_hasOneInconsistentFk(t *testing.T) {
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

func TestModify_Update_hasManyNotLoaded(t *testing.T) {
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

func TestModify_Update_hasManyInconsistentPk(t *testing.T) {
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

func TestModify_Update_hasManyInconsistentFk(t *testing.T) {
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

func TestModify_Update_unknownField(t *testing.T) {
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

func TestModify_Update_notUnique(t *testing.T) {
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
