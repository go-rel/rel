package grimoire

import (
	"testing"

	"github.com/Fs02/grimoire/change"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/query"
	"github.com/Fs02/grimoire/schema"
	"github.com/Fs02/grimoire/where"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var repo = Repo{}

func TestNew(t *testing.T) {
	adapter := new(TestAdapter)
	repo := New(adapter)

	assert.NotNil(t, repo)
	assert.Equal(t, adapter, repo.Adapter())
}

func TestRepo_SetLogger(t *testing.T) {
	repo := Repo{}
	assert.Nil(t, repo.logger)
	repo.SetLogger(DefaultLogger)
	assert.NotNil(t, repo.logger)
}

func TestRepo_Aggregate(t *testing.T) {
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users")
	mode := "COUNT"
	field := "*"

	var out struct {
		Count int
	}

	adapter.On("Aggregate", query, &out, mode, field).Return(nil)

	err := repo.Aggregate(User{}, "COUNT", "*", &out, query)
	assert.Nil(t, err)

	assert.NotPanics(t, func() {
		repo.MustAggregate(User{}, "COUNT", "*", &out, query)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_Count(t *testing.T) {
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users")

	adapter.On("Aggregate", query, mock.Anything, "COUNT", "*").Return(nil)

	count, err := repo.Count(User{}, query)
	assert.Nil(t, err)
	assert.Equal(t, 0, count)

	assert.NotPanics(t, func() {
		assert.Equal(t, 0, repo.MustCount(User{}, query))
	})

	adapter.AssertExpectations(t)
}

func TestRepo_One(t *testing.T) {
	user := User{}
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users").Limit(1)

	adapter.On("All", query, &user).Return(1, nil)

	assert.Nil(t, repo.One(&user, query))
	assert.NotPanics(t, func() { repo.MustOne(&user, query) })
	adapter.AssertExpectations(t)
}

func TestRepo_One_unexpectedError(t *testing.T) {
	user := User{}
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users").Limit(1)

	adapter.On("All", query, &user).Return(1, errors.NewUnexpected("error"))

	assert.NotNil(t, repo.One(&user, query))
	assert.Panics(t, func() { repo.MustOne(&user, query) })
	adapter.AssertExpectations(t)
}

func TestRepo_One_notFound(t *testing.T) {
	user := User{}
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users").Limit(1)

	adapter.On("All", query, &user).Return(0, nil)

	assert.Equal(t, errors.New("no result found", "", errors.NotFound), repo.One(&user, query))
	assert.Panics(t, func() { repo.MustOne(&user, query) })
	adapter.AssertExpectations(t)
}

func TestRepo_All(t *testing.T) {
	user := User{}
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}
	query := query.From("users").Limit(1)

	adapter.On("All", query, &user).Return(1, nil)

	assert.Nil(t, repo.All(&user, query))
	assert.NotPanics(t, func() { repo.MustAll(&user, query) })
	adapter.AssertExpectations(t)
}

func TestRepo_Insert(t *testing.T) {
	var (
		user      User
		adapter   = &TestAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
	)

	adapter.
		On("Insert", query.From("users"), changes).Return(1, nil).
		On("All", query.From("users").Where(where.Eq("id", 1)).Limit(1), &user).Return(1, nil)

	assert.Nil(t, repo.Insert(&user, cbuilders...))
	assert.NotPanics(t, func() { repo.MustInsert(&user, cbuilders...) })
	adapter.AssertExpectations(t)
}

func TestRepo_Insert_error(t *testing.T) {
	var (
		user      User
		adapter   = &TestAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
	)

	adapter.
		On("Insert", query.From("users"), changes).Return(0, errors.NewUnexpected("error"))

	assert.NotNil(t, repo.Insert(&user, cbuilders...))
	assert.Panics(t, func() { repo.MustInsert(&user, cbuilders...) })

	adapter.AssertExpectations(t)
}

func TestRepo_InsertAll(t *testing.T) {
	var (
		users   []User
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		changes = []change.Changes{
			change.Build(change.Set("name", "name1")),
			change.Build(change.Set("name", "name2")),
		}
	)

	adapter.
		On("InsertAll", query.From("users"), changes).Return([]interface{}{1, 2}, nil).
		On("All", query.From("users").Where(where.In("id", 1, 2)), &users).Return(2, nil)

	assert.Nil(t, repo.InsertAll(&users, changes))
	assert.NotPanics(t, func() { repo.MustInsertAll(&users, changes) })
	adapter.AssertExpectations(t)
}

func TestRepo_Update(t *testing.T) {
	var (
		user      = User{ID: 1}
		adapter   = &TestAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
		queries = query.From("users").Where(where.Eq("id", user.ID))
	)

	adapter.
		On("Update", queries, changes).Return(nil).
		On("All", queries.Limit(1), &user).Return(1, nil)

	assert.Nil(t, repo.Update(&user, cbuilders...))
	assert.NotPanics(t, func() { repo.MustUpdate(&user, cbuilders...) })
	adapter.AssertExpectations(t)
}

func TestRepo_Update_nothing(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
	)

	assert.Nil(t, repo.Update(nil))
	assert.NotPanics(t, func() { repo.MustUpdate(nil) })

	adapter.AssertExpectations(t)
}

func TestRepo_Update_unchanged(t *testing.T) {
	var (
		user    = User{ID: 1}
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
	)

	assert.Nil(t, repo.Update(&user))
	assert.NotPanics(t, func() { repo.MustUpdate(&user) })

	adapter.AssertExpectations(t)
}

func TestRepo_Update_error(t *testing.T) {
	var (
		user      = User{ID: 1}
		adapter   = &TestAdapter{}
		repo      = Repo{adapter: adapter}
		cbuilders = []change.Builder{
			change.Set("name", "name"),
		}
		changes = change.Build(cbuilders...)
		queries = query.From("users").Where(where.Eq("id", user.ID))
	)

	adapter.
		On("Update", queries, changes).Return(errors.NewUnexpected("error"))

	assert.NotNil(t, repo.Update(&user, cbuilders...))
	assert.Panics(t, func() { repo.MustUpdate(&user, cbuilders...) })
	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_update(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &Transaction{Buyer: User{ID: 1}}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users", where.Eq("id", 1))
		buyer, _ = changes.GetAssoc("Buyer")
	)

	adapter.
		On("Update", q, buyer[0]).Return(nil).
		On("All", q.Limit(1), &record.Buyer).Return(1, nil)

	err := repo.upsertBelongsTo(assocs, &changes)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_updateError(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &Transaction{Buyer: User{ID: 1}}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users", where.Eq("id", 1))
		buyer, _ = changes.GetAssoc("Buyer")
	)

	adapter.
		On("Update", q, buyer[0]).Return(errors.NewUnexpected("update error"))

	err := repo.upsertBelongsTo(assocs, &changes)
	assert.Equal(t, errors.NewUnexpected("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_updateInconsistentPrimaryKey(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &Transaction{Buyer: User{ID: 1}}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Buyer": change.Map{
					"id":   2,
					"name": "buyer1",
					"age":  20,
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.upsertBelongsTo(assocs, &changes)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_insertNew(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &Transaction{}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users")
		buyer, _ = changes.GetAssoc("Buyer")
	)

	adapter.
		On("Insert", q, buyer[0]).Return(1, nil).
		On("All", q.Where(where.Eq("id", 1)).Limit(1), &record.Buyer).Return(1, nil).
		Run(func(args mock.Arguments) {
			user := args.Get(1).(*User)
			user.ID = 1
		})

	err := repo.upsertBelongsTo(assocs, &changes)
	assert.Nil(t, err)

	ref, ok := changes.Get("user_id")
	assert.True(t, ok)
	assert.Equal(t, change.Set("user_id", 1), ref)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_insertNewError(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &Transaction{}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Buyer": change.Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		q        = query.Build("users")
		buyer, _ = changes.GetAssoc("Buyer")
	)

	adapter.
		On("Insert", q, buyer[0]).Return(0, errors.NewUnexpected("insert error"))

	err := repo.upsertBelongsTo(assocs, &changes)
	assert.Equal(t, errors.NewUnexpected("insert error"), err)

	_, ok := changes.Get("user_id")
	assert.False(t, ok)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertBelongsTo_notChanged(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &Transaction{}
		assocs  = schema.InferAssociations(record)
		changes = change.Build()
	)

	err := repo.upsertBelongsTo(assocs, &changes)
	assert.Nil(t, err)
	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_update(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &User{ID: 1, Address: Address{ID: 2}}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q            = query.Build("addresses").Where(where.Eq("id", 2).AndEq("user_id", 1))
		addresses, _ = changes.GetAssoc("Address")
	)

	adapter.
		On("Update", q, addresses[0]).Return(nil).
		On("All", q.Limit(1), &record.Address).Return(1, nil)

	err := repo.upsertHasOne(assocs, &changes, nil)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_updateError(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &User{ID: 1, Address: Address{ID: 2}}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q            = query.Build("addresses").Where(where.Eq("id", 2).AndEq("user_id", 1))
		addresses, _ = changes.GetAssoc("Address")
	)

	adapter.
		On("Update", q, addresses[0]).Return(errors.NewUnexpected("update error"))

	err := repo.upsertHasOne(assocs, &changes, nil)
	assert.Equal(t, errors.NewUnexpected("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_updateInconsistentPrimaryKey(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &User{ID: 1, Address: Address{ID: 2}}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"id":     1,
					"street": "street1",
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.upsertHasOne(assocs, &changes, nil)
	})

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_insertNew(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &User{}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q       = query.Build("addresses")
		address = change.Build(change.Set("street", "street1"))
	)

	// foreign value set after associations infered
	record.ID = 1
	address.SetValue("user_id", record.ID)

	adapter.
		On("Insert", q, address).Return(2, nil).
		On("All", q.Where(where.Eq("id", 2)).Limit(1), &record.Address).Return(1, nil)

	err := repo.upsertHasOne(assocs, &changes, nil)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepo_upsertHasOne_insertNewError(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		record  = &User{}
		assocs  = schema.InferAssociations(record)
		changes = change.Build(
			change.Map{
				"Address": change.Map{
					"street": "street1",
				},
			},
		)
		q       = query.Build("addresses")
		address = change.Build(change.Set("street", "street1"))
	)

	// foreign value set after associations infered
	record.ID = 1
	address.SetValue("user_id", record.ID)

	adapter.
		On("Insert", q, address).Return(nil, errors.NewUnexpected("insert error"))

	err := repo.upsertHasOne(assocs, &changes, nil)
	assert.Equal(t, errors.NewUnexpected("insert error"), err)

	adapter.AssertExpectations(t)
}

func TestRepo_Delete(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		user    = User{ID: 1}
	)

	adapter.
		On("Delete", query.From("users").Where(where.In("id", user.ID))).Return(nil)

	assert.Nil(t, repo.Delete(user))
	assert.NotPanics(t, func() { repo.MustDelete(user) })
	adapter.AssertExpectations(t)
}

func TestRepo_Delete_slice(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		users   = []User{
			{ID: 1},
			{ID: 2},
		}
	)

	adapter.
		On("Delete", query.From("users").Where(where.In("id", 1, 2))).Return(nil)

	assert.Nil(t, repo.Delete(users))
	assert.NotPanics(t, func() { repo.MustDelete(users) })
	adapter.AssertExpectations(t)
}

func TestRepo_Delete_emptySlice(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		users   = []User{}
	)

	assert.Nil(t, repo.Delete(users))
	assert.NotPanics(t, func() { repo.MustDelete(users) })
	adapter.AssertExpectations(t)
}

func TestRepo_DeleteAll(t *testing.T) {
	var (
		adapter = &TestAdapter{}
		repo    = Repo{adapter: adapter}
		queries = query.From("logs").Where(where.Eq("user_id", 1))
	)

	adapter.
		On("Delete", query.From("logs").Where(where.Eq("user_id", 1))).Return(nil)

	assert.Nil(t, repo.DeleteAll(queries))
	assert.NotPanics(t, func() { repo.MustDeleteAll(queries) })
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction(t *testing.T) {
	adapter := new(TestAdapter)
	adapter.On("Begin").Return(nil).
		On("Commit").Return(nil)

	repo := Repo{adapter: adapter}

	err := repo.Transaction(func(repo Repo) error {
		assert.True(t, repo.inTransaction)
		return nil
	})

	assert.False(t, repo.inTransaction)
	assert.Nil(t, err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_beginError(t *testing.T) {
	adapter := new(TestAdapter)
	adapter.On("Begin").Return(errors.NewUnexpected("error"))

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_commitError(t *testing.T) {
	adapter := new(TestAdapter)
	adapter.On("Begin").Return(nil).
		On("Commit").Return(errors.NewUnexpected("error"))

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_returnErrorAndRollback(t *testing.T) {
	adapter := new(TestAdapter)
	adapter.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		return errors.NewUnexpected("error")
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_panicWithKnownErrorAndRollback(t *testing.T) {
	adapter := new(TestAdapter)
	adapter.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	err := Repo{adapter: adapter}.Transaction(func(r Repo) error {
		// doing good things
		panic(errors.New("error", "", errors.NotFound))
	})

	assert.Equal(t, errors.New("error", "", errors.NotFound), err)
	adapter.AssertExpectations(t)
}

func TestRepo_Transaction_panicAndRollback(t *testing.T) {
	adapter := new(TestAdapter)
	adapter.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	assert.Panics(t, func() {
		Repo{adapter: adapter}.Transaction(func(r Repo) error {
			// doing good things
			panic(errors.NewUnexpected("error"))
		})
	})

	adapter.AssertExpectations(t)
}
