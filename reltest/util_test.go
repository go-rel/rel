package reltest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCsprint(t *testing.T) {
	assert.Equal(t, "reltest.Book{}", csprint(Book{}, true))
	assert.Equal(t, "&reltest.Book{}", csprint(&Book{}, true))
	assert.Equal(t, "[]reltest.Book{}", csprint([]Book{}, true))
	assert.Equal(t, "&[]reltest.Book{}", csprint(&[]Book{}, true))

	assert.Equal(t, "reltest.Book{Title: book}", csprint(Book{Title: "book"}, true))
	assert.Equal(t, "reltest.Book{ID: 1, Title: book}", csprint(Book{ID: 1, Title: "book"}, true))
	assert.Equal(t, "reltest.Book{Ratings: []reltest.Rating{reltest.Rating{Score: 10}}}", csprint(Book{Ratings: []Rating{{Score: 10}}}, true))

	assert.Equal(t, "[]reltest.Book{reltest.Book{}}", csprint([]Book{{}}, true))
	assert.Equal(t, "[]reltest.Book{reltest.Book{}, reltest.Book{}}", csprint([]Book{{}, {}}, true))
	assert.Equal(t, "[]reltest.Book{reltest.Book{Title: book}}", csprint([]Book{{Title: "book"}}, true))
}
