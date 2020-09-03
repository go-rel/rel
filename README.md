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
- Supports nested transactions.
- Composite Primary Key.
- Multi adapter.
- Soft Deletion.
- Pagination.
- Schema Migration.

## Install

```bash
go get github.com/Fs02/rel
```

## Getting Started

- Guides [https://fs02.github.io/rel](https://fs02.github.io/rel)

## Examples

- [go-todo-backend](https://github.com/Fs02/go-todo-backend) - Todo Backend

## License

Released under the [MIT License](https://github.com/Fs02/rel/blob/master/LICENSE)
