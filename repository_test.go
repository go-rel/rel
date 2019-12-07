package rel

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	now = func() time.Time {
		return time.Time{}
	}
}

var repo = repository{}

func createCursor(row int) *testCursor {
	cur := &testCursor{}

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id"}, nil).Once()

	if row > 0 {
		cur.On("Next").Return(true).Times(row)
		cur.MockScan(10).Times(row)
	}

	cur.On("Next").Return(false).Once()

	return cur
}

func TestNew(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	assert.NotNil(t, repo)
	assert.Equal(t, adapter, repo.Adapter())
}

func TestRepository_SetLogger(t *testing.T) {
	var (
		repo = repository{}
	)

	assert.Nil(t, repo.logger)
	repo.SetLogger(DefaultLogger)
	assert.NotNil(t, repo.logger)
}

func TestRepository_Aggregate(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		query     = From("users")
		aggregate = "count"
		field     = "*"
	)

	adapter.On("Aggregate", query, aggregate, field).Return(1, nil).Once()

	count, err := repo.Aggregate(query, "count", "*")
	assert.Equal(t, 1, count)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepository_MustAggregate(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		query     = From("users")
		aggregate = "count"
		field     = "*"
	)

	adapter.On("Aggregate", query, aggregate, field).Return(1, nil).Once()

	assert.NotPanics(t, func() {
		count := repo.MustAggregate(query, "count", "*")
		assert.Equal(t, 1, count)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Count(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users")
	)

	adapter.On("Aggregate", query, "count", "*").Return(1, nil).Once()

	count, err := repo.Count("users")
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	adapter.AssertExpectations(t)
}

func TestRepository_MustCount(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users")
	)

	adapter.On("Aggregate", query, "count", "*").Return(1, nil).Once()

	assert.NotPanics(t, func() {
		count := repo.MustCount("users")
		assert.Equal(t, 1, count)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Find(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users").Limit(1)
		cur     = createCursor(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.Find(&user, query))
	assert.Equal(t, 10, user.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_queryError(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		cur     = &testCursor{}
		query   = From("users").Limit(1)
	)

	adapter.On("Query", query).Return(cur, errors.New("error")).Once()

	assert.NotNil(t, repo.Find(&user, query))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_notFound(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		cur     = createCursor(0)
		query   = From("users").Limit(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	err := repo.Find(&user, query)
	assert.Equal(t, NoResultError{}, err)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_MustFind(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users").Limit(1)
		cur     = createCursor(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustFind(&user, query)
	})

	assert.Equal(t, 10, user.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAll(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users").Limit(1)
		cur     = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.FindAll(&users, query))
	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAll_error(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users").Limit(1)
		err     = errors.New("error")
	)

	adapter.On("Query", query).Return(&testCursor{}, err).Once()

	assert.Equal(t, err, repo.FindAll(&users, query))

	adapter.AssertExpectations(t)
}

func TestRepository_MustFindAll(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		query   = From("users").Limit(1)
		cur     = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustFindAll(&users, query)
	})

	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert(t *testing.T) {
	var (
		user      User
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		cbuilders = []Changer{
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		changes = BuildChanges(cbuilders...)
		cur     = createCursor(1)
	)

	adapter.On("Insert", From("users"), changes).Return(1, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 1)).Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Insert(&user, cbuilders...))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert_oneError(t *testing.T) {
	var (
		user      User
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		cbuilders = []Changer{
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		changes = BuildChanges(cbuilders...)
		cur     = &testCursor{}
		err     = errors.New("error")
	)

	adapter.On("Insert", From("users"), changes).Return(1, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 1)).Limit(1)).Return(cur, err).Once()

	assert.Equal(t, err, repo.Insert(&user, cbuilders...))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert_saveBelongsToError(t *testing.T) {
	var (
		address = Address{
			Street: "street",
			User:   &User{Name: "name"},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(&address))

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasOneError(t *testing.T) {
	var (
		user = User{
			Name: "name",
			Address: Address{
				Street: "street",
			},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		cur     = createCursor(1)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 1)).Limit(1)).Return(cur, nil).Once()
	adapter.On("Insert", From("addresses"), mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(&user))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert_saveHasManyError(t *testing.T) {
	var (
		user = User{
			Name: "name",
			Transactions: []Transaction{
				{Item: "soap"},
			},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		cur     = createCursor(1)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 1)).Limit(1)).Return(cur, nil).Once()
	adapter.On("InsertAll", From("transactions"), []string{"item", "user_id"}, mock.Anything).Return([]interface{}{}, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(&user))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert_error(t *testing.T) {
	var (
		user      User
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		cbuilders = []Changer{
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		changes = BuildChanges(cbuilders...)
	)

	adapter.On("Insert", From("users"), changes).Return(0, errors.New("error")).Once()

	assert.NotNil(t, repo.Insert(&user, cbuilders...))
	assert.Panics(t, func() { repo.MustInsert(&user, cbuilders...) })

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
	)

	assert.Nil(t, repo.Insert(nil))
	assert.NotPanics(t, func() { repo.MustInsert(nil) })

	adapter.AssertExpectations(t)
}

func TestRepository_InsertAll(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		changes = []Changes{
			BuildChanges(Set("name", "name1")),
			BuildChanges(Set("name", "name2")),
		}
		cur = createCursor(2)
	)

	adapter.On("InsertAll", From("users"), []string{"name"}, changes).Return([]interface{}{1, 2}, nil).Once()
	adapter.On("Query", From("users").Where(In("id", 1, 2))).Return(cur, nil).Once()

	assert.Nil(t, repo.InsertAll(&users, changes...))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_InsertAll_collection(t *testing.T) {
	var (
		users = []User{
			{Name: "name1"},
			{Name: "name2"},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		changes = []Changes{
			BuildChanges(
				Set("name", "name1"),
				Set("created_at", now()),
				Set("updated_at", now()),
			),
			BuildChanges(
				Set("name", "name2"),
				Set("created_at", now()),
				Set("updated_at", now()),
			),
		}
		cur = createCursor(2)
	)

	adapter.On("InsertAll", From("users"), []string{"name", "created_at", "updated_at"}, changes).Return([]interface{}{1, 2}, nil).Once()
	adapter.On("Query", From("users").Where(In("id", 1, 2))).Return(cur, nil).Once()

	assert.Nil(t, repo.InsertAll(&users))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_InsertAll_unchanged(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
	)

	assert.Nil(t, repo.InsertAll(&users))

	adapter.AssertExpectations(t)
}

func TestRepository_InsertAll_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
	)

	assert.Nil(t, repo.InsertAll(nil))
	assert.NotPanics(t, func() { repo.MustInsertAll(nil) })

	adapter.AssertExpectations(t)
}

func TestRepository_Update(t *testing.T) {
	var (
		user      = User{ID: 1}
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		cbuilders = []Changer{
			Set("name", "name"),
			Set("updated_at", now()),
		}
		changes = BuildChanges(cbuilders...)
		queries = From("users").Where(Eq("id", user.ID))
		cur     = createCursor(1)
	)

	adapter.On("Update", queries, changes).Return(nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Update(&user, cbuilders...))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Update_oneError(t *testing.T) {
	var (
		user      = User{ID: 1}
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		cbuilders = []Changer{
			Set("name", "name"),
			Set("updated_at", now()),
		}
		changes = BuildChanges(cbuilders...)
		queries = From("users").Where(Eq("id", user.ID))
		cur     = &testCursor{}
		err     = errors.New("error")
	)

	adapter.On("Update", queries, changes).Return(nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, err).Once()

	assert.Equal(t, err, repo.Update(&user, cbuilders...))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsToError(t *testing.T) {
	var (
		address = Address{
			ID: 1,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		queries = From("users").Where(Eq("id", address.ID))
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", queries, mock.Anything).Return(err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(&address))

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasOneError(t *testing.T) {
	var (
		user = User{
			ID: 10,
			Address: Address{
				ID:     1,
				Street: "street",
			},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		cur     = createCursor(1)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 10)).Limit(1)).Return(cur, nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10)), mock.Anything).Return(err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(&user))

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasManyError(t *testing.T) {
	var (
		user = User{
			ID: 10,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "soap",
				},
			},
		}
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		cur     = createCursor(1)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 10)).Limit(1)).Return(cur, nil).Once()
	adapter.On("Delete", From("transactions").Where(Eq("user_id", 10))).Return(err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(&user))

	adapter.AssertExpectations(t)
}

func TestRepository_Update_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
	)

	assert.Nil(t, repo.Update(nil))
	assert.NotPanics(t, func() { repo.MustUpdate(nil) })

	adapter.AssertExpectations(t)
}

func TestRepository_Update_error(t *testing.T) {
	var (
		user      = User{ID: 1}
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		cbuilders = []Changer{
			Set("name", "name"),
			Set("updated_at", now()),
		}
		changes = BuildChanges(cbuilders...)
		queries = From("users").Where(Eq("id", user.ID))
	)

	adapter.On("Update", queries, changes).Return(errors.New("error")).Once()

	assert.NotNil(t, repo.Update(&user, cbuilders...))
	assert.Panics(t, func() { repo.MustUpdate(&user, cbuilders...) })
	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_update(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = &Transaction{Buyer: User{ID: 1}}
		doc         = NewDocument(transaction)
		changes     = BuildChanges(
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"updated_at": now(),
				},
			},
		)
		q        = BuildQuery("users", Eq("id", 1))
		buyer, _ = changes.GetAssoc("buyer")
		cur      = createCursor(1)
	)

	adapter.On("Update", q, buyer.Changes[0]).Return(nil).Once()
	adapter.On("Query", q.Limit(1)).Return(cur, nil).Once()

	err := repo.saveBelongsTo(doc, &changes)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_updateError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = &Transaction{Buyer: User{ID: 1}}
		doc         = NewDocument(transaction)
		changes     = BuildChanges(
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"updated_at": now(),
				},
			},
		)
		q        = BuildQuery("users", Eq("id", 1))
		buyer, _ = changes.GetAssoc("buyer")
	)

	adapter.On("Update", q, buyer.Changes[0]).Return(errors.New("update error")).Once()

	err := repo.saveBelongsTo(doc, &changes)
	assert.Equal(t, errors.New("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_updateInconsistentPrimaryKey(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = &Transaction{Buyer: User{ID: 1}}
		doc         = NewDocument(transaction)
		changes     = BuildChanges(
			Map{
				"buyer": Map{
					"id":   2,
					"name": "buyer1",
					"age":  20,
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.saveBelongsTo(doc, &changes)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_insertNew(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = &Transaction{}
		doc         = NewDocument(transaction)
		changes     = BuildChanges(
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"created_at": now(),
					"updated_at": now(),
				},
			},
		)
		q        = BuildQuery("users")
		buyer, _ = changes.GetAssoc("buyer")
		cur      = createCursor(1)
	)

	adapter.On("Insert", q, buyer.Changes[0]).Return(1, nil).Once()
	adapter.On("Query", q.Where(Eq("id", 1)).Limit(1)).Return(cur, nil).Once()

	err := repo.saveBelongsTo(doc, &changes)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	ref, ok := changes.Get("user_id")
	assert.True(t, ok)
	assert.Equal(t, Set("user_id", 10), ref)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_insertNewError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = &Transaction{}
		doc         = NewDocument(transaction)
		changes     = BuildChanges(
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"created_at": now(),
					"updated_at": now(),
				},
			},
		)
		q        = BuildQuery("users")
		buyer, _ = changes.GetAssoc("buyer")
	)

	adapter.On("Insert", q, buyer.Changes[0]).Return(0, errors.New("insert error")).Once()

	err := repo.saveBelongsTo(doc, &changes)
	assert.Equal(t, errors.New("insert error"), err)

	_, ok := changes.Get("user_id")
	assert.False(t, ok)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_notChanged(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = &Transaction{}
		doc         = NewDocument(transaction)
		changes     = BuildChanges()
	)

	err := repo.saveBelongsTo(doc, &changes)
	assert.Nil(t, err)
	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_update(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{ID: 1, Address: Address{ID: 2}}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		q            = BuildQuery("addresses").Where(Eq("id", 2).AndEq("user_id", 1))
		addresses, _ = changes.GetAssoc("address")
		cur          = createCursor(1)
	)

	adapter.On("Update", q, addresses.Changes[0]).Return(nil).Once()
	adapter.On("Query", q.Limit(1)).Return(cur, nil).Once()

	err := repo.saveHasOne(doc, &changes)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveHasOne_updateError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{ID: 1, Address: Address{ID: 2}}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		q            = BuildQuery("addresses").Where(Eq("id", 2).AndEq("user_id", 1))
		addresses, _ = changes.GetAssoc("address")
	)

	adapter.On("Update", q, addresses.Changes[0]).Return(errors.New("update error")).Once()

	err := repo.saveHasOne(doc, &changes)
	assert.Equal(t, errors.New("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_updateInconsistentPrimaryKey(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{ID: 1, Address: Address{ID: 2}}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"address": Map{
					"id":     1,
					"street": "street1",
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.saveHasOne(doc, &changes)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_insertNew(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		q       = BuildQuery("addresses")
		address = BuildChanges(Set("street", "street1"))
		cur     = createCursor(1)
	)

	// foreign value set after associations infered
	user.ID = 1
	address.SetValue("user_id", user.ID)

	adapter.On("Insert", q, address).Return(2, nil).Once()
	adapter.On("Query", q.Where(Eq("id", 2)).Limit(1)).Return(cur, nil).Once()

	err := repo.saveHasOne(doc, &changes)
	assert.Nil(t, err)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveHasOne_insertNewError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		q       = BuildQuery("addresses")
		address = BuildChanges(Set("street", "street1"))
	)

	// foreign value set after associations infered
	user.ID = 1
	address.SetValue("user_id", user.ID)

	adapter.On("Insert", q, address).Return(nil, errors.New("insert error")).Once()

	err := repo.saveHasOne(doc, &changes)
	assert.Equal(t, errors.New("insert error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_insert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{ID: 1}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "item1",
					},
					{
						"item": "item2",
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		cur             = createCursor(2)
	)

	adapter.On("InsertAll", q, []string{"item", "user_id"}, transactions.Changes).Return(nil).Return([]interface{}{2, 3}, nil).Once()
	adapter.On("Query", q.Where(In("id", 2, 3))).Return(cur, nil).Once()

	err := repo.saveHasMany(doc, &changes, true)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(transactions.Changes))

	for i := range transactions.Changes {
		id, ok := transactions.Changes[i].GetValue("user_id")
		assert.True(t, ok)
		assert.Equal(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveHasMany_insertError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{ID: 1}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "item1",
					},
					{
						"item": "item2",
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		rerr            = errors.New("insert all error")
	)

	adapter.On("InsertAll", q, []string{"item", "user_id"}, transactions.Changes).Return(nil).Return([]interface{}{}, rerr).Once()

	err := repo.saveHasMany(doc, &changes, true)
	assert.Equal(t, rerr, err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_update(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
				{
					ID:   2,
					Item: "item2",
				},
				{
					ID:   3,
					Item: "item3",
				},
			},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"id":   1,
						"item": "item1 updated",
					},
					{
						"id":   2,
						"item": "item2 updated",
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		cur1            = createCursor(1)
		cur2            = createCursor(1)
	)

	changes.assocChanges[0].StaleIDs = []interface{}{3}

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 3))).Return(nil).Once()
	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), transactions.Changes[0]).Return(nil).Once()
	adapter.On("Update", q.Where(Eq("id", 2).AndEq("user_id", 1)), transactions.Changes[1]).Return(nil).Once()
	adapter.On("Query", q.Where(Eq("id", 1).AndEq("user_id", 1)).Limit(1)).Return(cur1, nil).Once()
	adapter.On("Query", q.Where(Eq("id", 2).AndEq("user_id", 1)).Limit(1)).Return(cur2, nil).Once()

	err := repo.saveHasMany(doc, &changes, false)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(transactions.Changes))
	assert.False(t, cur1.Next())
	assert.False(t, cur2.Next())

	for i := range transactions.Changes {
		id, ok := transactions.Changes[i].GetValue("user_id")
		assert.False(t, ok)
		assert.NotEqual(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur1.AssertExpectations(t)
	cur2.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateWithInsert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
			},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"id":   1,
						"item": "item1 updated",
					},
					{
						"item":    "new item",
						"user_id": 1,
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		updateCur       = createCursor(1)
		insertAllCur    = createCursor(1)
	)

	changes.assocChanges[0].StaleIDs = []interface{}{3}

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 3))).Return(nil).Once()
	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), transactions.Changes[0]).Return(nil).Once()
	adapter.On("Query", q.Where(Eq("id", 1).AndEq("user_id", 1)).Limit(1)).Return(updateCur, nil).Once()
	adapter.On("InsertAll", q, []string{"item", "user_id"}, transactions.Changes[1:]).Return(nil).Return([]interface{}{2}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "item"}, transactions.Changes[1:]).Return(nil).Return([]interface{}{2}, nil).Maybe()
	adapter.On("Query", q.Where(In("id", 2))).Return(insertAllCur, nil).Once()

	err := repo.saveHasMany(doc, &changes, false)
	assert.Nil(t, err)

	assert.Equal(t, 2, len(transactions.Changes))
	assert.False(t, updateCur.Next())

	adapter.AssertExpectations(t)
	updateCur.AssertExpectations(t)
	insertAllCur.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateDeleteError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
			},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "new item",
					},
				},
			},
		)
		q   = BuildQuery("transactions")
		err = errors.New("error")
	)

	changes.assocChanges[0].StaleIDs = []interface{}{1}

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 1))).Return(err).Once()

	assert.Equal(t, err, repo.saveHasMany(doc, &changes, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
			},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"id":   1,
						"item": "item1 updated",
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		err             = errors.New("error")
	)

	changes.assocChanges[0].StaleIDs = []interface{}{}

	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), transactions.Changes[0]).Return(err).Once()

	assert.Equal(t, err, repo.saveHasMany(doc, &changes, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_replace(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
				{
					ID:   2,
					Item: "item2",
				},
			},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "item3",
					},
					{
						"item": "item4",
					},
					{
						"item": "item5",
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		cur             = createCursor(3)
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(nil).Once()
	adapter.On("InsertAll", q, []string{"item", "user_id"}, transactions.Changes).Return(nil).Return([]interface{}{3, 4, 5}, nil).Once()
	adapter.On("Query", q.Where(In("id", 3, 4, 5))).Return(cur, nil).Once()

	err := repo.saveHasMany(doc, &changes, false)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(transactions.Changes))

	for i := range transactions.Changes {
		id, ok := transactions.Changes[i].GetValue("user_id")
		assert.True(t, ok)
		assert.Equal(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveHasMany_replaceEmptyAssoc(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID:           1,
			Transactions: []Transaction{},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "item3",
					},
					{
						"item": "item4",
					},
					{
						"item": "item5",
					},
				},
			},
		)
		q               = BuildQuery("transactions")
		transactions, _ = changes.GetAssoc("transactions")
		cur             = createCursor(3)
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(nil).Once()
	adapter.On("InsertAll", q, []string{"item", "user_id"}, transactions.Changes).Return(nil).Return([]interface{}{3, 4, 5}, nil).Once()
	adapter.On("Query", q.Where(In("id", 3, 4, 5))).Return(cur, nil).Once()

	err := repo.saveHasMany(doc, &changes, false)
	assert.Nil(t, err)

	assert.Equal(t, 3, len(transactions.Changes))

	for i := range transactions.Changes {
		id, ok := transactions.Changes[i].GetValue("user_id")
		assert.True(t, ok)
		assert.Equal(t, user.ID, id)
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_saveHasMany_replaceDeleteAllError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{
			ID: 1,
			Transactions: []Transaction{
				{
					ID:   1,
					Item: "item1",
				},
				{
					ID:   2,
					Item: "item2",
				},
			},
		}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "item3",
					},
				},
			},
		)
		q    = BuildQuery("transactions")
		rerr = errors.New("delete all error")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(rerr).Once()

	err := repo.saveHasMany(doc, &changes, false)
	assert.Equal(t, rerr, err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_replaceAssocNotLoaded(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = &User{ID: 1}
		doc     = NewDocument(user)
		changes = BuildChanges(
			Map{
				"transactions": []Map{
					{
						"item": "item3",
					},
				},
			},
		)
	)

	assert.Panics(t, func() {
		repo.saveHasMany(doc, &changes, false)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Save_insert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = User{Name: "name"}
		cur     = createCursor(1)
	)

	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 1)).Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Save(&user))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Save_update(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = User{ID: 1, Name: "name"}
		queries = From("users").Where(Eq("id", 1))
		cur     = createCursor(1)
	)

	adapter.On("Update", queries, mock.Anything).Return(nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Save(&user))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Save_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
	)

	assert.Nil(t, repo.Save(nil))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = User{ID: 1}
	)

	adapter.On("Delete", From("users").Where(Eq("id", user.ID))).Return(nil).Once()

	assert.Nil(t, repo.Delete(&user))

	adapter.AssertExpectations(t)
}

