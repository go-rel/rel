package main

import (
	"context"
	"io"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/join"
	"github.com/Fs02/rel/sort"
	"github.com/Fs02/rel/where"
)

// QueryFind docs example.
func QueryFind(ctx context.Context, repo rel.Repository) error {
	/// [find]
	var book Book
	err := repo.Find(ctx, &book, where.Eq("id", 1))
	/// [find]

	return err
}

// QueryFindAll docs example.
func QueryFindAll(ctx context.Context, repo rel.Repository) error {
	/// [find-all]
	var books []Book
	err := repo.FindAll(ctx, &books)
	/// [find-all]

	return err
}

// QueryCondition docs example.
func QueryCondition(ctx context.Context, repo rel.Repository) error {
	/// [condition]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Eq("available", true))
	/// [condition]

	return err
}

// QueryConditionAlias docs example.
func QueryConditionAlias(ctx context.Context, repo rel.Repository) error {
	/// [condition-alias]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Eq("available", true))
	/// [condition-alias]

	return err
}

// QueryConditionFragment docs example.
func QueryConditionFragment(ctx context.Context, repo rel.Repository) error {
	/// [condition-fragment]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Fragment("available=?", true))
	/// [condition-fragment]

	return err
}

// QueryConditionAdvanced docs example.
func QueryConditionAdvanced(ctx context.Context, repo rel.Repository) error {
	/// [condition-advanced]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.And(rel.Eq("available", true), rel.Or(rel.Gte("price", 100), rel.Eq("discount", true))))
	/// [condition-advanced]

	return err
}

// QueryConditionAdvancedChain docs example.
func QueryConditionAdvancedChain(ctx context.Context, repo rel.Repository) error {
	/// [condition-advanced-chain]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Eq("available", true).And(rel.Gte("price", 100).OrEq("discount", true)))
	/// [condition-advanced-chain]

	return err
}

// QueryConditionAdvancedAlias docs example.
func QueryConditionAdvancedAlias(ctx context.Context, repo rel.Repository) error {
	/// [condition-advanced-alias]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Eq("available", true).And(where.Gte("price", 100).OrEq("discount", true)))
	/// [condition-advanced-alias]

	return err
}

// QuerySorting docs example.
func QuerySorting(ctx context.Context, repo rel.Repository) error {
	/// [sorting]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.NewSortAsc("updated_at"))
	/// [sorting]

	return err
}

// QuerySortingAlias docs example.
func QuerySortingAlias(ctx context.Context, repo rel.Repository) error {
	/// [sorting-alias]
	var books []Book
	err := repo.FindAll(ctx, &books, sort.Asc("updated_at"))
	/// [sorting-alias]

	return err
}

// QuerySortingWithCondition docs example.
func QuerySortingWithCondition(ctx context.Context, repo rel.Repository) error {
	/// [sorting-with-condition]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Where(where.Eq("available", true)).SortAsc("updated_at"))
	/// [sorting-with-condition]

	return err
}

// QuerySortingWithConditionVariadic docs example.
func QuerySortingWithConditionVariadic(ctx context.Context, repo rel.Repository) error {
	/// [sorting-with-condition-variadic]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Eq("available", true), sort.Asc("updated_at"))
	/// [sorting-with-condition-variadic]

	return err
}

// QuerySelect docs example.
func QuerySelect(ctx context.Context, repo rel.Repository) error {
	/// [select]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Select("id", "title"))
	/// [select]

	return err
}

// QueryTable docs example.
func QueryTable(ctx context.Context, repo rel.Repository) error {
	/// [table]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.From("ebooks"))
	/// [table]

	return err
}

// QueryTableChained docs example.
func QueryTableChained(ctx context.Context, repo rel.Repository) error {
	/// [table-chained]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Select("id", "title").From("ebooks"))
	/// [table-chained]

	return err
}

// QueryLimitOffset docs example.
func QueryLimitOffset(ctx context.Context, repo rel.Repository) error {
	/// [limit-offset]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Limit(10), rel.Offset(20))
	/// [limit-offset]

	return err
}

// QueryLimitOffsetChained docs example.
func QueryLimitOffsetChained(ctx context.Context, repo rel.Repository) error {
	/// [limit-offset-chained]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Select().Limit(10).Offset(20))
	/// [limit-offset-chained]

	return err
}

// QueryGroup docs example.
func QueryGroup(ctx context.Context, repo rel.Repository) error {
	/// [group]
	// custom struct to store the result.
	var results []struct {
		Category string
		Total    int
	}

	// we need to explicitly specify table name since we are using an anonymous struct.
	err := repo.FindAll(ctx, &results, rel.Select("category", "COUNT(id) as total").From("books").Group("category"))
	/// [group]

	return err
}

