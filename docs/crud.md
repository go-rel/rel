## Reading and Writing Data

### Create

A new record can be inserted to database using a struct, map or set function. To insert a new record using a struct, simply pass the pointer to the instance as the only argment. Insertion using struct will update `created_at` and `updated_at` field if any.

<!-- tabs:start -->

#### **main.go**

```golang
user := Book{
    Name:     "Rel for dummies",
    Category: "education",
}

// Insert directly using struct.
if err := repo.Insert(&user); err != nil {
    // handle error
}
```

#### **main_test.go**

> reltest.Repository will automatically sets any primary key value to be 1.

```golang
// Expect any insert called.
repo.ExpectInsert()

// OR: Expect insertion for a specific type.
repo.ExpectInsert().ForType("main.Book")

// OR: Expect insertion for a specific record.
repo.ExpectInsert().For(&Book{
    Name:     "Rel for dummies",
    Category: "education",
})

// OR: Expect it to return an error.
repo.ExpectInsert().ForType("main.Book").Error(errors.New("oops!"))

// Assert all expectation is called.
repo.AssertExpectations(t)
```

<!-- tabs:end -->

To insert a new record using a map, simply pass a `rel.Map` as the second argument, changes defined in the map will be applied to the struct passed as the first argument. Insertion using map wont update `created_at` asnd `updated_at` field.

<!-- tabs:start -->

#### **main.go**

```golang
var user User
data := rel.Map{
    "name":     "Rel for dummies",
    "category": "education",
}

// Insert using map.
if err := repo.Insert(&user, data); err != nil {
    // handle error
}
```

#### **main_test.go**

> reltest.Repository will automatically populate record using value provided by map.

```golang
// Expect insertion with given changer.
repo.ExpectInsert(rel.Map{
    "name":     "Rel for dummies",
    "category": "education",
}).ForType("main.Book")

// Assert all expectation is called.
repo.AssertExpectations(t)
```

<!-- tabs:end -->

It's also possible to insert a new record manually using `rel.Set`, which is a very basic type of `changer`.

<!-- tabs:start -->

#### **main.go**

```golang
// Insert using set.
if err := repo.Insert(&user, rel.Set("name", "Rel for dummies"), rel.Set("category", "education")); err != nil {
    // handle error
}
```

#### **main_test.go**

```golang
// Expect insertion with given changer.
repo.ExpectInsert(
    rel.Set("name", "Rel for dummies"),
    rel.Set("category", "education"),
).ForType("main.Book")

// Assert all expectation is called.
repo.AssertExpectations(t)
```

<!-- tabs:end -->
