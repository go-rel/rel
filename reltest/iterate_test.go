package reltest

import (
	"context"
	"io"
	"testing"

	"github.com/go-rel/rel"
	"github.com/stretchr/testify/assert"
)

func TestIterate(t *testing.T) {
	tests := []struct {
		name   string
		result interface{}
		count  int
	}{
		{
			name:   "struct",
			result: Book{ID: 1},
			count:  1,
		},
		{
			name:   "struct pointer",
			result: &Book{ID: 1},
			count:  1,
		},
		{
			name:   "slice",
			result: []Book{{ID: 1}, {ID: 2}},
			count:  2,
		},
		{
			name:   "slice pointer",
			result: &[]Book{{ID: 1}, {ID: 2}},
			count:  2,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				book  Book
				repo  = New()
				query = rel.From("users")
			)

			repo.ExpectIterate(query, rel.BatchSize(500)).Result(test.result)

			var (
				count = 0
				it    = repo.Iterate(context.TODO(), query, rel.BatchSize(500))
			)

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

			assert.Equal(t, test.count, count)
			repo.AssertExpectations(t)
		})
	}
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

func TestIterate_assert(t *testing.T) {
	var (
		repo = New()
	)

	repo.ExpectIterate(rel.From("users"))

	assert.Panics(t, func() {
		repo.Iterate(context.TODO(), rel.From("books"))
	})
	assert.False(t, repo.AssertExpectations(nt))
	assert.Equal(t, "FAIL: The code you are testing needs to call:\n\tIterate(ctx, query todo)", nt.lastLog)
}

func TestIterate_String(t *testing.T) {
	var (
		mockIterate = MockIterate{assert: &Assert{}, argQuery: rel.From("users")}
	)

	assert.Equal(t, "Iterate(ctx, query todo)", mockIterate.String())
	assert.Equal(t, "ExpectIterate(query todo)", mockIterate.ExpectString())
}
