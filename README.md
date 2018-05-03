# grimoire
[![GoDoc](https://godoc.org/github.com/Fs02/grimoire?status.svg)](https://godoc.org/github.com/Fs02/grimoire) [![Build Status](https://travis-ci.org/Fs02/grimoire.svg?branch=master)](https://travis-ci.org/Fs02/grimoire) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/grimoire)](https://goreportcard.com/report/github.com/Fs02/grimoire) [![Maintainability](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/maintainability)](https://codeclimate.com/github/Fs02/grimoire/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/test_coverage)](https://codeclimate.com/github/Fs02/grimoire/test_coverage)

Grimoire is a flexible ORM for golang. It features flexible query API and builtin validation. It currently supports MySQL, PostgreSQL and SQLite3 but a custom adapter can be implemented easily using the Adapter interface.

Features:

- Query Builder
- Association Preloading
- Struct style create and update
- Changeset Style create and update
- Builtin validation using changeset
- Multi adapter support
- Logger

## Install

```bash
go get github.com/Fs02/grimoire
```

## Quick Start

```golang
package main

import (
	"time"
	"github.com/Fs02/grimoire"
	. "github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/adapter/mysql"
	"github.com/Fs02/grimoire/changeset"
)

type Product struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ProductChangeset(product interface{}, params map[string]interface{}) *changeset.Changeset {
	ch := changeset.Cast(product, params, []string{"name", "price"})
	changeset.ValidateRequired(ch, []string{"name", "price"})
	changeset.ValidateMin(ch, "price", 100)
	return ch
}

func main() {
	// initialize mysql adapter.
	adapter, err := mysql.Open("root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize grimoire's repo.
	repo := grimoire.New(adapter)

	var product Product

	// Changeset is used when creating or updating your data.
	ch := ProductChangeset(product, map[string]interface{}{
		"name": "shampoo",
		"price": 1000
	})

	if ch.Error() != nil {
		// do something
	}

	// Create products with changeset and return the result to &product,
	repo.From("products").MustCreate(&product, ch)

	// Find a product with id 1.
	repo.From("products").Find(1).MustOne(&product)

	// Update products with id=1.
	repo.From("products").Find(1).MustUpdate(&product, ch)

	// Delete Product with id=1.
	repo.From("products").Find(1).MustDelete()
}
```

## Documentation

See: [https://fs02.github.io/grimoire](https://fs02.github.io/grimoire)

## License

Released under the [MIT License](https://github.com/Fs02/grimoire/blob/master/LICENSE)
