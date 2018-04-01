// Package grimoire is a data access layer and validation for go.
//
// Quick Start:
//   package main
//
//   import (
//   	"time"
//   	"github.com/Fs02/grimoire"
//   	. "github.com/Fs02/grimoire/c"
//   	"github.com/Fs02/grimoire/adapter/mysql"
//   	"github.com/Fs02/grimoire/changeset"
//   )
//
//   type Product struct {
//   	ID        int
//   	Name      string
//   	Price     int
//   	CreatedAt time.Time
//   	UpdatedAt time.Time
//   }
//
//   func ProductChangeset(product interface{}, params map[string]interface{}) *changeset.Changeset {
//   	ch := changeset.Cast(product, params, []string{"name", "price"})
//   	changeset.ValidateRequired(ch, []string{"name", price})
//   	changeset.ValidateMin(ch, "price", 100)
//   	return ch
//   }
//
//   func main() {
//   	// initialize mysql adapter
//   	adapter, err := mysql.Open("root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
//   	if err != nil {
//   		panic(err)
//   	}
//   	defer adapter.Close()
//
//   	// initialize grimoire's repo
//   	repo := grimoire.New(adapter)
//
//   	var product Product
//
//   	// Changeset is used when creating or updating your data
//   	ch := ProductChangeset(product, map[string]interface{}{
//   		"name": "shampoo",
//   		"price": 1000
//   	})
//
//   	if ch.Error() != nil {
//   		// do something
//   	}
//
//   	// Create
//
//   	// Create products with changeset and return the result to &product
//   	repo.From("products").MustCreate(&product, ch)
//
//   	// or create without returning
//   	repo.From("products").MustCreate(nil, ch)
//
//   	// or create without changeset
//   	repo.From("products").Set("name", "shampoo").Set("price", 1000).MustCreate(nil)
//
//   	// Read
//
//   	// Find a product with id 1
//   	repo.From("products").Find(1).MustOne(&product)
//
//   	// Find() is a shortcut for this
//   	// this equal: SELECT * FROM products WHERE id=1 LIMIT 1;
//   	const id = I("id")
//   	repo.From("products").Where(Eq(id, 1)).MustOne(&product)
//
//   	// More advanced query that returns array of results
//   	// this equal: SELECT * FROM products WHERE (name="shampoo" AND price<1000) OR (name<>"shampoo" AND price>1000);
//   	var products []Product
//   	const name = I("name")
//   	const price = I("price")
//   	repo.From("products").Where(Eq(name, "shampoo"), Lt(price, 1000)).
//   		OrWhere(Ne(name, "shampoo"), Gt(price, 1000)).
//   		All(products)
//
//   	// Update
//
//   	// Update products with id=1
//   	repo.From("products").Find(1).MustUpdate(&product, ch)
//
//   	// Update products without returning
//   	repo.From("products").Find(1).MustUpdate(nil, ch)
//
//   	// Or Update products without changset
//   	repo.From("products").Find(1).Set("name", "shampoo").Set("price", 1000).MustUpdate(nil)
//
//   	// If no condition is specified, all records will be updated
//   	repo.From("products").MustUpdate(&product, ch)
//
//   	// Delete
//   	// Delete Product with id=1
//   	repo.From("products").Find(1).MustDelete()
//   	// If no condition is specified, all records will be deleted
//   	repo.From("products").MustDelete()
//   }
package grimoire
