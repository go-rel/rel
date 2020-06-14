package rel

import (
	"context"
	"errors"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func init() {
	t := time.Now().Truncate(time.Second)
	now = func() time.Time {
		return t
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

func TestRepository_Instrumentation(t *testing.T) {
	var (
		repo = repository{adapter: &testAdapter{}}
	)

	assert.Nil(t, repo.instrumenter)
	assert.NotPanics(t, func() {
		repo.instrument(context.TODO(), "test", "test")(nil)
	})

	repo.Instrumentation(DefaultLogger)
	assert.NotNil(t, repo.instrumenter)
	assert.NotPanics(t, func() {
		repo.instrument(context.TODO(), "test", "test")(nil)
	})
}

func TestRepository_Ping(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Ping").Return(nil).Once()

	assert.Nil(t, repo.Ping(context.TODO()))
	adapter.AssertExpectations(t)
}

func TestRepository_Iterate(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users")
		cur     = createCursor(1)
	)

	adapter.On("Query", query.SortAsc("id").Limit(1000)).Return(cur, nil).Once()

	it := repo.Iterate(context.TODO(), query)
	defer it.Close()
	for {
		if err := it.Next(&user); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
		}

		assert.NotEqual(t, 0, user.ID)
	}

	adapter.AssertExpectations(t)
}

func TestRepository_Aggregate(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = New(adapter)
		query     = From("users")
		aggregate = "count"
		field     = "*"
	)

	adapter.On("Aggregate", query, aggregate, field).Return(1, nil).Once()

	count, err := repo.Aggregate(context.TODO(), query, "count", "*")
	assert.Equal(t, 1, count)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepository_MustAggregate(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = New(adapter)
		query     = From("users")
		aggregate = "count"
		field     = "*"
	)

	adapter.On("Aggregate", query, aggregate, field).Return(1, nil).Once()

	assert.NotPanics(t, func() {
		count := repo.MustAggregate(context.TODO(), query, "count", "*")
		assert.Equal(t, 1, count)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Count(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users")
	)

	adapter.On("Aggregate", query, "count", "*").Return(1, nil).Once()

	count, err := repo.Count(context.TODO(), "users")
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	adapter.AssertExpectations(t)
}

func TestRepository_MustCount(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users")
	)

	adapter.On("Aggregate", query, "count", "*").Return(1, nil).Once()

	assert.NotPanics(t, func() {
		count := repo.MustCount(context.TODO(), "users")
		assert.Equal(t, 1, count)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Find(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(1)
		cur     = createCursor(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.Find(context.TODO(), &user, query))
	assert.Equal(t, 10, user.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_softDelete(t *testing.T) {
	var (
		address Address
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("addresses").Limit(1)
		cur     = createCursor(1)
	)

	adapter.On("Query", query.Where(Nil("deleted_at"))).Return(cur, nil).Once()

	assert.Nil(t, repo.Find(context.TODO(), &address, query))
	assert.Equal(t, 10, address.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_softDeleteUnscoped(t *testing.T) {
	var (
		address Address
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("addresses").Limit(1).Unscoped()
		cur     = createCursor(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.Find(context.TODO(), &address, query))
	assert.Equal(t, 10, address.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_queryError(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		cur     = &testCursor{}
		query   = From("users").Limit(1)
	)

	adapter.On("Query", query).Return(cur, errors.New("error")).Once()

	assert.NotNil(t, repo.Find(context.TODO(), &user, query))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_notFound(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		cur     = createCursor(0)
		query   = From("users").Limit(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	err := repo.Find(context.TODO(), &user, query)
	assert.Equal(t, NotFoundError{}, err)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_MustFind(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(1)
		cur     = createCursor(1)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustFind(context.TODO(), &user, query)
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
		repo    = New(adapter)
		query   = From("users").Limit(1)
		cur     = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.FindAll(context.TODO(), &users, query))
	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAll_softDelete(t *testing.T) {
	var (
		addresses []Address
		adapter   = &testAdapter{}
		repo      = New(adapter)
		query     = From("addresses").Limit(1)
		cur       = createCursor(2)
	)

	adapter.On("Query", query.Where(Nil("deleted_at"))).Return(cur, nil).Once()

	assert.Nil(t, repo.FindAll(context.TODO(), &addresses, query))
	assert.Len(t, addresses, 2)
	assert.Equal(t, 10, addresses[0].ID)
	assert.Equal(t, 10, addresses[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAll_softDeleteUnscoped(t *testing.T) {
	var (
		addresses []Address
		adapter   = &testAdapter{}
		repo      = New(adapter)
		query     = From("addresses").Limit(1).Unscoped()
		cur       = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Nil(t, repo.FindAll(context.TODO(), &addresses, query))
	assert.Len(t, addresses, 2)
	assert.Equal(t, 10, addresses[0].ID)
	assert.Equal(t, 10, addresses[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAll_error(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(1)
		err     = errors.New("error")
	)

	adapter.On("Query", query).Return(&testCursor{}, err).Once()

	assert.Equal(t, err, repo.FindAll(context.TODO(), &users, query))

	adapter.AssertExpectations(t)
}

func TestRepository_MustFindAll(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(1)
		cur     = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustFindAll(context.TODO(), &users, query)
	})

	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAndCountAll(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(10)
		cur     = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()
	adapter.On("Aggregate", query.Limit(0), "count", "*").Return(2, nil).Once()

	count, err := repo.FindAndCountAll(context.TODO(), &users, query)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAndCountAll_queryError(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(10)
		err     = errors.New("error")
	)

	adapter.On("Query", query).Return(&testCursor{}, err).Once()

	count, ferr := repo.FindAndCountAll(context.TODO(), &users, query)
	assert.Equal(t, err, ferr)
	assert.Equal(t, 0, count)

	adapter.AssertExpectations(t)
}

func TestRepository_MustFindAndCountAll(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(10)
		cur     = createCursor(2)
	)

	adapter.On("Query", query).Return(cur, nil).Once()
	adapter.On("Aggregate", query.Limit(0), "count", "*").Return(2, nil).Once()

	assert.NotPanics(t, func() {
		count := repo.MustFindAndCountAll(context.TODO(), &users, query)
		assert.Equal(t, 2, count)
	})

	assert.Len(t, users, 2)
	assert.Equal(t, 10, users[0].ID)
	assert.Equal(t, 10, users[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert(t *testing.T) {
	var (
		user     User
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
	)

	adapter.On("Insert", From("users"), mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user, mutators...))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_reload(t *testing.T) {
	var (
		user     User
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Reload(true),
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
		cur = createCursor(1)
	)

	adapter.On("Insert", From("users"), mutates).Return(10, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 10)).Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user, mutators...))
	assert.Equal(t, User{
		ID:        10,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
	}, user)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert_reloadError(t *testing.T) {
	var (
		user     User
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Reload(true),
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
		cur = &testCursor{}
		err = errors.New("error")
	)

	adapter.On("Insert", From("users"), mutates).Return(10, nil).Once()
	adapter.On("Query", From("users").Where(Eq("id", 10)).Limit(1)).Return(cur, err).Once()

	assert.Equal(t, err, repo.Insert(context.TODO(), &user, mutators...))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Insert_saveBelongsTo(t *testing.T) {
	var (
		userID    = 1
		addressID = 2
		address   = Address{
			Street: "street",
			User:   &User{Name: "name"},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(userID, nil).Once()
	adapter.On("Insert", From("addresses"), mock.Anything).Return(addressID, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &address))
	assert.Equal(t, Address{
		ID:     addressID,
		Street: "street",
		UserID: &userID,
		User: &User{
			ID:        userID,
			Name:      "name",
			CreatedAt: now(),
			UpdatedAt: now(),
		},
	}, address)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveBelongsToCascadeDisabled(t *testing.T) {
	var (
		address = Address{
			Street: "street",
			User:   &User{Name: "name"},
		}
		adapter      = &testAdapter{}
		repo         = New(adapter)
		newAddressID = 2
	)

	adapter.On("Insert", From("addresses"), mock.Anything).Return(newAddressID, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &address, Cascade(false)))
	assert.Equal(t, Address{
		ID:     newAddressID,
		Street: "street",
		User:   &User{Name: "name"},
	}, address)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveBelongsToError(t *testing.T) {
	var (
		address = Address{
			Street: "street",
			User:   &User{Name: "name"},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(context.TODO(), &address))

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasOne(t *testing.T) {
	var (
		userID    = 1
		addressID = 2
		user      = User{
			Name: "name",
			Address: Address{
				Street: "street",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(userID, nil).Once()
	adapter.On("Insert", From("addresses"), mock.Anything).Return(addressID, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        userID,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
		Address: Address{
			ID:     addressID,
			UserID: &userID,
			Street: "street",
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasOneCascadeDisabled(t *testing.T) {
	var (
		userID = 1
		user   = User{
			Name: "name",
			Address: Address{
				Street: "street",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Insert", From("users"), mock.Anything).Return(userID, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user, Cascade(false)))
	assert.Equal(t, User{
		ID:        userID,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
		Address: Address{
			Street: "street",
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasOneError(t *testing.T) {
	var (
		userID = 1
		user   = User{
			Name: "name",
			Address: Address{
				Street: "street",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(userID, nil).Once()
	adapter.On("Insert", From("addresses"), mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(context.TODO(), &user))

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasMany(t *testing.T) {
	var (
		user = User{
			Name: "name",
			Transactions: []Transaction{
				{Item: "soap"},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("InsertAll", From("transactions"), mock.Anything, mock.Anything).Return([]interface{}{2}, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
		Transactions: []Transaction{
			{ID: 2, BuyerID: 1, Item: "soap"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasManyCascadeDisabled(t *testing.T) {
	var (
		user = User{
			Name: "name",
			Transactions: []Transaction{
				{Item: "soap"},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user, Cascade(false)))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
		Transactions: []Transaction{
			{Item: "soap"},
		},
	}, user)

	adapter.AssertExpectations(t)
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
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("InsertAll", From("transactions"), mock.Anything, mock.Anything).Return([]interface{}{}, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(context.TODO(), &user))

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_error(t *testing.T) {
	var (
		user     User
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("name", "name"),
			Set("created_at", now()),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
	)

	adapter.On("Insert", From("users"), mutates).Return(0, errors.New("error")).Once()

	assert.NotNil(t, repo.Insert(context.TODO(), &user, mutators...))
	assert.Panics(t, func() { repo.MustInsert(context.TODO(), &user, mutators...) })

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_customError(t *testing.T) {
	var (
		user     User
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("name", "name"),
			ErrorFunc(func(err error) error { return errors.New("custom error") }),
		}
		mutates = map[string]Mutate{
			"name": Set("name", "name"),
		}
	)

	adapter.On("Insert", From("users"), mutates).Return(0, errors.New("error")).Once()

	assert.Equal(t, errors.New("custom error"), repo.Insert(context.TODO(), &user, mutators...))
	assert.Panics(t, func() { repo.MustInsert(context.TODO(), &user, mutators...) })

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_customErrorNested(t *testing.T) {
	var (
		address = Address{
			Street: "street",
			User: &User{
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, errors.New("error")).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, errors.New("error"), repo.Insert(context.TODO(), &address,
		NewStructset(&address, false),
		ErrorFunc(func(err error) error { return errors.New("custom error") }), // should not transform any errors of its children.
	))
	adapter.AssertExpectations(t)
}

func TestRepository_Insert_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	assert.Nil(t, repo.Insert(context.TODO(), nil))
	assert.NotPanics(t, func() { repo.MustInsert(context.TODO(), nil) })

	adapter.AssertExpectations(t)
}

func TestRepository_InsertAll(t *testing.T) {
	var (
		users = []User{
			{Name: "name1"},
			{Name: "name2", Age: 12},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		mutates = []map[string]Mutate{
			{
				"name":       Set("name", "name1"),
				"age":        Set("age", 0),
				"created_at": Set("created_at", now()),
				"updated_at": Set("updated_at", now()),
			},
			{
				"name":       Set("name", "name2"),
				"age":        Set("age", 12),
				"created_at": Set("created_at", now()),
				"updated_at": Set("updated_at", now()),
			},
		}
	)

	adapter.On("InsertAll", From("users"), mock.Anything, mutates).Return([]interface{}{1, 2}, nil).Once()

	assert.Nil(t, repo.InsertAll(context.TODO(), &users))
	assert.Equal(t, []User{
		{ID: 1, Name: "name1", Age: 0, CreatedAt: now(), UpdatedAt: now()},
		{ID: 2, Name: "name2", Age: 12, CreatedAt: now(), UpdatedAt: now()},
	}, users)

	adapter.AssertExpectations(t)
}

func TestRepository_InsertAll_empty(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	assert.Nil(t, repo.InsertAll(context.TODO(), &users))

	adapter.AssertExpectations(t)
}

func TestRepository_InsertAll_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	assert.Nil(t, repo.InsertAll(context.TODO(), nil))
	assert.NotPanics(t, func() { repo.MustInsertAll(context.TODO(), nil) })

	adapter.AssertExpectations(t)
}

func TestRepository_Update(t *testing.T) {
	var (
		user     = User{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("name", "name"),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"updated_at": Set("updated_at", now()),
		}
		queries = From("users").Where(Eq("id", user.ID))
	)

	adapter.On("Update", queries, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user, mutators...))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		UpdatedAt: now(),
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_softDelete(t *testing.T) {
	var (
		address  = Address{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("street", "street"),
		}
		mutates = map[string]Mutate{
			"street": Set("street", "street"),
		}
		queries = From("addresses").Where(Eq("id", address.ID))
	)

	adapter.On("Update", queries.Where(Nil("deleted_at")), mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &address, mutators...))
	assert.Equal(t, Address{
		ID:     1,
		Street: "street",
	}, address)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_softDeleteUnscoped(t *testing.T) {
	var (
		address  = Address{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Unscoped(true),
			Set("street", "street"),
		}
		mutates = map[string]Mutate{
			"street": Set("street", "street"),
		}
		queries = From("addresses").Where(Eq("id", address.ID)).Unscoped()
	)

	adapter.On("Update", queries, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &address, mutators...))
	assert.Equal(t, Address{
		ID:     1,
		Street: "street",
	}, address)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_notFound(t *testing.T) {
	var (
		user     = User{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("name", "name"),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"updated_at": Set("updated_at", now()),
		}
		queries = From("users").Where(Eq("id", user.ID))
	)

	adapter.On("Update", queries, mutates).Return(0, nil).Once()

	assert.Equal(t, NotFoundError{}, repo.Update(context.TODO(), &user, mutators...))

	adapter.AssertExpectations(t)
}

func TestRepository_Update_reload(t *testing.T) {
	var (
		user     = User{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			SetFragment("name=?", "name"),
		}
		mutates = map[string]Mutate{
			"name=?": SetFragment("name=?", "name"),
		}
		queries = From("users").Where(Eq("id", user.ID))
		cur     = createCursor(1)
	)

	adapter.On("Update", queries, mutates).Return(1, nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user, mutators...))
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Update_reloadError(t *testing.T) {
	var (
		user     = User{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			SetFragment("name=?", "name"),
		}
		mutates = map[string]Mutate{
			"name=?": SetFragment("name=?", "name"),
		}
		queries = From("users").Where(Eq("id", user.ID))
		cur     = &testCursor{}
		err     = errors.New("error")
	)

	adapter.On("Update", queries, mutates).Return(1, nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, err).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &user, mutators...))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsTo(t *testing.T) {
	var (
		userID  = 1
		address = Address{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", *address.UserID)), mock.Anything).Return(1, nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", address.ID).AndNil("deleted_at")), mock.Anything).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &address))
	assert.Equal(t, Address{
		ID:     1,
		UserID: &userID,
		User: &User{
			ID:        1,
			Name:      "name",
			UpdatedAt: now(),
			CreatedAt: now(),
		},
	}, address)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsToCascadeDisabled(t *testing.T) {
	var (
		userID  = 1
		address = Address{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Update", From("addresses").Where(Eq("id", address.ID).AndNil("deleted_at")), mock.Anything).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &address, Cascade(false)))
	assert.Equal(t, Address{
		ID:     1,
		UserID: &userID,
		User: &User{
			ID:   1,
			Name: "name",
		},
	}, address)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsToError(t *testing.T) {
	var (
		userID  = 1
		address = Address{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		queries = From("users").Where(Eq("id", address.ID))
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", queries, mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &address))

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasOne(t *testing.T) {
	var (
		userID = 10
		user   = User{
			ID: userID,
			Address: Address{
				ID:     1,
				Street: "street",
				UserID: &userID,
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(1, nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10).AndNil("deleted_at")), mock.Anything).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        userID,
		UpdatedAt: now(),
		CreatedAt: now(),
		Address: Address{
			ID:     1,
			Street: "street",
			UserID: &userID,
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasOneCascadeDisabled(t *testing.T) {
	var (
		userID = 10
		user   = User{
			ID: userID,
			Address: Address{
				ID:     1,
				Street: "street",
				UserID: &userID,
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user, Cascade(false)))
	assert.Equal(t, User{
		ID:        userID,
		UpdatedAt: now(),
		CreatedAt: now(),
		Address: Address{
			ID:     1,
			Street: "street",
			UserID: &userID,
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasOneError(t *testing.T) {
	var (
		userID = 10
		user   = User{
			ID: userID,
			Address: Address{
				ID:     1,
				Street: "street",
				UserID: &userID,
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(1, nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10).AndNil("deleted_at")), mock.Anything).Return(1, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &user))
	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasMany(t *testing.T) {
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
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(1, nil).Once()
	adapter.On("Delete", From("transactions").Where(Eq("user_id", 10))).Return(1, nil).Once()
	adapter.On("InsertAll", From("transactions"), mock.Anything, mock.Anything).Return([]interface{}{1}, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        10,
		CreatedAt: now(),
		UpdatedAt: now(),
		Transactions: []Transaction{
			{
				ID:      1,
				BuyerID: 10,
				Item:    "soap",
			},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasManyCascadeDisabled(t *testing.T) {
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
		repo    = New(adapter)
	)

	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user, Cascade(false)))
	assert.Equal(t, User{
		ID:        10,
		CreatedAt: now(),
		UpdatedAt: now(),
		Transactions: []Transaction{
			{
				ID:   1,
				Item: "soap",
			},
		},
	}, user)

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
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), mock.Anything).Return(1, nil).Once()
	adapter.On("Delete", From("transactions").Where(Eq("user_id", 10))).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &user))
	adapter.AssertExpectations(t)
}

func TestRepository_Update_nothing(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	assert.Nil(t, repo.Update(context.TODO(), nil))
	assert.NotPanics(t, func() { repo.MustUpdate(context.TODO(), nil) })

	adapter.AssertExpectations(t)
}

func TestRepository_Update_error(t *testing.T) {
	var (
		user     = User{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		mutators = []Mutator{
			Set("name", "name"),
			Set("updated_at", now()),
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"updated_at": Set("updated_at", now()),
		}
		queries = From("users").Where(Eq("id", user.ID))
	)

	adapter.On("Update", queries, mutates).Return(0, errors.New("error")).Once()

	assert.NotNil(t, repo.Update(context.TODO(), &user, mutators...))
	assert.Panics(t, func() { repo.MustUpdate(context.TODO(), &user, mutators...) })
	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_update(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{BuyerID: 1, Buyer: User{ID: 1}}
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc,
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"updated_at": now(),
				},
			},
		)
		mutates = map[string]Mutate{
			"name":       Set("name", "buyer1"),
			"age":        Set("age", 20),
			"updated_at": Set("updated_at", now()),
		}
		q = Build("users", Eq("id", 1))
	)

	adapter.On("Update", q, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveBelongsTo(context.TODO(), doc, &mutation))
	assert.Equal(t, Transaction{
		BuyerID: 1,
		Buyer: User{
			ID:        1,
			Name:      "buyer1",
			Age:       20,
			UpdatedAt: now(),
		},
	}, transaction)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_updateError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{BuyerID: 1, Buyer: User{ID: 1}}
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc,
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"updated_at": now(),
				},
			},
		)
		mutates = map[string]Mutate{
			"name":       Set("name", "buyer1"),
			"age":        Set("age", 20),
			"updated_at": Set("updated_at", now()),
		}
		q = Build("users", Eq("id", 1))
	)

	adapter.On("Update", q, mutates).Return(0, errors.New("update error")).Once()

	err := repo.(*repository).saveBelongsTo(context.TODO(), doc, &mutation)
	assert.Equal(t, errors.New("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_updateInconsistentAssoc(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{Buyer: User{ID: 1}}
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc,
			Map{
				"buyer": Map{
					"id":   1,
					"name": "buyer1",
					"age":  20,
				},
			},
		)
	)

	assert.Equal(t, ConstraintError{
		Key:  "user_id",
		Type: ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent belongs to ref and fk"),
	}, repo.(*repository).saveBelongsTo(context.TODO(), doc, &mutation))

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_insertNew(t *testing.T) {
	var (
		transaction Transaction
		adapter     = &testAdapter{}
		repo        = New(adapter)
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc,
			Map{
				"buyer": Map{
					"name": "buyer1",
					"age":  20,
				},
			},
		)
		mutates = map[string]Mutate{
			"name": Set("name", "buyer1"),
			"age":  Set("age", 20),
		}
		q = Build("users")
	)

	adapter.On("Insert", q, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveBelongsTo(context.TODO(), doc, &mutation))
	assert.Equal(t, Set("user_id", 1), mutation.Mutates["user_id"])
	assert.Equal(t, Transaction{
		Buyer: User{
			ID:   1,
			Name: "buyer1",
			Age:  20,
		},
		BuyerID: 1,
	}, transaction)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_insertNewError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{}
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc,
			Map{
				"buyer": Map{
					"name":       "buyer1",
					"age":        20,
					"created_at": now(),
					"updated_at": now(),
				},
			},
		)
		mutates = map[string]Mutate{
			"name":       Set("name", "buyer1"),
			"age":        Set("age", 20),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
		q = Build("users")
	)

	adapter.On("Insert", q, mutates).Return(0, errors.New("insert error")).Once()

	assert.Equal(t, errors.New("insert error"), repo.(*repository).saveBelongsTo(context.TODO(), doc, &mutation))
	assert.Zero(t, mutation.Mutates["user_id"])

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_notChanged(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{}
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc)
	)

	err := repo.(*repository).saveBelongsTo(context.TODO(), doc, &mutation)
	assert.Nil(t, err)
	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_update(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		userID   = 1
		user     = User{ID: userID, Address: Address{ID: 2, UserID: &userID}}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		mutates = map[string]Mutate{
			"street": Set("street", "street1"),
		}
		q = Build("addresses").Where(Eq("id", 2).AndEq("user_id", 1).AndNil("deleted_at"))
	)

	adapter.On("Update", q, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasOne(context.TODO(), doc, &mutation))
	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_updateError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		userID   = 1
		user     = User{ID: userID, Address: Address{ID: 2, UserID: &userID}}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		mutates = map[string]Mutate{
			"street": Set("street", "street1"),
		}
		q = Build("addresses").Where(Eq("id", 2).AndEq("user_id", 1).AndNil("deleted_at"))
	)

	adapter.On("Update", q, mutates).Return(0, errors.New("update error")).Once()

	err := repo.(*repository).saveHasOne(context.TODO(), doc, &mutation)
	assert.Equal(t, errors.New("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_updateInconsistentAssoc(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		user     = User{ID: 1, Address: Address{ID: 2}}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"address": Map{
					"id":     2,
					"street": "street1",
				},
			},
		)
	)

	assert.Equal(t, ConstraintError{
		Key:  "user_id",
		Type: ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent has one ref and fk"),
	}, repo.(*repository).saveHasOne(context.TODO(), doc, &mutation))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_insertNew(t *testing.T) {
	var (
		user     = User{ID: 1}
		adapter  = &testAdapter{}
		repo     = New(adapter)
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		mutates = map[string]Mutate{
			"street":  Set("street", "street1"),
			"user_id": Set("user_id", 1),
		}
		q = Build("addresses")
	)

	adapter.On("Insert", q, mutates).Return(2, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasOne(context.TODO(), doc, &mutation))
	assert.Equal(t, User{
		ID: 1,
		Address: Address{
			ID:     2,
			Street: "street1",
			UserID: &user.ID,
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_insertNewError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"address": Map{
					"street": "street1",
				},
			},
		)
		mutates = map[string]Mutate{
			"street":  Set("street", "street1"),
			"user_id": Set("user_id", 1),
		}
		q = Build("addresses")
	)

	adapter.On("Insert", q, mutates).Return(nil, errors.New("insert error")).Once()

	assert.Equal(t, errors.New("insert error"), repo.(*repository).saveHasOne(context.TODO(), doc, &mutation))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_insert(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"transactions": []Map{
					{"item": "item1"},
					{"item": "item2"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item1")},
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item2")},
		}
		q = Build("transactions")
	)

	adapter.On("InsertAll", q, []string{"item", "user_id"}, mutates).Return(nil).Return([]interface{}{2, 3}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "item"}, mutates).Return(nil).Return([]interface{}{2, 3}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, true))
	assert.Equal(t, User{
		ID: 1,
		Transactions: []Transaction{
			{ID: 2, BuyerID: 1, Item: "item1"},
			{ID: 3, BuyerID: 1, Item: "item2"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_insertError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"transactions": []Map{
					{"item": "item1"},
					{"item": "item2"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item1")},
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item2")},
		}
		q   = Build("transactions")
		err = errors.New("insert all error")
	)

	adapter.On("InsertAll", q, []string{"item", "user_id"}, mutates).Return(nil).Return([]interface{}{}, err).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "item"}, mutates).Return(nil).Return([]interface{}{}, err).Maybe()

	assert.Equal(t, err, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, true))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_update(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 1, BuyerID: 1, Item: "item1"},
				{ID: 2, BuyerID: 1, Item: "item2"},
				{ID: 3, BuyerID: 1, Item: "item3"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"transactions": []Map{
					{"id": 1, "item": "item1 updated"},
					{"id": 2, "item": "item2 updated"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"item": Set("item", "item1 updated")},
			{"item": Set("item", "item2 updated")},
		}
		q = Build("transactions")
	)

	mutation.SetDeletedIDs("transactions", []interface{}{3})

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 3))).Return(1, nil).Once()
	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), mutates[0]).Return(1, nil).Once()
	adapter.On("Update", q.Where(Eq("id", 2).AndEq("user_id", 1)), mutates[1]).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Transactions: []Transaction{
			{ID: 1, BuyerID: 1, Item: "item1 updated"},
			{ID: 2, BuyerID: 1, Item: "item2 updated"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateWithInsert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 1, BuyerID: 1, Item: "item1"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"transactions": []Map{
					{"id": 1, "item": "item1 updated"},
					{"item": "new item", "user_id": 1},
				},
			},
		)
		q       = Build("transactions")
		mutates = []map[string]Mutate{
			{"item": Set("item", "item1 updated")},
			{"user_id": Set("user_id", user.ID), "item": Set("item", "new item")},
		}
	)

	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), mutates[0]).Return(1, nil).Once()
	adapter.On("InsertAll", q, []string{"item", "user_id"}, mutates[1:]).Return(nil).Return([]interface{}{2}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "item"}, mutates[1:]).Return(nil).Return([]interface{}{2}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Transactions: []Transaction{
			{ID: 1, BuyerID: 1, Item: "item1 updated"},
			{ID: 2, BuyerID: 1, Item: "new item"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_deleteWithInsert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 1, Item: "item1"},
				{ID: 2, Item: "item2"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"transactions": []Map{
					{"item": "item3"},
					{"item": "item4"},
					{"item": "item5"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item3")},
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item4")},
			{"user_id": Set("user_id", user.ID), "item": Set("item", "item5")},
		}
		q = Build("transactions")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 1, 2))).Return(1, nil).Once()
	adapter.On("InsertAll", q, []string{"item", "user_id"}, mutates).Return(nil).Return([]interface{}{3, 4, 5}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "item"}, mutates).Return(nil).Return([]interface{}{3, 4, 5}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Transactions: []Transaction{
			{ID: 3, BuyerID: 1, Item: "item3"},
			{ID: 4, BuyerID: 1, Item: "item4"},
			{ID: 5, BuyerID: 1, Item: "item5"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_deleteError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 1, Item: "item1"},
				{ID: 2, Item: "item2"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"transactions": []Map{
					{"item": "item3"},
					{"item": "item4"},
					{"item": "item5"},
				},
			},
		)
		q   = Build("transactions")
		err = errors.New("delete all error")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 1, 2))).Return(0, err).Once()

	assert.Equal(t, err, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_replace(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{Item: "item3"},
				{Item: "item4"},
				{Item: "item5"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc, NewStructset(doc, false))
		mutates  = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "address_id": Set("address_id", 0), "status": Set("status", Status("")), "item": Set("item", "item3")},
			{"user_id": Set("user_id", user.ID), "address_id": Set("address_id", 0), "status": Set("status", Status("")), "item": Set("item", "item4")},
			{"user_id": Set("user_id", user.ID), "address_id": Set("address_id", 0), "status": Set("status", Status("")), "item": Set("item", "item5")},
		}
		q = Build("transactions")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(1, nil).Once()
	adapter.On("InsertAll", q, mock.Anything, mutates).Return(nil).Return([]interface{}{3, 4, 5}, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false))
	assert.Equal(t, User{
		ID:        1,
		CreatedAt: now(),
		UpdatedAt: now(),
		Transactions: []Transaction{
			{ID: 3, BuyerID: 1, Item: "item3"},
			{ID: 4, BuyerID: 1, Item: "item4"},
			{ID: 5, BuyerID: 1, Item: "item5"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_replaceDeleteAllError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Transactions: []Transaction{
				{ID: 1, Item: "item1"},
				{ID: 2, Item: "item2"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc, NewStructset(doc, false))
		q        = Build("transactions")
		err      = errors.New("delete all error")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(0, err).Once()

	assert.Equal(t, err, repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_invalidMutator(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(NewDocument(&User{}),
			Map{
				"transactions": []Map{
					{"item": "item3"},
				},
			},
		)
	)

	assert.PanicsWithValue(t, "rel: invalid mutator", func() {
		repo.(*repository).saveHasMany(context.TODO(), doc, &mutation, false)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_UpdateAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("addresses").Where(Eq("user_id", 1))
		mutates = map[string]Mutate{
			"notes": Set("notes", "notes"),
		}
	)

	adapter.On("Update", query, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.UpdateAll(context.TODO(), query, Set("notes", "notes")))

	adapter.AssertExpectations(t)
}

func TestRepository_MustUpdateAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("addresses").Where(Eq("user_id", 1))
		mutates = map[string]Mutate{
			"notes": Set("notes", "notes"),
		}
	)

	adapter.On("Update", query, mutates).Return(1, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustUpdateAll(context.TODO(), query, Set("notes", "notes"))
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Delete(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{ID: 1}
	)

	adapter.On("Delete", From("users").Where(Eq("id", user.ID))).Return(1, nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &user))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_notFound(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{ID: 1}
	)

	adapter.On("Delete", From("users").Where(Eq("id", user.ID))).Return(0, nil).Once()

	assert.Equal(t, NotFoundError{}, repo.Delete(context.TODO(), &user))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_softDelete(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		address = Address{ID: 1}
		query   = From("addresses").Where(Eq("id", address.ID))
		mutates = map[string]Mutate{
			"deleted_at": Set("deleted_at", now()),
		}
	)

	adapter.On("Update", query, mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &address))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_belongsTo(t *testing.T) {
	var (
		userID  = 1
		address = Address{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		addressMut = map[string]Mutate{
			"deleted_at": Set("deleted_at", now()),
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	// if db constraint enabled, deleting users where it referenced by soft deleted address will raise error.
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", address.ID)), addressMut).Return(1, nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", *address.UserID))).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &address, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_belongsToInconsistentAssoc(t *testing.T) {
	var (
		userID  = 2
		address = Address{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		addressMut = map[string]Mutate{
			"deleted_at": Set("deleted_at", now()),
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", address.ID)), addressMut).Return(1, nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, ConstraintError{
		Key:  "user_id",
		Type: ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent belongs to ref and fk"),
	}, repo.Delete(context.TODO(), &address, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_belongsToError(t *testing.T) {
	var (
		userID  = 1
		address = Address{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		addressMut = map[string]Mutate{
			"deleted_at": Set("deleted_at", now()),
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", *address.UserID))).Return(1, err).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", address.ID)), addressMut).Return(1, nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Delete(context.TODO(), &address, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasOne(t *testing.T) {
	var (
		userID = 10
		user   = User{
			ID: userID,
			Address: Address{
				ID:     1,
				Street: "street",
				UserID: &userID,
			},
		}
		addressMut = map[string]Mutate{
			"deleted_at": Set("deleted_at", now()),
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10)), addressMut).Return(1, nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", 10))).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasOneInconsistentAssoc(t *testing.T) {
	var (
		userID = 10
		user   = User{
			ID: 5,
			Address: Address{
				ID:     1,
				Street: "street",
				UserID: &userID,
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, ConstraintError{
		Key:  "user_id",
		Type: ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent has one ref and fk"),
	}, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasOneError(t *testing.T) {
	var (
		userID = 10
		user   = User{
			ID: userID,
			Address: Address{
				ID:     1,
				Street: "street",
				UserID: &userID,
			},
		}
		addressMut = map[string]Mutate{
			"deleted_at": Set("deleted_at", now()),
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10)), addressMut).Return(1, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasMany(t *testing.T) {
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
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Delete", From("transactions").Where(Eq("user_id", 10).AndIn("id", 1))).Return(1, nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", 10))).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasManyError(t *testing.T) {
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
		repo    = New(adapter)
		err     = errors.New("err")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Delete", From("transactions").Where(Eq("user_id", 10).AndIn("id", 1))).Return(1, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_MustDelete(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{ID: 1}
	)

	adapter.On("Delete", From("users").Where(Eq("id", user.ID))).Return(1, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustDelete(context.TODO(), &user)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_DeleteAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		queries = From("logs").Where(Eq("user_id", 1))
	)

	adapter.On("Delete", From("logs").Where(Eq("user_id", 1))).Return(1, nil).Once()

	assert.Nil(t, repo.DeleteAll(context.TODO(), queries))

	adapter.AssertExpectations(t)
}

func TestRepository_MustDeleteAll(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		queries = From("logs").Where(Eq("user_id", 1))
	)

	adapter.On("Delete", From("logs").Where(Eq("user_id", 1))).Return(1, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustDeleteAll(context.TODO(), queries)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Preload_hasOne(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{ID: 10}
		address = Address{ID: 100, UserID: &user.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", From("addresses").Where(In("user_id", 10).AndNil("deleted_at"))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(address.ID, *address.UserID).Times(2)
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(context.TODO(), &user, "address"))
	assert.Equal(t, address, user.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceHasOne(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = New(adapter)
		users     = []User{{ID: 10}, {ID: 20}}
		addresses = []Address{
			{ID: 100, UserID: &users[0].ID},
			{ID: 200, UserID: &users[1].ID},
		}
		cur = &testCursor{}
	)

	// one of these, because of map ordering
	adapter.On("Query", From("addresses").Where(In("user_id", 10, 20).AndNil("deleted_at"))).Return(cur, nil).Maybe()
	adapter.On("Query", From("addresses").Where(In("user_id", 20, 10).AndNil("deleted_at"))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(addresses[0].ID, *addresses[0].UserID).Twice()
	cur.MockScan(addresses[1].ID, *addresses[1].UserID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(context.TODO(), &users, "address"))
	assert.Equal(t, addresses[0], users[0].Address)
	assert.Equal(t, addresses[1], users[1].Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedHasOne(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{
			Buyer: User{ID: 10},
		}
		address = Address{ID: 100, UserID: &transaction.Buyer.ID}
		cur     = &testCursor{}
	)

	adapter.On("Query", From("addresses").Where(In("user_id", 10).AndNil("deleted_at"))).Return(cur, nil).Once()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(address.ID, *address.UserID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(context.TODO(), &transaction, "buyer.address"))
	assert.Equal(t, address, transaction.Buyer.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceNestedHasOne(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = New(adapter)
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
	adapter.On("Query", From("addresses").Where(In("user_id", 10, 20).AndNil("deleted_at"))).Return(cur, nil).Maybe()
	adapter.On("Query", From("addresses").Where(In("user_id", 20, 10).AndNil("deleted_at"))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Twice()
	cur.MockScan(addresses[0].ID, *addresses[0].UserID).Twice()
	cur.MockScan(addresses[1].ID, *addresses[1].UserID).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(context.TODO(), &transactions, "buyer.address"))
	assert.Equal(t, addresses[0], transactions[0].Buyer.Address)
	assert.Equal(t, addresses[1], transactions[1].Buyer.Address)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_hasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &user, "transactions"))
	assert.Equal(t, transactions, user.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceHasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &users, "transactions"))
	assert.Equal(t, transactions[:2], users[0].Transactions)
	assert.Equal(t, transactions[2:], users[1].Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedHasMany(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &address, "user.transactions"))
	assert.Equal(t, transactions, address.User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedNullHasMany(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		address = Address{User: nil}
	)

	assert.Nil(t, repo.Preload(context.TODO(), &address, "user.transactions"))

	adapter.AssertExpectations(t)
}

func TestRepository_Preload_nestedSliceHasMany(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &addresses, "user.transactions"))
	assert.Equal(t, transactions[:2], addresses[0].User.Transactions)
	assert.Equal(t, transactions[2:], addresses[1].User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nestedNullSliceHasMany(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &addresses, "user.transactions"))
	assert.Equal(t, transactions[:2], addresses[0].User.Transactions)
	assert.Equal(t, []Transaction(nil), addresses[1].User.Transactions)
	assert.Equal(t, transactions[2:], addresses[2].User.Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_belongsTo(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &transaction, "buyer"))
	assert.Equal(t, user, transaction.Buyer)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_ptrBelongsTo(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &address, "user"))
	assert.Equal(t, user, *address.User)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_nullBelongsTo(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		address = Address{}
	)

	assert.Nil(t, repo.Preload(context.TODO(), &address, "user"))
	assert.Nil(t, address.User)

	adapter.AssertExpectations(t)
}

func TestRepository_Preload_sliceBelongsTo(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &transactions, "buyer"))
	assert.Equal(t, users[0], transactions[0].Buyer)
	assert.Equal(t, users[1], transactions[1].Buyer)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_ptrSliceBelongsTo(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
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

	assert.Nil(t, repo.Preload(context.TODO(), &addresses, "user"))
	assert.Equal(t, users[0], *addresses[0].User)
	assert.Equal(t, users[1], *addresses[1].User)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_sliceNestedBelongsTo(t *testing.T) {
	var (
		adapter      = &testAdapter{}
		repo         = New(adapter)
		address      = Address{ID: 10, Street: "Continassa"}
		transactions = []Transaction{
			{AddressID: 10},
		}
		users = []User{
			{ID: 10, Name: "Del Piero", Transactions: transactions},
			{ID: 20, Name: "Nedved"},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("addresses").Where(In("id", 10).AndNil("deleted_at"))).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "street"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(address.ID, address.Street).Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(context.TODO(), &users, "transactions.address"))
	assert.Equal(t, []Transaction{
		{AddressID: 10, Address: address},
	}, users[0].Transactions)
	assert.Equal(t, []Transaction{}, users[1].Transactions)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Preload_emptySlice(t *testing.T) {
	var (
		adapter   = &testAdapter{}
		repo      = New(adapter)
		addresses = []Address{}
	)

	assert.Nil(t, repo.Preload(context.TODO(), &addresses, "user.transactions"))
}

func TestQuery_Preload_notPointerPanic(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{}
	)

	assert.Panics(t, func() { repo.Preload(context.TODO(), transaction, "User") })
}

func TestRepository_Preload_queryError(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{BuyerID: 10}
		cur         = &testCursor{}
		err         = errors.New("error")
	)

	adapter.On("Query", From("users").Where(In("id", 10))).Return(cur, err).Once()

	assert.Equal(t, err, repo.Preload(context.TODO(), &transaction, "buyer"))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_MustPreload(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		repo        = New(adapter)
		transaction = Transaction{BuyerID: 10}
		cur         = createCursor(0)
	)

	adapter.On("Query", From("users").Where(In("id", 10))).Return(cur, nil).Once()

	assert.NotPanics(t, func() {
		repo.MustPreload(context.TODO(), &transaction, "buyer")
	})

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Transaction(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).On("Commit").Return(nil).Once()

	repo := New(adapter)

	err := repo.Transaction(context.TODO(), func(repo Repository) error {
		assert.True(t, repo.(*repository).inTransaction)
		return nil
	})

	assert.False(t, repo.(*repository).inTransaction)
	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_beginError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(errors.New("error")).Once()

	err := New(adapter).Transaction(context.TODO(), func(r Repository) error {
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

	err := New(adapter).Transaction(context.TODO(), func(r Repository) error {
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

	err := New(adapter).Transaction(context.TODO(), func(r Repository) error {
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

	err := New(adapter).Transaction(context.TODO(), func(r Repository) error {
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
		_ = New(adapter).Transaction(context.TODO(), func(r Repository) error {
			// doing good things
			panic("error")
		})
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_runtimeError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	var user *User
	assert.Panics(t, func() {
		_ = New(adapter).Transaction(context.TODO(), func(r Repository) error {
			_ = user.ID
			return nil
		})
	})

	adapter.AssertExpectations(t)
}
