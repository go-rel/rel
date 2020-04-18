package main

import (
	"context"
	"io"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/join"
	"github.com/Fs02/rel/sort"
	"github.com/Fs02/rel/where"
)

// Find docs example.
func Find(ctx context.Context, repo rel.Repository) error {
	/// [find]
	var book Book
	err := repo.Find(ctx, &book)
	/// [find]

	return err
}

// FindAll docs example.
func FindAll(ctx context.Context, repo rel.Repository) error {
	/// [find-all]
	var books []Book
	err := repo.FindAll(ctx, &books)
	/// [find-all]

	return err
}

// Condition docs example.
func Condition(ctx context.Context, repo rel.Repository) error {
	/// [condition]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Eq("available", true))
	/// [condition]

	return err
}

// ConditionAlias docs example.
func ConditionAlias(ctx context.Context, repo rel.Repository) error {
	/// [condition-alias]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Eq("available", true))
	/// [condition-alias]

	return err
}

// ConditionFragment docs example.
func ConditionFragment(ctx context.Context, repo rel.Repository) error {
	/// [condition-fragment]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Fragment("available=?", true))
	/// [condition-fragment]

	return err
}

// ConditionAdvanced docs example.
func ConditionAdvanced(ctx context.Context, repo rel.Repository) error {
	/// [condition-advanced]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.And(rel.Eq("available", true), rel.Or(rel.Gte("price", 100), rel.Eq("discount", true))))
	/// [condition-advanced]

	return err
}

// ConditionAdvancedChain docs example.
func ConditionAdvancedChain(ctx context.Context, repo rel.Repository) error {
	/// [condition-advanced-chain]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Eq("available", true).And(rel.Gte("price", 100).OrEq("discount", true)))
	/// [condition-advanced-chain]

	return err
}

// ConditionAdvancedAlias docs example.
func ConditionAdvancedAlias(ctx context.Context, repo rel.Repository) error {
	/// [condition-advanced-alias]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Eq("available", true).And(where.Gte("price", 100).OrEq("discount", true)))
	/// [condition-advanced-alias]

	return err
}

// Sorting docs example.
func Sorting(ctx context.Context, repo rel.Repository) error {
	/// [sorting]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.NewSortAsc("updated_at"))
	/// [sorting]

	return err
}

// SortingAlias docs example.
func SortingAlias(ctx context.Context, repo rel.Repository) error {
	/// [sorting-alias]
	var books []Book
	err := repo.FindAll(ctx, &books, sort.Asc("updated_at"))
	/// [sorting-alias]

	return err
}

// SortingWithCondition docs example.
func SortingWithCondition(ctx context.Context, repo rel.Repository) error {
	/// [sorting-with-condition]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Where(where.Eq("available", true)).SortAsc("updated_at"))
	/// [sorting-with-condition]

	return err
}

// SortingWithConditionVariadic docs example.
func SortingWithConditionVariadic(ctx context.Context, repo rel.Repository) error {
	/// [sorting-with-condition-variadic]
	var books []Book
	err := repo.FindAll(ctx, &books, where.Eq("available", true), sort.Asc("updated_at"))
	/// [sorting-with-condition-variadic]

	return err
}

// Select docs example.
func Select(ctx context.Context, repo rel.Repository) error {
	/// [select]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Select("id", "title"))
	/// [select]

	return err
}

// Table docs example.
func Table(ctx context.Context, repo rel.Repository) error {
	/// [table]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.From("ebooks"))
	/// [table]

	return err
}

// TableChained docs example.
func TableChained(ctx context.Context, repo rel.Repository) error {
	/// [table-chained]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Select("id", "title").From("ebooks"))
	/// [table-chained]

	return err
}

// LimitOffset docs example.
func LimitOffset(ctx context.Context, repo rel.Repository) error {
	/// [limit-offset]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Limit(10), rel.Offset(20))
	/// [limit-offset]

	return err
}

// LimitOffsetChained docs example.
func LimitOffsetChained(ctx context.Context, repo rel.Repository) error {
	/// [limit-offset-chained]
	var books []Book
	err := repo.FindAll(ctx, &books, rel.Select().Limit(10).Offset(20))
	/// [limit-offset-chained]

	return err
}

// Group docs example.
func Group(ctx context.Context, repo rel.Repository) error {
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

// Join docs example.
func Join(ctx context.Context, repo rel.Repository) error {
	/// [join]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.Join("books").Where(where.Eq("books.name", "REL for Dummies")))
	/// [join]

	return err
}

// JoinOn docs example.
func JoinOn(ctx context.Context, repo rel.Repository) error {
	/// [join-on]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.JoinOn("books", "transactions.book_id", "books.id"))
	/// [join-on]

	return err
}

// JoinAlias docs example.
func JoinAlias(ctx context.Context, repo rel.Repository) error {
	/// [join-alias]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, join.On("books", "transactions.book_id", "books.id"))
	/// [join-alias]

	return err
}

// JoinWith docs example.
func JoinWith(ctx context.Context, repo rel.Repository) error {
	/// [join-with]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.JoinWith("LEFT JOIN", "books", "transactions.book_id", "books.id"))
	/// [join-with]

	return err
}

// JoinFragment docs example.
func JoinFragment(ctx context.Context, repo rel.Repository) error {
	/// [join-fragment]
	var transactions []Transaction
	err := repo.FindAll(ctx, &transactions, rel.Joinf("JOIN `books` ON `transactions`.`book_id`=`books`.`id`"))
	/// [join-fragment]

	return err
}

// SendPromotionEmail tp demonstrate Iteration.
func SendPromotionEmail(*User) {}

// Iteration docs example.
func Iteration(ctx context.Context, repo rel.Repository) error {
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
