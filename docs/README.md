# REL

> Golang SQL Database Layer for Layered Architecture.

REL is golang orm-ish database layer for layered architecture. It's testable and comes with it's own test library. REL also features extendable query builder that allows you to write query using builder or plain sql.

## Features

- Testable repository with builtin reltest package.
- Elegant, yet extendable query builder.
- Supports Eager loading.
- Multi adapter.
- Soft Deletion.

## Install

```
go get github.com/Fs02/rel
go get github.com/Fs02/rel/reltest
```

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
type userRepository struct {
	db *DB
}

func (ur userRepository) Find(user *User, int id) error {
	return db.Where("id = ?", 1).First(&user)
}
```

Compared to other orm, REL api is built with testability in mind. REL uses interface to define contract of every database query or execution, all while making a chainable query possible. The ultimate **goal of REL is to be your database package** without the needs of making your own wrapper.

**Learn More: [Basics](basics.md)**
