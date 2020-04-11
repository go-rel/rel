# Query Iterface

## Retrieving Data

REL provides two basic finders method, `Find` for retrieving single record, and `FindAll` for retrieving multiple record.

`Find` only accepts struct as the first argument, and always return the first result from the query.

<!-- tabs:start -->

### **Example**

```go
repo.Find(ctx, &book)
```

### **Mock**

```go
repo.ExpectFind().Result(book)
```

<!-- tabs:end -->

`FindAll` only accepts slice as the first argument, and always return all result from the query.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, &books)
```

### **Mock**

```go
repo.ExpectFindAll().Result(books)
```

<!-- tabs:end -->

## Conditions

To retrieve filtered recods from database, you can use filter api to specify coondition. For example, to filter all books that available, you can use `rel.Eq` in the query builder.

<!-- tabs:start -->

### **Example**

```go
// Retrieve all available books
repo.FindAll(ctx, &books, rel.Eq("available", true))

// or use alias: github.com/Fs02/rel/where
repo.FindAll(ctx, &books, where.Eq("available", true))

// or use raw query
repo.FindAll(ctx, &books, where.Fragment("available=?", true))
```

### **Mock**

```go
// mock and returns books
repo.ExpectFindAll(rel.Eq("available", true)).Result(book)

// with alias
repo.ExpectFindAll(where.Eq("available", true)).Result(book)

// with raw query
repo.ExpectFindAll(&books, where.Fragment("available=?", true)).Result(book)
```

<!-- tabs:end -->

You can use `rel.And` or `rel.Or` to specify more conditions.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, rel.And(rel.Eq("available", true), rel.Or(rel.Gte("price", 100), rel.Eq("discount", true))))

// or use filter chain
repo.FindAll(ctx, rel.Eq("available", true).And(rel.Gte("price", 100).OrEq("discount", true)))

// or use alias: github.com/Fs02/rel/where
repo.FindAll(ctx, where.Eq("available", true).And(where.Gte("price", 100).OrEq("discount", true)))
```

### **Mock**

```go
repo.ExpectFindAll(rel.And(rel.Eq("available", true), rel.Or(rel.Gte("price", 100), rel.Eq("discount", true)))).Result(book)

// or use filter chain
repo.ExpectFindAll(rel.Eq("available", true).And(rel.Gte("price", 100).OrEq("discount", true))).Result(book)

// or use alias: github.com/Fs02/rel/where
repo.ExpectFindAll(where.Eq("available", true).And(where.Gte("price", 100).OrEq("discount", true))).Result(book)
```

<!-- tabs:end -->

## Sorting

To retrieve records from database in a specific order, you can use the sort api.

<!-- tabs:start -->
### **Example**

```go
repo.FindAll(ctx, &books, rel.NewSortAsc("updated_at"))

// or use alias: github.com/Fs02/rel/sort
repo.FindAll(ctx, &books, sort.Asc("updated_at"))
```

### **Mock**

```go
repo.ExpectFindAll(rel.NewSortAsc("updated_at")).Result(book)

// or use alias: github.com/Fs02/rel/sort
repo.ExpectFindAll(sort.Asc("updated_at")).Result(book)
```

<!-- tabs:end -->

You can also chain sort with other query.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, &books, rel.Where(where.Eq("available", true).SortAsc("updated_at")))

// which is equal to:
repo.FindAll(ctx, &books, where.Eq("available", true), sort.Asc("updated_at"))
```

### **Mock**

```go
repo.ExpectFindAll(rel.Where(where.Eq("available", true).SortAsc("updated_at"))).Result(books)

// which is equal to:
repo.ExpectFindAll(where.Eq("available", true), sort.Asc("updated_at")).Result(books)
```

<!-- tabs:end -->

## Selecting Specific Fields

To select specific fields, you can use `Select` method, this way only specificied field will be mapped to books.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, &books, rel.Select("id", "title"))
```

### **Mock**

```go
repo.ExpectFindAll(rel.Select("id", "title")).Result(books)
```

<!-- tabs:end -->

## Using Specific Table

By default, REL will use pluralized-snakecase struct name as the table name. To select from specific table, you can use `From` method.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, &books, rel.From("ebooks"))

// chain it with select
repo.FindAll(ctx, &books, rel.Select("id", "title").From("ebooks"))
```

### **Mock**

```go
repo.ExpectFindAll(rel.From("ebooks")).Result(books)

// chain it with select
repo.ExpectFindAll(rel.Select("id", "title").From("ebooks")).Result(books)
```

<!-- tabs:end -->

## Limit and Offset

To set the limit and offset of query, use `Limit` and `Offset` api. `Offset` will be ignored if `Limit` is not specified.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, &books, rel.Limit(10), rel.Offset(20))

// as chainable query.
repo.FindAll(ctx, &books, rel.Select().Limit(10).Offset(20))
```

