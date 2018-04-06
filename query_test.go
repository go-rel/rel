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
	Name      string
	Age       int
	CreatedAt time.Time
	UpdatedAt time.Time
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
		Fields:     []string{"*"},
		AsDistinct: true,
	})
}

func TestQueryJoin(t *testing.T) {
	assert.Equal(t, repo.From("users").Join("transactions"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"*"},
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
		Fields:     []string{"*"},
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

func TestQueryJoinWith(t *testing.T) {
	t.Skip("PENDING")
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
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).Where(Ne("active", false)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"*"},
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
				Fields:     []string{"*"},
				Condition:  And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			repo.From("users").Where(Eq("id", 1)).OrWhere(Nil("deleted_at")),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  Or(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  Or(And(Eq("id", 1), Nil("deleted_at")), Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false), Gte("score", 80)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"*"},
				Condition:  Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			repo.From("users").Where(Eq("id", 1), Nil("deleted_at")).OrWhere(Ne("active", false), Gte("score", 80)).Where(Lt("price", 10000)),
			Query{
				repo:       &repo,
				Collection: "users",
				Fields:     []string{"*"},
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
		Fields:      []string{"*"},
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
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 AND deleted_at IS NIL AND active<>false`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).Having(Ne("active", false)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"*"},
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
				Fields:          []string{"*"},
				HavingCondition: And(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`id=1 OR deleted_at IS NIL`,
			repo.From("users").Having(Eq("id", 1)).OrHaving(Nil("deleted_at")),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: Or(Eq("id", 1), Nil("deleted_at")),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR active<>true`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: Or(And(Eq("id", 1), Nil("deleted_at")), Ne("active", false)),
			},
		},
		{
			`(id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false), Gte("score", 80)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"*"},
				HavingCondition: Or(And(Eq("id", 1), Nil("deleted_at")), And(Ne("active", false), Gte("score", 80))),
			},
		},
		{
			`((id=1 AND deleted_at IS NIL) OR (active<>true AND score>=80)) AND price<10000`,
			repo.From("users").Having(Eq("id", 1), Nil("deleted_at")).OrHaving(Ne("active", false), Gte("score", 80)).Having(Lt("price", 10000)),
			Query{
				repo:            &repo,
				Collection:      "users",
				Fields:          []string{"*"},
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
		Fields:     []string{"*"},
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
		Fields:       []string{"*"},
		OffsetResult: 10,
	})
}

func TestQueryLimit(t *testing.T) {
	assert.Equal(t, repo.From("users").Limit(10), Query{
		repo:        &repo,
		Collection:  "users",
		Fields:      []string{"*"},
		LimitResult: 10,
	})
}

func TestQueryFind(t *testing.T) {
	assert.Equal(t, repo.From("users").Find(1), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"*"},
		Condition:  And(Eq(I("users.id"), 1)),
	})

	assert.Equal(t, repo.From("users").Find("abc123"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"*"},
		Condition:  And(Eq(I("users.id"), "abc123")),
	})
}

func TestQuerySet(t *testing.T) {
	assert.Equal(t, repo.From("users").Set("field", 1), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"*"},
		Changes: map[string]interface{}{
			"field": 1,
		},
	})

	assert.Equal(t, repo.From("users").Set("field1", 1).Set("field2", "2"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"*"},
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

	mock.On("Insert", query, changes).Return(0, nil).
		On("All", query.Find(0).Limit(1), &user).Return(1, nil)

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

	mock.On("Insert", query, changes).Return(0, nil).
		On("All", query.Where(In(I("id"), 0, 0)), &users).Return(2, nil)

	assert.Nil(t, query.Insert(&users, ch1, ch2))
	assert.NotPanics(t, func() { query.MustInsert(&users, ch1, ch2) })
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

	mock.On("Insert", query, changes).Return(0, nil)

	assert.Nil(t, query.Insert(nil, ch))
	assert.NotPanics(t, func() { query.MustInsert(nil, ch) })
	mock.AssertExpectations(t)
}

func TestQueryInsert_withSet(t *testing.T) {
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

	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("groups")

	changes := map[string]interface{}{
		"name": "name",
	}

	mock.On("Insert", query, changes).Return(0, nil).
		On("All", query.Find(0).Limit(1), &group).Return(1, nil)

	assert.Nil(t, query.Insert(&group, ch))
	assert.NotPanics(t, func() { query.MustInsert(&group, ch) })
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

func TestQueryDelete(t *testing.T) {
	mock := new(TestAdapter)
	query := Repo{adapter: mock}.From("users")

	mock.On("Delete", query).Return(nil)

	assert.Nil(t, query.Delete())
	assert.NotPanics(t, func() { query.MustDelete() })
	mock.AssertExpectations(t)
}
