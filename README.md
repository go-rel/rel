# grimoire
[![GoDoc](https://godoc.org/github.com/Fs02/grimoire?status.svg)](https://godoc.org/github.com/Fs02/grimoire) [![Build Status](https://travis-ci.org/Fs02/grimoire.svg?branch=master)](https://travis-ci.org/Fs02/grimoire) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/grimoire)](https://goreportcard.com/report/github.com/Fs02/grimoire) [![Maintainability](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/maintainability)](https://codeclimate.com/github/Fs02/grimoire/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/test_coverage)](https://codeclimate.com/github/Fs02/grimoire/test_coverage)

Grimoire is a database access layer inspired by Ecto. It features a flexible query API and built-in validation. It currently supports MySQL, PostgreSQL, and SQLite3 but a custom adapter can be implemented easily using the Adapter interface.

Features:

- Query Builder
- Association Preloading
- Struct style create and update
- Changeset Style create and update
- Builtin validation using changeset
- Multi adapter support
- Logger

## Motivation

Common go ORM accepts struct as a value for modifying records which has a problem of unable to differentiate between an empty, nil, or undefined value. It's a tricky problem especially when you want to have an endpoint that supports partial updates. Grimoire attempts to solve that problem by integrating Changeset system inspired from Elixir's Ecto. Changeset is a form like entity which allows us to not only solve that problem but also help us with casting, validations, and constraints check.

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
	"github.com/Fs02/grimoire/adapter/mysql"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/params"
)

type Product struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChangeProduct prepares data before database operation.
// Such as casting value to appropriate types and perform validations.
func ChangeProduct(product interface{}, params params.Params) *changeset.Changeset {
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

	// Inserting Products.
	// Changeset is used when creating or updating your data.
	ch := ChangeProduct(product, params.Map{
		"name":  "shampoo",
		"price": 1000,
	})

	if ch.Error() != nil {
		// handle error
	}

	// Changeset can also be created directly from json string.
	jsonch := ChangeProduct(product, params.ParseJSON(`{
		"name":  "soap",
		"price": 2000,
	}`))

	// Create products with changeset and return the result to &product,
	if err = repo.From("products").Insert(&product, ch); err != nil {
		// handle error
	}

	// or panic when insertion pailed
	repo.From("products").MustInsert(&product, jsonch)

	// Querying Products.
	// Find a product with id 1.
	repo.From("products").Find(1).MustOne(&product)

	// Updating Products.
	// Update products with id=1.
	repo.From("products").Find(1).MustUpdate(&product, ch)

	// Deleting Products.
	// Delete Product with id=1.
	repo.From("products").Find(1).MustDelete()
}
```

## Examples

- [Todo API](https://github.com/Fs02/grimoire-todo-example)

## Documentation

Guides: [https://fs02.github.io/grimoire](https://fs02.github.io/grimoire)

API Documentation: [https://godoc.org/github.com/Fs02/grimoire](https://godoc.org/github.com/Fs02/grimoire)

## License

Released under the [MIT License](https://github.com/Fs02/grimoire/blob/master/LICENSE)
