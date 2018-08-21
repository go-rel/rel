package grimoire_test

import (
	"time"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/mysql"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/params"
)

type Product struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// ChangeProduct prepares data before database operation.
// Such as casting value to appropriate types and perform validations.
func ChangeProduct(product interface{}, params params.Params) *changeset.Changeset {
	ch := changeset.Cast(product, params, []string{"name", "price"})
	changeset.ValidateRequired(ch, []string{"name", "price"})
	changeset.ValidateMin(ch, "price", 100)
	return ch
}

func Example() {
	// initialize mysql adapter.
	adapter, err := mysql.Open("root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize grimoire's repo.
	repo := grimoire.New(adapter)

	var product Product

	// Inserting Products.
	// Changeset is used when creating or updating your data.
	ch := ChangeProduct(product, params.Map{
		"name":  "shampoo",
		"price": 1000,
	})

	if ch.Error() != nil {
		// handle error
	}

	// Changeset can also be created directly from json string.
	jsonch := ChangeProduct(product, params.ParseJSON(`{
		"name":  "soap",
		"price": 2000,
	}`))

	// Create products with changeset and return the result to &product,
	if err = repo.From("products").Insert(&product, ch); err != nil {
		// handle error
	}

	// or panic when insertion pailed
	repo.From("products").MustInsert(&product, jsonch)

	// Querying Products.
	// Find a product with id 1.
	repo.From("products").Find(1).MustOne(&product)

	// Updating Products.
	// Update products with id=1.
	repo.From("products").Find(1).MustUpdate(&product, ch)

	// Deleting Products.
	// Delete Product with id=1.
	repo.From("products").Find(1).MustDelete()
}
