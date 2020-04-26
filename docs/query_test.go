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

func TestQueryFind(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find]
	book := Book{ID: 1, Title: "REL for dummies"}
	repo.ExpectFind(where.Eq("id", 1)).Result(book)
	/// [find]

	assert.Nil(t, QueryFind(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryFindAll(t *testing.T) {
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

	assert.Nil(t, QueryFindAll(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryCondition(t *testing.T) {
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

	assert.Nil(t, QueryCondition(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryConditionAlias(t *testing.T) {
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

	assert.Nil(t, QueryConditionAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryConditionFragment(t *testing.T) {
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

	assert.Nil(t, QueryConditionFragment(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryConditionAdvanced(t *testing.T) {
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

	assert.Nil(t, QueryConditionAdvanced(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryConditionAdvancedChain(t *testing.T) {
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

	assert.Nil(t, QueryConditionAdvancedChain(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryConditionAdvancedAlias(t *testing.T) {
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

	assert.Nil(t, QueryConditionAdvancedAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQuerySorting(t *testing.T) {
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

	assert.Nil(t, QuerySorting(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQuerySortingAlias(t *testing.T) {
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

	assert.Nil(t, QuerySortingAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQuerySortingWithCondition(t *testing.T) {
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

	assert.Nil(t, QuerySortingWithCondition(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQuerySortingWithConditionVariadic(t *testing.T) {
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

	assert.Nil(t, QuerySortingWithConditionVariadic(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQuerySelect(t *testing.T) {
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

	assert.Nil(t, QuerySelect(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryTable(t *testing.T) {
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

	assert.Nil(t, QueryTable(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryTableChained(t *testing.T) {
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

	assert.Nil(t, QueryTableChained(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryLimitOffset(t *testing.T) {
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

	assert.Nil(t, QueryLimitOffset(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryLimitOffsetChained(t *testing.T) {
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

	assert.Nil(t, QueryLimitOffsetChained(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryGroup(t *testing.T) {
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

	assert.Nil(t, QueryGroup(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryJoin(t *testing.T) {
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

	assert.Nil(t, QueryJoin(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryJoinOn(t *testing.T) {
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

	assert.Nil(t, QueryJoinOn(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryJoinAlias(t *testing.T) {
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

	assert.Nil(t, QueryJoinAlias(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryJoinWith(t *testing.T) {
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

	assert.Nil(t, QueryJoinWith(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryJoinFragment(t *testing.T) {
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

	assert.Nil(t, QueryJoinFragment(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryLock(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [lock]
	var book Book
	repo.ExpectFind(where.Eq("id", 1), rel.Lock("FOR UPDATE")).Result(book)
	/// [lock]

	assert.Nil(t, QueryLock(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryLockForUpdate(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [lock-for-update]
	var book Book
	repo.ExpectFind(where.Eq("id", 1), rel.ForUpdate()).Result(book)
	/// [lock-for-update]

	assert.Nil(t, QueryLockForUpdate(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryLockChained(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [lock-chained]
	var book Book
	repo.ExpectFind(rel.Where(where.Eq("id", 1)).Lock("FOR UPDATE")).Result(book)
	/// [lock-chained]

	assert.Nil(t, QueryLockChained(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryAggregate(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [aggregate]
	repo.ExpectAggregate(rel.From("books").Where(where.Eq("available", true)), "count", "id").Result(5)
	/// [aggregate]

	assert.Nil(t, QueryAggregate(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryCount(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [count]
	repo.ExpectCount("books").Result(7)
	/// [count]

	assert.Nil(t, QueryCount(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryCountWithCondition(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [count-with-condition]
	repo.ExpectCount("books", where.Eq("available", true)).Result(5)
	/// [count-with-condition]

	assert.Nil(t, QueryCountWithCondition(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryFindAndCountAll(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [find-and-count-all]
	books := []Book{
		{ID: 1, Title: "REL for dummies"},
	}
	repo.ExpectFindAndCountAll(rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10)).Result(books, 12)
	/// [find-and-count-all]

	assert.Nil(t, QueryFindAndCountAll(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryIteration(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [batch-iteration]
	users := make([]User, 5)
	repo.ExpectIterate(rel.From("users"), rel.BatchSize(500)).Result(users)
	/// [batch-iteration]

	assert.Nil(t, QueryIteration(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQueryIteration_error(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [batch-iteration-connection-error]
	repo.ExpectIterate(rel.From("users"), rel.BatchSize(500)).ConnectionClosed()
	/// [batch-iteration-connection-error]

	assert.Equal(t, reltest.ErrConnectionClosed, QueryIteration(ctx, repo))
	repo.AssertExpectations(t)
}

func TestQuerySQL(t *testing.T) {
	var (
		ctx  = context.TODO()
		repo = reltest.New()
	)

	/// [sql]
	var book Book
	sql := rel.SQL("SELECT id, title, price, orders = (SELECT COUNT(t.id) FROM [transactions] t WHERE t.book_id = b.id) FROM books b where b.id=?", 1)
	repo.ExpectFind(sql).Result(book)
	/// [sql]

	assert.Nil(t, QuerySQL(ctx, repo))
	repo.AssertExpectations(t)
}