// QueryJoin docs example.
func QueryJoin(ctx context.Context, repo rel.Repository) error {
	/// [join]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.Join("books").Where(where.Eq("books.name", "REL for Dummies")))
	/// [join]

	return err
}

// QueryJoinOn docs example.
func QueryJoinOn(ctx context.Context, repo rel.Repository) error {
	/// [join-on]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.JoinOn("books", "transactions.book_id", "books.id"))
	/// [join-on]

	return err
}

// QueryJoinAlias docs example.
func QueryJoinAlias(ctx context.Context, repo rel.Repository) error {
	/// [join-alias]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, join.On("books", "transactions.book_id", "books.id"))
	/// [join-alias]

	return err
}

// QueryJoinWith docs example.
func QueryJoinWith(ctx context.Context, repo rel.Repository) error {
	/// [join-with]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.JoinWith("LEFT JOIN", "books", "transactions.book_id", "books.id"))
	/// [join-with]

	return err
}

// QueryJoinFragment docs example.
func QueryJoinFragment(ctx context.Context, repo rel.Repository) error {
	/// [join-fragment]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.Joinf("JOIN `books` ON `transactions`.`book_id`=`books`.`id`"))
	/// [join-fragment]

	return err
}

// QueryLock docs example.
func QueryLock(ctx context.Context, repo rel.Repository) error {
	/// [lock]
	var book Book
	err := repo.Find(ctx, &book, where.Eq("id", 1), rel.Lock("FOR UPDATE"))
	/// [lock]

	return err
}

// QueryLockForUpdate docs example.
func QueryLockForUpdate(ctx context.Context, repo rel.Repository) error {
	/// [lock-for-update]
	var book Book
	err := repo.Find(ctx, &book, where.Eq("id", 1), rel.ForUpdate())
	/// [lock-for-update]

	return err
}

// QueryLockChained docs example.
func QueryLockChained(ctx context.Context, repo rel.Repository) error {
	/// [lock-chained]
	var book Book
	err := repo.Find(ctx, &book, rel.Where(where.Eq("id", 1)).Lock("FOR UPDATE"))
	/// [lock-chained]

	return err
}

// QueryAggregate docs example.
func QueryAggregate(ctx context.Context, repo rel.Repository) error {
	/// [aggregate]
	count, err := repo.Aggregate(ctx, rel.From("books").Where(where.Eq("available", true)), "count", "id")
	/// [aggregate]

	_ = count
	return err
}

// QueryCount docs example.
func QueryCount(ctx context.Context, repo rel.Repository) error {
	/// [count]
	count, err := repo.Count(ctx, "books")
	/// [count]

	_ = count
	return err
}

// QueryCountWithCondition docs example.
func QueryCountWithCondition(ctx context.Context, repo rel.Repository) error {
	/// [count-with-condition]
	count, err := repo.Count(ctx, "books", where.Eq("available", true))
	/// [count-with-condition]

	_ = count
	return err
}

// QueryFindAndCountAll docs example.
func QueryFindAndCountAll(ctx context.Context, repo rel.Repository) error {
	/// [find-and-count-all]
	var books []Book
	count, err := repo.FindAndCountAll(ctx, &books, rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10))
	/// [find-and-count-all]

	_ = count
	return err
}

// SendPromotionEmail tp demonstrate Iteration.
func SendPromotionEmail(*User) {}

// QueryIteration docs example.
func QueryIteration(ctx context.Context, repo rel.Repository) error {
	/// [batch-iteration]
	var (
		user User
		iter = repo.Iterate(ctx, rel.From("users"), rel.BatchSize(500))
	)

	// make sure iterator is closed after process is finish.
	defer iter.Close()
	for {
		// retrieve next user.
		if err := iter.Next(&user); err != nil {
			if err == io.EOF {
				break
			}

			// handle error
			return err
		}

		// process user
		SendPromotionEmail(&user)
	}
	/// [batch-iteration]

	return nil
}

// QuerySQL natively.
func QuerySQL(ctx context.Context, repo rel.Repository) error {
	/// [sql]
	var book Book
	sql := rel.SQL("SELECT id, title, price, orders = (SELECT COUNT(t.id) FROM [transactions] t WHERE t.book_id = b.id) FROM books b where b.id=?", 1)
	err := repo.Find(ctx, &book, sql)
	/// [sql]

	return err
}
