package grimoire

import (
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
	Transactions []Transaction `assoc:"has_many" ref:"ID" fk:"UserID"` // Map User.ID to Transaction.Address
	Address      Address       `assoc:"has_one" ref:"ID" fk:"UserID"`  // Map User.ID to Address.UserID
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Transaction struct {
	ID     int
	UserID int
	User   User `assoc:"belongs_to" ref:"UserID" fk:"ID"` // Map Transaction.UserID to User.ID
}

type Address struct {
	ID     int
	UserID int
	User   *User `assoc:"belongs_to" ref:"UserID" fk:"ID"` // Map Address.UserID to User.ID
}

func TestQuerySelect(t *testing.T) {
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

func TestQueryDistinct(t *testing.T) {
	assert.Equal(t, repo.From("users").Distinct(), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
		AsDistinct: true,
	})
}

func TestQueryJoin(t *testing.T) {
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

func TestQueryWhere(t *testing.T) {
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

func TestQueryOrWhere(t *testing.T) {
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

func TestQueryGroup(t *testing.T) {
	assert.Equal(t, repo.From("users").Group("active", "plan"), Query{
		repo:        &repo,
		Collection:  "users",
		Fields:      []string{"users.*"},
		GroupFields: []string{"active", "plan"},
	})
}

func TestQueryHaving(t *testing.T) {
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

func TestQueryOrHaving(t *testing.T) {
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

func TestQueryOrderBy(t *testing.T) {
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

func TestQueryOffset(t *testing.T) {
	assert.Equal(t, repo.From("users").Offset(10), Query{
		repo:         &repo,
		Collection:   "users",
		Fields:       []string{"users.*"},
		OffsetResult: 10,
	})
}

func TestQueryLimit(t *testing.T) {
	assert.Equal(t, repo.From("users").Limit(10), Query{
		repo:        &repo,
		Collection:  "users",
		Fields:      []string{"users.*"},
		LimitResult: 10,
	})
}

func TestQueryFind(t *testing.T) {
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

func TestQuerySet(t *testing.T) {
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

func TestQueryOne(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.One(&user))
	assert.NotPanics(t, func() { query.MustOne(&user) })
	mock.AssertExpectations(t)
}

func TestQueryOneUnexpectedError(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(1, errors.UnexpectedError("error"))

	assert.NotNil(t, query.One(&user))
	assert.Panics(t, func() { query.MustOne(&user) })
	mock.AssertExpectations(t)
}

func TestQueryOneNotFound(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(0, nil)

	assert.Equal(t, errors.NotFoundError("no result found"), query.One(&user))
	assert.Panics(t, func() { query.MustOne(&user) })
	mock.AssertExpectations(t)
}

func TestQueryAll(t *testing.T) {
	user := User{}
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Limit(1)

	mock.On("All", query, &user).Return(1, nil)

	assert.Nil(t, query.All(&user))
	assert.NotPanics(t, func() { query.MustAll(&user) })
	mock.AssertExpectations(t)
}

func TestQueryCount(t *testing.T) {
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

func TestQueryInsert(t *testing.T) {
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

func TestQueryInsertMultiple(t *testing.T) {
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

func TestQueryInsertMultipleWithSet(t *testing.T) {
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

func TestQueryInsertMultipleError(t *testing.T) {
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

	mock.On("InsertAll", query, allchanges).Return([]interface{}{1, 2}, errors.UnexpectedError("error"))

	assert.NotNil(t, query.Insert(&users, ch1, ch2))
	assert.Panics(t, func() { query.MustInsert(&users, ch1, ch2) })
	mock.AssertExpectations(t)
}

func TestQueryInsertNotReturning(t *testing.T) {
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

func TestQueryInsertWithSet(t *testing.T) {
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

func TestQueryInsertOnlySet(t *testing.T) {
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

func TestQueryInsertOnlySetError(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users").Set("age", 10)

	changes := map[string]interface{}{
		"age": 10,
	}

	mock.On("Insert", query, changes).Return(0, errors.UnexpectedError("error"))

	assert.NotNil(t, query.Insert(nil))
	assert.Panics(t, func() { query.MustInsert(nil) })
	mock.AssertExpectations(t)
}

func TestQueryInsertAssocOne(t *testing.T) {
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

func TestQueryInsertAssocMany(t *testing.T) {
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
func TestQueryInsertError(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Insert", query, changes).Return(0, errors.UnexpectedError("error"))

	assert.NotNil(t, query.Insert(&user, ch))
	assert.Panics(t, func() { query.MustInsert(&user, ch) })
	mock.AssertExpectations(t)
}

func TestQueryUpdate(t *testing.T) {
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

func TestUpdateWithSet(t *testing.T) {
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

func TestUpdateOnlySet(t *testing.T) {
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

func TestUpdateNotReturning(t *testing.T) {
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

func TestUpdateNothing(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	assert.Nil(t, query.Update(nil))
	assert.NotPanics(t, func() { query.MustUpdate(nil) })
	mock.AssertExpectations(t)
}

func TestQueryUpdateAssocOne(t *testing.T) {
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

func TestUpdateAssocMany(t *testing.T) {
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

func TestQueryUpdateError(t *testing.T) {
	ch, user := createChangeset()
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	changes := map[string]interface{}{
		"name":       "name",
		"updated_at": time.Now().Round(time.Second),
	}

	mock.On("Update", query, changes).Return(errors.UnexpectedError("error"))

	assert.NotNil(t, query.Update(&user, ch))
	assert.Panics(t, func() { query.MustUpdate(&user, ch) })
	mock.AssertExpectations(t)
}

func TestPutInsert(t *testing.T) {
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

func TestPutInsertMultiple(t *testing.T) {
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

func TestPutUpdate(t *testing.T) {
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

func TestPutUpdateMultiple(t *testing.T) {
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

func TestPutSliceEmpty(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")
	users := []User{}

	assert.Nil(t, query.Save(&users))
	assert.NotPanics(t, func() { query.MustSave(&users) })
}

func TestQueryDelete(t *testing.T) {
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

func TestPreloadHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	user := User{ID: 10}
	result := []Address{
		{ID: 100, UserID: 10},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Address{}).Return(1, nil)

	assert.Nil(t, query.Preload(&user, "Address"))
	assert.Equal(t, result[0], user.Address)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadSliceHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	users := []User{{ID: 10}, {ID: 20}}
	result := []Address{
		{ID: 100, UserID: 10},
		{ID: 200, UserID: 20},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Address{}).Return(2, nil)

	assert.Nil(t, query.Preload(&users, "Address"))
	assert.Equal(t, result[0], users[0].Address)
	assert.Equal(t, result[1], users[1].Address)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadNestedHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transaction := Transaction{
		User: User{ID: 10},
	}

	result := []Address{
		{ID: 100, UserID: 10},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Address{}).Return(1, nil)

	assert.Nil(t, query.Preload(&transaction, "User.Address"))
	assert.Equal(t, result[0], transaction.User.Address)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadSliceNestedHasOne(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transactions := []Transaction{
		{User: User{ID: 10}},
		{User: User{ID: 20}},
	}

	result := []Address{
		{ID: 100, UserID: 10},
		{ID: 200, UserID: 20},
	}

	query := repo.From("addresses")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Address{}).Return(2, nil)

	assert.Nil(t, query.Preload(&transactions, "User.Address"))
	assert.Equal(t, result[0], transactions[0].User.Address)
	assert.Equal(t, result[1], transactions[1].User.Address)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	user := User{ID: 10}
	result := []Transaction{
		{ID: 5, UserID: 10},
		{ID: 10, UserID: 10},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Transaction{}).Return(2, nil)

	assert.Nil(t, query.Preload(&user, "Transactions"))
	assert.Equal(t, result, user.Transactions)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadSliceHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	users := []User{{ID: 10}, {ID: 20}}
	result := []Transaction{
		{ID: 5, UserID: 10},
		{ID: 10, UserID: 10},
		{ID: 15, UserID: 20},
		{ID: 20, UserID: 20},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Transaction{}).Return(4, nil)

	assert.Nil(t, query.Preload(&users, "Transactions"))
	assert.Equal(t, result[:2], users[0].Transactions)
	assert.Equal(t, result[2:], users[1].Transactions)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadNestedHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	address := Address{User: &User{ID: 10}}
	result := []Transaction{
		{ID: 5, UserID: 10},
		{ID: 10, UserID: 10},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10)), &[]Transaction{}).Return(2, nil)

	assert.Nil(t, query.Preload(&address, "User.Transactions"))
	assert.Equal(t, result, address.User.Transactions)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadNestedSliceHasMany(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	addresses := []Address{
		{User: &User{ID: 10}},
		{User: &User{ID: 20}},
	}

	result := []Transaction{
		{ID: 5, UserID: 10},
		{ID: 10, UserID: 10},
		{ID: 15, UserID: 20},
		{ID: 20, UserID: 20},
	}

	query := repo.From("transactions")

	mock.Result(result).On("All", query.Where(In("user_id", 10, 20)), &[]Transaction{}).Return(4, nil)

	assert.Nil(t, query.Preload(&addresses, "User.Transactions"))
	assert.Equal(t, result[:2], addresses[0].User.Transactions)
	assert.Equal(t, result[2:], addresses[1].User.Transactions)
	// assert.NotPanics(t, func() { query.MustSave(&user) })
	mock.AssertExpectations(t)
}

func TestPreloadBelongsTo(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transaction := Transaction{UserID: 10}
	address := Address{UserID: 10}
	result := []User{{ID: 10}}

	query := repo.From("users")

	mock.Result(result).On("All", query.Where(In("id", 10)), &[]User{}).Return(1, nil)

	assert.Nil(t, query.Preload(&transaction, "User"))
	assert.Equal(t, result[0], transaction.User)
	// assert.NotPanics(t, func() { query.MustSave(&user) })

	assert.Nil(t, query.Preload(&address, "User"))
	assert.Equal(t, result[0], *address.User)
	// assert.NotPanics(t, func() { query.MustSave(&user) })

	mock.AssertExpectations(t)
}

func TestPreloadSliceBelongsTo(t *testing.T) {
	mock := new(TestAdapter)
	repo := Repo{adapter: mock}

	transactions := []Transaction{{UserID: 10}, {UserID: 20}}
	addresses := []Address{{UserID: 10}, {UserID: 20}}
	result := []User{{ID: 10}, {ID: 20}}

	query := repo.From("users")

	mock.Result(result).On("All", query.Where(In("id", 10, 20)), &[]User{}).Return(2, nil)

	assert.Nil(t, query.Preload(&transactions, "User"))
	assert.Equal(t, result[0], transactions[0].User)
	assert.Equal(t, result[1], transactions[1].User)
	// assert.NotPanics(t, func() { query.MustSave(&user) })

	assert.Nil(t, query.Preload(&addresses, "User"))
	assert.Equal(t, result[0], *addresses[0].User)
	assert.Equal(t, result[1], *addresses[1].User)

	mock.AssertExpectations(t)
}

func TestPreloadEmptySlice(t *testing.T) {
	repo := Repo{}
	addresses := []Address{}

	assert.Nil(t, repo.From("transactions").Preload(&addresses, "User.Transactions"))
}
