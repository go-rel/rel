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

## Learn More

- [Basic Usage](quickstart.md)
- [Comparison with other ORMs](comparison.md)
- [Package documentation](https://godoc.org/github.com/Fs02/rel)
- [Test package documentation](https://godoc.org/github.com/Fs02/rel/reltest)
