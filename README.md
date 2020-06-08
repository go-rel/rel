# REL

[![GoDoc](https://godoc.org/github.com/Fs02/rel?status.svg)](https://godoc.org/github.com/Fs02/rel)
[![Build Status](https://travis-ci.com/Fs02/rel.svg?branch=master)](https://travis-ci.com/Fs02/rel)
[![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/rel)](https://goreportcard.com/report/github.com/Fs02/rel)
[![Maintainability](https://api.codeclimate.com/v1/badges/194611cc82f02edcda6e/maintainability)](https://codeclimate.com/github/Fs02/rel/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/194611cc82f02edcda6e/test_coverage)](https://codeclimate.com/github/Fs02/rel/test_coverage)

> Golang SQL Database Layer for Layered Architecture.

REL is golang orm-ish database layer for layered architecture. It's testable and comes with its own test library. REL also features extendable query builder that allows you to write query using builder or plain sql.

## Features

- Testable repository with builtin reltest package.
- Elegant, yet extendable query builder with mix of syntactic sugar.
- Supports Eager loading.
- Multi adapter.
- Soft Deletion.
- Pagination.

## Quick Example

REL api might be slightly different than other Golang ORM/DAL. But it's developer friendly and doesn't requires a steep learning curve.

```go
book := Book{Title: "REL for Dummies"}

// Insert a Book.
if err := repo.Insert(ctx, &book); err != nil {
    return err
}

// Find a Book with id 1.
if err := repo.Find(ctx, &book, where.Eq("id", 1)); err != nil {
    return err
}

// Update a Book.
book.Title = "REL for Dummies 2nd Edition"
if err := repo.Update(ctx, &book); err != nil {
    return err
}

// Delete a Book.
if err := repo.Delete(ctx, &book); err != nil {
    return err
}
```

## Install

```bash
go get github.com/Fs02/rel
```

## Getting Started

- Guides [https://fs02.github.io/rel](https://fs02.github.io/rel)

## License

Released under the [MIT License](https://github.com/Fs02/rel/blob/master/LICENSE)
