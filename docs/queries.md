# Queries

## Retrieving Data

REL provides two basic finders method, `Find` for retrieving single record, and `FindAll` for retrieving multiple record.

`Find` only accepts struct as the first argument, and always return the first result from the query.

=== "Example"

    Retrieve a book where id=1.

    {{ embed_code("docs/queries.go", "find", "\t") }}

=== "Mock"

    Mock retrieve a book where id=1.

    {{ embed_code("docs/queries_test.go", "find", "\t") }}


`FindAll` only accepts slice as the first argument, and always return all result from the query.

=== "Example"

    Retrieve all books.

    {{ embed_code("docs/queries.go", "find-all", "\t") }}

=== "Mock"

    Mock retrieve all books.

    {{ embed_code("docs/queries_test.go", "find-all", "\t") }}

## Conditions

To retrieve filtered recods from database, you can use filter api to specify [condition](https://pkg.go.dev/github.com/Fs02/rel/where). For example, to filter all books that available, you can use `rel.Eq` in the query builder.

=== "Example"

    Retrieve all available books using filter query.

    {{ embed_code("docs/queries.go", "condition", "\t") }}

    Alias can be used to boost readability when dealing with short query.

    {{ embed_code("docs/queries.go", "condition-alias", "\t") }}

    Use fragment to specify custom SQL query.

    {{ embed_code("docs/queries.go", "condition-fragment", "\t") }}

=== "Mock"

    Mock retrieve all available books.

    {{ embed_code("docs/queries_test.go", "condition", "\t") }}

    Mock retrieve all using alias.

    {{ embed_code("docs/queries_test.go", "condition-alias", "\t") }}

    Mock retrieve all using fragment to specify custom SQL query.

    {{ embed_code("docs/queries_test.go", "condition-fragment", "\t") }}


You can use `rel.And` or `rel.Or` to specify more conditions.

=== "Example"

    Retrieve all available books where price is at least 100 or in discount using filter query.

    {{ embed_code("docs/queries.go", "condition-advanced", "\t") }}

    Retrieve all available books where price is at least 100 or in discount using chained filter query.

    {{ embed_code("docs/queries.go", "condition-advanced-chain", "\t") }}

    Retrieve all available books where price is at least 100 or in discount using alias: github.com/Fs02/rel/where

    {{ embed_code("docs/queries.go", "condition-advanced-alias", "\t") }}

=== "Mock"

    Mock retrieve all available books where price is at least 100 or in discount using filter query.

    {{ embed_code("docs/queries_test.go", "condition-advanced", "\t") }}

    Mock retrieve all available books where price is at least 100 or in discount using chained filter query.

    {{ embed_code("docs/queries_test.go", "condition-advanced-chain", "\t") }}

    Mock retrieve all available books where price is at least 100 or in discount using alias: github.com/Fs02/rel/where

    {{ embed_code("docs/queries_test.go", "condition-advanced-alias", "\t") }}


## Sorting

To retrieve records from database in a specific order, you can use the sort api.

=== "Example"

    Sort books ascending by updated_at field.

    {{ embed_code("docs/queries.go", "sorting", "\t") }}

    Using alias if you need more syntactic sugar.

    {{ embed_code("docs/queries.go", "sorting-alias", "\t") }}

=== "Mock"

    Mock sort books ascending by updated_at field.

    {{ embed_code("docs/queries_test.go", "sorting", "\t") }}

    Mock sort using alias.

    {{ embed_code("docs/queries_test.go", "sorting-alias", "\t") }}


Combining with other query is fairly easy.

=== "Example"

    Chain where and sort using [query builder](https://pkg.go.dev/github.com/Fs02/rel?tab=doc#Query).

    {{ embed_code("docs/queries.go", "sorting-with-condition", "\t") }}

    It's also possible to use variadic arguments to combine multiple queries.

    {{ embed_code("docs/queries.go", "sorting-with-condition-variadic", "\t") }}

=== "Mock"

    Mock chain where and sort using [query builder](https://pkg.go.dev/github.com/Fs02/rel?tab=doc#Query).

    {{ embed_code("docs/queries_test.go", "sorting-with-condition", "\t") }}

    Mock query that uses variadic arguments to combine multiple queries.

    {{ embed_code("docs/queries_test.go", "sorting-with-condition-variadic", "\t") }}


## Selecting Specific Fields

To select specific fields, you can use `Select` method, this way only specificied field will be mapped to books.

?> Specifying select without argument (`rel.Select()`) will automatically load all fields. This is helpful when used as query builder entry point (compared to using `rel.From`), because you can let REL to infer the table name.

=== "Example"

    Load only id and title.

    {{ embed_code("docs/queries.go", "select", "\t") }}

=== "Mock"

    Mock find all to load only id and title.

    {{ embed_code("docs/queries_test.go", "select", "\t") }}


## Using Specific Table

By default, REL will use pluralized-snakecase struct name as the table name. To select from specific table, you can use `From` method.

=== "Example"

    Load from `ebooks` table.

    {{ embed_code("docs/queries.go", "table", "\t") }}

    Chain the query with select.

    {{ embed_code("docs/queries.go", "table-chained", "\t") }}

=== "Mock"

    Mock load from `ebooks` table.

    {{ embed_code("docs/queries_test.go", "table", "\t") }}

    Mock chain the query with select.

    {{ embed_code("docs/queries_test.go", "table-chained", "\t") }}


## Limit and Offset

To set the limit and offset of query, use `Limit` and `Offset` api. `Offset` will be ignored if `Limit` is not specified.

=== "Example"

    Specify limit and offset.

    {{ embed_code("docs/queries.go", "limit-offset", "\t") }}

    As a chainable api.

    {{ embed_code("docs/queries.go", "limit-offset-chained", "\t") }}

=== "Mock"

    Mock limit and offset.

    {{ embed_code("docs/queries_test.go", "limit-offset", "\t") }}

    Mock using chainable api.

    {{ embed_code("docs/queries_test.go", "limit-offset-chained", "\t") }}


## Group

To use group by query, you can use `Group` method.

=== "Example"

    Retrieve count of books for every category.

    {{ embed_code("docs/queries.go", "group", "\t") }}

=== "Mock"

    Mock retrieve count of books for every category.

    {{ embed_code("docs/queries_test.go", "group", "\t") }}


## Joining Tables

To join tables, you can use `join` api.

?> Joining table won't load the association to struct. If you want to load association on a struct, use [preload](associations.md#preload) instead.

=== "Example"

    Join transaction and book table, then filter only transaction that have specified book name. This methods assumes belongs to relation, which means it'll try to join using `transactions.book_id=books.id`.

    {{ embed_code("docs/queries.go", "join", "\t") }}

    Specifying which column to join using JoinOn.

    {{ embed_code("docs/queries.go", "join-on", "\t") }}

    Syntactic sugar also available for join.

    {{ embed_code("docs/queries.go", "join-alias", "\t") }}

    Joining table with custom join mode.

    {{ embed_code("docs/queries.go", "join-with", "\t") }}

    Use fragment for more complex join query.

    {{ embed_code("docs/queries.go", "join-fragment", "\t") }}

=== "Mock"

    Mock join transaction and book table, then filter only transaction that have specified book name.

    {{ embed_code("docs/queries_test.go", "join", "\t") }}

    Specifying which column to join using JoinOn.

    {{ embed_code("docs/queries_test.go", "join-on", "\t") }}

    Syntactic sugar also available for join.

    {{ embed_code("docs/queries_test.go", "join-alias", "\t") }}

    Joining table with custom join mode.

    {{ embed_code("docs/queries_test.go", "join-with", "\t") }}

    Use fragment for more complex join query.

    {{ embed_code("docs/queries_test.go", "join-fragment", "\t") }}


## Pessimistic Locking

REL supports pessimistic locking by using mechanism provided by the underlying database. `Lock` can be only used only inside transaction.

=== "Example"

    Retrieve and lock a row for update.

    {{ embed_code("docs/queries.go", "lock", "\t") }}

    Retrieve and lock a row using predefined lock alias.

    {{ embed_code("docs/queries.go", "lock-for-update", "\t") }}

    Retrieve and lock a row using chained query.

    {{ embed_code("docs/queries.go", "lock-chained", "\t") }}

=== "Mock"

    Mock retrieve and lock a row for update.

    {{ embed_code("docs/queries_test.go", "lock", "\t") }}

    Mock retrieve and lock a row using predefined lock alias.

    {{ embed_code("docs/queries_test.go", "lock-for-update", "\t") }}

    Mock retrieve and lock a row using chained query.

    {{ embed_code("docs/queries_test.go", "lock-chained", "\t") }}

## Aggregation

REL provides a very basic `Aggregate` method which can be used to count, sum, max etc.

=== "Example"

    Count all available books using aggregate.

    {{ embed_code("docs/queries.go", "aggregate", "\t") }}

    Count all available books using count.

    {{ embed_code("docs/queries.go", "count", "\t") }}

    Count all available books using count.

    {{ embed_code("docs/queries.go", "count-with-condition", "\t") }}

=== "Mock"

    Mock count all available books using aggregate.

    {{ embed_code("docs/queries_test.go", "aggregate", "\t") }}

    Mock count all available books using count.

    {{ embed_code("docs/queries_test.go", "count", "\t") }}

    Mock count all available books using count.

    {{ embed_code("docs/queries_test.go", "count-with-condition", "\t") }}

## Pagination

REL provides a convenient `FindAndCountAll` methods that is useful for pagination, It's a combination of `FindAll` and `Count` method.
FindAndCountAll returns count of records (ignoring limit and offset query) and an error.

=== "Example"

    Retrieve all books within limit and offset and also count of all books.

    {{ embed_code("docs/queries.go", "find-and-count-all", "\t") }}

=== "Mock"

    Mock retrieve all books within limit and offset and also count of all books.

    {{ embed_code("docs/queries_test.go", "find-and-count-all", "\t") }}


## Batch Iteration

REL provides records iterator that can be use for perform batch processing of large amounts of records.

Options:

- `BatchSize` - The size of batches (default 1000).
- `Start` - The primary value (ID) to start from (inclusive).
- `Finish` - The primary value (ID) to finish at (inclusive).

=== "Example"

    {{ embed_code("docs/queries.go", "batch-iteration", "\t") }}

=== "Mock"

    Mock and returns users.

    {{ embed_code("docs/queries_test.go", "batch-iteration", "\t") }}

    Mock and retuns `reltest.ErrConnectionClosed` (`sql.ErrConnDone`).

    {{ embed_code("docs/queries_test.go", "batch-iteration-connection-error", "\t") }}


## Native SQL Query

REL allows querying using native SQL query, this is especially useful when using complex query that cannot be covered with the query builder.

=== "Example"

    Retrieve a book using native sql query.

    {{ embed_code("docs/queries.go", "sql", "\t") }}

=== "Mock"

    Mock and retrieve a book using native sql query.

    {{ embed_code("docs/queries_test.go", "sql", "\t") }}
