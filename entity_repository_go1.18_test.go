//go:build go1.18
// +build go1.18

package rel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type testRepository struct {
	mock.Mock
}

func (tr *testRepository) Adapter(ctx context.Context) Adapter {
	return nil
}

func (tr *testRepository) Instrumentation(instrumenter Instrumenter) {}

func (tr *testRepository) Ping(ctx context.Context) error {
	return nil
}

func (tr *testRepository) Iterate(ctx context.Context, query Query, option ...IteratorOption) Iterator {
	args := tr.Called(query, option)
	return args.Get(0).(Iterator)
}

func (tr *testRepository) Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	args := tr.Called(query, aggregate, field)
	return args.Int(0), args.Error(1)
}

func (tr *testRepository) MustAggregate(ctx context.Context, query Query, aggregate string, field string) int {
	args := tr.Called(query, aggregate, query, field)
	return args.Int(0)
}

func (tr *testRepository) Count(ctx context.Context, collection string, queriers ...Querier) (int, error) {
	args := tr.Called(collection, queriers)
	return args.Int(0), args.Error(1)
}

func (tr *testRepository) MustCount(ctx context.Context, collection string, queriers ...Querier) int {
	args := tr.Called(collection, queriers)
	return args.Int(0)
}

func (tr *testRepository) Find(ctx context.Context, record any, queriers ...Querier) error {
	args := tr.Called(record, queriers)
	return args.Error(0)
}

func (tr *testRepository) MustFind(ctx context.Context, record any, queriers ...Querier) {
	tr.Called(record, queriers)
}

func (tr *testRepository) FindAll(ctx context.Context, records any, queriers ...Querier) error {
	args := tr.Called(records, queriers)
	return args.Error(0)
}

func (tr *testRepository) MustFindAll(ctx context.Context, records any, queriers ...Querier) {
	tr.Called(records, queriers)
}

func (tr *testRepository) FindAndCountAll(ctx context.Context, records any, queriers ...Querier) (int, error) {
	args := tr.Called(records, queriers)
	return args.Int(0), args.Error(1)
}

func (tr *testRepository) MustFindAndCountAll(ctx context.Context, records any, queriers ...Querier) int {
	args := tr.Called(records, queriers)
	return args.Int(0)
}

func (tr *testRepository) Insert(ctx context.Context, record any, mutators ...Mutator) error {
	args := tr.Called(record, mutators)
	return args.Error(0)
}

func (tr *testRepository) MustInsert(ctx context.Context, record any, mutators ...Mutator) {
	tr.Called(record, mutators)
}

func (tr *testRepository) InsertAll(ctx context.Context, records any, mutators ...Mutator) error {
	args := tr.Called(records, mutators)
	return args.Error(0)
}

func (tr *testRepository) MustInsertAll(ctx context.Context, records any, mutators ...Mutator) {
	tr.Called(records, mutators)
}

func (tr *testRepository) Update(ctx context.Context, record any, mutators ...Mutator) error {
	args := tr.Called(record, mutators)
	return args.Error(0)
}

func (tr *testRepository) MustUpdate(ctx context.Context, record any, mutators ...Mutator) {
	tr.Called(record, mutators)
}

func (tr *testRepository) UpdateAny(ctx context.Context, query Query, mutates ...Mutate) (int, error) {
	args := tr.Called(query, mutates)
	return args.Int(0), args.Error(1)
}

func (tr *testRepository) MustUpdateAny(ctx context.Context, query Query, mutates ...Mutate) int {
	args := tr.Called(query, mutates)
	return args.Int(0)
}

func (tr *testRepository) Delete(ctx context.Context, record any, mutators ...Mutator) error {
	args := tr.Called(record, mutators)
	return args.Error(0)
}

func (tr *testRepository) MustDelete(ctx context.Context, record any, mutators ...Mutator) {
	tr.Called(record, mutators)
}

func (tr *testRepository) DeleteAll(ctx context.Context, records any) error {
	args := tr.Called(records)
	return args.Error(0)
}

func (tr *testRepository) MustDeleteAll(ctx context.Context, records any) {
	tr.Called(records)
}

func (tr *testRepository) DeleteAny(ctx context.Context, query Query) (int, error) {
	args := tr.Called(query)
	return args.Int(0), args.Error(1)
}

func (tr *testRepository) MustDeleteAny(ctx context.Context, query Query) int {
	args := tr.Called(query)
	return args.Int(0)
}

func (tr *testRepository) Preload(ctx context.Context, records any, field string, queriers ...Querier) error {
	args := tr.Called(records, field, queriers)
	return args.Error(0)
}

func (tr *testRepository) MustPreload(ctx context.Context, records any, field string, queriers ...Querier) {
	tr.Called(records, field, queriers)
}

func (tr *testRepository) Exec(ctx context.Context, statement string, arg ...any) (int, int, error) {
	args := tr.Called(statement, statement, arg)
	return args.Int(0), args.Int(1), args.Error(2)
}

func (tr *testRepository) MustExec(ctx context.Context, statement string, arg ...any) (int, int) {
	args := tr.Called(statement, statement, arg)
	return args.Int(0), args.Int(1)
}

func (tr *testRepository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	tr.Called()
	return fn(ctx)
}

