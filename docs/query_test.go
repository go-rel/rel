package main

import (
	"context"
	"testing"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/reltest"
	"github.com/Fs02/rel/sort"
	"github.com/Fs02/rel/where"
	"github.com/stretchr/testify/assert"
)

func TestFind(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find]
	book := Book{ID: 1, Title: "REL for dummies"}
	repo.ExpectFind().Result(book)
	/// [find]

	assert.Nil(t, Find(ctx, repo))
	repo.AssertExpectations(t)
}

func TestFindAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find-all]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll().Result(books)
	/// [find-all]

	assert.Nil(t, FindAll(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCondition(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [condition]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(rel.Eq("available", true)).Result(books)
	/// [condition]

	assert.Nil(t, Condition(ctx, repo))
	repo.AssertExpectations(t)
}

func TestConditionAlias(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [condition-alias]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(where.Eq("available", true)).Result(books)
	/// [condition-alias]

	assert.Nil(t, ConditionAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestConditionFragment(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [condition-fragment]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(where.Fragment("available=?", true)).Result(books)
	/// [condition-fragment]

	assert.Nil(t, ConditionFragment(ctx, repo))
	repo.AssertExpectations(t)
}

func TestConditionAdvanced(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [condition-advanced]
	books := []Book{
		{ID: 1, Title: "REL for dummies", Price: 100},
		{ID: 2, Title: "REL for dummies", Price: 50, Discount: true},
	}
	repo.ExpectFindAll(rel.And(rel.Eq("available", true), rel.Or(rel.Gte("price", 100), rel.Eq("discount", true)))).Result(books)
	/// [condition-advanced]

	assert.Nil(t, ConditionAdvanced(ctx, repo))
	repo.AssertExpectations(t)
}

func TestConditionAdvancedChain(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [condition-advanced-chain]
	books := []Book{
		{ID: 1, Title: "REL for dummies", Price: 100},
		{ID: 2, Title: "REL for dummies", Price: 50, Discount: true},
	}
	repo.ExpectFindAll(rel.Eq("available", true).And(rel.Gte("price", 100).OrEq("discount", true))).Result(books)
	/// [condition-advanced-chain]

	assert.Nil(t, ConditionAdvancedChain(ctx, repo))
	repo.AssertExpectations(t)
}

func TestConditionAdvancedAlias(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [condition-advanced-alias]
	books := []Book{
		{ID: 1, Title: "REL for dummies", Price: 100},
		{ID: 2, Title: "REL for dummies", Price: 50, Discount: true},
	}
	repo.ExpectFindAll(where.Eq("available", true).And(where.Gte("price", 100).OrEq("discount", true))).Result(books)
	/// [condition-advanced-alias]

	assert.Nil(t, ConditionAdvancedAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestSorting(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [sorting]
	books := []Book{
		{ID: 1, Title: "REL for dummies", UpdatedAt: time.Now()},
	}
	repo.ExpectFindAll(rel.NewSortAsc("updated_at")).Result(books)
	/// [sorting]

	assert.Nil(t, Sorting(ctx, repo))
	repo.AssertExpectations(t)
}

func TestSortingAlias(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [sorting-alias]
	books := []Book{
		{ID: 1, Title: "REL for dummies", UpdatedAt: time.Now()},
	}
	repo.ExpectFindAll(sort.Asc("updated_at")).Result(books)
	/// [sorting-alias]

	assert.Nil(t, SortingAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestSortingWithCondition(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [sorting-with-condition]
	books := []Book{
		{ID: 1, Title: "REL for dummies", UpdatedAt: time.Now()},
	}
	repo.ExpectFindAll(rel.Where(where.Eq("available", true)).SortAsc("updated_at")).Result(books)
	/// [sorting-with-condition]

	assert.Nil(t, SortingWithCondition(ctx, repo))
	repo.AssertExpectations(t)
}

func TestSortingWithConditionVariadic(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [sorting-with-condition-variadic]
	books := []Book{
		{ID: 1, Title: "REL for dummies", UpdatedAt: time.Now()},
	}
	repo.ExpectFindAll(where.Eq("available", true), sort.Asc("updated_at")).Result(books)
	/// [sorting-with-condition-variadic]

	assert.Nil(t, SortingWithConditionVariadic(ctx, repo))
	repo.AssertExpectations(t)
}

func TestSelect(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [select]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(rel.Select("id", "title")).Result(books)
	/// [select]

	assert.Nil(t, Select(ctx, repo))
	repo.AssertExpectations(t)
}

func TestIteration(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [batch-iteration]
	users := make([]User, 5)
	repo.ExpectIterate(rel.From("users"), rel.BatchSize(500)).Result(users)
	/// [batch-iteration]

	assert.Nil(t, Iteration(ctx, repo))
	repo.AssertExpectations(t)
}

func TestIteration_error(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [batch-iteration-connection-error]
	repo.ExpectIterate(rel.From("users"), rel.BatchSize(500)).ConnectionClosed()
	/// [batch-iteration-connection-error]

	assert.Equal(t, reltest.ErrConnectionClosed, Iteration(ctx, repo))
	repo.AssertExpectations(t)
}
