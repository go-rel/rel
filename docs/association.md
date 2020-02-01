# Association

## Defining Association

Association in rel can be declared by ensuring that each association have an association field, reference id field and foreign id field.
Association field is a field with the type of another struct.
Reference id is an id field that can be mapped to the foreign id field in another struct.
By following that convention, rel currently supports `belongs to`, `has one` and `has many` association.

```go
type User struct {
	ID        int
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time

	// has many transactions.
	// with custom reference and foreign field declaration.
	// ref: id refers to User.ID field.
	// fk: buyer_id refers to Transaction.BuyerID
	Transactions []Transaction `ref:"id" fk:"buyer_id"`

	// has one address.
	// doesn't contains primary key of other struct.
	// rel can guess the reference and foreign field if it's not specified.
	Address Address
}

type Transaction struct {
	ID   int
	Paid bool

	// belongs to user.
	// contains primary key of other struct.
	Buyer   User `ref:"buyer_id" fk:"id"`
	BuyerID int
}

type Address struct {
	ID   int
	City string

	// belongs to user.
	User   *User
	UserID *int
}
```

## Preloading Association

Preload will load association to structs. To preload association, use `Preload`.

<!-- tabs:start -->

### **main.go**

```go
// preload transaction `belongs to` buyer association.
repo.Preload(ctx, &transaction, "buyer")

// preload user `has one` address association
repo.Preload(ctx, &user, "address")

// preload user `has many` transactions association
repo.Preload(ctx, &user, "transactions")

// preload paid transactions from user.
repo.Preload(ctx, &user, "transactions", where.Eq("paid", true))

// preload every buyer's address in transactions.
// note: buyer needs to be preloaded before preloading buyer's address.
repo.Preload(ctx, &transactions, "buyer.address")
```

### **main_test.go**

```go
// preload transaction `belongs to` buyer association.
repo.ExpectPreload("buyer").Result(user)

// preload user `has one` address association
repo.ExpectPreload("address").Result(address)

// preload user `has many` transactions association
repo.ExpectPreload("transactions").Result(transactions)

// preload paid transactions from user.
repo.ExpectPreload(&user, "transactions", where.Eq("paid", true)).Result(transactions)

// preload every buyer's address in transactions.
// note: buyer needs to be preloaded before preloading buyer's address.
repo.ExpectPreload("buyer.address").Result(addresses)
```

<!-- tabs:end -->

## Modifying Association

rel will automatically creates or updates association by using `Insert` or `Update` method. If `ID` of association struct is not a zero value, rel will try to update the association, else it'll create a new association.

rel will try to create a new record for association if `ID` is a zero value.

<!-- tabs:start -->

### **main.go**

```go
user := User{
    Name: "rel",
    Address: Address{
        City: "Bandung",
    },
}

// Inserts a new record to users and address table.
// Result: User{ID: 1, Name: "rel", Address: Address{ID: 1, City: "Bandung", UserID: 1}}
repo.Insert(ctx, &user)
```

### **main_tesst.go**

```go
repo.ExpectInsert().For(&user)
```

<!-- tabs:end -->

rel will try to update a new record for association if `ID` is a zero value. To update association, it first needs to be preloaded.

<!-- tabs:start -->

### **main.go**

```go
userId := 1
user := User{
    ID:   1,
    Name: "rel",
    // association is loaded when the primary key (id) is not zero.
    Address: Address{
        ID:     1,
        UserID: &userId,
        City:   "Bandung",
    },
}

// Update user record with id 1.
// Update address record with id 1.
repo.Update(ctx, &user)
```

### **main_tesst.go**

```go
repo.ExpectUpdate().For(&user)
```

<!-- tabs:end -->

To selectively update only specific fields or association, `use rel.Map`.

<!-- tabs:start -->

### **main.go**

```go
modification := rel.Map{
    "address": rel.Map{
        "city": "bandung",
    },
}

// Update address record with id 1, only set city to bandung.
repo.Update(ctx, &user, modification)
```

### **main_test.go**

```go
modification := rel.Map{
    "address": rel.Map{
        "city": "bandung",
    },
}

// Update address record with id 1, only set city to bandung.
repo.ExpectUpdate(modification).For(&user)
```

<!-- tabs:end -->

**Next: [Transactions](transactions.md)**
