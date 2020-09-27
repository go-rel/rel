# Association

## Defining Association

Association in REL can be declared by ensuring that each association have an association field, reference id field and foreign id field.
Association field is a field with the type of another struct.
Reference id is an id field that can be mapped to the foreign id field in another struct.
By following that convention, REL currently supports `belongs to`, `has one` and `has many` association.

{{ embed_code("examples/association.go","association-schema") }}

## Preloading Association

Preload will load association to structs. To preload association, use `Preload`.

*Preload Transaction's Buyer (`belongs to` association):*

=== "Example"
    {{ embed_code("examples/association.go","preload-belongs-to", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "preload-belongs-to", "\t") }}

*Preload User's Address (`has one` association):*

=== "Example"
    {{ embed_code("examples/association.go","preload-has-one", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "preload-has-one", "\t") }}

*Preload User's Transactions (`has many` association):*

=== "Example"
    {{ embed_code("examples/association.go","preload-has-many", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "preload-has-many", "\t") }}

*Preload only paid Transactions from users:*

=== "Example"
    {{ embed_code("examples/association.go","preload-has-many-filter", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "preload-has-many-filter", "\t") }}

*Preload every Buyer's Address in Transactions (Buyer needs to be preloaded before preloading Buyer's Address):*

=== "Example"
    {{ embed_code("examples/association.go","preload-nested", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "preload-nested", "\t") }}

## Inserting and Updating Association

REL will automatically creates or updates association by using `Insert` or `Update` method. If `ID` of association struct is not a zero value, REL will try to update the association, else it'll create a new association.

!!! note
    REL will try to create a new record for association if Primary Value (`ID`) is a zero value.

=== "Example"
    {{ embed_code("examples/association.go","insert-association", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "insert-association", "\t") }}


REL will try to update a new record for association if `ID` is a zero value. To update association, it first needs to be preloaded.

=== "Example"
    {{ embed_code("examples/association.go","update-association", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "update-association", "\t") }}


To selectively update only specific fields or association, `use rel.Map`.

=== "Example"
    {{ embed_code("examples/association.go","update-association-with-map", "\t") }}
=== "Mock"
    {{ embed_code("examples/association_test.go", "update-association-with-map", "\t") }}
