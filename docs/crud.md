# Reading and Writing Record

## Create

A new record can be inserted to database using a struct, map or set function. To insert a new record using a struct, simply pass the pointer to the instance as the only argment. Insertion using struct will update `created_at` and `updated_at` field if any.

<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=insert')

### **Mock**

?> reltest.Repository will automatically sets any primary key value to be 1.

Expect any insert called.

[crud_test.go](crud_test.go ':include :fragment=insert')

Expect insertion only for a specific record.

[crud_test.go](crud_test.go ':include :fragment=insert-for')

Expect insertion only for a specific type.

[crud_test.go](crud_test.go ':include :fragment=insert-for-type')

Expect insertion to to return an error.

[crud_test.go](crud_test.go ':include :fragment=insert-error')

<!-- tabs:end -->

To inserts multiple records at once, use `InsertAll`.

<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=insert-all')

### **Mock**

[crud_test.go](crud_test.go ':include :fragment=insert-all')

<!-- tabs:end -->

## Read

REL provides a powerful API for querying record from database. To query a single record, simply use the Find method, it's accept the returned result as the first argument, and the conditions for the rest arguments.


<!-- tabs:start -->

### **Example**

Retrieve a book with id 1.

[crud.go](crud.go ':include :fragment=find')

Retrieve a book with iid 1 using syntactic sugar.

[crud.go](crud.go ':include :fragment=find-alias')


### **Mock**

Mock retrieve a book with id 1.

[crud_test.go](crud_test.go ':include :fragment=find')

Mock retrieve a book with id 1 using syntactic sugar and returns error.

[crud_test.go](crud_test.go ':include :fragment=find-alias-error')

<!-- tabs:end -->

To query multiple records, use `FindAll` method.


<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=find-all')


### **Mock**

[crud_test.go](crud_test.go ':include :fragment=find-all')

<!-- tabs:end -->

REL also support chainable query api for a more complex query use case.


<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=find-all-chained')

### **Mock**

[crud_test.go](crud_test.go ':include :fragment=find-all-chained')

<!-- tabs:end -->

## Update

Similar to create, updating a record in REL can also be done using struct, map or set function. Updating using struct will also update `updated_at` field if any.

An update using struct will cause all fields and association to be saved to database, regardless of whether it's been updated or not. Use `rel.Map`, `rel.Set` or `rel.Changeset` to update only specific fields.

?> When updating belongs to association, it's recommended to not expose reference key (`[other]_id`) for updates directly from user, since there's no way to validate belongs to association using query.

<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=update')

### **Mock**

[crud_test.go](crud_test.go ':include :fragment=update')

<!-- tabs:end -->

## Delete

To delete a record in rel, simply pass the record to be deleted.

?> REL will automatically apply soft-delete if `DeletedAt time.Time` field exists in a struct. To query soft-deleted records, use `rel.Unscoped(true)` when querying.

<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=delete')

### **Mock**

[crud_test.go](crud_test.go ':include :fragment=delete')

<!-- tabs:end -->

Deleting multiple records is possible using `DeleteAll`.


<!-- tabs:start -->

### **Example**

[crud.go](crud.go ':include :fragment=delete-all')

### **Mock**

[crud_test.go](crud_test.go ':include :fragment=delete-all')

<!-- tabs:end -->


**Next: [Queries](queries.md)**
