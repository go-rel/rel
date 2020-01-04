package main

import (
	"testing"

	"github.com/Fs02/rel/reltest"
	"github.com/Fs02/rel/where"
)

func TestExample(t *testing.T) {
	// create a mocked repository.
	repo := reltest.New()

	// mocks insert
	repo.ExpectInsert()

	// mock find and return other result
	repo.ExpectFind(where.Eq("id", 1)).Result(Book{
		ID:       1,
		Title:    "Go for dummies",
		Category: "learning",
	})

	// run
	Example(repo)

	// asserts
	repo.AssertExpectations(t)
}
