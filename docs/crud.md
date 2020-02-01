# Reading and Writing Record

## Create

A new record can be inserted to database using a struct, map or set function. To insert a new record using a struct, simply pass the pointer to the instance as the only argment. Insertion using struct will update `created_at` and `updated_at` field if any.

<!-- tabs:start -->

### **main.go**

```go
book := Book{
    Title:    "Rel for dummies",
    Category: "education",
}

// Insert directly using struct.
if err := repo.Insert(ctx, &book); err != nil {
    // handle error
}
```

### **main_test.go**

> reltest.Repository will automatically sets any primary key value to be 1.

```go
// Expect any insert called.
repo.ExpectInsert()

// OR: Expect insertion for a specific type.
repo.ExpectInsert().ForType("main.Book")

// OR: Expect insertion for a specific record.
repo.ExpectInsert().For(&Book{
    Title:    "Rel for dummies",
    Category: "education",
})

// OR: Expect it to return an error.
repo.ExpectInsert().ForType("main.Book").Error(errors.New("oops!"))

// Assert all expectation is called.
repo.AssertExpectations(t)
```

<!-- tabs:end -->

To insert a new record using a map, simply pass a `rel.Map` as the second argument, modification defined in the map will be applied to the struct passed as the first argument. Insertion using map wont update `created_at` or `updated_at` field.

<!-- tabs:start -->

### **main.go**

```go
var book Book
data := rel.Map{
    "title":    "Rel for dummies",
    "category": "education",
}

// Insert using map.
repo.Insert(ctx, &book, data)
```

### **main_test.go**

> reltest.Repository will automatically populate record using value provided by map.

```go
// Expect insertion with given modifier.
repo.ExpectInsert(rel.Map{
    "title":    "Rel for dummies",
    "category": "education",
}).ForType("main.Book")
```

<!-- tabs:end -->

It's also possible to insert a new record manually using `rel.Set`, which is a very basic type of `modifier`.

<!-- tabs:start -->

### **main.go**

```go
// Insert using set.
repo.Insert(ctx, &book, rel.Set("title", "Rel for dummies"), rel.Set("category", "education"))
```

### **main_test.go**

```go
// Expect insertion with given modifier.
repo.ExpectInsert(
    rel.Set("title", "Rel for dummies"),
    rel.Set("category", "education"),
).ForType("main.Book")
```

<!-- tabs:end -->

To inserts multiple records at once, use `InsertAll`.


<!-- tabs:start -->

### **main.go**

```go
// InsertAll books.
repo.InsertAll(ctx, &books)
```

### **main_test.go**

```go
// Expect any insert all.
repo.ExpectInsertAll()
```

<!-- tabs:end -->


## Read

rel provides a powerful API for querying record from database. To query a single record, simply use the Find method, it's accept the returned result as the first argument, and the conditions for the rest arguments.


<!-- tabs:start -->

### **main.go**

```go
repo.Find(ctx, &book, rel.Eq("id", 1))

// OR: with sugar alias
repo.Find(ctx, &book, where.Eq("id", 1))
```

### **main_test.go**

```go
// Expect a find query and mock the result.
repo.ExpectFind(rel.Eq("id", 1)).Result(book)

// OR: Expect a find query and returns rel.NotFoundError
repo.ExpectFind(where.Eq("id", 1)).NotFound()
```

<!-- tabs:end -->

To query multiple records, use `FindAll` method.


<!-- tabs:start -->

### **main.go**

```go
repo.FindAll(ctx, &books, where.Like("title", "%dummies%").AndEq("category", "education"), rel.Limit(10))
```

### **main_test.go**

```go
// Expect a find all query and mock the result.
repo.ExpectFindAll(where.Like("title", "%dummies%").AndEq("category", "education"), rel.Limit(10))).Result(books)
```

<!-- tabs:end -->

rel also support chainable query api for a more complex query use case.


<!-- tabs:start -->

### **main.go**

```go
query := rel.Select("title", "category").Where(where.Eq("category", "education")).SortAsc("title")
repo.FindAll(ctx, &books, query)
```

### **main_test.go**

```go
// Expect a find all query and mock the result.
query := rel.Select("title", "category").Where(where.Eq("category", "education")).SortAsc("title")
repo.ExpectFindAll(query).Result(books)
```

<!-- tabs:end -->

## Update

Similar to create, updating a record in rel can also be done using struct, map or set function. Updating using struct will also update `updated_at` field if any.

> An update using struct will cause all fields to be saved to database, regardless of whether it's been updated or not. Use `rel.Map`, `rel.Set` or `rel.Structset` to update only specified fields.

<!-- tabs:start -->

### **main.go**

```go
// Update directly using struct.
repo.Update(ctx, &book)
```

### **main_test.go**

```go
// Expect any update is called.
repo.ExpectUpdate()
```

<!-- tabs:end -->

Besides struct, map and set function. There's also increment and decrement modifier to atomically increment/decrement any value in database.

<!-- tabs:start -->

### **main.go**

```go
// Update directly using struct.
repo.Update(ctx, &book, rel.Inc("views"))
```

### **main_test.go**

```go
// Expect any update is called.
repo.ExpectUpdate(rel.Inc("views"))
```

<!-- tabs:end -->

## Delete

To delete a record in rel, simply pass the record to be deleted.

<!-- tabs:start -->

### **main.go**

```go
// Delete a record.
repo.Delete(ctx, &book)
```

### **main_test.go**

```go
// Expect book to be deleted.
repo.ExpectDelete().For(&book)
```

<!-- tabs:end -->

Deleting multiple records is possible using `DeleteAll`.


<!-- tabs:start -->

### **main.go**

```go
// We have manually define the table here.
repo.DeleteAll(ctx, rel.From("books").Where(where.Eq("id", 1)))
```

### **main_test.go**

```go
// Expect books to be deleted.
repo.ExpectDeleteAll(rel.From("books").Where(where.Eq("id", 1)))
```

<!-- tabs:end -->


**Next: [Query Interface](query.md)**
