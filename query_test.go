package grimoire

import (
	"fmt"
	"testing"
	"time"

	. "github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
)

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
	BuyerID int  `db:"user_id"`
	Buyer   User `references:"BuyerID" foreign_key:"ID"`
}

type Address struct {
	ID     int
	UserID *int
	User   *User
}

type Owner struct {
	User   *User
	UserID *int
}

func TestQuery_Select(t *testing.T) {
	assert.Equal(t, repo.From("users").Select("*"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"*"},
	})

	assert.Equal(t, repo.From("users").Select("id", "name", "email"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"id", "name", "email"},
	})
}

func TestQuery_Distinct(t *testing.T) {
	assert.Equal(t, repo.From("users").Distinct(), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		AsDistinct: true,
	})
}

func TestQuery_Join(t *testing.T) {
	assert.Equal(t, repo.From("users").Join("transactions"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		JoinClause: []Join{
			{
				Mode:       "JOIN",
				Collection: "transactions",
				Condition: And(Eq(
					I("users.transaction_id"),
					I("transactions.id"),
				)),
			},
		},
	})

	assert.Equal(t, repo.From("users").Join("transactions", Eq(
		I("users.transaction_id"),
		I("transactions.id"),
	)), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		JoinClause: []Join{
			{
				Mode:       "JOIN",
				Collection: "transactions",
				Condition: And(Eq(
					I("users.transaction_id"),
					I("transactions.id"),
				)),
			},
		},
	})
}

func TestQuery_Where(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).Where(Ne("active", false)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at"), Ne("active", false)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_OrWhere(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").OrWhere(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			repo.From("users").Where(Eq("id", 1)).OrWhere(Nil("deleted_at")),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  Or(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  Or(And(Eq("id", 1), Nil("deleted_at")), Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false), Gte("score", 80)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false), Gte("score", 80)).Where(Lt("price", 10000)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"users.*"},
				Condition:  And(Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))), Lt("price", 10000)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_Group(t *testing.T) {
	assert.Equal(t, repo.From("users").Group("active", "plan"), Query{
		repo:        &repo,
		Collection:  "users",
		Fields:      []string{"users.*"},
		GroupFields: []string{"active", "plan"},
	})
}

func TestQuery_Having(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).Having(Ne("active", false)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at"), Ne("active", false)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_OrHaving(t *testing.T) {
	tests := []struct {
		Case     string
		Build    Query
		Expected Query
	}{
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").OrHaving(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			repo.From("users").Having(Eq("id", 1)).OrHaving(Nil("deleted_at")),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: Or(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: Or(And(Eq("id", 1), Nil("deleted_at")), Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false), Gte("score", 80)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false), Gte("score", 80)).Having(Lt("price", 10000)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"users.*"},
				HavingCondition: And(Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))), Lt("price", 10000)),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			assert.Equal(t, tt.Expected, tt.Build)
		})
	}
}

func TestQuery_OrderBy(t *testing.T) {
	assert.Equal(t, repo.From("users").Order(Asc("id")), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		OrderClause: []Order{
			{
				Field: "id",
				Order: 1,
			},
		},
	})
}

func TestQuery_Offset(t *testing.T) {
	assert.Equal(t, repo.From("users").Offset(10), Query{
		repo:         &repo,
		Collection:   "users",
		Fields:       []string{"users.*"},
		OffsetResult: 10,
	})
}

func TestQuery_Limit(t *testing.T) {
	assert.Equal(t, repo.From("users").Limit(10), Query{
		repo:        &repo,
		Collection:  "users",
		Fields:      []string{"users.*"},
		LimitResult: 10,
	})
}

func TestQuery_Find(t *testing.T) {
	assert.Equal(t, repo.From("users").Find(1), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		Condition:  And(Eq(I("users.id"), 1)),
	})

	assert.Equal(t, repo.From("users").Find("abc123"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		Condition:  And(Eq(I("users.id"), "abc123")),
	})
}

func TestQuery_FindBy(t *testing.T) {
	assert.Equal(t, repo.From("users").FindBy("email", "user@email.com"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		Condition:  And(Eq(I("users.email"), "user@email.com")),
	})
}

