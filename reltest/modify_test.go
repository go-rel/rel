package reltest

import (
	"context"
	"errors"
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestModify_Insert(t *testing.T) {
	var (
		repo   = New()
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectInsert()
	assert.Nil(t, repo.Insert(context.TODO(), &result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert()
	assert.NotPanics(t, func() {
		repo.MustInsert(context.TODO(), &result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Insert_nested(t *testing.T) {
	var (
		repo     = New()
		authorID = 1
		result   = Book{
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
			AuthorID: &authorID,
			Ratings: []Rating{
				{ID: 1, Score: 9, BookID: 1},
				{ID: 2, Score: 10, BookID: 1},
			},
			Poster: Poster{ID: 1, BookID: 1, Image: "http://image.url"},
		}
	)

	repo.ExpectInsert()
	assert.Nil(t, repo.Insert(context.TODO(), &result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert()
	assert.NotPanics(t, func() {
		repo.MustInsert(context.TODO(), &result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Insert_record(t *testing.T) {
	var (
		repo   = New()
		result = Book{Title: "Golang for dummies"}
		book   = Book{ID: 1, Title: "Golang for dummies"}
	)

	repo.ExpectInsert().ForType("reltest.Book")
	assert.Nil(t, repo.Insert(context.TODO(), &result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert().ForType("reltest.Book")
	assert.NotPanics(t, func() {
		repo.MustInsert(context.TODO(), &result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Insert_set(t *testing.T) {
	var (
		repo   = New()
		result Book
		book   = Book{ID: 1, Title: "Rel for dummies"}
	)

	repo.ExpectInsert(rel.Set("title", "Rel for dummies"))
	assert.Nil(t, repo.Insert(context.TODO(), &result, rel.Set("title", "Rel for dummies")))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert(rel.Set("title", "Rel for dummies"))
	assert.NotPanics(t, func() {
		repo.MustInsert(context.TODO(), &result, rel.Set("title", "Rel for dummies"))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Insert_map(t *testing.T) {
	var (
		repo     = New()
		result   Book
		authorID = 1
		book     = Book{
			ID:       1,
			Title:    "Rel for dummies",
			Author:   Author{ID: 1, Name: "Kia"},
			AuthorID: &authorID,
			Ratings: []Rating{
				{ID: 1, Score: 9, BookID: 1},
				{ID: 2, Score: 10, BookID: 1},
			},
			Poster: Poster{ID: 1, BookID: 1, Image: "http://image.url"},
		}
		mod = rel.Map{
			"title": "Rel for dummies",
			"author": rel.Map{
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

	repo.ExpectInsert(mod)
	assert.Nil(t, repo.Insert(context.TODO(), &result, mod))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectInsert(mod)
	assert.NotPanics(t, func() {
		repo.MustInsert(context.TODO(), &result, mod)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Insert_unknownField(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectInsert(rel.Set("titles", "Rel for dummies"))
	assert.Panics(t, func() {
		_ = repo.Insert(context.TODO(), &result, rel.Set("titles", "Rel for dummies"))
	})
	repo.AssertExpectations(t)
}

func TestModify_Insert_notUnique(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectInsert(rel.Set("title", "Rel for dummies")).NotUnique("title")
	assert.Equal(t,
		rel.ConstraintError{Key: "title", Type: rel.UniqueConstraint},
		repo.Insert(context.TODO(), &result, rel.Set("title", "Rel for dummies")),
	)
	repo.AssertExpectations(t)
}

func TestModify_InsertAll(t *testing.T) {
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

func TestModify_Update(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate()
	assert.Nil(t, repo.Update(context.TODO(), &result))
	repo.AssertExpectations(t)

	repo.ExpectUpdate()
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Update_nested(t *testing.T) {
	var (
		repo     = New()
		authorID = 2
		result   = Book{
			ID:       2,
			Title:    "Rel for dummies",
			Author:   Author{ID: 2, Name: "Kia"},
			AuthorID: &authorID,
			Ratings: []Rating{
				{ID: 1, BookID: 2, Score: 9},
				{ID: 2, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
		book = Book{
			ID:       2,
			Title:    "Rel for dummies",
			Author:   Author{ID: 2, Name: "Kia"},
			AuthorID: &authorID,
			Ratings: []Rating{
				{ID: 1, BookID: 2, Score: 9},
				{ID: 2, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
	)

	repo.ExpectUpdate()
	assert.Nil(t, repo.Update(context.TODO(), &result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate()
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Update_nestedInsert(t *testing.T) {
	var (
		repo     = New()
		authorID = 1
		result   = Book{
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
			AuthorID: &authorID,
			Ratings: []Rating{
				{ID: 1, BookID: 2, Score: 9},
				{ID: 2, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
	)

	repo.ExpectUpdate()
	assert.Nil(t, repo.Update(context.TODO(), &result))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate()
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Update_record(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate().For(&result)
	assert.Nil(t, repo.Update(context.TODO(), &result))
	repo.AssertExpectations(t)

	repo.ExpectUpdate().For(&result)
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result)
	})
	repo.AssertExpectations(t)
}

// func TestModify_Update_withoutPrimaryValue(t *testing.T) {
// 	var (
// 		repo   = New()
// 		result = Book{Title: "Golang for dummies"}
// 	)

// 	repo.ExpectUpdate().For(&result)
// 	assert.Panics(t, func() {
// 		_ = repo.Update(context.TODO(),&result)
// 	})
// 	repo.AssertExpectations(t)
// }

func TestModify_Update_set(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
		book   = Book{ID: 2, Title: "Rel for dummies"}
	)

	repo.ExpectUpdate(rel.Set("title", "Rel for dummies"))
	assert.Nil(t, repo.Update(context.TODO(), &result, rel.Set("title", "Rel for dummies")))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate(rel.Set("title", "Rel for dummies"))
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result, rel.Set("title", "Rel for dummies"))
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

// func TestModify_Update_inc(t *testing.T) {
// 	var (
// 		repo   = New()
// 		result = Book{ID: 2, Views: 10}
// 		book   = Book{ID: 2, Views: 11}
// 	)

// 	repo.ExpectUpdate(rel.Inc("views"))
// 	assert.Nil(t, repo.Update(context.TODO(),&result, rel.Inc("views")))
// 	assert.Equal(t, book, result)
// 	repo.AssertExpectations(t)
// }

// func TestModify_Update_dec(t *testing.T) {
// 	var (
// 		repo   = New()
// 		result = Book{ID: 2, Views: 10}
// 		book   = Book{ID: 2, Views: 9}
// 	)

// 	repo.ExpectUpdate(rel.Dec("views"))
// 	assert.Nil(t, repo.Update(context.TODO(),&result, rel.Dec("views")))
// 	assert.Equal(t, book, result)
// 	repo.AssertExpectations(t)
// }

func TestModify_Update_incOrDecFieldNotExists(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Views: 10}
	)

	repo.ExpectUpdate(rel.Inc("watistis"))
	assert.Panics(t, func() {
		assert.Nil(t, repo.Update(context.TODO(), &result, rel.Inc("watistis")))
	})

	repo.AssertExpectations(t)
}

func TestModify_Update_incOrDecFieldInvalid(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Views: 10}
	)

	repo.ExpectUpdate(rel.Inc("title"))
	assert.Panics(t, func() {
		assert.Nil(t, repo.Update(context.TODO(), &result, rel.Inc("title")))
	})

	repo.AssertExpectations(t)
}

func TestModify_Update_setNil(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
		book   = Book{ID: 2, Title: ""}
	)

	repo.ExpectUpdate(rel.Set("title", nil))
	assert.Nil(t, repo.Update(context.TODO(), &result, rel.Set("title", nil)))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)
}

func TestModify_Update_map(t *testing.T) {
	var (
		repo     = New()
		authorID = 2
		result   = Book{
			ID:       2,
			Title:    "Golang for dummies",
			Author:   Author{ID: 2, Name: "unknown"},
			AuthorID: &authorID,
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
			AuthorID: &authorID,
			Ratings: []Rating{
				{ID: 2, BookID: 2, Score: 9},
				{ID: 1, BookID: 2, Score: 10},
			},
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
		mod = rel.Map{
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

	repo.ExpectUpdate(mod)
	assert.Nil(t, repo.Update(context.TODO(), &result, mod))
	assert.Equal(t, book, result)
	repo.AssertExpectations(t)

	repo.ExpectUpdate(mod)
	assert.NotPanics(t, func() {
		repo.MustUpdate(context.TODO(), &result, mod)
		assert.Equal(t, book, result)
	})
	repo.AssertExpectations(t)
}

func TestModify_Update_belongsToInconsistentFk(t *testing.T) {
	var (
		repo     = New()
		authorID = 1
		result   = Book{
			ID:       2,
			Title:    "Golang for dummies",
			AuthorID: &authorID,
			Author:   Author{ID: 2, Name: "Kia"},
		}
		mod = rel.Map{
			"author": rel.Map{
				"id":   2,
				"name": "Koa",
			},
		}
	)

	repo.ExpectUpdate(mod)
	assert.Equal(t, rel.ConstraintError{
		Key:  "author_id",
		Type: rel.ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent belongs to ref and fk"),
	}, repo.Update(context.TODO(), &result, mod))
	repo.AssertExpectations(t)
}

func TestModify_Update_hasOneInconsistentPk(t *testing.T) {
	var (
		repo   = New()
		result = Book{
			ID:     2,
			Title:  "Golang for dummies",
			Poster: Poster{ID: 1, BookID: 2, Image: "http://image.url"},
		}
		mod = rel.Map{
			"poster": rel.Map{
				"id":    2,
				"image": "http://image.url/other",
			},
		}
	)

	repo.ExpectUpdate(mod)
	assert.Panics(t, func() {
		_ = repo.Update(context.TODO(), &result, mod)
	})
	repo.AssertExpectations(t)
}

func TestModify_Update_hasOneInconsistentFk(t *testing.T) {
	var (
		repo   = New()
		result = Book{
			ID:     2,
			Title:  "Golang for dummies",
			Poster: Poster{ID: 1, BookID: 1, Image: "http://image.url"},
		}
		mod = rel.Map{
			"poster": rel.Map{
				"id":    1,
				"image": "http://image.url/other",
			},
		}
	)

	repo.ExpectUpdate(mod)
	assert.Equal(t, rel.ConstraintError{
		Key:  "book_id",
		Type: rel.ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent has one ref and fk"),
	}, repo.Update(context.TODO(), &result, mod))
	repo.AssertExpectations(t)
}

func TestModify_Update_hasManyInconsistentFk(t *testing.T) {
	var (
		repo   = New()
		result = Book{
			ID:    2,
			Title: "Golang for dummies",
			Ratings: []Rating{
				{ID: 2, BookID: 1, Score: 5},
			},
		}
		mod = rel.Map{
			"ratings": []rel.Map{
				{"id": 2, "score": 9},
			},
		}
	)

	repo.ExpectUpdate(mod)
	assert.Equal(t, rel.ConstraintError{
		Key:  "book_id",
		Type: rel.ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent has many ref and fk"),
	}, repo.Update(context.TODO(), &result, mod))
	repo.AssertExpectations(t)
}

func TestModify_Update_unknownField(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate(rel.Set("titles", "Rel for dummies"))
	assert.Panics(t, func() {
		_ = repo.Update(context.TODO(), &result, rel.Set("titles", "Rel for dummies"))
	})
	repo.AssertExpectations(t)
}

func TestModify_Update_notUnique(t *testing.T) {
	var (
		repo   = New()
		result = Book{ID: 2, Title: "Golang for dummies"}
	)

	repo.ExpectUpdate(rel.Set("title", "Rel for dummies")).NotUnique("title")
	assert.Equal(t,
		rel.ConstraintError{Key: "title", Type: rel.UniqueConstraint},
		repo.Update(context.TODO(), &result, rel.Set("title", "Rel for dummies")),
	)
	repo.AssertExpectations(t)
}
