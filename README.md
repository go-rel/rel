# grimoire
[![GoDoc](https://godoc.org/github.com/Fs02/grimoire?status.svg)](https://godoc.org/github.com/Fs02/grimoire) [![Build Status](https://travis-ci.org/Fs02/grimoire.svg?branch=master)](https://travis-ci.org/Fs02/grimoire) [![Go Report Card](https://goreportcard.com/badge/github.com/Fs02/grimoire)](https://goreportcard.com/report/github.com/Fs02/grimoire) [![Maintainability](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/maintainability)](https://codeclimate.com/github/Fs02/grimoire/maintainability) [![Test Coverage](https://api.codeclimate.com/v1/badges/d487e2be0ed7b0b1fed1/test_coverage)](https://codeclimate.com/github/Fs02/grimoire/test_coverage)

Grimoire is a database access layer and validation for go. Grimoire is not an ORM but it gives a nice way to deal with relations.

Features:

- Query Builder
- Struct style create and update
- Changeset Style create and update
- Builtin validation using changeset
- Multi adapter support
- Logger

## Table of contents

<!--ts-->
   * [Install](#install)
   * [Quick Start](#quick-start)
   * [Connecting to a database](#connecting-to-database)
   * [CRUD Interface](#crud-interface)
      * [Create](#create)
      * [Query](#query)
      * [Update](#update)
      * [Delete](#delete)
   * [Transaction](#transaction)
   * [Logger](#logger)
   * [Field Mapping](#field-mapping)
<!--te-->

## Install

```bash
go get github.com/Fs02/grimoire
```

## Quick Start

```golang
package main

import (
	"time"
	"github.com/Fs02/grimoire"
	. "github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/adapter/mysql"
	"github.com/Fs02/grimoire/changeset"
)

type Product struct {
	ID        int
	Name      string
	Price     int
	CreatedAt time.Time
	UpdatedAt time.Time
}

func ProductChangeset(product interface{}, params map[string]interface{}) *changeset.Changeset {
	ch := changeset.Cast(product, params, []string{"name", "price"})
	changeset.ValidateRequired(ch, []string{"name", "price"})
	changeset.ValidateMin(ch, "price", 100)
	return ch
}

func main() {
	// initialize mysql adapter.
	adapter, err := mysql.Open("root@(127.0.0.1:3306)/db?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize grimoire's repo.
	repo := grimoire.New(adapter)

	var product Product

	// Changeset is used when creating or updating your data.
	ch := ProductChangeset(product, map[string]interface{}{
		"name": "shampoo",
		"price": 1000
	})

	if ch.Error() != nil {
		// do something
	}

	// Create products with changeset and return the result to &product,
	repo.From("products").MustCreate(&product, ch)

	// Find a product with id 1.
	repo.From("products").Find(1).MustOne(&product)

	// Update products with id=1.
	repo.From("products").Find(1).MustUpdate(&product, ch)

	// Delete Product with id=1.
	repo.From("products").Find(1).MustDelete()
}
```

### Connecting to database

In order to connect to database, first you need to initialize adapter and then create a grimoire's repo using the adapter instance.

```golang
import (
	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/mysql" // use mysql adapter
)

func main() {
	// open mysql connection.
	adapter, err := mysql.Open("root@(127.0.0.1:3306)/grimoire_test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic(err)
	}
	defer adapter.Close()

	// initialize grimoire's repo.
	repo := grimoire.New(adapter)
}
```

## CRUD Interface

### Create

There's three alternatives on how you can insert records to a database depending on your needs. The easiest way is by using struct directly like in most other Golang ORM.

```golang
user := User{Name: "Alice", Age: 18}
// Create a new record in `users`.
err := repo.From("users").Save(&user)

// Create multiple records at once,
err := repo.From("users").Save(&[]User{user, user})
```

The other way is by using changeset, most of the time you might want to use this way especially when handling data from user.
The advantages of using changeset is you can validates and pre-process your data before presisting to database. Using changeset also solves the problem of dealing with `zero values`, `null` and `undefined` fields where it's usually tricky to handle in `patch` request.

```golang
user := User{} // this also hold the schema

// prepare the changes
ch := changeset.Cast(user, map[string]interface{}{
	"name": "Alice",
	"age": 18,
	"address": "world",
}, []string{"name", "age"}) // this will filter `address`.
changeset.ValidateRequired(ch, []string{"name", "age"}) // validate `name` and `age` field exists.

// Insert changes to `users` table and return the result to `&user`.
err := repo.From("users").Insert(&user, ch)

// It's also possible to insert multiple changes at once.
users := []User{}
err := repo.From("users").Insert(&users, ch, ch)

// If you don't care about the return value, you can pass nil.
err := repo.From("users").Insert(nil, ch)
```

It's also possible to insert using query builder directly. Inserting without using changeset or `Save` method won't set `created_at` and `updated_at` fields.

```golang
// Insert a record to users.
err := repo.From("users").Set("name", "Alice").Set("age", 18).Insert(&user)

// If you don't care about the return value, you can pass nil.
err := repo.From("users").Set("name", "Alice").Set("age", 18).Insert(nil)

// When used alongside Changeset or using `Save` function, it'll replace value defined in changeset or struct.
// This behaviour especially useful when dealing with relation.
err := repo.From("users").Set("crew_id", 10).Insert(&user, ch, ch, ch)
err := repo.From("users").Set("crew_id", 10).Save(&users)
```

### Query

In general Grimoire's use query builder to perform select, insert, update and delete query.

#### Where

Grimoire's support various condition in query builder. to see a full list of supported condition check [this](https://godoc.org/github.com/Fs02/grimoire/c) out.

```golang
// declare some variables.
users     := "users"
id        := c.I("id") // identify column in a table.
name      := c.I("name")
age       := c.I("age")

user := User{}

// Get a record from users table using id.
// One will automatically add limit=1 to query.
err := repo.From(users).Find(1).One(&user)

// Which equal to:
err := repo.From(users).Where(c.Eq(id, user.ID)).One(&user)

alluser := []User{}

// Retrieves all record where name=Alice
err := repo.From(users).Where(c.Eq(name, "Alice")).All(&alluser)

// You can add more condition just by adding more arguments.
// This will retrieve all record where name=Alice and age=10.
// All condition inside where will automatically joined using and expression.
err := repo.From(users).Where(c.Eq(name, "Alice"), c.Ne(age, 10)).All(&alluser)

// Use or:
// This will retrieve all record where age=10 or age=15.
err := repo.From(users).Where(c.Eq(age, 10)).OrWhere(c.Eq(age, 15)).All(&alluser)

// Multiple condition inside `OrWhere` will be joined using and expression.
// This will retrieve all record where age=10 or (name=Alice and age=15).
err := repo.From(users).Where(c.Eq(age, 10)).OrWhere(c.Eq(name, "Alice"), c.Eq(age, 15)).All(&alluser)
```

#### Selecting Fields

```golang
// Get one record and only select only it's name and age.
er := repo.From(users).Select("name", "age").Find(1).One(&user)
```

#### Offset and Limit

```golang
// Limit query by 5.
err := repo.From(users).Limit(5).All(&alluser)

// Limit by 5 and offset by 5.
err := repo.From(users).Limit(5).Offset(5).All(&alluser)
```

#### Order

```golang
// Order by name ascending.
err := repo.From(users).Order(c.Asc(name)).All(&alluser)

// Order by name descending.
err := repo.From(users).Order(c.Desc(name)).All(&alluser)
```

#### Group By

```golang
// Group by gender and retrieves it's count.
err := repo.From(users).Group("gender").Select("COUNT(id)").All(&alluser)

// Group by with having.
err := repo.From(users).Group("age").Having(c.Gt(age, 10)).Select("COUNT(id)").All(&alluser)
```

#### Join

```golang
// Join addresses table with users on addresses.user_id=users.id
err := repo.From(addresses).Join(users).All(&alluser)

// Which equal to this if the join condition is manually specified.
err := repo.From(addresses).Join(users, c.Eq(c.I("addresses.user_id"), c.I("users.id"))).All(&alluser)

// It's also possible to join using other join mode using `JoinWith`
err := repo.From(addresses).JoinWith("LEFT JOIN", users).All(&alluser)
err := repo.From(addresses).JoinWith("LEFT OUTER JOIN", users).All(&alluser)
```

### Update

There's also three alternatives on how you can update records to a database. The easiest way is by using struct directly.


```golang
user := User{Name: "Alice", Age: 18}

// Update a record from `users` where id=1.
// Notice updating a record using `Save` function is similar to creating, but you will need to specify condition.
err := repo.From("users").Find(1).Save(&user)

// It's also possible to update multiple record (where age=18) at once and retrieves all the results.
// The following will update all record matches the condition and return it to array.
// Only the first item from slice will be used as update value.
users := []User{user}
err := repo.From("users").Where(Eq(I("age"), 18)).Save(&users)
```

Updating using changeset is similar to inserting.

```golang
user := User{} // this also hold the schema

// prepare the changes
ch := changeset.Cast(user, map[string]interface{}{
	"name": "Alice",
	"age": 18,
	"address": "world",
}, []string{"name", "age"}) // this will filter `address`.
changeset.ValidateRequired(ch, []string{"name", "age"}) // validate `name` and `age` field exists.

// Update changes to `users` where id=1 table and return the result to `&user`.
err := repo.From("users").Find(1).Update(&user, ch)

// If you don't care about the return value, you can pass nil.
err := repo.From("users").Find(1).Update(nil, ch)

// The following will update all record matches the condition and return it to array.
users := []User{user}
err := repo.From("users").Where(Eq(I("age"), 18)).Update(&user, ch)

// If no condition is used, grimoire will update all records.
err := repo.From("users").Update(nil, ch)
```

It's also possible to update using query builder directly. Update a record without using changeset or `Save` method won't set `updated_at` fields.

```golang
// Update a record where id=1.
err := repo.From("users").Find(1).Set("name", "Alice").Set("age", 18).Update(&user)

// If you don't care about the return value, you can pass nil.
err := repo.From("users").Find(1).Set("name", "Alice").Set("age", 18).Update(nil)

// When used alongside Changeset or using `Save` function, it'll replace value defined in changeset or struct.
// This behaviour especially useful when dealing with relation.
err := repo.From("users").Find(1).Set("crew_id", 10).Update(&user, ch, ch, ch)
err := repo.From("users").Find(1).Set("crew_id", 10).Save(&users)
```

### Delete

Deleting one or more records is simple.

```golang
// Delete a record with id=1.
err := repo.From("products").Find(1).Delete()

// Delete records where age=18
err := repo.From("users").Where(Eq(I("age"), 18)).Delete()

// Delete all records.
err := repo.From("users").Delete()
```


## Transaction

Transactions in grimoire are run inside a transaction function that returns an error.
Commit and Rollback are handled automatically by grimoire, transaction will rollback when the function returns an error or throws a panic.
if panic is not caused by grimoire's error or it's an grimoire's `UnexpectedError()`, the function will repanic after recover.
If no error returned or function did not panic, then the transaction will be commited.

```golang
user := User{}

// cast user changes alongside addreses
ch := changeUser(user, params)
if ch.Error() {
	// do something
}

// declare and execute transaction
err := repo.Transaction(func repo grimoire.Repo) error {
	// MustInsert similar to Insert, but this will panic if any error occured.
	// If it's panic, transaction will automatically rolled back
	// and the panic cause will be returned as an error as long as it's grimoire's error.
	repo.From("users").MustInsert(&user, ch)

	// Get array of addresses changeset.
	addresses := ch.Changes()["addresses"].([]*changeset.Changeset)

	// Insert multiple addresses changeset at once.
	// Set("user_id", user.ID) will ensure it's user_id refer to previous inserted user.
	repo.From("addresses").Set("user_id", user.ID).MustInsert(&user.Addresses, addresses...)

	// commit transaction
	return nil
})

if err != nil {
	// do something
}
```

## Logger

Grimoire's default logger can be replaced with repo's `SetLogger()` function.

```golang
// replace default logger with custom logger function.
repo.SetLogger(func(query string, duration time.Duration, err error) {
	log.Print("[", duration, "] - ", query)
})
```

## Field Mapping

By default Grimoire's will map struct fields by converting field's name to snake case.
But field name can be defined manually using struct tags.

```golang
type Entity struct {
	ALLCAP		bool 	`db:"all_cap"` 	// map `ALLCAP` as `all_cap`
	Virtual		bool 	`db:"-"` 	// using `-` will treat field as virtual.
	FullName	string			// by default it'll map to `full_name`
}
```
