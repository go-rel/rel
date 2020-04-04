package reltest

import (
	"context"
	"io"
	"testing"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	var (
		book  Book
		repo  = New()
		query = rel.From("users")
	)

	repo.ExpectIterate(query).Result([]Book{{ID: 1}, {ID: 2}})

	count := 0
	it := repo.Iterate(context.TODO(), query)
	defer it.Close()
	for {
		if err := it.Next(&book); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
		}

		assert.NotEqual(t, 0, book.ID)
		count++
	}

	assert.Equal(t, 2, count)
	repo.AssertExpectations(t)
}

func TestIterate_single(t *testing.T) {
	var (
		book  Book
		repo  = New()
		query = rel.From("users")
	)

	repo.ExpectIterate(query).Result(Book{ID: 1})

	count := 0
	it := repo.Iterate(context.TODO(), query)
	defer it.Close()
	for {
		if err := it.Next(&book); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
		}

		assert.NotEqual(t, 0, book.ID)
		count++
	}

	assert.Equal(t, 1, count)
	repo.AssertExpectations(t)
}

func TestIterate_error(t *testing.T) {
	var (
		book  Book
		repo  = New()
		query = rel.From("users")
	)

	repo.ExpectIterate(query).ConnectionClosed()

	it := repo.Iterate(context.TODO(), query)
	defer it.Close()

	assert.Equal(t, ErrConnectionClosed, it.Next(&book))
	repo.AssertExpectations(t)
}
