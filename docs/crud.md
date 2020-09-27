# Reading and Writing Record

## Create

A new record can be inserted to database using a struct, map or set function. To insert a new record using a struct, simply pass the pointer to the instance as the only argment. Insertion using struct will update `created_at` and `updated_at` field if any.

!!! note
    reltest.Repository will automatically sets any primary key value to be 1.

*Inserting a record:*

=== "Example"
    {{ embed_code("examples/crud.go", "insert", "\t") }}
=== "Mock Any"
    {{ embed_code("examples/crud_test.go", "insert", "\t") }}
=== "Mock by Record"
    {{ embed_code("examples/crud_test.go", "insert-for", "\t") }}
=== "Mock by Type"
    {{ embed_code("examples/crud_test.go", "insert-for-type", "\t") }}
=== "Mock Error"
    {{ embed_code("examples/crud_test.go", "insert-error", "\t") }}


*To inserts multiple records at once, use `InsertAll`:*

=== "Example"
    {{ embed_code("examples/crud.go", "insert-all", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "insert-all", "\t") }}

## Read

REL provides a powerful API for querying record from database. To query a record, simply use the Find method, it's accept the returned result as the first argument, and the conditions for the rest arguments.

*Retrieve a book with id 1:*

=== "Example"
    {{ embed_code("examples/crud.go", "find", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "find", "\t") }}

*Retrieve a book with id 1 using syntactic sugar:*

=== "Example"
    {{ embed_code("examples/crud.go", "find-alias", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "find-alias-error", "\t") }}

*Querying multiple records using `FindAll` method:*

=== "Example"
    {{ embed_code("examples/crud.go", "find-all", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "find-all", "\t") }}


*Using chainable query api for a more complex query use case:*

=== "Example"
    {{ embed_code("examples/crud.go", "find-all-chained", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "find-all-chained", "\t") }}

## Update

Similar to create, updating a record in REL can also be done using struct, map or set function. Updating using struct will also update `updated_at` field if any.

An update using struct will cause all fields and association to be saved to database, regardless of whether it's been updated or not. Use `rel.Map`, `rel.Set` or `rel.Changeset` to update only specific fields.

!!! note
    When updating belongs to association, it's recommended to not expose reference key (`[other]_id`) for updates directly from user, since there's no way to validate belongs to association using query.

*Updating a record:*

=== "Example"
    {{ embed_code("examples/crud.go", "update", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "update", "\t") }}

*Updating multiple records is possible using `UpdateAll`:*

=== "Example"
    {{ embed_code("examples/crud.go", "update-all", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "update-all", "\t") }}

## Delete

To delete a record in rel, simply pass the record to be deleted.

!!! note
    REL will automatically apply soft-delete if `DeletedAt time.Time` field exists in a struct. To query soft-deleted records, use `rel.Unscoped(true)` when querying.

*Deleting a record:*

=== "Example"
    {{ embed_code("examples/crud.go", "delete", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "delete", "\t") }}


*Deleting multiple records:*

=== "Example"
    {{ embed_code("examples/crud.go", "delete-all", "\t") }}
=== "Mock"
    {{ embed_code("examples/crud_test.go", "delete-all", "\t") }}
