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

func TestRepoSetLogger(t *testing.T) {
	repo := Repo{}
	assert.Nil(t, repo.logger)
	repo.SetLogger(DefaultLogger)
	assert.NotNil(t, repo.logger)
}

func TestRepoFrom(t *testing.T) {
	assert.Equal(t, repo.From("users"), Query{
		repo:       &repo,
		Collection: "users",
		Fields:     []string{"users.*"},
	})
}

func TestRepoTransaction(t *testing.T) {
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

func TestTransactionBeginError(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(errors.UnexpectedError("error"))

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.UnexpectedError("error"), err)
	mock.AssertExpectations(t)
}

func TestTransactionCommitError(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Commit").Return(errors.UnexpectedError("error"))

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return nil
	})

	assert.Equal(t, errors.UnexpectedError("error"), err)
	mock.AssertExpectations(t)
}

func TestTransactionReturnErrorAndRollback(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		return errors.UnexpectedError("error")
	})

	assert.Equal(t, errors.UnexpectedError("error"), err)
	mock.AssertExpectations(t)
}

func TestTransactionPanicWithKnownErrorAndRollback(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	err := Repo{adapter: mock}.Transaction(func(r Repo) error {
		// doing good things
		panic(errors.NotFoundError("error"))
	})

	assert.Equal(t, errors.NotFoundError("error"), err)
	mock.AssertExpectations(t)
}

func TestTransactionPanicAndRollback(t *testing.T) {
	mock := new(TestAdapter)
	mock.On("Begin").Return(nil).
		On("Rollback").Return(nil)

	assert.Panics(t, func() {
		Repo{adapter: mock}.Transaction(func(r Repo) error {
			// doing good things
			panic(errors.UnexpectedError("error"))
		})
	})

	mock.AssertExpectations(t)
}
