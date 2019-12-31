## Basics

### Adapters

Rel uses adapter in order to generate and execute query to a database, below is the list of available adapters supported by rel out of the box.

| Adapter    | Package                              | Godoc                                                                                                                                 |
|------------|--------------------------------------|---------------------------------------------------------------------------------------------------------------------------------------|
| MySQL      | github.com/Fs02/rel/adapter/mysql    | [![GoDoc](https://godoc.org/github.com/Fs02/rel/adapter/mysql?status.svg)](https://godoc.org/github.com/Fs02/rel/adapter/mysql)       |
| PostgreSQL | github.com/Fs02/rel/adapter/postgres | [![GoDoc](https://godoc.org/github.com/Fs02/rel/adapter/postgres?status.svg)](https://godoc.org/github.com/Fs02/rel/adapter/postgres) |
| SQLite3    | github.com/Fs02/rel/adapter/sqlite3  | [![GoDoc](https://godoc.org/github.com/Fs02/rel/adapter/sqlite3?status.svg)](https://godoc.org/github.com/Fs02/rel/adapter/sqlite3)   |

### Example

Below is a very basic example on how to utilize rel using mysql adapter.
Testing database query using rel can be done using [reltest](https://godoc.org/github.com/Fs02/rel/reltest) package.

<!-- tabs:start -->

#### **main.go**

```go
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
func run(repo rel.Repository) {
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

#### **main_test.go**

```go
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
	// create a mocked repository.
	repo := reltest.New()
	
	// prepare mocks
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

<!-- tabs:end -->

### Conventions

#### Schema Definition

rel uses a struct as the schema to infer `table name`, `columns` and `primary field`.

```go
// Table name: books
type Book struct {
	ID        int       // id
	Title     string    // title
	Category  string    // category
	CreatedAt time.Time // created_at
	UpdatedAt time.Time // updated_at
}
```

#### Table Name

Table name will be the pluralized struct name in snake case, you may create a `Table() string` method to override the default table name.

```go
// Default table name is `books`
type Book struct {}

// Override table name to be `ebooks`
func (b Book) Table() string {
	return "ebooks"
}
```

#### Column Name

Column name will be the struct field name in snake case, you may override the column name by using using `db` tag.

```go
type Book struct {
	ID       int                // this field will be mapped to `id` column.
	Title    string `db:"name"` // this field will be mapped to `name` column.
	Category string `db:"-"`    // this field will be skipped
}
```

#### Primary Key

rel requires every struct to have at least `primary` key. by default field named `id` will be used as primary key. to use other field as primary key. you may define it as `primary` using `db` tag.


```go
type Book struct {
	UUID string `db:"uuid,primary"` // or just `db:",primary"`
}
```

#### Timestamp

rel automatically track created and updated time of each struct if `CreatedAt` or `UpdatedAt` field exists.

**Next: [Reading and Writing Data](crud.md)**