func TestEntityRepository_Repository(t *testing.T) {
	var (
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	assert.Equal(t, repo, entityRepo.Repository())
	repo.AssertExpectations(t)
}

func TestEntityRepository_Aggregate(t *testing.T) {
	var (
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Aggregate", From("users"), "max", "score").Return(1, nil)

	result, err := entityRepo.Aggregate(context.TODO(), "max", "score")
	assert.Nil(t, err)
	assert.Equal(t, 1, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustAggregate(t *testing.T) {
	var (
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Aggregate", From("users"), "max", "score").Return(1, nil)

	result := entityRepo.MustAggregate(context.TODO(), "max", "score")
	assert.Equal(t, 1, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_Count(t *testing.T) {
	var (
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Count", "users", []Querier(nil)).Return(1, nil)

	result, err := entityRepo.Count(context.TODO())
	assert.Nil(t, err)
	assert.Equal(t, 1, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustCount(t *testing.T) {
	var (
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Count", "users", []Querier(nil)).Return(1, nil)

	result := entityRepo.MustCount(context.TODO())
	assert.Equal(t, 1, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_Find(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
		query      = From("users").Limit(1)
	)

	repo.On("Find", &user, []Querier{query}).Return(nil)

	result, err := entityRepo.Find(context.TODO(), query)
	assert.Nil(t, err)
	assert.Equal(t, user, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustFind(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
		query      = From("users").Limit(1)
	)

	repo.On("Find", &user, []Querier{query}).Return(nil)

	result := entityRepo.MustFind(context.TODO(), query)
	assert.Equal(t, user, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_FindAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
		query      = From("users").Limit(1)
	)

	repo.On("FindAll", &users, []Querier{query}).Return(nil)

	result, err := entityRepo.FindAll(context.TODO(), query)
	assert.Nil(t, err)
	assert.Equal(t, users, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustFindAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
		query      = From("users").Limit(1)
	)

	repo.On("FindAll", &users, []Querier{query}).Return(nil)

	result := entityRepo.MustFindAll(context.TODO(), query)
	assert.Equal(t, users, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_FindAndCountAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
		query      = From("users").Limit(1)
	)

	repo.On("FindAndCountAll", &users, []Querier{query}).Return(1, nil)

	result, count, err := entityRepo.FindAndCountAll(context.TODO(), query)
	assert.Nil(t, err)
	assert.Equal(t, 1, count)
	assert.Equal(t, users, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustFindAndCountAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
		query      = From("users").Limit(1)
	)

	repo.On("FindAndCountAll", &users, []Querier{query}).Return(1, nil)

	result, count := entityRepo.MustFindAndCountAll(context.TODO(), query)
	assert.Equal(t, 1, count)
	assert.Equal(t, users, result)

	repo.AssertExpectations(t)
}

func TestEntityRepository_Insert(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Insert", &user, []Mutator(nil)).Return(nil)

	err := entityRepo.Insert(context.TODO(), &user)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustInsert(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustInsert", &user, []Mutator(nil))

	entityRepo.MustInsert(context.TODO(), &user)

	repo.AssertExpectations(t)
}

func TestEntityRepository_InsertAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("InsertAll", &users, []Mutator(nil)).Return(nil)

	err := entityRepo.InsertAll(context.TODO(), &users)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustInsertAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustInsertAll", &users, []Mutator(nil))

	entityRepo.MustInsertAll(context.TODO(), &users)

	repo.AssertExpectations(t)
}

func TestEntityRepository_Update(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Update", &user, []Mutator(nil)).Return(nil)

	err := entityRepo.Update(context.TODO(), &user)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustUpdate(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustUpdate", &user, []Mutator(nil))

	entityRepo.MustUpdate(context.TODO(), &user)

	repo.AssertExpectations(t)
}

func TestEntityRepository_Delete(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Delete", &user, []Mutator(nil)).Return(nil)

	err := entityRepo.Delete(context.TODO(), &user)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustDelete(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustDelete", &user, []Mutator(nil))

	entityRepo.MustDelete(context.TODO(), &user)

	repo.AssertExpectations(t)
}

func TestEntityRepository_DeleteAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("DeleteAll", &users).Return(nil)

	err := entityRepo.DeleteAll(context.TODO(), &users)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustDeleteAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustDeleteAll", &users)

	entityRepo.MustDeleteAll(context.TODO(), &users)

	repo.AssertExpectations(t)
}

func TestEntityRepository_Preload(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Preload", &user, "address", []Querier(nil)).Return(nil)

	err := entityRepo.Preload(context.TODO(), &user, "address")
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustPreload(t *testing.T) {
	var (
		user       User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustPreload", &user, "address", []Querier(nil))

	entityRepo.MustPreload(context.TODO(), &user, "address")

	repo.AssertExpectations(t)
}

func TestEntityRepository_PreloadAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Preload", &users, "address", []Querier(nil)).Return(nil)

	err := entityRepo.PreloadAll(context.TODO(), &users, "address")
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestEntityRepository_MustPreloadAll(t *testing.T) {
	var (
		users      []User
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("MustPreload", &users, "address", []Querier(nil))

	entityRepo.MustPreloadAll(context.TODO(), &users, "address")

	repo.AssertExpectations(t)
}

func TestEntityRepository_Transaction(t *testing.T) {
	var (
		repo       = &testRepository{}
		entityRepo = NewEntityRepository[User](repo)
	)

	repo.On("Transaction")

	err := entityRepo.Transaction(context.TODO(), func(ctx context.Context) error {
		return nil
	})

	assert.Nil(t, err)

	repo.AssertExpectations(t)
}
