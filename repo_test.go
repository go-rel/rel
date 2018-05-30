package grimoire

import (
	"testing"

	"github.com/Fs02/grimoire/errors"
	"github.com/stretchr/testify/assert"
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

func TestRepo_From(t *testing.T) {
	assert.Equal(t, repo.From("users"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
	})
}

func TestRepo_Transaction(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Commit").Return(nil)

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Nil(t, err)
	mock.AssertExpectations(t)
}

func TestRepo_Transaction_beginError(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(errors.NewUnexpected("error"))

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	mock.AssertExpectations(t)
}

func TestRepo_Transaction_commitError(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Commit").Return(errors.NewUnexpected("error"))

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	mock.AssertExpectations(t)
}

func TestRepo_Transaction_returnErrorAndRollback(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return errors.NewUnexpected("error")
	})

	assert.Equal(t, errors.NewUnexpected("error"), err)
	mock.AssertExpectations(t)
}

func TestRepo_Transaction_panicWithKnownErrorAndRollback(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		panic(errors.New("error", "", errors.NotFound))
	})

	assert.Equal(t, errors.New("error", "", errors.NotFound), err)
	mock.AssertExpectations(t)
}

func TestRepo_Transaction_panicAndRollback(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	assert.Panics(t, func() {
		Repo{adapter: mock}.Transaction(func(r Repo) error {
			// doing good things
			panic(errors.NewUnexpected("error"))
		})
	})

	mock.AssertExpectations(t)
}