func TestQuery_Set(t *testing.T) {
	assert.Equal(t, repo.From("users").Set("field", 1), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		Changes: map[string]interface{}{
			"field": 1,
		},
	})

	assert.Equal(t, repo.From("users").Set("field1", 1).Set("field2", "2"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		Changes: map[string]interface{}{
			"field1": 1,
			"field2": "2",
		},
	})
}

func TestQuery_One(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.One(&user))
	assert.NotPanics(t, func() { query.MustOne(&user) })
	mock.AssertExpectations(t)
}

func TestQuery_One_unexpectedError(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(1, errors.NewUnexpected("error"))

	assert.NotNil(t, query.One(&user))
	assert.Panics(t, func() { query.MustOne(&user) })
	mock.AssertExpectations(t)
}

func TestQuery_One_notFound(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(0, nil)

	assert.Equal(t, errors.New("no result found", "", errors.NotFound), query.One(&user))
	assert.Panics(t, func() { query.MustOne(&user) })
	mock.AssertExpectations(t)
}

func TestQuery_All(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.All(&user))
	assert.NotPanics(t, func() { query.MustAll(&user) })
	mock.AssertExpectations(t)
}

func TestQuery_Count(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	mock.On("Count", query).Return(10, nil)

	count, err := query.Count()
	assert.Nil(t, err)
	assert.Equal(t, 10, count)

	assert.NotPanics(t, func() {
		assert.Equal(t, 10, query.MustCount())
	})

	mock.AssertExpectations(t)
}

func createChangeset() (*changeset.Changeset, User) {
	user := User{}
	ch := changeset.Cast(user, map[string]interface{}{
		"name": "name",
	}, []string{"name"})

	if ch.Error() != nil {
		panic(ch.Error())
	}

	return ch, user
}

