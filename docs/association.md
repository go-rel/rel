# Association

## Defining Association

Association in REL can be declared by ensuring that each association have an association field, reference id field and foreign id field.
Association field is a field with the type of another struct.
Reference id is an id field that can be mapped to the foreign id field in another struct.
By following that convention, REL currently supports `belongs to`, `has one` and `has many` association.

[association.go](association.go ':include :fragment=association-schema')

## Preloading Association

Preload will load association to structs. To preload association, use `Preload`.

<!-- tabs:start -->

### **Example**

Preload Transaction's Buyer (`belongs to` association).

[association.go](association.go ':include :fragment=preload-belongs-to')

Preload User's Address (`has one` association).

[association.go](association.go ':include :fragment=preload-has-one')

Preload User's Transactions (`has many` association).

[association.go](association.go ':include :fragment=preload-has-many')

Preload only paid Transactions from users.

[association.go](association.go ':include :fragment=preload-has-many-filter')

Preload every Buyer's Address in Transactions.

**Note:** Buyer needs to be preloaded before preloading Buyer's Address.

[association.go](association.go ':include :fragment=preload-nested')

### **Mock**

Mock preload Transaction's Buyer (`belongs to` association).

[association_test.go](association_test.go ':include :fragment=preload-belongs-to')

Mock preload User's Address (`has one` association).

[association_test.go](association_test.go ':include :fragment=preload-has-one')

Preload User's Transactions (`has many` association).

[association_test.go](association_test.go ':include :fragment=preload-has-many')

Mock preload only paid Transactions from users.

[association_test.go](association_test.go ':include :fragment=preload-has-many-filter')

Mock preload every Buyer's Address in Transactions.

**Note:** Address will be assigned based on `UserID` (association key).

[association_test.go](association_test.go ':include :fragment=preload-nested')

<!-- tabs:end -->

## Inserting and Updating Association

REL will automatically creates or updates association by using `Insert` or `Update` method. If `ID` of association struct is not a zero value, REL will try to update the association, else it'll create a new association.

> REL will try to create a new record for association if Primary Value (`ID`) is a zero value.

<!-- tabs:start -->

### **Example**

[association.go](association.go ':include :fragment=insert-association')

### **Mock**

[association_test.go](association_test.go ':include :fragment=insert-association')

<!-- tabs:end -->

REL will try to update a new record for association if `ID` is a zero value. To update association, it first needs to be preloaded.

<!-- tabs:start -->

### **Example**

[association.go](association.go ':include :fragment=update-association')

### **Mock**

[association_test.go](association_test.go ':include :fragment=update-association')

<!-- tabs:end -->

To selectively update only specific fields or association, `use rel.Map`.

<!-- tabs:start -->

### **Example**

[association.go](association.go ':include :fragment=update-association-with-map')

### **Mock**

[association_test.go](association_test.go ':include :fragment=update-association-with-map')

<!-- tabs:end -->

**Next: [Transactions](transactions.md)**
