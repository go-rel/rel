## Query Iterface

### Retrieving Data

TODO:

### Conditions



### Sorting

To retrieve a record from database in a specific order, you can use the sort api.

<!-- tabs:start -->
#### **main.go**

```go
repo.MustFindAll(&books, rel.NewSortAsc("updated_at"))

// or use alias: github.com/Fs02/rel/sort 
repo.MustFindAll(&books, sort.Asc("updated_at"))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.NewSortAsc("updated_at")).Result(book)

// or use alias: github.com/Fs02/rel/sort 
repo.ExpectFindAll(sort.Asc("updated_at")).Result(book)
```

<!-- tabs:end -->

You can also chain sort with other query.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustFindAll(&books, rel.Where(where.Eq("available", true).SortAsc("updated_at")))

// which is equal to:
repo.MustFindAll(&books, where.Eq("available", true), sort.Asc("updated_at"))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.Where(where.Eq("available", true).SortAsc("updated_at"))).Result(books)

// which is equal to:
repo.ExpectFindAll(where.Eq("available", true), sort.Asc("updated_at")).Result(books)
```

<!-- tabs:end -->

### Selecting Specific Fields

To select specific fields, you can use `Select` method, this way only specificied field will be mapped to books.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustFindAll(&books, rel.Select("id", "title"))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.Select("id", "title")).Result(books)
```

<!-- tabs:end -->

### Using Specific Table

By default, rel will use pluralized-snakecase struct name as the table name. To select from specific table, you can use `From` method.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustFindAll(&books, rel.From("ebooks"))

// chain it with select
repo.MustFindAll(&books, rel.Select("id", "title").From("ebooks"))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.From("ebooks")).Result(books)

// chain it with select
repo.ExpectFindAll(rel.Select("id", "title").From("ebooks")).Result(books)
```

<!-- tabs:end -->

### Limit and Offset

To set the limit and offset of query, use `Limit` and `Offset` api. `Offset` will be ignored if `Limit` is not specified.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustFindAll(&books, rel.Limit(10), rel.Offset(20))

// as chainable query.
repo.MustFindAll(&books, rel.Select().Limit(10).Offset(20))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.Limit(10), rel.Offset(20)).Result(books)

// as chainable query.
repo.ExpectFindAll(rel.Select().Limit(10).Offset(20)).Result(books)
```

<!-- tabs:end -->

### Group

To use group by query, you can use `Group` method.

<!-- tabs:start -->

#### **main.go**

```go
// custom struct to store the result.
var results []struct {
    Category string
    Count    int
}

// we need to explicitly specify table name since we are using an anonymous struct. 
repo.MustFindAll(&results, rel.Select("category", "COUNT(id) as id").From("books").Group("category"))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.Select("category", "COUNT(id) as id").From("books").Group("category")).Result(results)
```

<!-- tabs:end -->

### Joining Tables

To join tables, you can use `join` api.

> Joining table won't load the association to struct. If you want to load association on a struct, use [preload](associations.md#preload) instead.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustFindAll(&books, rel.Join("users"))
// or with explicit columns
repo.MustFindAll(&books, rel.JoinOn("users", "addresses.users_id", "users.id"))
// or use alias: github.com/Fs02/rel/join
repo.MustFindAll(&books, join.On("users", "addresses.users_id", "users.id"))
// or with custom join mode.
repo.MustFindAll(&books, rel.JoinWith("LEFT JOIN", "users", "addresses.users_id", "users.id"))
// or with raw sql
repo.MustFindAll(&books, rel.Joinf("JOIN `users` ON `addresses`.`user_id`=`users`.`id`"))
```

#### **main_test.go**

```go
repo.ExpectFindAll(rel.Join("users")).Result(books)
// or with explicit columns
repo.ExpectFindAll(rel.JoinOn("users", "addresses.users_id", "users.id")).Result(books)
// or use alias: github.com/Fs02/rel/join
repo.ExpectFindAll(join.On("users", "addresses.users_id", "users.id")).Result(books)
// or with custom join mode.
repo.ExpectFindAll(rel.JoinWith("LEFT JOIN", "users", "addresses.users_id", "users.id")).Result(books)
// or with raw sql
repo.ExpectFindAll(rel.Joinf("JOIN `users` ON `addresses`.`user_id`=`users`.`id`")).Result(books)
```

<!-- tabs:end -->

### Pessimistic Locking

rel supports pessimistic locking by using mechanism provided by the underlying database. `Lock` can be only used only inside transaction.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustFind(&book, where.Eq("id", 1), rel.Lock("FOR UPDATE"))
// or
repo.MustFind(&book, where.Eq("id", 1), rel.ForUpdate())
// or
repo.MustFind(&book, query.Where(where.Eq("id", 1)).Lock("FOR UPDATE"))
```

#### **main_test.go**

```go
repo.ExpectFind(where.Eq("id", 1), rel.Lock("FOR UPDATE")).Result(book)
// or
repo.ExpectFind(where.Eq("id", 1), rel.ForUpdate()).Result(book)
// or
repo.ExpectFind(query.Where(where.Eq("id", 1)).Lock("FOR UPDATE")).Result(book)
```

<!-- tabs:end -->

### Aggregation

rel provides a very basic `Aggregate` method which can be used to count, sum, max etc.

<!-- tabs:start -->

#### **main.go**

```go
repo.MustAggregaate(rel.From("books").Where(where.Eq("available", true)), "count", "id")
// or
repo.MustCount("books", where.Eq("available", true))
// or just count all books.
repo.MustCount("books")
```

#### **main_test.go**

```go
repo.ExpectAggregaate(rel.From("books").Where(where.Eq("available", true)), "count", "id").Result(5)
// or
repo.ExpectCount("books", where.Eq("available", true)).Result(5)
// or just count all books.
repo.ExpectCount("books").Result(7)
```

<!-- tabs:end -->
