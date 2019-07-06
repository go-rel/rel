package grimoire

import (
	"testing"
	"time"

	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/query"
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
	ch, user := createChangeset()
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	adapter.On("Insert", query.From("users"), changes).Return(1, nil).
		On("All", query.From("users").Where(where.Eq("id", 1)).Limit(1), &user).Return(1, nil)

	assert.Nil(t, repo.Insert(&user, ch))
	assert.NotPanics(t, func() { repo.MustInsert(&user, ch) })
	adapter.AssertExpectations(t)
}

func TestRepo_Insert_multiple(t *testing.T) {
	ch1, user1 := createChangeset()
	ch2, user2 := createChangeset()
	users := []User{user1, user2}

	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	allchanges := []map[string]interface{}{changes, changes}

	adapter.On("InsertAll", query.From("users"), allchanges).Return([]interface{}{1, 2}, nil).
		On("All", query.From("users").Where(where.In("id", 1, 2)), &users).Return(2, nil)

	assert.Nil(t, repo.Insert(&users, ch1, ch2))
	assert.NotPanics(t, func() { repo.MustInsert(&users, ch1, ch2) })
	adapter.AssertExpectations(t)
}

func TestRepo_Insert_multipleError(t *testing.T) {
	ch1, user1 := createChangeset()
	ch2, user2 := createChangeset()
	users := []User{user1, user2}

	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	allchanges := []map[string]interface{}{changes, changes}

	adapter.On("InsertAll", query.From("users"), allchanges).
		Return([]interface{}{1, 2}, errors.NewUnexpected("error"))

	assert.NotNil(t, repo.Insert(&users, ch1, ch2))
	assert.Panics(t, func() { repo.MustInsert(&users, ch1, ch2) })
	adapter.AssertExpectations(t)
}

// func TestRepo_Insert_notReturning(t *testing.T) {
// 	ch, _ := createChangeset()
// 	adapter := &TestAdapter{}
// 	repo := Repo{adapter: adapter}

// 	changes := map[string]interface{}{
// 		"name":       "name",
// 		"created_at": time.Now().Round(time.Second),
// 		"updated_at": time.Now().Round(time.Second),
// 	}

// 	adapter.On("Insert", query.From("users"), changes).Return(1, nil)

// 	assert.Nil(t, repo.Insert(nil, ch))
// 	assert.NotPanics(t, func() { repo.MustInsert(nil, ch) })
// 	adapter.AssertExpectations(t)
// }

func TestRepo_Insert_error(t *testing.T) {
	ch, user := createChangeset()
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	changes := map[string]interface{}{
		"name":       "name",
		"created_at": time.Now().Round(time.Second),
		"updated_at": time.Now().Round(time.Second),
	}

	adapter.On("Insert", query.From("users"), changes).Return(0, errors.NewUnexpected("error"))

	assert.NotNil(t, repo.Insert(&user, ch))
	assert.Panics(t, func() { repo.MustInsert(&user, ch) })
	adapter.AssertExpectations(t)
}

func TestRepp_Update(t *testing.T) {
	ch, user := createChangeset()
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	changes := map[string]interface{}{
		"name":       "name",
		"updated_at": time.Now().Round(time.Second),
	}

	query := query.From("users").Where(where.Eq("id", user.ID))
	adapter.On("Update", query, changes).Return(nil).
		On("All", query.Limit(1), &user).Return(1, nil)

	assert.Nil(t, repo.Update(&user, ch))
	assert.NotPanics(t, func() { repo.MustUpdate(&user, ch) })
	adapter.AssertExpectations(t)
}

func TestRepo_Update_nothing(t *testing.T) {
	ch, _ := createChangeset()
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	assert.Nil(t, repo.Update(nil, ch))
	assert.NotPanics(t, func() { repo.MustUpdate(nil, ch) })

	adapter.AssertExpectations(t)
}

func TestRepo_Update_unchanged(t *testing.T) {
	ch := &changeset.Changeset{}
	user := User{}
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	assert.Nil(t, repo.Update(&user, ch))
	assert.NotPanics(t, func() { repo.MustUpdate(&user, ch) })

	adapter.AssertExpectations(t)
}

func TestRepo_Update_error(t *testing.T) {
	ch, user := createChangeset()
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	changes := map[string]interface{}{
		"name":       "name",
		"updated_at": time.Now().Round(time.Second),
	}

	adapter.On("Update", query.From("users").Where(where.Eq("id", 0)), changes).
		Return(errors.NewUnexpected("error"))

	assert.NotNil(t, repo.Update(&user, ch))
	assert.Panics(t, func() { repo.MustUpdate(&user, ch) })
	adapter.AssertExpectations(t)
}

func TestRepo_Delete(t *testing.T) {
	adapter := &TestAdapter{}
	repo := Repo{adapter: adapter}

	adapter.On("Delete", query.From("users").Where(where.Eq("id", 0))).Return(nil)

	assert.Nil(t, repo.Delete(User{}))
	assert.NotPanics(t, func() { repo.MustDelete(User{}) })
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
