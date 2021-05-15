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
	t := now().Truncate(time.Second)
	now = func() time.Time {
		return t
	}
}

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
		ctx     = context.TODO()
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	assert.NotNil(t, repo)
	assert.Equal(t, adapter, repo.Adapter(ctx))
}

func TestRepository_Instrumentation(t *testing.T) {
	var (
		repo = repository{rootAdapter: &testAdapter{}}
	)

	assert.Nil(t, repo.instrumenter)
	assert.NotPanics(t, func() {
		repo.instrumenter.Observe(context.TODO(), "test", "test")(nil)
	})

	repo.Instrumentation(DefaultLogger)
	assert.NotNil(t, repo.instrumenter)
	assert.NotPanics(t, func() {
		repo.instrumenter.Observe(context.TODO(), "test", "test")(nil)
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

		assert.NotZero(t, user.ID)
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

func TestRepository_Find_withCascade(t *testing.T) {
	var (
		trx        Transaction
		adapter    = &testAdapter{}
		repo       = New(adapter)
		query      = From("transactions").Limit(1).Cascade(true)
		cur        = createCursor(1)
		curPreload = createCursor(0)
	)

	adapter.On("Query", query.Preload("buyer")).Return(cur, nil).Once()
	adapter.On("Query", From("users").Where(In("id", 0))).Return(curPreload, nil).Once()

	assert.Nil(t, repo.Find(context.TODO(), &trx, query))
	assert.Equal(t, 10, trx.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_withPreload(t *testing.T) {
	var (
		user       User
		adapter    = &testAdapter{}
		repo       = New(adapter)
		query      = From("users").Limit(1).Preload("address")
		cur        = createCursor(1)
		curPreload = &testCursor{}
	)

	adapter.On("Query", query).Return(cur, nil).Once()
	adapter.On("Query", From("addresses").Where(In("user_id", 10).AndNil("deleted_at"))).
		Return(curPreload, nil).Once()

	curPreload.On("Close").Return(nil).Once()
	curPreload.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	curPreload.On("Next").Return(false).Once()

	assert.Nil(t, repo.Find(context.TODO(), &user, query))
	assert.Equal(t, 10, user.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
	curPreload.AssertExpectations(t)
}

func TestRepository_Find_withPreloadError(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(1).Preload("address")
		cur     = createCursor(1)
		err     = errors.New("error")
	)

	adapter.On("Query", query).Return(cur, nil).Once()
	adapter.On("Query", From("addresses").Where(In("user_id", 10).AndNil("deleted_at"))).
		Return(&testCursor{}, err).Once()

	assert.Equal(t, err, repo.Find(context.TODO(), &user, query))
	assert.Equal(t, 10, user.ID)
	assert.False(t, cur.Next())

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_scanError(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		cur     = &testCursor{}
		query   = From("users").Limit(1)
		err     = errors.New("error")
	)

	adapter.On("Query", query).Return(cur, nil).Once()

	cur.On("Fields").Return([]string(nil), err).Once()
	cur.On("Close").Return(nil).Once()

	assert.Equal(t, err, repo.Find(context.TODO(), &user, query))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Find_error(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		repo    = New(adapter)
		cur     = &testCursor{}
		query   = From("users").Limit(1)
		err     = errors.New("error")
	)

	adapter.On("Query", query).Return(cur, err).Once()

	assert.Equal(t, err, repo.Find(context.TODO(), &user, query))

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

func TestRepository_FindAll_withCascade(t *testing.T) {
	var (
		trxs       []Transaction
		adapter    = &testAdapter{}
		repo       = New(adapter)
		query      = From("transactions")
		cur        = createCursor(2)
		curPreload = createCursor(0)
	)

	adapter.On("Query", query.Preload("buyer")).Return(cur, nil).Once()
	adapter.On("Query", From("users").Where(In("id", 0))).Return(curPreload, nil)

	assert.Nil(t, repo.FindAll(context.TODO(), &trxs, query))
	assert.Len(t, trxs, 2)
	assert.Equal(t, 10, trxs[0].ID)
	assert.Equal(t, 10, trxs[1].ID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_FindAll_withPreload(t *testing.T) {
	var (
		addresses  []Address
		adapter    = &testAdapter{}
		repo       = New(adapter)
		query      = From("addresses").Preload("user")
		cur        = &testCursor{}
		curPreload = createCursor(0)
	)

	adapter.On("Query", query.Where(Nil("deleted_at"))).Return(cur, nil).Once()
	adapter.On("Query", From("users").Where(In("id", 20))).
		Return(curPreload, nil).Once()

	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(10, 20)
	cur.On("Next").Return(false).Once()
	cur.On("Close").Return(nil).Once()

	assert.Nil(t, repo.FindAll(context.TODO(), &addresses, query))
	assert.Len(t, addresses, 1)
	assert.Equal(t, 10, addresses[0].ID)
	assert.Equal(t, 20, *addresses[0].UserID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
	curPreload.AssertExpectations(t)
}

func TestRepository_FindAll_withPreloadError(t *testing.T) {
	var (
		addresses  []Address
		adapter    = &testAdapter{}
		repo       = New(adapter)
		query      = From("addresses").Preload("user")
		cur        = &testCursor{}
		curPreload = &testCursor{}
		err        = errors.New("error")
	)

	adapter.On("Query", query.Where(Nil("deleted_at"))).Return(cur, nil).Once()
	adapter.On("Query", From("users").Where(In("id", 20))).
		Return(curPreload, err).Once()

	cur.On("Fields").Return([]string{"id", "user_id"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(10, 20)
	cur.On("Next").Return(false).Once()
	cur.On("Close").Return(nil).Once()

	assert.Equal(t, err, repo.FindAll(context.TODO(), &addresses, query))
	assert.Len(t, addresses, 1)
	assert.Equal(t, 10, addresses[0].ID)
	assert.Equal(t, 20, *addresses[0].UserID)

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
	curPreload.AssertExpectations(t)
}

func TestRepository_FindAll_scanError(t *testing.T) {
	var (
		users   []User
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = From("users").Limit(1)
		cur     = &testCursor{}
		err     = errors.New("error")
	)

	cur.On("Fields").Return([]string(nil), err).Once()
	cur.On("Close").Return(nil).Once()

	adapter.On("Query", query).Return(cur, nil).Once()

	assert.Equal(t, err, repo.FindAll(context.TODO(), &users, query))

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

func TestRepository_FindAndCountAll_error(t *testing.T) {
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
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			Name: "name",
		}
		mutates = map[string]Mutate{
			"name":       Set("name", "name"),
			"age":        Set("age", 0),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
	)

	adapter.On("Insert", From("users"), mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_compositePrimaryFields(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		userRole = UserRole{
			UserID: 1,
			RoleID: 2,
		}
		mutates = map[string]Mutate{
			"user_id": Set("user_id", 1),
			"role_id": Set("role_id", 2),
		}
	)

	adapter.On("Insert", From("user_roles"), mutates).Return(0, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &userRole))
	assert.Equal(t, UserRole{
		UserID: 1,
		RoleID: 2,
	}, userRole)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_sets(t *testing.T) {
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

func TestRepository_Insert_saveBelongsTo(t *testing.T) {
	var (
		userID    = 1
		profileID = 2
		profile   = Profile{
			Name: "profile name",
			User: &User{Name: "name"},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(userID, nil).Once()
	adapter.On("Insert", From("profiles"), mock.Anything).Return(profileID, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &profile))
	assert.Equal(t, Profile{
		ID:     profileID,
		Name:   "profile name",
		UserID: &userID,
		User: &User{
			ID:        userID,
			Name:      "name",
			CreatedAt: now(),
			UpdatedAt: now(),
		},
	}, profile)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveBelongsToCascadeDisabled(t *testing.T) {
	var (
		profile = Profile{
			Name: "profile name",
			User: &User{Name: "name"},
		}
		adapter   = &testAdapter{}
		repo      = New(adapter)
		addressID = 2
	)

	adapter.On("Insert", From("profiles"), mock.Anything).Return(addressID, nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &profile, Cascade(false)))
	assert.Equal(t, Profile{
		ID:   addressID,
		Name: "profile name",
		User: &User{Name: "name"},
	}, profile)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveBelongsToError(t *testing.T) {
	var (
		profile = Profile{
			Name: "profile name",
			User: &User{Name: "name"},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Insert(context.TODO(), &profile))

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
			UserRoles: []UserRole{
				{RoleID: 2},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("InsertAll", From("user_roles"), mock.Anything, mock.Anything).Return([]interface{}(nil), nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Insert(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
		UserRoles: []UserRole{
			{UserID: 1, RoleID: 2},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasManyCascadeDisabled(t *testing.T) {
	var (
		user = User{
			Name: "name",
			UserRoles: []UserRole{
				{RoleID: 2},
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
		UserRoles: []UserRole{
			{RoleID: 2},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Insert_saveHasManyError(t *testing.T) {
	var (
		user = User{
			Name: "name",
			UserRoles: []UserRole{
				{RoleID: 2},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Insert", From("users"), mock.Anything).Return(1, nil).Once()
	adapter.On("InsertAll", From("user_roles"), mock.Anything, mock.Anything).Return([]interface{}{}, err).Once()
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
		err = errors.New("error")
	)

	adapter.On("Insert", From("users"), mutates).Return(0, err).Once()

	assert.Equal(t, err, repo.Insert(context.TODO(), &user, mutators...))
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
		profile = Profile{
			Name: "profile name",
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

	assert.Equal(t, errors.New("error"), repo.Insert(context.TODO(), &profile,
		NewStructset(&profile, false),
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

func TestRepository_InsertAll_compositePrimaryFields(t *testing.T) {
	var (
		userRoles = []UserRole{
			{UserID: 1, RoleID: 2},
			{UserID: 1, RoleID: 3},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		mutates = []map[string]Mutate{
			{
				"user_id": Set("user_id", 1),
				"role_id": Set("role_id", 2),
			},
			{
				"user_id": Set("user_id", 1),
				"role_id": Set("role_id", 3),
			},
		}
	)

	adapter.On("InsertAll", From("user_roles"), mock.Anything, mutates).Return([]interface{}{0, 0}, nil).Once()

	assert.Nil(t, repo.InsertAll(context.TODO(), &userRoles))
	assert.Equal(t, []UserRole{
		{UserID: 1, RoleID: 2},
		{UserID: 1, RoleID: 3},
	}, userRoles)

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
		adapter = &testAdapter{}
		repo    = New(adapter)
		user    = User{
			ID:        1,
			Name:      "name",
			CreatedAt: now(),
			UpdatedAt: now(),
		}
		mutates = map[string]Mutate{
			"id":         Set("id", 1),
			"name":       Set("name", "name"),
			"age":        Set("age", 0),
			"created_at": Set("created_at", now()),
			"updated_at": Set("updated_at", now()),
		}
		queries = From("users").Where(Eq("id", user.ID))
	)

	adapter.On("Update", queries, "id", mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        1,
		Name:      "name",
		CreatedAt: now(),
		UpdatedAt: now(),
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_compositePrimaryKeys(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		userRole = UserRole{
			UserID: 1,
			RoleID: 2,
		}
		mutates = map[string]Mutate{
			"user_id": Set("user_id", 1),
			"role_id": Set("role_id", 2),
		}
		queries = From("user_roles").Where(Eq("user_id", userRole.UserID).AndEq("role_id", userRole.RoleID))
	)

	adapter.On("Update", queries, "", mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &userRole))
	assert.Equal(t, UserRole{
		UserID: 1,
		RoleID: 2,
	}, userRole)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_sets(t *testing.T) {
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

	adapter.On("Update", queries, "id", mutates).Return(1, nil).Once()

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

	adapter.On("Update", queries.Where(Nil("deleted_at")), "id", mutates).Return(1, nil).Once()

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

	adapter.On("Update", queries, "id", mutates).Return(1, nil).Once()

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

	adapter.On("Update", queries, "id", mutates).Return(0, nil).Once()

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

	adapter.On("Update", queries, "id", mutates).Return(1, nil).Once()
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

	adapter.On("Update", queries, "id", mutates).Return(1, nil).Once()
	adapter.On("Query", queries.Limit(1)).Return(cur, err).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &user, mutators...))

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsTo(t *testing.T) {
	var (
		userID  = 1
		profile = Profile{
			ID:     1,
			Name:   "profile name",
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
	adapter.On("Update", From("users").Where(Eq("id", *profile.UserID)), "id", mock.Anything).Return(1, nil).Once()
	adapter.On("Update", From("profiles").Where(Eq("id", profile.ID)), "id", mock.Anything).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &profile))
	assert.Equal(t, Profile{
		ID:     1,
		Name:   "profile name",
		UserID: &userID,
		User: &User{
			ID:        1,
			Name:      "name",
			UpdatedAt: now(),
			CreatedAt: now(),
		},
	}, profile)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsToCascadeDisabled(t *testing.T) {
	var (
		userID  = 1
		profile = Profile{
			ID:     1,
			Name:   "profile name",
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Update", From("profiles").Where(Eq("id", profile.ID)).Cascade(false), "id", mock.Anything).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &profile, Cascade(false)))
	assert.Equal(t, Profile{
		ID:     1,
		Name:   "profile name",
		UserID: &userID,
		User: &User{
			ID:   1,
			Name: "name",
		},
	}, profile)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveBelongsToError(t *testing.T) {
	var (
		userID  = 1
		profile = Profile{
			ID:     1,
			Name:   "profile name",
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		queries = From("users").Where(Eq("id", profile.ID))
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", queries, "id", mock.Anything).Return(0, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &profile))

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
	adapter.On("Update", From("users").Where(Eq("id", 10)), "id", mock.Anything).Return(1, nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10).AndNil("deleted_at")), "id", mock.Anything).Return(1, nil).Once()
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

	adapter.On("Update", From("users").Where(Eq("id", 10)).Cascade(false), "id", mock.Anything).Return(1, nil).Once()

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
	adapter.On("Update", From("users").Where(Eq("id", 10)), "id", mock.Anything).Return(1, nil).Once()
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10).AndNil("deleted_at")), "id", mock.Anything).Return(1, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Update(context.TODO(), &user))
	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasMany(t *testing.T) {
	var (
		user = User{
			ID: 10,
			UserRoles: []UserRole{
				{RoleID: 2},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), "id", mock.Anything).Return(1, nil).Once()
	adapter.On("Delete", From("user_roles").Where(Eq("user_id", 10))).Return(1, nil).Once()
	adapter.On("InsertAll", From("user_roles"), mock.Anything, mock.Anything).Return([]interface{}(nil), nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user))
	assert.Equal(t, User{
		ID:        10,
		CreatedAt: now(),
		UpdatedAt: now(),
		UserRoles: []UserRole{
			{UserID: 10, RoleID: 2},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasManyCascadeDisabled(t *testing.T) {
	var (
		user = User{
			ID: 10,
			UserRoles: []UserRole{
				{RoleID: 2},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Update", From("users").Where(Eq("id", 10)).Cascade(false), "id", mock.Anything).Return(1, nil).Once()

	assert.Nil(t, repo.Update(context.TODO(), &user, Cascade(false)))
	assert.Equal(t, User{
		ID:        10,
		CreatedAt: now(),
		UpdatedAt: now(),
		UserRoles: []UserRole{
			{RoleID: 2},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_Update_saveHasManyError(t *testing.T) {
	var (
		user = User{
			ID: 10,
			UserRoles: []UserRole{
				{RoleID: 2},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Update", From("users").Where(Eq("id", 10)), "id", mock.Anything).Return(1, nil).Once()
	adapter.On("Delete", From("user_roles").Where(Eq("user_id", 10))).Return(0, err).Once()
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

	adapter.On("Update", queries, "id", mutates).Return(0, errors.New("error")).Once()

	assert.NotNil(t, repo.Update(context.TODO(), &user, mutators...))
	assert.Panics(t, func() { repo.MustUpdate(context.TODO(), &user, mutators...) })
	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_update(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		userID   = 1
		profile  = Profile{UserID: &userID, User: &User{ID: userID}}
		doc      = NewDocument(&profile)
		mutation = Apply(doc,
			Map{
				"user": Map{
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

	adapter.On("Update", q, "id", mutates).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveBelongsTo(cw, doc, &mutation))
	assert.Equal(t, Profile{
		UserID: &userID,
		User: &User{
			ID:        userID,
			Name:      "buyer1",
			Age:       20,
			UpdatedAt: now(),
		},
	}, profile)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_updateError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		userID   = 1
		profile  = Profile{UserID: &userID, User: &User{ID: userID}}
		doc      = NewDocument(&profile)
		mutation = Apply(doc,
			Map{
				"user": Map{
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

	adapter.On("Update", q, "id", mutates).Return(0, errors.New("update error")).Once()

	err := repo.(*repository).saveBelongsTo(cw, doc, &mutation)
	assert.Equal(t, errors.New("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_updateInconsistentAssoc(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		profile  = Profile{User: &User{ID: 1}}
		doc      = NewDocument(&profile)
		mutation = Apply(doc,
			Map{
				"user": Map{
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
	}, repo.(*repository).saveBelongsTo(cw, doc, &mutation))

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_insertNew(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		userID   = 1
		profile  = Profile{}
		doc      = NewDocument(&profile)
		mutation = Apply(doc,
			Map{
				"user": Map{
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

	assert.Nil(t, repo.(*repository).saveBelongsTo(cw, doc, &mutation))
	assert.Equal(t, Set("user_id", 1), mutation.Mutates["user_id"])
	assert.Equal(t, Profile{
		User: &User{
			ID:   1,
			Name: "buyer1",
			Age:  20,
		},
		UserID: &userID,
	}, profile)

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_insertNewError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		profile  = Profile{}
		doc      = NewDocument(&profile)
		mutation = Apply(doc,
			Map{
				"user": Map{
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

	assert.Equal(t, errors.New("insert error"), repo.(*repository).saveBelongsTo(cw, doc, &mutation))
	assert.Zero(t, mutation.Mutates["user_id"])

	adapter.AssertExpectations(t)
}

func TestRepository_saveBelongsTo_notChanged(t *testing.T) {
	var (
		adapter     = &testAdapter{}
		cw          = fetchContext(context.TODO(), adapter)
		repo        = New(adapter)
		transaction = Transaction{}
		doc         = NewDocument(&transaction)
		mutation    = Apply(doc)
	)

	err := repo.(*repository).saveBelongsTo(cw, doc, &mutation)
	assert.Nil(t, err)
	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_update(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
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

	adapter.On("Update", q, "id", mutates).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasOne(cw, doc, &mutation))
	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_updateError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
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

	adapter.On("Update", q, "id", mutates).Return(0, errors.New("update error")).Once()

	err := repo.(*repository).saveHasOne(cw, doc, &mutation)
	assert.Equal(t, errors.New("update error"), err)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_updateInconsistentAssoc(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
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
	}, repo.(*repository).saveHasOne(cw, doc, &mutation))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasOne_insertNew(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
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

	adapter.On("Insert", q, mutates).Return(2, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasOne(cw, doc, &mutation))
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
		cw       = fetchContext(context.TODO(), adapter)
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

	assert.Equal(t, errors.New("insert error"), repo.(*repository).saveHasOne(cw, doc, &mutation))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_insert(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"email": "email1@gmail.com"},
					{"email": "email2@gmail.com"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email1@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email2@gmail.com")},
		}
		q = Build("emails")
	)

	adapter.On("InsertAll", q, []string{"email", "user_id"}, mutates).Return([]interface{}{2, 3}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "email"}, mutates).Return([]interface{}{2, 3}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(cw, doc, &mutation, true))
	assert.Equal(t, User{
		ID: 1,
		Emails: []Email{
			{ID: 2, Email: "email1@gmail.com", UserID: 1},
			{ID: 3, Email: "email2@gmail.com", UserID: 1},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_insertError(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"email": "email1@gmail.com"},
					{"email": "email2@gmail.com"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email1@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email2@gmail.com")},
		}
		q   = Build("emails")
		err = errors.New("insert all error")
	)

	adapter.On("InsertAll", q, []string{"email", "user_id"}, mutates).Return([]interface{}{}, err).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "email"}, mutates).Return([]interface{}{}, err).Maybe()

	assert.Equal(t, err, repo.(*repository).saveHasMany(cw, doc, &mutation, true))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_update(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, UserID: 1, Email: "email1@gmail.com"},
				{ID: 2, UserID: 1, Email: "email2@gmail.com"},
				{ID: 3, UserID: 1, Email: "email3@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"id": 1, "email": "email1updated@gmail.com"},
					{"id": 3, "email": "email3updated@gmail.com"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"email": Set("email", "email1updated@gmail.com")},
			{"email": Set("email", "email3updated@gmail.com")},
		}
		q = Build("emails")
	)

	mutation.SetDeletedIDs("emails", []interface{}{2})

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 2))).Return(1, nil).Once()
	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), "id", mutates[0]).Return(1, nil).Once()
	adapter.On("Update", q.Where(Eq("id", 3).AndEq("user_id", 1)), "id", mutates[1]).Return(1, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasMany(cw, doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Emails: []Email{
			{ID: 1, UserID: 1, Email: "email1updated@gmail.com"},
			{ID: 3, UserID: 1, Email: "email3updated@gmail.com"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateInconsistentReferences(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, UserID: 2, Email: "email1@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"id": 1, "email": "email1updated@gmail.com"},
				},
			},
		)
	)

	mutation.SetDeletedIDs("emails", []interface{}{})

	assert.Equal(t, ConstraintError{
		Key:  "user_id",
		Type: ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent has many ref and fk"),
	}, repo.(*repository).saveHasMany(cw, doc, &mutation, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, UserID: 1, Email: "email1@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"id": 1, "email": "email1updated@gmail.com"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"email": Set("email", "email1updated@gmail.com")},
		}
		q   = Build("emails")
		err = errors.New("update error")
	)

	mutation.SetDeletedIDs("emails", []interface{}{})

	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), "id", mutates[0]).Return(0, err).Once()

	assert.Equal(t, err, repo.(*repository).saveHasMany(cw, doc, &mutation, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateWithInsert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, UserID: 1, Email: "email1@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"email": "newemail@gmail.com", "user_id": 1},
					{"id": 1, "email": "email1updated@gmail.com"},
				},
			},
		)
		q       = Build("emails")
		mutates = []map[string]Mutate{
			{"email": Set("email", "email1updated@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "newemail@gmail.com")},
		}
	)

	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), "id", mutates[0]).Return(1, nil).Once()
	adapter.On("InsertAll", q, []string{"email", "user_id"}, mutates[1:]).Return([]interface{}{2}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "email"}, mutates[1:]).Return([]interface{}{2}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(cw, doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Emails: []Email{
			{ID: 1, UserID: 1, Email: "email1updated@gmail.com"},
			{ID: 2, UserID: 1, Email: "newemail@gmail.com"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_updateWithReorderInsert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{Email: "email1@gmail.com"}, // new record not appended, but prepended/inserted
				{ID: 1, UserID: 1, Email: "email2@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Mutation{}
		q        = Build("emails")
		mutates  = []map[string]Mutate{
			{"email": Set("email", "update@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "new@gmail.com")},
		}
	)

	// insert first, so internally rel needs to reorder the assoc.
	mutation.SetAssoc("emails",
		Apply(NewDocument(&user.Emails[0]), Set("email", "new@gmail.com")),
		Apply(NewDocument(&user.Emails[1]), Set("email", "update@gmail.com")),
	)
	mutation.SetDeletedIDs("emails", []interface{}{})

	adapter.On("Update", q.Where(Eq("id", 1).AndEq("user_id", 1)), "id", mutates[0]).Return(1, nil).Once()
	adapter.On("InsertAll", q, []string{"email", "user_id"}, mutates[1:]).Return([]interface{}{2}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "email"}, mutates[1:]).Return([]interface{}{2}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(cw, doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Emails: []Email{
			{ID: 1, UserID: 1, Email: "update@gmail.com"},
			{ID: 2, UserID: 1, Email: "new@gmail.com"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_deleteWithInsert(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, Email: "email1@gmail.com"},
				{ID: 2, Email: "email2@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"email": "email3@gmail.com"},
					{"email": "email4@gmail.com"},
					{"email": "email5@gmail.com"},
				},
			},
		)
		mutates = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email3@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email4@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email5@gmail.com")},
		}
		q = Build("emails")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 1, 2))).Return(1, nil).Once()
	adapter.On("InsertAll", q, []string{"email", "user_id"}, mutates).Return([]interface{}{3, 4, 5}, nil).Maybe()
	adapter.On("InsertAll", q, []string{"user_id", "email"}, mutates).Return([]interface{}{3, 4, 5}, nil).Maybe()

	assert.Nil(t, repo.(*repository).saveHasMany(cw, doc, &mutation, false))
	assert.Equal(t, User{
		ID: 1,
		Emails: []Email{
			{ID: 3, UserID: 1, Email: "email3@gmail.com"},
			{ID: 4, UserID: 1, Email: "email4@gmail.com"},
			{ID: 5, UserID: 1, Email: "email5@gmail.com"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_deleteError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, Email: "email1@gmail.com"},
				{ID: 2, Email: "email2@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc,
			Map{
				"emails": []Map{
					{"email": "email3@gmail.com"},
					{"email": "email4@gmail.com"},
					{"email": "email5@gmail.com"},
				},
			},
		)
		q   = Build("emails")
		err = errors.New("delete all error")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1).AndIn("id", 1, 2))).Return(0, err).Once()

	assert.Equal(t, err, repo.(*repository).saveHasMany(cw, doc, &mutation, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_replace(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{Email: "email3@gmail.com"},
				{Email: "email4@gmail.com"},
				{Email: "email5@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc, NewStructset(doc, false))
		mutates  = []map[string]Mutate{
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email3@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email4@gmail.com")},
			{"user_id": Set("user_id", user.ID), "email": Set("email", "email5@gmail.com")},
		}
		q = Build("emails")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(1, nil).Once()
	adapter.On("InsertAll", q, mock.Anything, mutates).Return([]interface{}{3, 4, 5}, nil).Once()

	assert.Nil(t, repo.(*repository).saveHasMany(cw, doc, &mutation, false))
	assert.Equal(t, User{
		ID:        1,
		CreatedAt: now(),
		UpdatedAt: now(),
		Emails: []Email{
			{ID: 3, UserID: 1, Email: "email3@gmail.com"},
			{ID: 4, UserID: 1, Email: "email4@gmail.com"},
			{ID: 5, UserID: 1, Email: "email5@gmail.com"},
		},
	}, user)

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_replaceDeleteAllError(t *testing.T) {
	var (
		adapter = &testAdapter{}
		cw      = fetchContext(context.TODO(), adapter)
		repo    = New(adapter)
		user    = User{
			ID: 1,
			Emails: []Email{
				{ID: 1, Email: "email1@gmail.com"},
				{ID: 2, Email: "email2@gmail.com"},
			},
		}
		doc      = NewDocument(&user)
		mutation = Apply(doc, NewStructset(doc, false))
		q        = Build("emails")
		err      = errors.New("delete all error")
	)

	adapter.On("Delete", q.Where(Eq("user_id", 1))).Return(0, err).Once()

	assert.Equal(t, err, repo.(*repository).saveHasMany(cw, doc, &mutation, false))

	adapter.AssertExpectations(t)
}

func TestRepository_saveHasMany_invalidMutator(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		cw       = fetchContext(context.TODO(), adapter)
		repo     = New(adapter)
		user     = User{ID: 1}
		doc      = NewDocument(&user)
		mutation = Apply(NewDocument(&User{}),
			Map{
				"emails": []Map{
					{"email": "email3@gmail.com"},
				},
			},
		)
	)

	assert.PanicsWithValue(t, "rel: invalid mutator", func() {
		repo.(*repository).saveHasMany(cw, doc, &mutation, false)
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

	adapter.On("Update", query, "", mutates).Return(1, nil).Once()

	updatedCount, err := repo.UpdateAll(context.TODO(), query, Set("notes", "notes"))
	assert.Nil(t, err)
	assert.Equal(t, 1, updatedCount)

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

	adapter.On("Update", query, "", mutates).Return(1, nil).Once()

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

func TestRepository_Delete_compositePrimaryKey(t *testing.T) {
	var (
		adapter  = &testAdapter{}
		repo     = New(adapter)
		userRole = UserRole{UserID: 1, RoleID: 2}
	)

	adapter.On("Delete", From("user_roles").Where(Eq("user_id", userRole.UserID).AndEq("role_id", userRole.RoleID))).Return(1, nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &userRole))

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

	adapter.On("Update", query, "", mutates).Return(1, nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &address))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_belongsTo(t *testing.T) {
	var (
		userID  = 1
		profile = Profile{
			ID:     2,
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
	adapter.On("Delete", From("profiles").Where(Eq("id", profile.ID))).Return(1, nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", *profile.UserID))).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &profile, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_belongsToInconsistentAssoc(t *testing.T) {
	var (
		userID  = 2
		profile = Profile{
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
	adapter.On("Delete", From("profiles").Where(Eq("id", profile.ID))).Return(1, nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, ConstraintError{
		Key:  "user_id",
		Type: ForeignKeyConstraint,
		Err:  errors.New("rel: inconsistent belongs to ref and fk"),
	}, repo.Delete(context.TODO(), &profile, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_belongsToError(t *testing.T) {
	var (
		userID  = 1
		profile = Profile{
			ID:     1,
			UserID: &userID,
			User: &User{
				ID:   1,
				Name: "name",
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("error")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", *profile.UserID))).Return(1, err).Once()
	adapter.On("Delete", From("profiles").Where(Eq("id", profile.ID))).Return(1, nil).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Delete(context.TODO(), &profile, Cascade(true)))

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
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10)), "", addressMut).Return(1, nil).Once()
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
	adapter.On("Update", From("addresses").Where(Eq("id", 1).AndEq("user_id", 10)), "", addressMut).Return(1, err).Once()
	adapter.On("Rollback").Return(nil).Once()

	assert.Equal(t, err, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasMany(t *testing.T) {
	var (
		user = User{
			ID: 10,
			UserRoles: []UserRole{
				{UserID: 10, RoleID: 1},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Delete", From("user_roles").Where(Eq("user_id", 10).And(Or(Eq("user_id", 10).AndEq("role_id", 1))))).Return(1, nil).Once()
	adapter.On("Delete", From("users").Where(Eq("id", 10))).Return(1, nil).Once()
	adapter.On("Commit").Return(nil).Once()

	assert.Nil(t, repo.Delete(context.TODO(), &user, Cascade(true)))

	adapter.AssertExpectations(t)
}

func TestRepository_Delete_hasManyError(t *testing.T) {
	var (
		user = User{
			ID: 10,
			UserRoles: []UserRole{
				{UserID: 10, RoleID: 1},
			},
		}
		adapter = &testAdapter{}
		repo    = New(adapter)
		err     = errors.New("err")
	)

	adapter.On("Begin").Return(nil).Once()
	adapter.On("Delete", From("user_roles").Where(Eq("user_id", 10).And(Or(Eq("user_id", 10).AndEq("role_id", 1))))).Return(1, err).Once()
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

	deletedCount, err := repo.DeleteAll(context.TODO(), queries)
	assert.Nil(t, err)
	assert.Equal(t, 1, deletedCount)

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
		deletedCount := repo.MustDeleteAll(context.TODO(), queries)
		assert.Equal(t, 1, deletedCount)
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

func TestRepository_Preload_alreadyLoaded(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		userID  = 1
		address = Address{
			UserID: &userID,
			User:   &User{ID: userID},
		}
	)

	assert.Nil(t, repo.Preload(context.TODO(), &address, "user"))
	adapter.AssertExpectations(t)
}

func TestRepository_Preload_alreadyLoadedForceReload(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		userID  = 1
		address = Address{
			UserID: &userID,
			User:   &User{ID: userID},
		}
		cur = &testCursor{}
	)

	adapter.On("Query", From("users").Where(In("id", 1)).Reload()).Return(cur, nil).Maybe()

	cur.On("Close").Return(nil).Once()
	cur.On("Fields").Return([]string{"id", "name"}, nil).Once()
	cur.On("Next").Return(true).Once()
	cur.MockScan(userID, "Del Piero").Twice()
	cur.On("Next").Return(false).Once()

	assert.Nil(t, repo.Preload(context.TODO(), &address, "user", Reload(true)))
	adapter.AssertExpectations(t)
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

func TestRepository_Exec(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = "UPDATE users SET something = ? WHERE something2 = ?;"
		args    = []interface{}{3, "sdfds"}
		rets    = []interface{}{1, 2, nil}
	)

	adapter.On("Exec", context.TODO(), query, args).Return(rets...).Once()

	lastInsertedId, rowsAffected, err := repo.Exec(context.TODO(), query, args...)
	assert.Equal(t, rets[0], lastInsertedId)
	assert.Equal(t, rets[1], rowsAffected)
	assert.Equal(t, rets[2], err)

	adapter.AssertExpectations(t)
}

func TestRepository_MustExec(t *testing.T) {
	var (
		adapter = &testAdapter{}
		repo    = New(adapter)
		query   = "UPDATE users SET something = ? WHERE something2 = ?;"
		args    = []interface{}{3, "sdfds"}
		rets    = []interface{}{1, 2, nil}
	)

	adapter.On("Exec", context.TODO(), query, args).Return(rets...).Once()

	assert.NotPanics(t, func() {
		lastInsertedId, rowsAffected, err := repo.Exec(context.TODO(), query, args...)
		assert.Equal(t, rets[0], lastInsertedId)
		assert.Equal(t, rets[1], rowsAffected)
		assert.Equal(t, rets[2], err)
	})

	adapter.AssertExpectations(t)
}

func TestRepository_Transaction(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(nil).On("Commit").Return(nil).Once()

	repo := New(adapter)

	err := repo.Transaction(context.TODO(), func(ctx context.Context) error {
		return nil
	})

	assert.Nil(t, err)

	adapter.AssertExpectations(t)
}

func TestRepository_Transaction_beginError(t *testing.T) {
	adapter := &testAdapter{}
	adapter.On("Begin").Return(errors.New("error")).Once()

	err := New(adapter).Transaction(context.TODO(), func(ctx context.Context) error {
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

	err := New(adapter).Transaction(context.TODO(), func(ctx context.Context) error {
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

	err := New(adapter).Transaction(context.TODO(), func(ctx context.Context) error {
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

	err := New(adapter).Transaction(context.TODO(), func(ctx context.Context) error {
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
		_ = New(adapter).Transaction(context.TODO(), func(ctx context.Context) error {
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
		_ = New(adapter).Transaction(context.TODO(), func(ctx context.Context) error {
			_ = user.ID
			return nil
		})
	})

	adapter.AssertExpectations(t)
}
