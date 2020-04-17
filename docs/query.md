# Query Iterface

## Retrieving Data

REL provides two basic finders method, `Find` for retrieving single record, and `FindAll` for retrieving multiple record.

`Find` only accepts struct as the first argument, and always return the first result from the query.

<!-- tabs:start -->

### **Example**

Retrieve a book.

[query.go](query.go ':include :fragment=find')

### **Mock**

Mock retrieve a book.

[query_test.go](query_test.go ':include :fragment=find')

<!-- tabs:end -->

`FindAll` only accepts slice as the first argument, and always return all result from the query.

<!-- tabs:start -->

### **Example**

Retrieve all books.

[query.go](query.go ':include :fragment=find-all')

### **Mock**

Mock retrieve all books.

[query_test.go](query_test.go ':include :fragment=find-all')

<!-- tabs:end -->

## Conditions

To retrieve filtered recods from database, you can use filter api to specify [condition](https://pkg.go.dev/github.com/Fs02/rel/where). For example, to filter all books that available, you can use `rel.Eq` in the query builder.

<!-- tabs:start -->

### **Example**

Retrieve all available books using filter query.

[query.go](query.go ':include :fragment=condition')

Using alias: github.com/Fs02/rel/where

[query.go](query.go ':include :fragment=condition-alias')

Using fragment to using fragment to specify custom SQL query.

[query.go](query.go ':include :fragment=condition-fragment')

### **Mock**

Mock retrieve all available books.

[query_test.go](query_test.go ':include :fragment=condition')

Mock retrieve all using alias: github.com/Fs02/rel/where

[query_test.go](query_test.go ':include :fragment=condition-alias')

Mock retrieve all using fragment to specify custom SQL query.

[query_test.go](query_test.go ':include :fragment=condition-fragment')

<!-- tabs:end -->

You can use `rel.And` or `rel.Or` to specify more conditions.

<!-- tabs:start -->

### **Example**

Retrieve all available books where price is at least 100 or in discount using filter query.

[query.go](query.go ':include :fragment=condition-advanced')

Retrieve all available books where price is at least 100 or in discount using chained filter query.

[query.go](query.go ':include :fragment=condition-advanced-chain')

Retrieve all available books where price is at least 100 or in discount using alias: github.com/Fs02/rel/where

[query.go](query.go ':include :fragment=condition-advanced-alias')

### **Mock**


Mock retrieve all available books where price is at least 100 or in discount using filter query.

[query_test.go](query_test.go ':include :fragment=condition-advanced')

Mock retrieve all available books where price is at least 100 or in discount using chained filter query.

[query_test.go](query_test.go ':include :fragment=condition-advanced-chain')

Mock retrieve all available books where price is at least 100 or in discount using alias: github.com/Fs02/rel/where

[query_test.go](query_test.go ':include :fragment=condition-advanced-alias')


<!-- tabs:end -->

## Sorting

To retrieve records from database in a specific order, you can use the sort api.

<!-- tabs:start -->
### **Example**

Sort books ascending by updated_at field.

[query.go](query.go ':include :fragment=sorting')

Using alias if you need more syntactic sugar.

[query.go](query.go ':include :fragment=sorting-alias')

### **Mock**

Mock sort books ascending by updated_at field.

[query_test.go](query_test.go ':include :fragment=sorting')

Mock sort using alias.

[query_test.go](query_test.go ':include :fragment=sorting-alias')

<!-- tabs:end -->

Combining with other query is fairly easy.

<!-- tabs:start -->

### **Example**

Chain where and sort using [query builder](https://pkg.go.dev/github.com/Fs02/rel?tab=doc#Query).

[query.go](query.go ':include :fragment=sorting-with-condition')

It's also possible to use variadic arguments to combine multiple queries.

[query.go](query.go ':include :fragment=sorting-with-condition-variadic')

### **Mock**

Mock chain where and sort using [query builder](https://pkg.go.dev/github.com/Fs02/rel?tab=doc#Query).

[query_test.go](query_test.go ':include :fragment=sorting-with-condition')

Mock query that uses variadic arguments to combine multiple queries.

[query_test.go](query_test.go ':include :fragment=sorting-with-condition-variadic')


<!-- tabs:end -->

## Selecting Specific Fields

To select specific fields, you can use `Select` method, this way only specificied field will be mapped to books.

?> Specifying select without argument (`rel.Select()`) will automatically load all fields. This is helpful when used as query builder entry point (compared to using `rel.From`), because you can let REL to infer the table name.

<!-- tabs:start -->

### **Example**

Load only id and title.

[query.go](query.go ':include :fragment=select')

### **Mock**

Mock find all to load only id and title.

[query_test.go](query_test.go ':include :fragment=select')

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

Specify limit and offset.

[query.go](query.go ':include :fragment=limit-offset')

As a chainable api.

[query.go](query.go ':include :fragment=limit-offset-chained')

### **Mock**

Mock limit and offset.

[query_test.go](query_test.go ':include :fragment=limit-offset')

Mock using chainable api.

[query_test.go](query_test.go ':include :fragment=limit-offset-chained')

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

?> Joining table won't load the association to struct. If you want to load association on a struct, use [preload](associations.md#preload) instead.

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
