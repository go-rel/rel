# Queries

## Retrieving Data

REL provides two basic finders method, `Find` for retrieving single record, and `FindAll` for retrieving multiple record.

!!! note
    - `Find` only accepts struct as the first argument, and always return the first result from the query.
    - `FindAll` only accepts slice as the first argument, and always return all result from the query.

*Retrieve a book where id=1:*

=== "Example"
    {{ embed_code("examples/queries.go", "find", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "find", "\t") }}

*Retrieve all books:*

=== "Example"
    {{ embed_code("examples/queries.go", "find-all", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "find-all", "\t") }}

## Conditions

To retrieve filtered recods from database, you can use filter api to specify [condition](https://pkg.go.dev/github.com/go-rel/rel/where). For example, to filter all books that available, you can use `rel.Eq` in the query builder.

*Retrieve all available books using filter query:*

=== "Example"
    {{ embed_code("examples/queries.go", "condition", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "condition", "\t") }}

*Alias can be used to boost readability when dealing with short query:*

=== "Example"
    {{ embed_code("examples/queries.go", "condition-alias", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "condition-alias", "\t") }}

*Use fragment to specify custom SQL query:*

=== "Example"
    {{ embed_code("examples/queries.go", "condition-fragment", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "condition-fragment", "\t") }}

You can use `rel.And` or `rel.Or` to specify more conditions.

*Retrieve all available books where price is at least 100 or in discount using filter query:*

=== "Example"
    {{ embed_code("examples/queries.go", "condition-advanced", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "condition-advanced", "\t") }}

*Retrieve all available books where price is at least 100 or in discount using chained filter query:*

=== "Example"
    {{ embed_code("examples/queries.go", "condition-advanced-chain", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "condition-advanced-chain", "\t") }}

*Retrieve all available books where price is at least 100 or in discount using alias (`github.com/go-rel/rel/where`):*

=== "Example"
    {{ embed_code("examples/queries.go", "condition-advanced-alias", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "condition-advanced-alias", "\t") }}

## Sorting

To retrieve records from database in a specific order, you can use the sort api.

*Sort books ascending by updated_at field:*

=== "Example"
    {{ embed_code("examples/queries.go", "sorting", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "sorting", "\t") }}

*Using alias if you need more syntactic sugar:*

=== "Example"
    {{ embed_code("examples/queries.go", "sorting-alias", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "sorting-alias", "\t") }}

Combining with other query is fairly easy.

*Chain where and sort using query builder:*

=== "Example"
    {{ embed_code("examples/queries.go", "sorting-with-condition", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "sorting-with-condition", "\t") }}

*It's also possible to use variadic arguments to combine multiple queries:*

=== "Example"
    {{ embed_code("examples/queries.go", "sorting-with-condition-variadic", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "sorting-with-condition-variadic", "\t") }}

## Selecting Specific Fields

To select specific fields, you can use `Select` method, this way only specificied field will be mapped to books.

!!! note
    Specifying select without argument (`rel.Select()`) will automatically load all fields. This is helpful when used as query builder entry point (compared to using `rel.From`), because you can let REL to infer the table name.

*Load only id and title:*

=== "Example"
    {{ embed_code("examples/queries.go", "select", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "select", "\t") }}

## Using Specific Table

By default, REL will use pluralized-snakecase struct name as the table name. To select from specific table, you can use `From` method.

*Load from `ebooks` table:*

=== "Example"
    {{ embed_code("examples/queries.go", "table", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "table", "\t") }}

*Chain the query with select:*

=== "Example"
    {{ embed_code("examples/queries.go", "table-chained", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "table-chained", "\t") }}

## Limit and Offset

To set the limit and offset of query, use `Limit` and `Offset` api. `Offset` will be ignored if `Limit` is not specified.

*Specify limit and offset:*

=== "Example"
    {{ embed_code("examples/queries.go", "limit-offset", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "limit-offset", "\t") }}

*As a chainable api:*

=== "Example"
    {{ embed_code("examples/queries.go", "limit-offset-chained", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "limit-offset-chained", "\t") }}

## Group

To use group by query, you can use `Group` method.

*Retrieve count of books for every category:*

=== "Example"
    {{ embed_code("examples/queries.go", "group", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "group", "\t") }}

## Joining Tables

To join tables, you can use `join` api.

!!! note
    Joining table won't load the association to struct. If you want to load association on a struct, use [preload](/association/#preloading-association) instead.

*Join transaction and book table, then filter only transaction that have specified book name. This methods assumes belongs to relation, which means it'll try to join using `transactions.book_id=books.id`:*

=== "Example"
    {{ embed_code("examples/queries.go", "join", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "join", "\t") }}

*Specifying which column to join using JoinOn:*

=== "Example"
    {{ embed_code("examples/queries.go", "join-on", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "join-on", "\t") }}

*Syntactic sugar also available for join:*

=== "Example"
    {{ embed_code("examples/queries.go", "join-alias", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "join-alias", "\t") }}

*Joining table with custom join mode:*

=== "Example"
    {{ embed_code("examples/queries.go", "join-with", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "join-with", "\t") }}

*Use fragment for more complex join query:*

=== "Example"
    {{ embed_code("examples/queries.go", "join-fragment", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "join-fragment", "\t") }}

## Pessimistic Locking

REL supports pessimistic locking by using mechanism provided by the underlying database. `Lock` can be only used only inside transaction.

*Retrieve and lock a row for update:*

=== "Example"
    {{ embed_code("examples/queries.go", "lock", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "lock", "\t") }}

*Retrieve and lock a row using predefined lock alias:*

=== "Example"
    {{ embed_code("examples/queries.go", "lock-for-update", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "lock-for-update", "\t") }}

*Retrieve and lock a row using chained query:*

=== "Example"
    {{ embed_code("examples/queries.go", "lock-chained", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "lock-chained", "\t") }}

## Aggregation

REL provides a very basic `Aggregate` method which can be used to count, sum, max etc.

*Count all available books using aggregate:*

=== "Example"
    {{ embed_code("examples/queries.go", "aggregate", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "aggregate", "\t") }}

*Count all available books using count:*

=== "Example"
    {{ embed_code("examples/queries.go", "count", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "count", "\t") }}

*Count all available books using count:*

=== "Example"
    {{ embed_code("examples/queries.go", "count-with-condition", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "count-with-condition", "\t") }}

## Pagination

REL provides a convenient `FindAndCountAll` methods that is useful for pagination, It's a combination of `FindAll` and `Count` method.
FindAndCountAll returns count of records (ignoring limit and offset query) and an error.

*Retrieve all books within limit and offset and also count of all books:*

=== "Example"
    {{ embed_code("examples/queries.go", "find-and-count-all", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "find-and-count-all", "\t") }}


## Batch Iteration

REL provides records iterator that can be use for perform batch processing of large amounts of records.

Options:

- `BatchSize` - The size of batches (default 1000).
- `Start` - The primary value (ID) to start from (inclusive).
- `Finish` - The primary value (ID) to finish at (inclusive).

=== "Example"
    {{ embed_code("examples/queries.go", "batch-iteration", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "batch-iteration", "\t") }}
=== "Mock Error"
    {{ embed_code("examples/queries_test.go", "batch-iteration-connection-error", "\t") }}


## Native SQL Query

REL allows querying using native SQL query, this is especially useful when using complex query that cannot be covered with the query builder.

*Retrieve a book using native sql query:*

=== "Example"
    {{ embed_code("examples/queries.go", "sql", "\t") }}
=== "Mock"
    {{ embed_code("examples/queries_test.go", "sql", "\t") }}
