package rel

import (
	"time"
)

type Status string

type User struct {
	ID           int
	Name         string
	Age          int
	Transactions []Transaction `ref:"id" fk:"user_id"`
	Address      Address
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID      int
	Item    string
	Status  Status
	BuyerID int  `db:"user_id"`
	Buyer   User `ref:"user_id" fk:"id"`
}

type Address struct {
	ID        int
	UserID    *int
	User      *User
	Street    string
	DeletedAt *time.Time
}

type Owner struct {
	User   *User
	UserID *int
}
