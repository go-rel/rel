package grimoire

import (
	"reflect"
	"strings"

	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/query"
	"github.com/Fs02/grimoire/schema"
	"github.com/Fs02/grimoire/where"
)

// Repo defines grimoire repository.
type Repo struct {
	adapter       Adapter
	logger        []Logger
	inTransaction bool
}

// New create new repo using adapter.
func New(adapter Adapter) Repo {
	return Repo{
		adapter: adapter,
		logger:  []Logger{DefaultLogger},
	}
}

// Adapter returns adapter of repo.
func (r *Repo) Adapter() Adapter {
	return r.adapter
}

// SetLogger replace default logger with custom logger.
func (r *Repo) SetLogger(logger ...Logger) {
	r.logger = logger
}

// Aggregate calculate aggregate over the given field.
func (r Repo) Aggregate(record interface{}, mode string, field string, out interface{}, queries ...query.Builder) error {
	table := schema.InferTableName(record)
	q := query.Build(table, queries...)
	return r.adapter.Aggregate(q, out, mode, field, r.logger...)
}

// MustAggregate calculate aggregate over the given field.
// It'll panic if any error eccured.
func (r Repo) MustAggregate(record interface{}, mode string, field string, out interface{}, queries ...query.Builder) {
	must(r.Aggregate(record, mode, field, out, queries...))
}

// Count retrieves count of results that match the query.
func (r Repo) Count(record interface{}, queries ...query.Builder) (int, error) {
	var out struct {
		Count int
	}

	err := r.Aggregate(record, "COUNT", "*", &out, queries...)
	return out.Count, err
}

// MustCount retrieves count of results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustCount(record interface{}, queries ...query.Builder) int {
	count, err := r.Count(record, queries...)
	must(err)
	return count
}

// One retrieves one result that match the query.
// If no result found, it'll return not found error.
func (r Repo) One(record interface{}, queries ...query.Builder) error {
	table := schema.InferTableName(record)
	q := query.Build(table, queries...).Limit(1)

	count, err := r.adapter.All(q, record, r.logger...)

	if err != nil {
		return transformError(err)
	} else if count == 0 {
		return errors.New("no result found", "", errors.NotFound)
	} else {
		return nil
	}
}

// MustOne retrieves one result that match the query.
// If no result found, it'll panic.
func (r Repo) MustOne(record interface{}, queries ...query.Builder) {
	must(r.One(record, queries...))
}

// All retrieves all results that match the query.
func (r Repo) All(record interface{}, queries ...query.Builder) error {
	table := schema.InferTableName(record)
	q := query.Build(table, queries...)
	_, err := r.adapter.All(q, record, r.logger...)
	return err
}

// MustAll retrieves all results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustAll(record interface{}, queries ...query.Builder) {
	must(r.All(record, queries...))
}

// Insert records to database.
func (r Repo) Insert(record interface{}, chs ...*changeset.Changeset) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil || len(chs) == 0 {
		return nil
	}

	var err error
	var ids []interface{}

	// if string means no returning
	table := schema.InferTableName(record)
	primaryKey, _ := schema.InferPrimaryKey(record, false)

	q := query.Build(table)

	if len(chs) == 1 {
		// single insert
		ch := chs[0]
		changes := ch.Changes()
		// cloneChangeset(changes, ch.Changes())
		putTimestamp(changes, "created_at", ch.Types())
		putTimestamp(changes, "updated_at", ch.Types())
		// cloneQuery(changes, query.Changes)

		var id interface{}
		id, err = r.adapter.Insert(q, changes, r.logger...)
		ids = append(ids, id)
	} else if len(chs) > 1 {
		// multiple insert
		fields := getFields(chs)

		allchanges := make([]map[string]interface{}, len(chs))
		for i, ch := range chs {
			changes := ch.Changes()
			// cloneChangeset(changes, ch.Changes())
			putTimestamp(changes, "created_at", ch.Types())
			putTimestamp(changes, "updated_at", ch.Types())
			// cloneQuery(changes, query.Changes)

			allchanges[i] = changes
		}

		ids, err = r.adapter.InsertAll(q, fields, allchanges, r.logger...)
	} else { //if len(query.Changes) > 0 {
		// set only
		// var id interface{}
		// id, err = r.adapter.Insert(query, query.Changes, r.logger...)
		// ids = append(ids, id)
	}

	if err != nil {
		return transformError(err, chs...)
	} else if record == nil || len(ids) == 0 {
		return nil
	} else if len(ids) == 1 {
		return transformError(r.One(record, where.Eq(primaryKey, ids[0])))
	}

	return transformError(r.All(record, where.In(primaryKey, ids...)))
}

