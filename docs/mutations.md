# Mutations

## Basic Mutator

REL uses mutator to define inserts and updates operation. Using basic mutator won't update `created_at` and `updated_at` fields.

- `Dec(field string)` - Decrement a field value by 1.
- `DecBy(field string, n int)` - Decrement a field value by n.
- `Inc(field string)` - Increase a field value by 1.
- `IncBy(field string, n int)` - Increase a field value by n.
- `Set(field string, value interface{})` - Set a value to a field.
- `SetFragment(raw string, args ...interface{})` - Set a value of a field using SQL fragment.

<!-- tabs:start -->

### **Example**

Set title and category values.

[mutations.go](mutations.go ':include :fragment=basic-set')

Decrement stock.

[mutations.go](mutations.go ':include :fragment=basic-dec')

Update title using SQL fragment.

[mutations.go](mutations.go ':include :fragment=basic-fragment')

### **Mock**

Mock set title and category values.

[mutations_test.go](mutations_test.go ':include :fragment=basic-set')

Mock decrement stock.

[mutations_test.go](mutations_test.go ':include :fragment=basic-dec')

Mock update title using SQL fragment.

[mutations_test.go](mutations_test.go ':include :fragment=basic-fragment')

<!-- tabs:end -->

## Structset

Structset is a mutator that generates list of `Set` mutators based on a struct value. Using Structset will result in replacing the intire record in the database using provided struct, It'll always clear a has many association and re-insert it on updates if it's loaded. Changeset can be used to avoid clearing has many association on updates.

?> `Structset` is the default mutator used when none is provided explicitly.

<!-- tabs:start -->

### **Example**

Insert a struct.

[mutations.go](mutations.go ':include :fragment=structset')

### **Mock**

Mock insert a struct.

[mutations_test.go](mutations_test.go ':include :fragment=structset')

<!-- tabs:end -->

## Changeset

Changeset allows you to track and update only updated values and asssociation to database. This is very efficient when dealing with a complex struct that contains a lot of fields and associations.

<!-- tabs:start -->

### **Example**

Update only price and discount field using changeset.

[mutations.go](mutations.go ':include :fragment=changeset')

### **Mock**

Mock update only price and discount field using changeset.

[mutations_test.go](mutations_test.go ':include :fragment=changeset')

<!-- tabs:end -->

## Map

Map allows to define group of `Set` mutator, this is intended to be use internally and not to be exposed directly to user.

<!-- tabs:start -->

### **Example**

Insert books and its author using `Map`.

[mutations.go](mutations.go ':include :fragment=map')

### **Mock**

Mock insert books and its author using `Map`.

[mutations_test.go](mutations_test.go ':include :fragment=map')

<!-- tabs:end -->

## Reloading Updated Struct

By default, only `Inc`, `IncBy`, `Dec`, `DecBy` and `SetFragment` will reload struct from database, `Reload` mutator can be used to manually trigger reload after inserts/update operations.

<!-- tabs:start -->

### **Example**

Update title and force reload.

[mutations.go](mutations.go ':include :fragment=reload')

### **Mock**

Mock update title and force reload.

[mutations_test.go](mutations_test.go ':include :fragment=reload')

<!-- tabs:end -->

## Cascade Operations

REL supports insert/update/delete record and it's associations.

<!-- tabs:start -->

### **Example**

Disable cascade insert (default enabled).

[mutations.go](mutations.go ':include :fragment=cascade')

Enable cascade delete (default disabled).

[mutations.go](mutations.go ':include :fragment=delete-cascade')

### **Mock**

Mock disable cascade insert (default enabled).

[mutations_test.go](mutations_test.go ':include :fragment=cascade')

Mock enable cascade delete (default disabled).

[mutations_test.go](mutations_test.go ':include :fragment=delete-cascade')

<!-- tabs:end -->
