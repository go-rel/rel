## Quick Start

### Adapters

Rel uses adapter in order to generate and execute query to a database, below is the list of available adapters supported by rel out of the box.

| Adapter    | Package                              | Godoc                                                                                                                                 |
|------------|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| MySQL      | github.com/Fs02/rel/adapter/mysql    | [![GoDoc](https://godoc.org/github.com/Fs02/rel/adapter/mysql?status.svg)](https://godoc.org/github.com/Fs02/rel/adapter/mysql)       |
| PostgreSQL | github.com/Fs02/rel/adapter/postgres | [![GoDoc](https://godoc.org/github.com/Fs02/rel/adapter/postgres?status.svg)](https://godoc.org/github.com/Fs02/rel/adapter/postgres) |
| SQLite3    | github.com/Fs02/rel/adapter/sqlite3  | [![GoDoc](https://godoc.org/github.com/Fs02/rel/adapter/sqlite3?status.svg)](https://godoc.org/github.com/Fs02/rel/adapter/sqlite3)   |

### Basic Usage

Below is a very basic example on how to utilize rel using mysql adapter.

```golang
// main.go
package main

import (
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/mysql"
	"github.com/Fs02/rel/where"
	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID        int
	Name      string
	Price     int
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

	run(repo)
}

// run is an actual service function that run a complex business package.
// beware: it's actually doing nonsense here.
func run(repo *rel.Repository) {
	// Inserting Products.
	product := Product{
		Name: "shampoo",
		Price: 1000,
	}
	repo.Insert(&product)

	// Querying Products.
	// Find a product with id 1.
	repo.Find(&product, where.Eq("id", 1))
}
```

### Testing

Testing database query using rel can be done using [reltest](https://godoc.org/github.com/Fs02/rel/reltest) package which is based on testify mock libarary.

```golang
// main_test.go
package main

import (
	"time"
	"testing"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/adapter/mysql"
	"github.com/Fs02/rel/where"
	"github.com/Fs02/reltest"
)

func TestInsert(t *testing.T) {
	repo := reltest.New()
	
	// expectations
	repo.ExpectInsert()
	repo.ExpectFind(where.Eq("id", 1)).Result(Product{
		ID: 5,
		Name: "soap",
		Price: 2000,
	})
	
	// run
	run(repo)
	
	// asserts
	repo.AssertExpectations(t)
}
```
