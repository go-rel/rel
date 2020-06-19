# Transactions

To declare a transaction, use `Transaction` method. It can be called recursively to define nested transaction.
It accepts a function with `context.Context` argument which is used to determine the transaction scope.

If any error occured within transaction, the transaction will be rolled back, and returns the error.
If the error is a runtime error or `panic` with string argument, it'll panic after rollback.

<!-- tabs:start -->

### **Example**

[transactions.go](transactions.go ':include :fragment=transactions')


### **Mock**

[transactions_test.go](transactions_test.go ':include :fragment=transactions')

<!-- tabs:end -->

**Next: [Instrumentation](instrumentation.md)**
