package main

import (
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/mysql"
	"github.com/Fs02/rel/where"
	_ "github.com/go-sql-driver/mysql"
)

// Book is a model that maps to books table.
type Book struct {
	ID        int
	Title     string
	Category  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	// initialize mysql adapter.
	adapter, err := mysql.Open("root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize rel's repo.
	repo := rel.New(adapter)

	Example(repo)
}

// Example is an actual service function that run a complex business package.
// beware: it's actually doing nonsense here.
func Example(repo rel.Repository) {
	// Inserting Books.
	book := Book{
		Title:    "rel for dummies",
		Category: "learning",
	}
	repo.Insert(&book)

	// Querying Books.
	// Find a book with id 1.
	repo.Find(&book, where.Eq("id", 1))
}
