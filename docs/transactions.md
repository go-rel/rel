# Transactions

To declare a transaction, use `Transaction` method. It accepts a function with `rel.Repository` argument and returns an error.

If any error occured within transaction, the transaction will be rolled back, and returns the error. If the error is a runtime error or `panic` with string argument, it'll panic after rollback.

<!-- tabs:start -->

### **main.go**

```go
if err := repo.Transaction(func(repo rel.Repository) error {
    repo.Update(&books, rel.Dec("stock"))
    return repo.Update(&transaction, rel.Set("paid", true))
}); err != nil {
    // handle error
}
```

### **main_test.go**

```go
repo.ExpectTransaction(func(repo *Repository) {
    repo.ExpectUpdate(rel.Dec("stock")).ForType("main.Book")
    repo.ExpectUpdate(rel.Set("paid", true)).ForType("main.Transaction")
})
```

<!-- tabs:end -->
