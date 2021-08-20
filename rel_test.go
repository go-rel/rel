package rel

import (
	"encoding/json"
	"time"
)

type Status string

type extendedUser struct {
	ID        int
	Password  []byte
	Metadata  json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}

type User struct {
	ID           int
	Name         string
	Age          int
	Transactions []Transaction `ref:"id" fk:"user_id"`
	Address      Address       `autosave:"true"`
	WorkAddress  *Address
	UserRoles    []UserRole `autosave:"true"`
	Emails       []Email    `autosave:"true"`

	// many to many
	// user:id <- user_id:user_roles:role_id -> role:id
	Roles []Role `through:"user_roles"`

	// self-referencing needs two intermediate reference to be set up.
	Follows   []Follow `ref:"id" fk:"following_id"`
	Followeds []Follow `ref:"id" fk:"follower_id"`

	// association through
	Followings []User `through:"follows"`
	Followers  []User `through:"followeds"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type Follow struct {
	FollowerID  int  `db:",primary"`
	FollowingID int  `db:",primary"`
	Accepted    bool // this way, it may contains additional data
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
	Buyer     User `ref:"user_id" fk:"id" autoload:"true"`
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
	Users []User `through:"user_roles"`
}

type UserRole struct {
	UserID int `db:",primary"`
	RoleID int `db:",primary"`
}