func TestRepository_MustDelete(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = User{ID: 1}
	)

	adapter.On("Delete", From("users").Where(Eq("id", user.ID))).Return(nil).Once()

	assert.NotPanics(t, func() {
		repo.MustDelete(&user)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_DeleteAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		queries = From("logs").Where(Eq("user_id", 1))
	)

	adapter.On("Delete", From("logs").Where(Eq("user_id", 1))).Return(nil).Once()

	assert.Nil(t, repo.DeleteAll(queries))

	adapter.AssertExpectations(t)
}

func TestRepository_MustDeleteAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		queries = From("logs").Where(Eq("user_id", 1))
	)

	adapter.On("Delete", From("logs").Where(Eq("user_id", 1))).Return(nil).Once()

	assert.NotPanics(t, func() {
		repo.MustDeleteAll(queries)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Preload_hasOne(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = User{ID: 10}
		address = Address{ID: 100, UserID: &user.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", From("addresses").Where(In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(address.ID, *address.UserID).Times(2)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&user, "address"))
	assert.Equal(t, address, user.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceHasOne(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		users     = []User{{ID: 10}, {ID: 20}}
		addresses = []Address{
			{ID: 100, UserID: &users[0].ID},
			{ID: 200, UserID: &users[1].ID},
		}
		cur = &testCursor{}
	)

	// one of these, because of map ordering
	adapter.On("Query", From("addresses").Where(In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", From("addresses").Where(In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(addresses[0].ID, *addresses[0].UserID).Twice()
	cur.MockScan(addresses[1].ID, *addresses[1].UserID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&users, "address"))
	assert.Equal(t, addresses[0], users[0].Address)
	assert.Equal(t, addresses[1], users[1].Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedHasOne(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = Transaction{
			Buyer: User{ID: 10},
		}
		address = Address{ID: 100, UserID: &transaction.Buyer.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", From("addresses").Where(In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(address.ID, *address.UserID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transaction, "buyer.address"))
	assert.Equal(t, address, transaction.Buyer.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceNestedHasOne(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = repository{adapter: adapter}
		transactions = []Transaction{
			{Buyer: User{ID: 10}},
			{Buyer: User{ID: 20}},
		}
		addresses = []Address{
			{ID: 100, UserID: &transactions[0].Buyer.ID},
			{ID: 200, UserID: &transactions[1].Buyer.ID},
		}
		cur = &testCursor{}
	)

	// one of these, because of map ordering
	adapter.On("Query", From("addresses").Where(In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", From("addresses").Where(In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(addresses[0].ID, *addresses[0].UserID).Twice()
	cur.MockScan(addresses[1].ID, *addresses[1].UserID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transactions, "buyer.address"))
	assert.Equal(t, addresses[0], transactions[0].Buyer.Address)
	assert.Equal(t, addresses[1], transactions[1].Buyer.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_hasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = repository{adapter: adapter}
		user         = User{ID: 10}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("transactions").Where(In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(transactions[0].ID, transactions[0].BuyerID).Twice()
	cur.MockScan(transactions[1].ID, transactions[1].BuyerID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&user, "transactions"))
	assert.Equal(t, transactions, user.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceHasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = repository{adapter: adapter}
		users        = []User{{ID: 10}, {ID: 20}}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
			{ID: 15, BuyerID: 20},
			{ID: 20, BuyerID: 20},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("transactions").Where(In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", From("transactions").Where(In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Times(4)
	cur.MockScan(transactions[0].ID, transactions[0].BuyerID).Twice()
	cur.MockScan(transactions[1].ID, transactions[1].BuyerID).Twice()
	cur.MockScan(transactions[2].ID, transactions[2].BuyerID).Twice()
	cur.MockScan(transactions[3].ID, transactions[3].BuyerID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&users, "transactions"))
	assert.Equal(t, transactions[:2], users[0].Transactions)
	assert.Equal(t, transactions[2:], users[1].Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedHasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = repository{adapter: adapter}
		address      = Address{User: &User{ID: 10}}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
		}

		cur = &testCursor{}
	)

	adapter.On("Query", From("transactions").Where(In("user_id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(transactions[0].ID, transactions[0].BuyerID).Twice()
	cur.MockScan(transactions[1].ID, transactions[1].BuyerID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&address, "user.transactions"))
	assert.Equal(t, transactions, address.User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedNullHasMany(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		address = Address{User: nil}
	)

	assert.Nil(t, repo.Preload(&address, "user.transactions"))

	adapter.AssertExpectations(t)
}

func TestRepository_Preload_nestedSliceHasMany(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		addresses = []Address{
			{User: &User{ID: 10}},
			{User: &User{ID: 20}},
		}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
			{ID: 15, BuyerID: 20},
			{ID: 20, BuyerID: 20},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("transactions").Where(In("user_id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", From("transactions").Where(In("user_id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Times(4)
	cur.MockScan(transactions[0].ID, transactions[0].BuyerID).Twice()
	cur.MockScan(transactions[1].ID, transactions[1].BuyerID).Twice()
	cur.MockScan(transactions[2].ID, transactions[2].BuyerID).Twice()
	cur.MockScan(transactions[3].ID, transactions[3].BuyerID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&addresses, "user.transactions"))
	assert.Equal(t, transactions[:2], addresses[0].User.Transactions)
	assert.Equal(t, transactions[2:], addresses[1].User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedNullSliceHasMany(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = repository{adapter: adapter}
		addresses = []Address{
			{User: &User{ID: 10}},
			{User: nil},
			{User: &User{ID: 15}},
		}
		transactions = []Transaction{
			{ID: 5, BuyerID: 10},
			{ID: 10, BuyerID: 10},
			{ID: 15, BuyerID: 15},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("transactions").Where(In("user_id", 10, 15))).Return(cur, nil).Maybe()
	adapter.On("Query", From("transactions").Where(In("user_id", 15, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Times(3)
	cur.MockScan(transactions[0].ID, transactions[0].BuyerID).Twice()
	cur.MockScan(transactions[1].ID, transactions[1].BuyerID).Twice()
	cur.MockScan(transactions[2].ID, transactions[2].BuyerID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&addresses, "user.transactions"))
	assert.Equal(t, transactions[:2], addresses[0].User.Transactions)
	assert.Equal(t, []Transaction(nil), addresses[1].User.Transactions)
	assert.Equal(t, transactions[2:], addresses[2].User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_belongsTo(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		user        = User{ID: 10, Name: "Del Piero"}
		transaction = Transaction{BuyerID: 10}
		cur         = &testCursor{}
	)

	adapter.On("Query", From("users").Where(In("id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(user.ID, user.Name).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transaction, "buyer"))
	assert.Equal(t, user, transaction.Buyer)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_ptrBelongsTo(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		user    = User{ID: 10, Name: "Del Piero"}
		address = Address{UserID: &user.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", From("users").Where(In("id", 10))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(user.ID, user.Name).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&address, "user"))
	assert.Equal(t, user, *address.User)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nullBelongsTo(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		address = Address{}
	)

	assert.Nil(t, repo.Preload(&address, "user"))
	assert.Nil(t, address.User)

	adapter.AssertExpectations(t)
}

func TestRepository_Preload_sliceBelongsTo(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = repository{adapter: adapter}
		transactions = []Transaction{
			{BuyerID: 10},
			{BuyerID: 20},
		}
		users = []User{
			{ID: 10, Name: "Del Piero"},
			{ID: 20, Name: "Nedved"},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("users").Where(In("id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", From("users").Where(In("id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(users[0].ID, users[0].Name).Twice()
	cur.MockScan(users[1].ID, users[1].Name).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&transactions, "buyer"))
	assert.Equal(t, users[0], transactions[0].Buyer)
	assert.Equal(t, users[1], transactions[1].Buyer)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_ptrSliceBelongsTo(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = repository{adapter: adapter}
		users   = []User{
			{ID: 10, Name: "Del Piero"},
			{ID: 20, Name: "Nedved"},
		}
		addresses = []Address{
			{UserID: &users[0].ID},
			{UserID: &users[1].ID},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("users").Where(In("id", 10, 20))).Return(cur, nil).Maybe()
	adapter.On("Query", From("users").Where(In("id", 20, 10))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(users[0].ID, users[0].Name).Twice()
	cur.MockScan(users[1].ID, users[1].Name).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(&addresses, "user"))
	assert.Equal(t, users[0], *addresses[0].User)
	assert.Equal(t, users[1], *addresses[1].User)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_emptySlice(t *testing.T) {
	var (
		repo      = repository{}
		addresses = []Address{}
	)

	assert.Nil(t, repo.Preload(&addresses, "user.transactions"))
}

func TestQuery_Preload_notPointerPanic(t *testing.T) {
	var (
		repo        = repository{}
		transaction = Transaction{}
	)

	assert.Panics(t, func() { repo.Preload(transaction, "User") })
}

func TestRepository_Preload_queryError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = Transaction{BuyerID: 10}
		cur         = &testCursor{}
		err         = errors.New("error")
	)

	adapter.On("Query", From("users").Where(In("id", 10))).Return(cur, err).Once()

	assert.Equal(t, err, repo.Preload(&transaction, "buyer"))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_MustPreload(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = repository{adapter: adapter}
		transaction = Transaction{BuyerID: 10}
		cur         = createCursor(0)
	)

	adapter.On("Query", From("users").Where(In("id", 10))).Return(cur, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustPreload(&transaction, "buyer")
	})

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Transaction(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).On("Commit").Return(nil).Once()

	repo := repository{adapter: adapter}

	err := repo.Transaction(func(repo Repository) error {
		assert.True(t, repo.(*repository).inTransaction)
		return nil
	})

	assert.False(t, repo.inTransaction)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_beginError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(errors.New("error")).Once()

	err := repository{adapter: adapter}.Transaction(func(r Repository) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.New("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_commitError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Commit").Return(errors.New("error")).Once()

	err := repository{adapter: adapter}.Transaction(func(r Repository) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.New("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_returnErrorAndRollback(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	err := repository{adapter: adapter}.Transaction(func(r Repository) error {
		// doing good things
		return errors.New("error")
	})

	assert.Equal(t, errors.New("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_panicWithErrorAndRollback(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	err := repository{adapter: adapter}.Transaction(func(r Repository) error {
		// doing good things
		panic(errors.New("error"))
	})

	assert.Equal(t, errors.New("error"), err)
	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_panicWithStringAndRollback(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Panics(t, func() {
		repository{adapter: adapter}.Transaction(func(r Repository) error {
			// doing good things
			panic("error")
		})
	})

	adapter.AssertExpectations(t)
}
