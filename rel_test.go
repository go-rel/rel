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
	Address      Address       `autosave:"true"`
	WorkAddress  Address
	UserRoles    []UserRole `autosave:"true"`
	Emails       []Email    `autosave:"true"`

	// many to many
	// user:id <- user_id:user_roles:role_id -> role:id
	Roles []Role `through:"user_roles"`

	// many to many: self-referencing with explicitly defined ref and fk.
	// omit mapped field.
	Followers  []User `ref:"id:following_id" fk:"id:follower_id" through:"followers"`
	Followings []User `ref:"id:follower_id" fk:"id:following_id" through:"followers"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Email struct {
	ID     int
	Email  string
	UserID int
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

type Profile struct {
	ID     int
	Name   string
	User   *User `autosave:"true"`
	UserID *int
}

type Role struct {
	ID        int
	Name      string
	UserRoles []UserRole

	// explicit many to many declaration:
	// role:id <- role_id:user_roles:user_id -> user:id
	Users []User `ref:"id:role_id" fk:"id:user_id" through:"user_roles"`
}

type UserRole struct {
	UserID int `db:",primary"`
	RoleID int `db:",primary"`
}
