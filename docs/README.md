## rel

> Golang SQL Repository Layer for Clean (Onion) Architecture.

rel is orm-ish library for golang that aims to be the repository layer of onion architecture. It's testable and comes with it's own test library. rel also features extendable query builder that allows you to write query using builder or plain sql.

See the [Quick start](quickstart.md) guide for more details.

## Features

- Testable repository with builtin reltest package.
- Elegant, yet extendable query builder.
- Supports Eager loading.
- Multi adapter support.

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

## Learn More

- [Basic Usage](quickstart.md)
- [Package documentation](https://godoc.org/github.com/Fs02/rel)
- [Test package documentation](https://godoc.org/github.com/Fs02/rel/reltest)
