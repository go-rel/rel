package main

import (
	"context"
	"testing"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/join"
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
	repo.ExpectFind(where.Eq("id", 1)).Result(book)
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

func TestTable(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [table]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(rel.From("ebooks")).Result(books)
	/// [table]

	assert.Nil(t, Table(ctx, repo))
	repo.AssertExpectations(t)
}

func TestTableChained(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [table-chained]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(rel.Select("id", "title").From("ebooks")).Result(books)
	/// [table-chained]

	assert.Nil(t, TableChained(ctx, repo))
	repo.AssertExpectations(t)
}

func TestLimitOffset(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [limit-offset]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(rel.Limit(10), rel.Offset(20)).Result(books)
	/// [limit-offset]

	assert.Nil(t, LimitOffset(ctx, repo))
	repo.AssertExpectations(t)
}

func TestLimitOffsetChained(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [limit-offset-chained]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAll(rel.Select().Limit(10).Offset(20)).Result(books)
	/// [limit-offset-chained]

	assert.Nil(t, LimitOffsetChained(ctx, repo))
	repo.AssertExpectations(t)
}

func TestGroup(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [group]
	results := []struct {
		Category string
		Total    int
	}{
		{Category: "education", Total: 100},
	}
	repo.ExpectFindAll(rel.Select("category", "COUNT(id) as total").From("books").Group("category")).Result(results)
	/// [group]

	assert.Nil(t, Group(ctx, repo))
	repo.AssertExpectations(t)
}

func TestJoin(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [join]
	transactions := []Transaction{
		{ID: 1, Status: "paid"},
	}
	repo.ExpectFindAll(rel.Join("books").Where(where.Eq("books.name", "REL for Dummies"))).Result(transactions)
	/// [join]

	assert.Nil(t, Join(ctx, repo))
	repo.AssertExpectations(t)
}

func TestJoinOn(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [join-on]
	transactions := []Transaction{
		{ID: 1, Status: "paid"},
	}
	repo.ExpectFindAll(rel.JoinOn("books", "transactions.book_id", "books.id")).Result(transactions)
	/// [join-on]

	assert.Nil(t, JoinOn(ctx, repo))
	repo.AssertExpectations(t)
}

func TestJoinAlias(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [join-alias]
	transactions := []Transaction{
		{ID: 1, Status: "paid"},
	}
	repo.ExpectFindAll(join.On("books", "transactions.book_id", "books.id")).Result(transactions)
	/// [join-alias]

	assert.Nil(t, JoinAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestJoinWith(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [join-with]
	transactions := []Transaction{
		{ID: 1, Status: "paid"},
	}
	repo.ExpectFindAll(rel.JoinWith("LEFT JOIN", "books", "transactions.book_id", "books.id")).Result(transactions)
	/// [join-with]

	assert.Nil(t, JoinWith(ctx, repo))
	repo.AssertExpectations(t)
}

func TestJoinFragment(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [join-fragment]
	transactions := []Transaction{
		{ID: 1, Status: "paid"},
	}
	repo.ExpectFindAll(rel.Joinf("JOIN `books` ON `transactions`.`book_id`=`books`.`id`")).Result(transactions)
	/// [join-fragment]

	assert.Nil(t, JoinFragment(ctx, repo))
	repo.AssertExpectations(t)
}

func TestLock(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [lock]
	var book Book
	repo.ExpectFind(where.Eq("id", 1), rel.Lock("FOR UPDATE")).Result(book)
	/// [lock]

	assert.Nil(t, Lock(ctx, repo))
	repo.AssertExpectations(t)
}

func TestLockForUpdate(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [lock-for-update]
	var book Book
	repo.ExpectFind(where.Eq("id", 1), rel.ForUpdate()).Result(book)
	/// [lock-for-update]

	assert.Nil(t, LockForUpdate(ctx, repo))
	repo.AssertExpectations(t)
}

func TestLockChained(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [lock-chained]
	var book Book
	repo.ExpectFind(rel.Where(where.Eq("id", 1)).Lock("FOR UPDATE")).Result(book)
	/// [lock-chained]

	assert.Nil(t, LockChained(ctx, repo))
	repo.AssertExpectations(t)
}

func TestAggregate(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [aggregate]
	repo.ExpectAggregate(rel.From("books").Where(where.Eq("available", true)), "count", "id").Result(5)
	/// [aggregate]

	assert.Nil(t, Aggregate(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCount(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [count]
	repo.ExpectCount("books").Result(7)
	/// [count]

	assert.Nil(t, Count(ctx, repo))
	repo.AssertExpectations(t)
}

func TestCountWithCondition(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [count-with-condition]
	repo.ExpectCount("books", where.Eq("available", true)).Result(5)
	/// [count-with-condition]

	assert.Nil(t, CountWithCondition(ctx, repo))
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
