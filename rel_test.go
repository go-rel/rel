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
	UserRoles    []UserRole
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID        int
	Item      string
	Status    Status
	BuyerID   int  `db:"user_id"`
	Buyer     User `ref:"user_id" fk:"id"`
	AddressID int
	Address   Address
	Histories *[]History
}

type History struct {
	ID            int
	TransactionID int
}

type Notes string

func (n Notes) Equal(other interface{}) bool {
	if o, ok := other.(Notes); ok {
		return n == o
	}

	return false
}

type Address struct {
	ID        int
	UserID    *int
	User      *User
	Street    string
	Notes     Notes
	DeletedAt *time.Time
}

type Owner struct {
	User   *User
	UserID *int
}

type Role struct {
	ID        int
	Name      string
	UserRoles []UserRole
}

type UserRole struct {
	UserID int `db:",primary"`
	RoleID int `db:",primary"`
}
