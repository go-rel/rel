# Association

## Defining Association

Association in rel can be declared by ensuring that each association have an association field, reference id field and foreign id field.
Association field is a field with the type of another struct.
Reference id is an id field that can be mapped to the foreign id field in another struct.
By following that convention, rel currently supports `belongs to`, `has one` and `has many` association.

```go
type User struct {
	ID           int
	Name         string
    Age          int
    // has many transactions.
    // with custom reference and foreign field declaration.
    // references: id refers to User.ID field.
    // foreign_key: buyer_id refers to Transaction.BuyerID
    Transactions []Transaction `references:"id" foreign_key:"buyer_id"`
    // has one address.
    // doesn't contains primary key of other struct.
    // rel can guess the reference and foreign field if it's not specified.
    Address      Address
    CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID      int
    // belongs to user.
    // contains primary key of other struct.
	Buyer   User `ref:"buyer_id" foreign_key:"id"`
    BuyerID int
    Paid    bool
}

type Address struct {
	ID     int
    City   string
    // belongs to user.
	User   *User
	UserID *int
}
```

## Preloading

Preload will load association to structs. To preload association, use `Preload`.

<!-- tabs:start -->

### **main.go**

```go
// preload transaction `belongs to` buyer association.
repo.Preload(&transaction, "buyer")

// preload user `has one` address association
repo.Preload(&user, "address")

// preload user `has many` transactions association
repo.Preload(&user, "transactions")

// preload every buyer's address in transactions.
// note: buyer needs to be preloaded before preloading buyer's address.
repo.Preload(&transactions, "buyer.address")
```

### **main_test.go**

```go
// preload transaction `belongs to` buyer association.
repo.ExpectPreload("buyer").Result(user)

// preload user `has one` address association
repo.ExpectPreload("address").Result(address)

// preload user `has many` transactions association
repo.ExpectPreload("transactions").Result(transactions)

// preload every buyer's address in transactions.
// note: buyer needs to be preloaded before preloading buyer's address.
repo.ExpectPreload("buyer.address").Result(addresses)
```

<!-- tabs:end -->

## Modifying Associations
