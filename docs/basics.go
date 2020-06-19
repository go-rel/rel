package main

import (
	"context"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/mysql"
	"github.com/Fs02/rel/where"
	_ "github.com/go-sql-driver/mysql"
)

// Author is a model that maps to authors table.
type Author struct {
	ID   int
	Name string
}

// Book is a model that maps to books table.
type Book struct {
	ID        int
	Title     string
	Category  string
	Price     int
	Discount  bool
	Stock     int
	AuthorID  int
	Author    Author
	CreatedAt time.Time
	UpdatedAt time.Time
}

var dsn = "root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local"

func main() {
	// initialize mysql adapter.
	adapter, _ := mysql.Open(dsn)
	defer adapter.Close()

	// initialize rel's repo.
	repo := rel.New(adapter)

	// run
	Example(context.Background(), repo)
}

// Example is an actual service function that run a complex business package.
// beware: it's actually doing nonsense here.
func Example(ctx context.Context, repo rel.Repository) error {
	var book Book

	// Quickly find a book with id 1 using short alias.
	if err := repo.Find(ctx, &book, where.Eq("id", 1)); err != nil {
		return err
	}

	// Or use chainable query builder.
	query := rel.Select().Where(where.Eq("id", 1)).Limit(1)
	if err := repo.Find(ctx, &book, query); err != nil {
		return err
	}

	// Convenient method to preload Book's Author.
	if err := repo.Preload(ctx, &book, "author"); err != nil {
		return err
	}

	// Performs updates inside a transaction.
	return repo.Transaction(ctx, func(ctx context.Context) error {
		// basic update using struct.
		book.Title = "REL for dummies"
		repo.MustUpdate(ctx, &book)

		// update only specific fields.
		repo.MustUpdate(ctx, &book, rel.Set("discount", false))

		// it even supports atomic inc/dec mutation.
		return repo.Update(ctx, &book, rel.Dec("stock"))
	})
}