// MustInsert records to database.
// It'll panic if any error occurred.
func (r Repo) MustInsert(record interface{}, chs ...*changeset.Changeset) {
	must(r.Insert(record, chs...))
}

// Update records in database.
// It'll panic if any error occurred.
func (r Repo) Update(record interface{}, ch *changeset.Changeset) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil || len(ch.Changes()) == 0 {
		return nil
	}

	changes := ch.Changes()
	putTimestamp(changes, "updated_at", ch.Types())

	table := schema.InferTableName(record)
	primaryKey, primaryValue := schema.InferPrimaryKey(record, true)

	q := query.Build(table, where.Eq(primaryKey, primaryValue))

	// perform update
	err := r.adapter.Update(q, changes, r.logger...)
	if err != nil {
		return transformError(err, ch)
	}

	return r.One(record, q)
}

// MustUpdate records in database.
// It'll panic if any error occurred.
func (r Repo) MustUpdate(record interface{}, chs *changeset.Changeset) {
	must(r.Update(record, chs))
}

// Delete deletes all results that match the query.
func (r Repo) Delete(record interface{}) error {
	table := schema.InferTableName(record)
	primaryKey, primaryValue := schema.InferPrimaryKey(record, true)

	q := query.Build(table, where.Eq(primaryKey, primaryValue))

	return transformError(r.adapter.Delete(q, r.logger...))
}

// MustDelete deletes all results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustDelete(record interface{}) {
	must(r.Delete(record))
}

// Preload loads association with given query.
func (r Repo) Preload(record interface{}, field string, queries ...query.Builder) error {
	var (
		path = strings.Split(field, ".")
		rv   = reflect.ValueOf(record)
	)

	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		panic("grimoire: record parameter must be a pointer.")
	}

	preload := traversePreloadTarget(rv.Elem(), path)
	if len(preload) == 0 {
		return nil
	}

	schemaType := preload[0].schema.Type()
	refIndex, fkIndex, column := schema.InferAssociation(schemaType, path[len(path)-1])

	addrs, ids := collectPreloadTarget(preload, refIndex)
	if len(ids) == 0 {
		return nil
	}

	// prepare temp result variable for querying
	rt := preload[0].field.Type()
	if rt.Kind() == reflect.Slice || rt.Kind() == reflect.Array || rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	slice := reflect.MakeSlice(reflect.SliceOf(rt), 0, 0)
	result := reflect.New(slice.Type())
	result.Elem().Set(slice)

	// query all records using collected ids.
	err := r.All(result.Interface(), where.In(column, ids...))
	if err != nil {
		return err
	}

	// map results.
	result = result.Elem()
	for i := 0; i < result.Len(); i++ {
		curr := result.Index(i)
		id := getPreloadID(curr.FieldByIndex(fkIndex))

		for _, addr := range addrs[id] {
			if addr.Kind() == reflect.Slice {
				addr.Set(reflect.Append(addr, curr))
			} else if addr.Kind() == reflect.Ptr {
				currP := reflect.New(curr.Type())
				currP.Elem().Set(curr)
				addr.Set(currP)
			} else {
				addr.Set(curr)
			}
		}
	}

	return nil
}

// MustPreload loads association with given query.
// It'll panic if any error occurred.
func (r Repo) MustPreload(record interface{}, field string, queries ...query.Builder) {
	must(r.Preload(record, field, queries...))
}

// Transaction performs transaction with given function argument.
func (r Repo) Transaction(fn func(Repo) error) error {
	adp, err := r.adapter.Begin()
	if err != nil {
		return err
	}

	txRepo := New(adp)
	txRepo.inTransaction = true

	func() {
		defer func() {
			if p := recover(); p != nil {
				txRepo.adapter.Rollback()

				if e, ok := p.(errors.Error); ok && e.Kind() != errors.Unexpected {
					err = e
				} else {
					panic(p) // re-throw panic after Rollback
				}
			} else if err != nil {
				txRepo.adapter.Rollback()
			} else {
				err = txRepo.adapter.Commit()
			}
		}()

		err = fn(txRepo)
	}()

	return err
}
