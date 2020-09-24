# Mutations

## Basic Mutator

REL uses mutator to define inserts and updates operation. Using basic mutator won't update `created_at` and `updated_at` fields.

- `Dec(field string)` - Decrement a field value by 1.
- `DecBy(field string, n int)` - Decrement a field value by n.
- `Inc(field string)` - Increase a field value by 1.
- `IncBy(field string, n int)` - Increase a field value by n.
- `Set(field string, value interface{})` - Set a value to a field.
- `SetFragment(raw string, args ...interface{})` - Set a value of a field using SQL fragment.

=== "Example"

    Set title and category values.

    {{ embed_code("docs/mutations.go", "basic-set", "\t") }}

    Decrement stock.

    {{ embed_code("docs/mutations.go", "basic-dec", "\t") }}

    Update title using SQL fragment.

    {{ embed_code("docs/mutations.go", "basic-fragment", "\t") }}

=== "Mock"

    Mock set title and category values.

    {{ embed_code("docs/mutations_test.go", "basic-set", "\t") }}

    Mock decrement stock.

    {{ embed_code("docs/mutations_test.go", "basic-dec", "\t") }}

    Mock update title using SQL fragment.

    {{ embed_code("docs/mutations_test.go", "basic-fragment", "\t") }}

## Structset

Structset is a mutator that generates list of `Set` mutators based on a struct value. Using Structset will result in replacing the intire record in the database using provided struct, It'll always clear a has many association and re-insert it on updates if it's loaded. Changeset can be used to avoid clearing has many association on updates.

?> `Structset` is the default mutator used when none is provided explicitly.

=== "Example"

    Insert a struct.

    {{ embed_code("docs/mutations.go", "structset", "\t") }}

=== "Mock"

    Mock insert a struct.

    {{ embed_code("docs/mutations_test.go", "structset", "\t") }}

## Changeset

Changeset allows you to track and update only updated values and asssociation to database. This is very efficient when dealing with a complex struct that contains a lot of fields and associations.

=== "Example"

    Update only price and discount field using changeset.

    {{ embed_code("docs/mutations.go", "changeset", "\t") }}

=== "Mock"

    Mock update only price and discount field using changeset.

    {{ embed_code("docs/mutations_test.go", "changeset", "\t") }}

## Map

Map allows to define group of `Set` mutator, this is intended to be use internally and not to be exposed directly to user. Mutation defined in the map will be applied to the struct passed as the first argument. Insert/Update using map wont update `created_at` or `updated_at` field.

=== "Example"

    Insert books and its author using `Map`.

    {{ embed_code("docs/mutations.go", "map", "\t") }}

=== "Mock"

    Mock insert books and its author using `Map`.

    {{ embed_code("docs/mutations_test.go", "map", "\t") }}

## Reloading Updated Struct

By default, only `Inc`, `IncBy`, `Dec`, `DecBy` and `SetFragment` will reload struct from database, `Reload` mutator can be used to manually trigger reload after inserts/update operations.

=== "Example"

    Update title and force reload.

    {{ embed_code("docs/mutations.go", "reload", "\t") }}

=== "Mock"

    Mock update title and force reload.

    {{ embed_code("docs/mutations_test.go", "reload", "\t") }}

## Cascade Operations

REL supports insert/update/delete record and it's associations.

=== "Example"

    Disable cascade insert (default enabled).

    {{ embed_code("docs/mutations.go", "cascade", "\t") }}

    Enable cascade delete (default disabled).

    {{ embed_code("docs/mutations.go", "delete-cascade", "\t") }}

=== "Mock"

    Mock disable cascade insert (default enabled).

    {{ embed_code("docs/mutations_test.go", "cascade", "\t") }}

    Mock enable cascade delete (default disabled).

    {{ embed_code("docs/mutations_test.go", "delete-cascade", "\t") }}
