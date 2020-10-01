# REL

[![GoDoc](https://godoc.org/github.com/go-rel/rel?status.svg)](https://godoc.org/github.com/go-rel/rel)
[![Build Status](https://github.com/go-rel/rel/workflows/Build/badge.svg)](https://github.com/go-rel/rel/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-rel/rel)](https://goreportcard.com/report/github.com/go-rel/rel)
[![Maintainability](https://api.codeclimate.com/v1/badges/194611cc82f02edcda6e/maintainability)](https://codeclimate.com/github/go-rel/rel/maintainability)
[![Codecov](https://codecov.io/gh/go-rel/rel/branch/master/graph/badge.svg?token=0P505E1IWB)](https://codecov.io/gh/go-rel/rel)
[![Gitter chat](https://badges.gitter.im/go-rel/rel.png)](https://gitter.im/go-rel/rel)

> Modern Database Access Layer for Golang.

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
go get github.com/go-rel/rel
```

## Getting Started

- Guides [https://go-rel.github.io](https://go-rel.github.io)

## Examples

- [go-todo-backend](https://github.com/Fs02/go-todo-backend) - Todo Backend

## License

Released under the [MIT License](https://github.com/go-rel/rel/blob/master/LICENSE)
