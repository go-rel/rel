# rel
[![GoDoc](https://godoc.org/github.com/Fs02/rel?status.svg)](https://godoc.org/github.com/Fs02/rel) [![Build Status](https://travis-ci.com/Fs02/rel.svg?branch=master)](https://travis-ci.com/Fs02/rel) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/rel)](https://goreportcard.com/report/github.com/Fs02/rel) [![Maintainability](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/maintainability)](https://codeclimate.com/github/Fs02/rel/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/test_coverage)](https://codeclimate.com/github/Fs02/rel/test_coverage)

> Golang SQL Repository Layer for Clean (Onion) Architecture.

rel is orm-ish library for golang that aims to be the repository layer of onion architecture. It's testable and comes with it's own test library. rel also features extendable query builder that allows you to write query using builder or plain sql.

## Features

- Testable repository with builtin reltest package.
- Elegant, yet extendable query builder.
- Supports Eager loading.
- Multi adapter.

## Install

```bash
go get github.com/Fs02/rel
```

## Quick Start

### Basic Usage
```golang
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

### Testing

```golang
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

**Learn More:** [https://fs02.github.io/rel](https://fs02.github.io/rel)

## Why rel

Most (if not all) orm for golang is written as a chainable API, meaning all of the query need to be called before performing actual action as a chain of method invocations. example:

```go
db.Where("id = ?", 1).First(&user)
```

Chainable api is very hard to be unit tested without writing a wrapper. One way to make it testable is to make an interface that also acts as a wrapper, which is usually ends up as its own repository package resides somewhere in your project:

```go
// mockable interface.
type UserRepository interface {
	Find(user *User, int id) error
}
// actual implementation
type userRepository struct{
	db *DB
}
func (ur userRepository) Find(user *User, int id) error {
	return db.Where("id = ?", 1).First(&user)
}
```

Compared to other orm, rel api is built with [testability](https://godoc.org/github.com/Fs02/rel/reltest) in mind. rel uses [interface](https://godoc.org/github.com/Fs02/rel#Repository) to define contract of every database query or execution, all while making a chainable query possible. The ultimate goal of rel is to be **your repository package without the needs of making your own wrapper**. example:

```go
// rel repository
repo.Find(&user, where.Eq("id", 1))
```

## License

Released under the [MIT License](https://github.com/Fs02/rel/blob/master/LICENSE)
