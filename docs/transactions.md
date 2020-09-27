# Transactions

To declare a transaction, use the `Transaction` method which can be called recursively to define nested transactions.

Rel accepts a function with `context.Context` argument that is used to determine the transaction scope.
Context makes it easier to call any function that involves db operation inside a transaction, because the scope of transaction is automatically passed by context.

If any error occured within transaction, the transaction will be rolled back, and returns the error.
If the error is a runtime error or `panic` with string argument, it'll panic after rollback.

=== "Example"
    {{ embed_code("examples/transactions.go", "transactions", "\t") }}
=== "Mock"
    {{ embed_code("examples/transactions_test.go", "transactions", "\t") }}
