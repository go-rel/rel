# Reading and Writing Record

## Create

A new record can be inserted to database using a struct, map or set function. To insert a new record using a struct, simply pass the pointer to the instance as the only argment. Insertion using struct will update `created_at` and `updated_at` field if any.

=== "Example"

    {{ embed_code("docs/crud.go", "insert", "\t") }}

=== "Mock"

    ?> reltest.Repository will automatically sets any primary key value to be 1.

    Expect any insert called.

    {{ embed_code("docs/crud_test.go", "insert", "\t") }}

    Expect insertion only for a specific record.

    {{ embed_code("docs/crud_test.go", "insert-for", "\t") }}

    Expect insertion only for a specific type.

    {{ embed_code("docs/crud_test.go", "insert-for-type", "\t") }}

    Expect insertion to to return an error.

    {{ embed_code("docs/crud_test.go", "insert-error", "\t") }}


To inserts multiple records at once, use `InsertAll`.

=== "Example"

    {{ embed_code("docs/crud.go", "insert-all", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "insert-all", "\t") }}

## Read

REL provides a powerful API for querying record from database. To query a single record, simply use the Find method, it's accept the returned result as the first argument, and the conditions for the rest arguments.

=== "Example"

    Retrieve a book with id 1.

    {{ embed_code("docs/crud.go", "find", "\t") }}

    Retrieve a book with id 1 using syntactic sugar.

    {{ embed_code("docs/crud.go", "find-alias", "\t") }}

=== "Mock"

    Mock retrieve a book with id 1.

    {{ embed_code("docs/crud_test.go", "find", "\t") }}

    Mock retrieve a book with id 1 using syntactic sugar and returns error.

    {{ embed_code("docs/crud_test.go", "find-alias-error", "\t") }}


To query multiple records, use `FindAll` method.

=== "Example"

    {{ embed_code("docs/crud.go", "find-all", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "find-all", "\t") }}


REL also support chainable query api for a more complex query use case.

=== "Example"

    {{ embed_code("docs/crud.go", "find-all-chained", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "find-all-chained", "\t") }}

## Update

Similar to create, updating a record in REL can also be done using struct, map or set function. Updating using struct will also update `updated_at` field if any.

An update using struct will cause all fields and association to be saved to database, regardless of whether it's been updated or not. Use `rel.Map`, `rel.Set` or `rel.Changeset` to update only specific fields.

?> When updating belongs to association, it's recommended to not expose reference key (`[other]_id`) for updates directly from user, since there's no way to validate belongs to association using query.

=== "Example"

    {{ embed_code("docs/crud.go", "update", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "update", "\t") }}


Updating multiple records is possible using `UpdateAll`.

=== "Example"

    {{ embed_code("docs/crud.go", "update-all", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "update-all", "\t") }}

## Delete

To delete a record in rel, simply pass the record to be deleted.

?> REL will automatically apply soft-delete if `DeletedAt time.Time` field exists in a struct. To query soft-deleted records, use `rel.Unscoped(true)` when querying.

=== "Example"

    {{ embed_code("docs/crud.go", "delete", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "delete", "\t") }}


Deleting multiple records is possible using `DeleteAll`.

=== "Example"

    {{ embed_code("docs/crud.go", "delete-all", "\t") }}

=== "Mock"

    {{ embed_code("docs/crud_test.go", "delete-all", "\t") }}
