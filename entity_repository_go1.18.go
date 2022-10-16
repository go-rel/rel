//go:build go1.18
// +build go1.18

package rel

import (
	"context"
	"reflect"
)

type EntityRepository[T any] interface {
	// Repository returns base Repository wrapped by this EntityRepository.
	Repository() Repository

	// // Iterate through a collection of records from database in batches.
	// // This function returns iterator that can be used to loop all records.
	// // Limit, Offset and Sort query is automatically ignored.
	// Iterate(ctx context.Context, query Query, option ...IteratorOption) Iterator

	// Aggregate over the given field.
	// Supported aggregate: count, sum, avg, max, min.
	// Any select, group, offset, limit and sort query will be ignored automatically.
	// If complex aggregation is needed, consider using All instead.
	Aggregate(ctx context.Context, aggregate string, field string, queriers ...Querier) (int, error)

	// MustAggregate over the given field.
	// Supported aggregate: count, sum, avg, max, min.
	// Any select, group, offset, limit and sort query will be ignored automatically.
	// If complex aggregation is needed, consider using All instead.
	// It'll panic if any error eccured.
	MustAggregate(ctx context.Context, aggregate string, field string, queriers ...Querier) int

	// Count records that match the query.
	Count(ctx context.Context, queriers ...Querier) (int, error)

	// MustCount records that match the query.
	// It'll panic if any error eccured.
	MustCount(ctx context.Context, queriers ...Querier) int

	// Find a record that match the query.
	// If no result found, it'll return not found error.
	Find(ctx context.Context, queriers ...Querier) (T, error)

	// MustFind a record that match the query.
	// If no result found, it'll panic.
	MustFind(ctx context.Context, queriers ...Querier) T

	// FindAll records that match the query.
	FindAll(ctx context.Context, queriers ...Querier) ([]T, error)

	// MustFindAll records that match the query.
	// It'll panic if any error eccured.
	MustFindAll(ctx context.Context, queriers ...Querier) []T

	// FindAndCountAll records that match the query.
	// This is a convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
	// Limit and Offset property will be ignored when performing count query.
	FindAndCountAll(ctx context.Context, queriers ...Querier) ([]T, int, error)

	// MustFindAndCountAll records that match the query.
	// This is a convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
	// Limit and Offset property will be ignored when performing count query.
	// It'll panic if any error eccured.
	MustFindAndCountAll(ctx context.Context, queriers ...Querier) ([]T, int)

	// Insert a record to database.
	Insert(ctx context.Context, record *T, mutators ...Mutator) error

	// MustInsert an record to database.
	// It'll panic if any error occurred.
	MustInsert(ctx context.Context, record *T, mutators ...Mutator)

	// InsertAll records.
	// Does not supports application cascade insert.
	InsertAll(ctx context.Context, records *[]T, mutators ...Mutator) error

	// MustInsertAll records.
	// It'll panic if any error occurred.
	// Does not supports application cascade insert.
	MustInsertAll(ctx context.Context, records *[]T, mutators ...Mutator)

	// Update a record in database.
	// It'll panic if any error occurred.
	Update(ctx context.Context, record *T, mutators ...Mutator) error

	// MustUpdate a record in database.
	// It'll panic if any error occurred.
	MustUpdate(ctx context.Context, record *T, mutators ...Mutator)

	// Delete a record.
	Delete(ctx context.Context, record *T, mutators ...Mutator) error

	// MustDelete a record.
	// It'll panic if any error eccured.
	MustDelete(ctx context.Context, record *T, mutators ...Mutator)

	// DeleteAll records.
	// Does not supports application cascade delete.
	DeleteAll(ctx context.Context, records *[]T) error

	// MustDeleteAll records.
	// It'll panic if any error occurred.
	// Does not supports application cascade delete.
	MustDeleteAll(ctx context.Context, records *[]T)

	// Preload association with given query.
	// If association is already loaded, this will do nothing.
	// To force preloading even though association is already loaeded, add `Reload(true)` as query.
	Preload(ctx context.Context, record *T, field string, queriers ...Querier) error

	// MustPreload association with given query.
	// It'll panic if any error occurred.
	MustPreload(ctx context.Context, record *T, field string, queriers ...Querier)

	// Preload association with given query.
	// If association is already loaded, this will do nothing.
	// To force preloading even though association is already loaeded, add `Reload(true)` as query.
	PreloadAll(ctx context.Context, records *[]T, field string, queriers ...Querier) error

	// MustPreload association with given query.
	// It'll panic if any error occurred.
	MustPreloadAll(ctx context.Context, records *[]T, field string, queriers ...Querier)

	// Transaction performs transaction with given function argument.
	// Transaction scope/connection is automatically passed using context.
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type entityRepository[T any] struct {
	repository Repository
}

func (er entityRepository[T]) Repository() Repository {
	return er.repository
}

func (er entityRepository[T]) Aggregate(ctx context.Context, aggregate string, field string, queriers ...Querier) (int, error) {
	var (
		entity       T
		documentMeta = getDocumentMeta(reflect.TypeOf(entity), true)
		query        = Build(documentMeta.table, queriers...)
	)

	return er.repository.Aggregate(ctx, query, aggregate, field)
}

func (er entityRepository[T]) MustAggregate(ctx context.Context, aggregate string, field string, queriers ...Querier) int {
	result, err := er.Aggregate(ctx, aggregate, field, queriers...)
	must(err)
	return result
}

func (er entityRepository[T]) Count(ctx context.Context, queriers ...Querier) (int, error) {
	var (
		entity       T
		documentMeta = getDocumentMeta(reflect.TypeOf(entity), true)
	)

	return er.repository.Count(ctx, documentMeta.Table(), queriers...)
}

func (er entityRepository[T]) MustCount(ctx context.Context, queriers ...Querier) int {
	result, err := er.Count(ctx, queriers...)
	must(err)
	return result
}

func (er entityRepository[T]) Find(ctx context.Context, queriers ...Querier) (T, error) {
	var entity T
	return entity, er.repository.Find(ctx, &entity, queriers...)
}

func (er entityRepository[T]) MustFind(ctx context.Context, queriers ...Querier) T {
	entity, err := er.Find(ctx, queriers...)
	must(err)
	return entity
}

func (er entityRepository[T]) FindAll(ctx context.Context, queriers ...Querier) ([]T, error) {
	var entities []T
	return entities, er.repository.FindAll(ctx, &entities, queriers...)
}

func (er entityRepository[T]) MustFindAll(ctx context.Context, queriers ...Querier) []T {
	entities, err := er.FindAll(ctx, queriers...)
	must(err)
	return entities
}

func (er entityRepository[T]) FindAndCountAll(ctx context.Context, queriers ...Querier) ([]T, int, error) {
	var entities []T
	count, err := er.repository.FindAndCountAll(ctx, &entities, queriers...)
	return entities, count, err
}

func (er entityRepository[T]) MustFindAndCountAll(ctx context.Context, queriers ...Querier) ([]T, int) {
	entities, count, err := er.FindAndCountAll(ctx, queriers...)
	must(err)
	return entities, count
}

func (er entityRepository[T]) Insert(ctx context.Context, record *T, mutators ...Mutator) error {
	return er.repository.Insert(ctx, record, mutators...)
}

func (er entityRepository[T]) MustInsert(ctx context.Context, record *T, mutators ...Mutator) {
	er.repository.MustInsert(ctx, record, mutators...)
}

func (er entityRepository[T]) InsertAll(ctx context.Context, records *[]T, mutators ...Mutator) error {
	return er.repository.InsertAll(ctx, records, mutators...)
}

func (er entityRepository[T]) MustInsertAll(ctx context.Context, records *[]T, mutators ...Mutator) {
	er.repository.MustInsertAll(ctx, records, mutators...)
}

func (er entityRepository[T]) Update(ctx context.Context, record *T, mutators ...Mutator) error {
	return er.repository.Update(ctx, record, mutators...)
}

func (er entityRepository[T]) MustUpdate(ctx context.Context, record *T, mutators ...Mutator) {
	er.repository.MustUpdate(ctx, record, mutators...)
}

func (er entityRepository[T]) Delete(ctx context.Context, record *T, mutators ...Mutator) error {
	return er.repository.Delete(ctx, record, mutators...)
}

func (er entityRepository[T]) MustDelete(ctx context.Context, record *T, mutators ...Mutator) {
	er.repository.MustDelete(ctx, record, mutators...)
}

func (er entityRepository[T]) DeleteAll(ctx context.Context, records *[]T) error {
	return er.repository.DeleteAll(ctx, records)
}

func (er entityRepository[T]) MustDeleteAll(ctx context.Context, records *[]T) {
	er.repository.MustDeleteAll(ctx, records)
}

func (er entityRepository[T]) Preload(ctx context.Context, record *T, field string, queriers ...Querier) error {
	return er.repository.Preload(ctx, record, field, queriers...)
}

func (er entityRepository[T]) MustPreload(ctx context.Context, record *T, field string, queriers ...Querier) {
	er.repository.MustPreload(ctx, record, field, queriers...)
}

func (er entityRepository[T]) PreloadAll(ctx context.Context, records *[]T, field string, queriers ...Querier) error {
	return er.repository.Preload(ctx, records, field, queriers...)
}

func (er entityRepository[T]) MustPreloadAll(ctx context.Context, records *[]T, field string, queriers ...Querier) {
	er.repository.MustPreload(ctx, records, field, queriers...)
}

func (er entityRepository[T]) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	return er.repository.Transaction(ctx, fn)
}

func NewEntityRepository[T any](repository Repository) EntityRepository[T] {
	return entityRepository[T]{
		repository: repository,
	}
}