### **Mock**

```go
repo.ExpectFindAll(rel.Limit(10), rel.Offset(20)).Result(books)

// as chainable query.
repo.ExpectFindAll(rel.Select().Limit(10).Offset(20)).Result(books)
```

<!-- tabs:end -->

## Group

To use group by query, you can use `Group` method.

<!-- tabs:start -->

### **Example**

```go
// custom struct to store the result.
var results []struct {
    Category string
    Count    int
}

// we need to explicitly specify table name since we are using an anonymous struct. 
repo.FindAll(ctx, &results, rel.Select("category", "COUNT(id) as id").From("books").Group("category"))
```

### **Mock**

```go
repo.ExpectFindAll(rel.Select("category", "COUNT(id) as id").From("books").Group("category")).Result(results)
```

<!-- tabs:end -->

## Joining Tables

To join tables, you can use `join` api.

> Joining table won't load the association to struct. If you want to load association on a struct, use [preload](associations.md#preload) instead.

<!-- tabs:start -->

### **Example**

```go
repo.FindAll(ctx, &books, rel.Join("users"))
// or with explicit columns
repo.FindAll(ctx, &books, rel.JoinOn("users", "addresses.users_id", "users.id"))
// or use alias: github.com/Fs02/rel/join
repo.FindAll(ctx, &books, join.On("users", "addresses.users_id", "users.id"))
// or with custom join mode.
repo.FindAll(ctx, &books, rel.JoinWith("LEFT JOIN", "users", "addresses.users_id", "users.id"))
// or with raw sql
repo.FindAll(ctx, &books, rel.Joinf("JOIN `users` ON `addresses`.`user_id`=`users`.`id`"))
```

### **Mock**

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

## Pessimistic Locking

REL supports pessimistic locking by using mechanism provided by the underlying database. `Lock` can be only used only inside transaction.

<!-- tabs:start -->

### **Example**

```go
repo.Find(ctx, &book, where.Eq("id", 1), rel.Lock("FOR UPDATE"))
// or
repo.Find(ctx, &book, where.Eq("id", 1), rel.ForUpdate())
// or
repo.Find(ctx, &book, query.Where(where.Eq("id", 1)).Lock("FOR UPDATE"))
```

### **Mock**

```go
repo.ExpectFind(where.Eq("id", 1), rel.Lock("FOR UPDATE")).Result(book)
// or
repo.ExpectFind(where.Eq("id", 1), rel.ForUpdate()).Result(book)
// or
repo.ExpectFind(query.Where(where.Eq("id", 1)).Lock("FOR UPDATE")).Result(book)
```

<!-- tabs:end -->

## Aggregation

REL provides a very basic `Aggregate` method which can be used to count, sum, max etc.

<!-- tabs:start -->

### **Example**

```go
count, err = repo.Aggregate(ctx, rel.From("books").Where(where.Eq("available", true)), "count", "id")
// or
count, err = repo.Count(ctx, "books", where.Eq("available", true))
// or just count all books.
count, err = repo.Count(ctx, "books")
```

### **Mock**

```go
repo.ExpectAggregate(rel.From("books").Where(where.Eq("available", true)), "count", "id").Result(5)
// or
repo.ExpectCount("books", where.Eq("available", true)).Result(5)
// or just count all books.
repo.ExpectCount("books").Result(7)
```

<!-- tabs:end -->

## Pagination

REL provides a convenient `FindAndCountAll` methods that is useful for pagination, It's a combination of `FindAll` and `Count` method.
FindAndCountAll returns count of records (ignoring limit and offset query) and an error.

<!-- tabs:start -->

### **Example**

```go
// Find and count total books in database.
count, err = repo.FindAndCountAll(ctx, &books, rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10))
```

### **Mock**

```go
// Expect count to returns books and with total count of 12.
repo.ExpectFindAndCountAll(rel.Where(where.Like("title", "%dummies%")).Limit(10).Offset(10)).Result(books, 12)
```

<!-- tabs:end -->

## Batch Iteration

REL provides records iterator that can be use for perform batch processing of large amounts of records.

Options:

- `BatchSize` - The size of batches (default 1000).
- `Start` - The primary value (ID) to start from (inclusive).
- `Finish` - The primary value (ID) to finish at (inclusive).

<!-- tabs:start -->

### **Example**

[query.go](query.go ':include :fragment=batch-iteration')

### **Mock**

Mock and returns users.

[query_test.go](query_test.go ':include :fragment=batch-iteration')

Mock and retuns `reltest.ErrConnectionClosed` (`sql.ErrConnDone`).

[query_test.go](query_test.go ':include :fragment=batch-iteration-connection-error')

<!-- tabs:end -->

**Next: [Association](association.md)**
