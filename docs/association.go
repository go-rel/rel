package main

import (
	"context"
	"time"

	"github.com/Fs02/rel"
	"github.com/Fs02/rel/where"
)

/// [association-schema]

// User schema.
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
	// REL can guess the reference and foreign field if it's not specified.
	Address Address
}

// Transaction schema.
type Transaction struct {
	ID     int
	Item   string
	Status string

	// belongs to user.
	// contains primary key of other struct.
	Buyer   User `ref:"buyer_id" fk:"id"`
	BuyerID int
}

// Address schema.
type Address struct {
	ID   int
	City string

	// belongs to user.
	User   *User
	UserID *int
}

/// [association-schema]

// PreloadBelongsTo docs example.
func PreloadBelongsTo(ctx context.Context, repo rel.Repository) error {
	var transaction Transaction

	/// [preload-belongs-to]
	err := repo.Preload(ctx, &transaction, "buyer")
	/// [preload-belongs-to]

	return err
}

// PreloadHasOne docs example.
func PreloadHasOne(ctx context.Context, repo rel.Repository) error {
	var user User

	/// [preload-has-one]
	err := repo.Preload(ctx, &user, "address")
	/// [preload-has-one]

	return err
}

// PreloadHasMany docs example.
func PreloadHasMany(ctx context.Context, repo rel.Repository) error {
	var user User

	/// [preload-has-many]
	err := repo.Preload(ctx, &user, "transactions")
	/// [preload-has-many]

	return err
}

// PreloadHasManyFilter docs example.
func PreloadHasManyFilter(ctx context.Context, repo rel.Repository) error {
	var user User

	/// [preload-has-many-filter]
	err := repo.Preload(ctx, &user, "transactions", where.Eq("status", "paid"))
	/// [preload-has-many-filter]

	return err
}

// PreloadNested docs example.
func PreloadNested(ctx context.Context, repo rel.Repository) error {
	var transaction Transaction

	/// [preload-nested]
	err := repo.Preload(ctx, &transaction, "buyer.address")
	/// [preload-nested]

	return err
}

// InsertAssociation docs example.
func InsertAssociation(ctx context.Context, repo rel.Repository) error {
	/// [insert-association]
	user := User{
		Name: "rel",
		Address: Address{
			City: "Bandung",
		},
	}

	// Inserts a new record to users and address table.
	// Result: User{ID: 1, Name: "rel", Address: Address{ID: 1, City: "Bandung", UserID: 1}}
	err := repo.Insert(ctx, &user)
	/// [insert-association]

	return err
}

// UpdateAssociation docs example.
func UpdateAssociation(ctx context.Context, repo rel.Repository) error {
	/// [update-association]
	userID := 1
	user := User{
		ID:   1,
		Name: "rel",
		// association is loaded when the primary key (id) is not zero.
		Address: Address{
			ID:     1,
			UserID: &userID,
			City:   "Bandung",
		},
	}

	// Update user record with id 1.
	// Update address record with id 1.
	err := repo.Update(ctx, &user)
	/// [update-association]

	return err
}

// UpdateAssociationWithMap docs example.
func UpdateAssociationWithMap(ctx context.Context, repo rel.Repository) error {
	var user User

	/// [update-association-with-map]
	mutation := rel.Map{
		"address": rel.Map{
			"city": "bandung",
		},
	}

	// Update address record if it's loaded else it'll creates a new address.
	// only set city to bandung.
	err := repo.Update(ctx, &user, mutation)
	/// [update-association-with-map]

	return err
}
