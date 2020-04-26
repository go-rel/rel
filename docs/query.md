# Query Iterface

## Retrieving Data

REL provides two basic finders method, `Find` for retrieving single record, and `FindAll` for retrieving multiple record.

`Find` only accepts struct as the first argument, and always return the first result from the query.

<!-- tabs:start -->

### **Example**

Retrieve a book where id=1.

[query.go](query.go ':include :fragment=find')

### **Mock**

Mock retrieve a book where id=1.

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

Alias can be used to boost readability when dealing with short query.

[query.go](query.go ':include :fragment=condition-alias')

Use fragment to specify custom SQL query.

[query.go](query.go ':include :fragment=condition-fragment')

### **Mock**

Mock retrieve all available books.

[query_test.go](query_test.go ':include :fragment=condition')

Mock retrieve all using alias.

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

Load from `ebooks` table.

[query.go](query.go ':include :fragment=table')

Chain the query with select.

[query.go](query.go ':include :fragment=table-chained')

### **Mock**

Mock load from `ebooks` table.

[query_test.go](query_test.go ':include :fragment=table')

Mock chain the query with select.

[query_test.go](query_test.go ':include :fragment=table-chained')

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

Retrieve count of books for every category.

[query.go](query.go ':include :fragment=group')

### **Mock**

Mock retrieve count of books for every category.

[query_test.go](query_test.go ':include :fragment=group')

<!-- tabs:end -->

## Joining Tables

To join tables, you can use `join` api.

?> Joining table won't load the association to struct. If you want to load association on a struct, use [preload](associations.md#preload) instead.

<!-- tabs:start -->

### **Example**

Join transaction and book table, then filter only transaction that have specified book name. This methods assumes belongs to relation, which means it'll try to join using `transactions.book_id=books.id`.

[query.go](query.go ':include :fragment=join')

Specifying which column to join using JoinOn.

[query.go](query.go ':include :fragment=join-on')

Syntactic sugar also available for join.

[query.go](query.go ':include :fragment=join-alias')

Joining table with custom join mode.

[query.go](query.go ':include :fragment=join-with')

Use fragment for more complex join query.

[query.go](query.go ':include :fragment=join-fragment')

### **Mock**

Mock join transaction and book table, then filter only transaction that have specified book name.

[query_test.go](query_test.go ':include :fragment=join')

Specifying which column to join using JoinOn.

[query_test.go](query_test.go ':include :fragment=join-on')

Syntactic sugar also available for join.

[query_test.go](query_test.go ':include :fragment=join-alias')

Joining table with custom join mode.

[query_test.go](query_test.go ':include :fragment=join-with')

Use fragment for more complex join query.

[query_test.go](query_test.go ':include :fragment=join-fragment')

<!-- tabs:end -->

## Pessimistic Locking

REL supports pessimistic locking by using mechanism provided by the underlying database. `Lock` can be only used only inside transaction.

<!-- tabs:start -->

### **Example**

Retrieve and lock a row for update.

[query.go](query.go ':include :fragment=lock')

Retrieve and lock a row using predefined lock alias.

[query.go](query.go ':include :fragment=lock-for-update')

Retrieve and lock a row using chained query.

[query.go](query.go ':include :fragment=lock-chained')

### **Mock**

Mock retrieve and lock a row for update.

[query_test.go](query_test.go ':include :fragment=lock')

Mock retrieve and lock a row using predefined lock alias.

[query_test.go](query_test.go ':include :fragment=lock-for-update')

Mock retrieve and lock a row using chained query.

[query_test.go](query_test.go ':include :fragment=lock-chained')

<!-- tabs:end -->

## Aggregation

REL provides a very basic `Aggregate` method which can be used to count, sum, max etc.

<!-- tabs:start -->

### **Example**

Count all available books using aggregate.

[query.go](query.go ':include :fragment=aggregate')

Count all available books using count.

[query.go](query.go ':include :fragment=count')

Count all available books using count.

[query.go](query.go ':include :fragment=count-with-condition')

### **Mock**

Mock count all available books using aggregate.

[query_test.go](query_test.go ':include :fragment=aggregate')

Mock count all available books using count.

[query_test.go](query_test.go ':include :fragment=count')

Mock count all available books using count.

[query_test.go](query_test.go ':include :fragment=count-with-condition')

<!-- tabs:end -->

## Pagination

REL provides a convenient `FindAndCountAll` methods that is useful for pagination, It's a combination of `FindAll` and `Count` method.
FindAndCountAll returns count of records (ignoring limit and offset query) and an error.

<!-- tabs:start -->

### **Example**

Retrieve all books within limit and offset and also count of all books.

[query.go](query.go ':include :fragment=find-and-count-all')

### **Mock**

Mock retrieve all books within limit and offset and also count of all books.

[query_test.go](query_test.go ':include :fragment=find-and-count-all')

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

## Native SQL Query

REL allows querying using native SQL query, this is especially useful when using complex query that cannot be covered with the query builder.

<!-- tabs:start -->

### **Example**

Retrieve a book using raw sql query.

[query.go](query.go ':include :fragment=sql')

### **Mock**

Mock and retrieve a book using raw sql query.

[query_test.go](query_test.go ':include :fragment=sql')

<!-- tabs:end -->

**Next: [Association](association.md)**
