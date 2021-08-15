package reltest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchContains(t *testing.T) {
	assert.True(t, matchContains(&Book{ID: 1}, &Book{ID: 1}))
	assert.True(t, matchContains(Book{ID: 1}, Book{ID: 1}))
	assert.True(t, matchContains(Book{ID: 1}, &Book{ID: 1}))
	assert.True(t, matchContains(&Book{ID: 1}, Book{ID: 1}))
	assert.True(t, matchContains(Book{}, Book{ID: 1, Title: "book"}))
	assert.True(t, matchContains(Book{ID: 1}, Book{ID: 1, Title: "book"}))
	assert.True(t, matchContains(Book{Title: "book"}, Book{ID: 1, Title: "book"}))
	assert.False(t, matchContains(Book{ID: 2}, Book{ID: 1, Title: "book"}))
	assert.False(t, matchContains(Book{Title: "paper"}, Book{ID: 1, Title: "book"}))
}