func TestQuery_Insert(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Insert", query, changes).Return(1, nil).
		On("All", query.Find(1).Limit(1), &user).Return(1, nil)

	assert.Nil(t, query.Insert(&user, ch))
	assert.NotPanics(t, func() { query.MustInsert(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_multiple(t *testing.T) {
	ch1, user1 := createChangeset()
	ch2, user2 := createChangeset()
	users := []User{user1, user2}

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	allchanges := []map[string]interface{}{changes, changes}

	mock.On("InsertAll", query, allchanges).Return([]interface{}{1, 2}, nil).
		On("All", query.Where(In(I("id"), 1, 2)), &users).Return(2, nil)

	assert.Nil(t, query.Insert(&users, ch1, ch2))
	assert.NotPanics(t, func() { query.MustInsert(&users, ch1, ch2) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_multipleWithSet(t *testing.T) {
	ch1, user1 := createChangeset()
	ch2, user2 := createChangeset()
	users := []User{user1, user2}

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 18)

	changes := map[string]interface{}{
		"name":       "name",
		"age":        18,
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	allchanges := []map[string]interface{}{changes, changes}

	mock.On("InsertAll", query, allchanges).Return([]interface{}{1, 2}, nil).
		On("All", query.Where(In(I("id"), 1, 2)), &users).Return(2, nil)

	assert.Nil(t, query.Insert(&users, ch1, ch2))
	assert.NotPanics(t, func() { query.MustInsert(&users, ch1, ch2) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_multipleError(t *testing.T) {
	ch1, user1 := createChangeset()
	ch2, user2 := createChangeset()
	users := []User{user1, user2}

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	allchanges := []map[string]interface{}{changes, changes}

	mock.On("InsertAll", query, allchanges).Return([]interface{}{1, 2}, errors.NewUnexpected("error"))

	assert.NotNil(t, query.Insert(&users, ch1, ch2))
	assert.Panics(t, func() { query.MustInsert(&users, ch1, ch2) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_notReturning(t *testing.T) {
	ch, _ := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Insert", query, changes).Return(1, nil)

	assert.Nil(t, query.Insert(nil, ch))
	assert.NotPanics(t, func() { query.MustInsert(nil, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_withSet(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 10)

	changes := map[string]interface{}{
		"name":       "name",
		"age":        10,
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Insert", query, changes).Return(0, nil).
		On("All", query.Find(0).Limit(1), &user).Return(1, nil)

	assert.Nil(t, query.Insert(&user, ch))
	assert.NotPanics(t, func() { query.MustInsert(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_onlySet(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 10)

	changes := map[string]interface{}{
		"age": 10,
	}

	mock.On("Insert", query, changes).Return(0, nil)

	assert.Nil(t, query.Insert(nil))
	assert.NotPanics(t, func() { query.MustInsert(nil) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_onlySetError(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 10)

	changes := map[string]interface{}{
		"age": 10,
	}

	mock.On("Insert", query, changes).Return(0, errors.NewUnexpected("error"))

	assert.NotNil(t, query.Insert(nil))
	assert.Panics(t, func() { query.MustInsert(nil) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_assocOne(t *testing.T) {
	var card struct {
		ID   int
		User User
	}

	params := map[string]interface{}{
		"id": 1,
		"user": map[string]interface{}{
			"name": "name",
		},
	}

	userChangeset := func(entity interface{}, params map[string]interface{}) *changeset.Changeset {
		ch := changeset.Cast(entity, params, []string{"name"})
		return ch
	}

	ch := changeset.Cast(card, params, []string{"id"})
	changeset.CastAssoc(ch, "user", userChangeset)

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("cards")

	changes := map[string]interface{}{
		"id": 1,
	}

	mock.On("Insert", query, changes).Return(0, nil).
		On("All", query.Find(0).Limit(1), &card).Return(1, nil)

	assert.Nil(t, query.Insert(&card, ch))
	assert.NotPanics(t, func() { query.MustInsert(&card, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Insert_assocMany(t *testing.T) {
	var group struct {
		Name  string
		Users []User
	}

	params := map[string]interface{}{
		"name": "name",
		"users": []map[string]interface{}{
			{
				"name": "name1",
			},
			{
				"name": "name2",
			},
		},
	}

	userChangeset := func(entity interface{}, params map[string]interface{}) *changeset.Changeset {
		ch := changeset.Cast(entity, params, []string{"name"})
		return ch
	}

	ch := changeset.Cast(group, params, []string{"name"})
	changeset.CastAssoc(ch, "users", userChangeset)

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("groups")

	allchanges := []map[string]interface{}{
		{
			"name":       "name1",
			"created_at": time.Now().Round(time.Second),
			"updated_at": time.Now().Round(time.Second),
		},
		{
			"name":       "name2",
			"created_at": time.Now().Round(time.Second),
			"updated_at": time.Now().Round(time.Second),
		},
	}

	userChs := ch.Changes()["users"].([]*changeset.Changeset)

	mock.On("InsertAll", query, allchanges).Return([]interface{}{1, 2}, nil).
		On("All", query.Where(In("id", 1, 2)), &group.Users).Return(1, nil)

	assert.Nil(t, query.Insert(&group.Users, userChs...))
	assert.NotPanics(t, func() { query.MustInsert(&group.Users, userChs...) })
	mock.AssertExpectations(t)
}
func TestQuery_Insert_error(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Insert", query, changes).Return(0, errors.NewUnexpected("error"))

	assert.NotNil(t, query.Insert(&user, ch))
	assert.Panics(t, func() { query.MustInsert(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Update(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(nil).
		On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.Update(&user, ch))
	assert.NotPanics(t, func() { query.MustUpdate(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_withSet(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 10)

	changes := map[string]interface{}{
		"name":       "name",
		"age":        10,
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(nil).
		On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.Update(&user, ch))
	assert.NotPanics(t, func() { query.MustUpdate(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_onlySet(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 10)

	changes := map[string]interface{}{
		"age": 10,
	}

	mock.On("Update", query, changes).Return(nil)

	assert.Nil(t, query.Update(nil))
	assert.NotPanics(t, func() { query.MustUpdate(nil) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_notReturning(t *testing.T) {
	ch, _ := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(nil)

	assert.Nil(t, query.Update(nil, ch))
	assert.NotPanics(t, func() { query.MustUpdate(nil, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_nothing(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	assert.Nil(t, query.Update(nil))
	assert.NotPanics(t, func() { query.MustUpdate(nil) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_assocOne(t *testing.T) {
	var card struct {
		ID   int
		User User
	}

	params := map[string]interface{}{
		"id": 1,
		"user": map[string]interface{}{
			"name": "name",
		},
	}

	userChangeset := func(entity interface{}, params map[string]interface{}) *changeset.Changeset {
		ch := changeset.Cast(entity, params, []string{"name"})
		return ch
	}

	ch := changeset.Cast(card, params, []string{"id"})
	changeset.CastAssoc(ch, "user", userChangeset)

	changes := map[string]interface{}{
		"id": 1,
	}

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("cards")

	mock.On("Update", query, changes).Return(nil).
		On("All", query, &card).Return(1, nil)

	assert.Nil(t, query.Update(&card, ch))
	assert.NotPanics(t, func() { query.MustUpdate(&card, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_assocMany(t *testing.T) {
	var group struct {
		Name  string
		Users []User
	}

	params := map[string]interface{}{
		"name": "name",
		"user": []map[string]interface{}{
			{
				"name": "name",
			},
		},
	}

	userChangeset := func(entity interface{}, params map[string]interface{}) *changeset.Changeset {
		ch := changeset.Cast(entity, params, []string{"name"})
		return ch
	}

	ch := changeset.Cast(group, params, []string{"name"})
	changeset.CastAssoc(ch, "users", userChangeset)

	changes := map[string]interface{}{
		"name": "name",
	}

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("groups")

	mock.On("Update", query, changes).Return(nil).
		On("All", query, &group).Return(1, nil)

	assert.Nil(t, query.Update(&group, ch))
	assert.NotPanics(t, func() { query.MustUpdate(&group, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Update_error(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(errors.NewUnexpected("error"))

	assert.NotNil(t, query.Update(&user, ch))
	assert.Panics(t, func() { query.MustUpdate(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQuery_Put_insert(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")
	user := User{}

	changes := map[string]interface{}{
		"name":       "",
		"age":        0,
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Insert", query, changes).Return(1, nil).
		On("All", query.Find(1).Limit(1), &user).Return(1, nil)

	assert.Nil(t, query.Save(&user))
	assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestQuery_Put_insertMultiple(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")
	users := []User{{}, {}}

	changes := map[string]interface{}{
		"name":       "",
		"age":        0,
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("InsertAll", query, []map[string]interface{}{changes, changes}).Return([]interface{}{1, 2}, nil).
		On("All", query.Where(In(I("id"), 1, 2)), &users).Return(1, nil)

	assert.Nil(t, query.Save(&users))
	assert.NotPanics(t, func() { query.MustSave(&users) })
	mock.AssertExpectations(t)
}

func TestQuery_Put_update(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Find(1)
	user := User{}

	changes := map[string]interface{}{
		"name":       "",
		"age":        0,
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(nil).
		On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.Save(&user))
	assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestQuery_Put_updateMultiple(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Where(Eq("name", "name"))
	users := []User{{}, {}}

	changes := map[string]interface{}{
		"name":       "",
		"age":        0,
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(nil).
		On("All", query, &users).Return(1, nil)

	assert.Nil(t, query.Save(&users))
	assert.NotPanics(t, func() { query.MustSave(&users) })
	mock.AssertExpectations(t)
}

func TestQuery_Put_sliceEmpty(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")
	users := []User{}

	assert.Nil(t, query.Save(&users))
	assert.NotPanics(t, func() { query.MustSave(&users) })
}

func TestQuery_Delete(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	mock.On("Delete", query).Return(nil)

	assert.Nil(t, query.Delete())
	assert.NotPanics(t, func() { query.MustDelete() })
	mock.AssertExpectations(t)
}

func TestGetFields(t *testing.T) {
	var group struct {
		Name  string
		Users []User
	}

	query := Repo{}.From("users")
	params := map[string]interface{}{
		"name": "name",
		"users": []map[string]interface{}{
			{
				"name": "name1",
			},
		},
	}

	userChangeset := func(entity interface{}, params map[string]interface{}) *changeset.Changeset {
		ch := changeset.Cast(entity, params, []string{"name"})
		return ch
	}

	ch := changeset.Cast(group, params, []string{"name"})
	changeset.CastAssoc(ch, "users", userChangeset)

	assert.Equal(t, []string{"name"}, getFields(query, []*changeset.Changeset{ch}))
}

func TestQuery_Preload_hasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	user := User{ID: 10}
	result := []Address{
		{ID: 100, UserID: &user.ID},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Address{}).Return(1, nil)

	assert.Nil(t, query.Preload(&user, "Address"))
	assert.Equal(t, result[0], user.Address)
	assert.NotPanics(t, func() { query.MustPreload(&user, "Address") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_sliceHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	users := []User{{ID: 10}, {ID: 20}}
	result := []Address{
		{ID: 100, UserID: &users[0].ID},
		{ID: 200, UserID: &users[1].ID},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Address{}).Return(2, nil)

	assert.Nil(t, query.Preload(&users, "Address"))
	assert.Equal(t, result[0], users[0].Address)
	assert.Equal(t, result[1], users[1].Address)
	assert.NotPanics(t, func() { query.MustPreload(&users, "Address") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_nestedHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transaction := Transaction{
		Buyer: User{ID: 10},
	}

	result := []Address{
		{ID: 100, UserID: &transaction.Buyer.ID},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Address{}).Return(1, nil)

	assert.Nil(t, query.Preload(&transaction, "Buyer.Address"))
	assert.Equal(t, result[0], transaction.Buyer.Address)
	assert.NotPanics(t, func() { query.MustPreload(&transaction, "Buyer.Address") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_sliceNestedHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transactions := []Transaction{
		{Buyer: User{ID: 10}},
		{Buyer: User{ID: 20}},
	}

	result := []Address{
		{ID: 100, UserID: &transactions[0].Buyer.ID},
		{ID: 200, UserID: &transactions[1].Buyer.ID},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Address{}).Return(2, nil)

	assert.Nil(t, query.Preload(&transactions, "Buyer.Address"))
	assert.Equal(t, result[0], transactions[0].Buyer.Address)
	assert.Equal(t, result[1], transactions[1].Buyer.Address)
	assert.NotPanics(t, func() { query.MustPreload(&transactions, "Buyer.Address") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_hasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	user := User{ID: 10}
	result := []Transaction{
		{ID: 5, BuyerID: 10},
		{ID: 10, BuyerID: 10},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Transaction{}).Return(2, nil)

	assert.Nil(t, query.Preload(&user, "Transactions"))
	assert.Equal(t, result, user.Transactions)
	assert.NotPanics(t, func() { query.MustPreload(&user, "Transactions") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_sliceHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	users := []User{{ID: 10}, {ID: 20}}
	result := []Transaction{
		{ID: 5, BuyerID: 10},
		{ID: 10, BuyerID: 10},
		{ID: 15, BuyerID: 20},
		{ID: 20, BuyerID: 20},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Transaction{}).Return(4, nil)

	assert.Nil(t, query.Preload(&users, "Transactions"))
	assert.Equal(t, result[:2], users[0].Transactions)
	assert.Equal(t, result[2:], users[1].Transactions)
	assert.NotPanics(t, func() { query.MustPreload(&users, "Transactions") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_nestedHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	address := Address{User: &User{ID: 10}}
	result := []Transaction{
		{ID: 5, BuyerID: 10},
		{ID: 10, BuyerID: 10},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Transaction{}).Return(2, nil)

	assert.Nil(t, query.Preload(&address, "User.Transactions"))
	assert.Equal(t, result, address.User.Transactions)
	assert.NotPanics(t, func() { query.MustPreload(&address, "User.Transactions") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_nestedSliceHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	addresses := []Address{
		{User: &User{ID: 10}},
		{User: &User{ID: 20}},
	}

	result := []Transaction{
		{ID: 5, BuyerID: 10},
		{ID: 10, BuyerID: 10},
		{ID: 15, BuyerID: 20},
		{ID: 20, BuyerID: 20},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Transaction{}).Return(4, nil)

	assert.Nil(t, query.Preload(&addresses, "User.Transactions"))
	assert.Equal(t, result[:2], addresses[0].User.Transactions)
	assert.Equal(t, result[2:], addresses[1].User.Transactions)
	assert.NotPanics(t, func() { query.MustPreload(&addresses, "User.Transactions") })
	mock.AssertExpectations(t)
}

func TestQuery_Preload_belongsTo(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transaction := Transaction{BuyerID: 10}
	address := Address{UserID: &transaction.BuyerID}
	result := []User{{ID: 10}}

	query := repo.From("users")

	mock.Result(result).On("All", query.Where(In("id", 10)), &[]User{}).Return(1, nil)

	assert.Nil(t, query.Preload(&transaction, "Buyer"))
	assert.Equal(t, result[0], transaction.Buyer)
	assert.NotPanics(t, func() { query.MustPreload(&transaction, "Buyer") })

	assert.Nil(t, query.Preload(&address, "User"))
	assert.Equal(t, result[0], *address.User)
	assert.NotPanics(t, func() { query.MustPreload(&address, "User") })

	mock.AssertExpectations(t)
}

func TestQuery_Preload_ptr(t *testing.T) {
	repo := Repo{}
	query := repo.From("owners")

	owner := Owner{}

	assert.Nil(t, query.Preload(&owner, "User"))
	assert.Nil(t, owner.User)
	assert.Nil(t, owner.UserID)
}

func TestQuery_Preload_slicePtr(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}
	id := 1

	owners := []Owner{
		{}, // nil
		{
			UserID: &id,
		},
	}

	result := []User{{ID: 1}}

	query := repo.From("owners")

	mock.Result(result).On("All", query.Where(In("id", 1)), &[]User{}).Return(1, nil)

	assert.Nil(t, query.Preload(&owners, "User"))
	assert.Nil(t, owners[0].User)
	assert.Nil(t, owners[0].UserID)

	assert.Equal(t, result[0], *owners[1].User)
	assert.Equal(t, result[0].ID, *owners[1].UserID)

	mock.AssertExpectations(t)
}

func TestQuery_Preload_sliceBelongsTo(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transactions := []Transaction{{BuyerID: 10}, {BuyerID: 20}}
	addresses := []Address{
		{UserID: &transactions[0].BuyerID},
		{UserID: &transactions[1].BuyerID},
	}

	result := []User{{ID: 10}, {ID: 20}}

	query := repo.From("users")

	mock.Result(result).On("All", query.Where(In("id", 10, 20)), &[]User{}).Return(2, nil)

	assert.Nil(t, query.Preload(&transactions, "Buyer"))
	assert.Equal(t, result[0], transactions[0].Buyer)
	assert.Equal(t, result[1], transactions[1].Buyer)
	assert.NotPanics(t, func() { query.MustPreload(&transactions, "Buyer") })

	assert.Nil(t, query.Preload(&addresses, "User"))
	assert.Equal(t, result[0], *addresses[0].User)
	assert.Equal(t, result[1], *addresses[1].User)
	assert.NotPanics(t, func() { query.MustPreload(&addresses, "User") })

	mock.AssertExpectations(t)
}

func TestQuery_Preload_emptySlice(t *testing.T) {
	repo := Repo{}
	addresses := []Address{}

	assert.Nil(t, repo.From("transactions").Preload(&addresses, "User.Transactions"))
}

func TestQuery_Preload_notPointerPanic(t *testing.T) {
	repo := Repo{}
	transaction := Transaction{}

	assert.Panics(t, func() { repo.From("users").Preload(transaction, "User") })
}

func TestQuery_Preload_notValidPanic(t *testing.T) {
	repo := Repo{}
	transaction := Transaction{}

	assert.Panics(t, func() { repo.From("users").Preload(&transaction, "ID") })
	assert.Panics(t, func() { repo.From("users").Preload(&transaction, "ID.User") })
}

func TestQuery_Preload_invalidKeyPanic(t *testing.T) {
	repo := Repo{}
	info := struct {
		ID        int
		User      User `references:"UID" foreign_key:"InfoID"`
		OtherUser User `references:"ID" foreign_key:"InfoID"`
	}{}

	assert.Panics(t, func() { repo.From("users").Preload(&info, "User") })
	assert.Panics(t, func() { repo.From("users").Preload(&info, "OtherUser") })
}

func TestQuery_Preload_queryError(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	user := User{ID: 10}
	query := repo.From("addresses")

	mock.On("All", query.Where(In("user_id", 10)), &[]Address{}).Return(1, errors.NewUnexpected("query error"))

	assert.NotNil(t, query.Preload(&user, "Address"))
	assert.Panics(t, func() { query.MustPreload(&user, "Address") })
	mock.AssertExpectations(t)
}

func TestTransformError_unknown(t *testing.T) {
	assert.Equal(t, errors.NewUnexpected("error"), transformError(fmt.Errorf("error")))
}
