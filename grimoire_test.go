package grimoire

import (
	"time"
)

type Status string

type User struct {
	ID           int
	Name         string
	Age          int
	Transactions []Transaction `references:"ID" foreign_key:"BuyerID"`
	Address      Address
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID      int
	Item    string
	Status  Status
	BuyerID int  `db:"user_id"`
	Buyer   User `references:"BuyerID" foreign_key:"ID"`
}

type Address struct {
	ID     int
	UserID *int
	User   *User
	Street string
}

type Owner struct {
	User   *User
	UserID *int
}
