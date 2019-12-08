# rel
[![GoDoc](https://godoc.org/github.com/Fs02/rel?status.svg)](https://godoc.org/github.com/Fs02/rel) [![Build Status](https://travis-ci.com/Fs02/rel.svg?branch=master)](https://travis-ci.com/Fs02/rel) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/rel)](https://goreportcard.com/report/github.com/Fs02/rel) [![Maintainability](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/maintainability)](https://codeclimate.com/github/Fs02/rel/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/test_coverage)](https://codeclimate.com/github/Fs02/rel/test_coverage)

rel is a testable repository layer for sql database. 

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

## License

Released under the [MIT License](https://github.com/Fs02/rel/blob/master/LICENSE)
