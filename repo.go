package grimoire

import (
	"reflect"
	"strings"

	"github.com/Fs02/grimoire/change"
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
	var (
		table  = schema.InferTableName(record)
		q      = query.Build(table, queries...)
		_, err = r.adapter.All(q, record, r.logger...)
	)

	return err
}

// MustAll retrieves all results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustAll(record interface{}, queries ...query.Builder) {
	must(r.All(record, queries...))
}

// Insert a record to database.
// TODO: insert all (multiple changes as multiple records)
func (r Repo) Insert(record interface{}, cbuilders ...change.Builder) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil || len(cbuilders) == 0 {
		return nil
	}

	// TODO: transform changeset error
	return transformError(r.insert(record, change.Build(cbuilders...)))
}

func (r Repo) insert(record interface{}, changes change.Changes) error {
	var (
		table         = schema.InferTableName(record)
		primaryKey, _ = schema.InferPrimaryKey(record, false)
		association   = schema.InferAssociations(record)
		queries       = query.Build(table)
	)

	if err := r.upsertBelongsTo(association, &changes); err != nil {
		return err
	}

	// TODO: put timestamp (updated_at, created_at) if those fields exist.
	id, err := r.Adapter().Insert(queries, changes, r.logger...)
	if err != nil {
		return err
	}

	// fetch record
	if err := r.One(record, where.Eq(primaryKey, id)); err != nil {
		return err
	}

	if err := r.upsertHasOne(association, &changes, id); err != nil {
		return err
	}

	if err := r.upsertHasMany(association, &changes, id); err != nil {
		return err
	}

	return nil
}

// MustInsert a record to database.
// It'll panic if any error occurred.
func (r Repo) MustInsert(record interface{}, cbuilders ...change.Builder) {
	must(r.Insert(record, cbuilders...))
}

func (r Repo) InsertAll(record interface{}, changes []change.Changes) error {
	return transformError(r.insertAll(record, changes))
}

func (r Repo) MustInsertAll(record interface{}, changes []change.Changes) {
	must(r.InsertAll(record, changes))
}

// TODO: support assocs
func (r Repo) insertAll(record interface{}, changes []change.Changes) error {
	if len(changes) == 0 {
		return nil
	}

	var (
		table         = schema.InferTableName(record)
		primaryKey, _ = schema.InferPrimaryKey(record, false)
		queries       = query.Build(table)
		fields        = make([]string, 0, len(changes[0].Fields))
		fieldMap      = make(map[string]struct{}, len(changes[0].Fields))
	)

	for i := range changes {
		for _, ch := range changes[i].Changes {
			if _, exist := fieldMap[ch.Field]; !exist {
				fieldMap[ch.Field] = struct{}{}
				fields = append(fields, ch.Field)
			}
		}
	}

	ids, err := r.adapter.InsertAll(queries, fields, changes, r.logger...)
	if err != nil {
		return err
	}

	_, err = r.adapter.All(queries.Where(where.In(primaryKey, ids...)), record, r.logger...)
	return err
}

// Update a record in database.
// It'll panic if any error occurred.
func (r Repo) Update(record interface{}, cbuilders ...change.Builder) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil || len(cbuilders) == 0 {
		return nil
	}

	var (
		pKey, pValues = schema.InferPrimaryKey(record, true)
		changes       = change.Build(cbuilders...)
	)

	if len(pValues) == 0 {
		panic("grimoire: must be a struct")
	}

	return r.update(record, changes, where.Eq(pKey, pValues[0]))
}

func (r Repo) update(record interface{}, changes change.Changes, filter query.FilterClause) error {
	if changes.Empty() {
		return nil
	}

	var (
		table   = schema.InferTableName(record)
		queries = query.Build(table, filter)
	)

	// TODO: update timestamp (updated_at)

	// perform update
	err := r.adapter.Update(queries, changes, r.logger...)
	if err != nil {
		// TODO: changeset error
		return transformError(err)
	}

	return r.One(record, queries)
}

// MustUpdate a record in database.
// It'll panic if any error occurred.
func (r Repo) MustUpdate(record interface{}, cbuilders ...change.Builder) {
	must(r.Update(record, cbuilders...))
}

func (r Repo) upsertBelongsTo(assocs schema.Associations, changes *change.Changes) error {
	for _, field := range assocs.BelongsTo() {
		allAssocChanges, changed := changes.GetAssoc(field)
		if !changed || len(allAssocChanges) == 0 {
			continue
		}

		var (
			assocChanges   = allAssocChanges[0]
			assoc          = assocs.Association(field)
			target, loaded = assoc.TargetAddr()
			foreignValue   = assoc.ForeignValue()
		)

		if loaded {
			var (
				pKey, pValues = schema.InferPrimaryKey(target, true)
			)

			if pch, exist := assocChanges.Get(pKey); exist && pch.Value != pValues[0] {
				panic("cannot update assoc: inconsistent primary key")
			}

			var (
				filter = where.Eq(assoc.ForeignColumn(), foreignValue)
			)

			if err := r.update(target, assocChanges, filter); err != nil {
				return err
			}
		} else {
			if err := r.insert(target, assocChanges); err != nil {
				return err
			}

			changes.SetValue(assoc.ReferenceColumn(), assoc.ForeignValue())
		}
	}

	return nil
}

func (r Repo) upsertHasOne(assocs schema.Associations, changes *change.Changes, id interface{}) error {
	for _, field := range assocs.HasOne() {
		allAssocChanges, changed := changes.GetAssoc(field)
		if !changed || len(allAssocChanges) == 0 {
			continue
		}

		var (
			assocChanges   = allAssocChanges[0]
			assoc          = assocs.Association(field)
			target, loaded = assoc.TargetAddr()
			referenceValue = assoc.ReferenceValue()
			pKey, pValues  = schema.InferPrimaryKey(target, true)
		)

		if loaded {
			if pch, exist := assocChanges.Get(pKey); exist && pch.Value != pValues[0] {
				panic("cannot update assoc: inconsistent primary key")
			}

			var (
				filter = where.Eq(pKey, pValues[0]).
					AndEq(assoc.ForeignColumn(), referenceValue)
			)

			if err := r.update(target, assocChanges, filter); err != nil {
				return err
			}
		} else {
			assocChanges.SetValue(assoc.ForeignColumn(), referenceValue)

			if err := r.insert(target, assocChanges); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r Repo) upsertHasMany(assocs schema.Associations, changes *change.Changes, id interface{}) error {
	for _, field := range assocs.HasMany() {
		changes, changed := changes.GetAssoc(field)
		if !changed {
			continue
		}

		// Check must be loaded if updated
		// Collect primary keys
		// Delete all based on primary keys
		// Insert all using assoc
		var (
			assoc          = assocs.Association(field)
			target, loaded = assoc.TargetAddr()
			table          = schema.InferTableName(target)
			pKey, pValues  = schema.InferPrimaryKey(target, true)
			referenceValue = assoc.ReferenceValue()
			filter         = where.Eq(assoc.ForeignColumn(), referenceValue).AndIn(pKey, pValues...)
		)

		if loaded && !filter.None() {
			if err := r.deleteAll(query.Build(table, filter)); err != nil {
				return err
			}
		}

		// set assocs
		for i := range changes {
			changes[i].SetValue(assoc.ForeignColumn(), referenceValue)
		}

		if err := r.insertAll(target, changes); err != nil {
			return err
		}

	}

	return nil
}

// Delete deletes all results that match the query.
// TODO: supports array
func (r Repo) Delete(record interface{}) error {
	var (
		table         = schema.InferTableName(record)
		pKey, pValues = schema.InferPrimaryKey(record, true)
		q             = query.Build(table, where.In(pKey, pValues...))
	)

	if len(pValues) == 0 {
		return nil
	}

	return transformError(r.adapter.Delete(q, r.logger...))
}

// MustDelete deletes all results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustDelete(record interface{}) {
	must(r.Delete(record))
}

func (r Repo) DeleteAll(queries ...query.Builder) error {
	var (
		q = query.Build("", queries...)
	)

	return transformError(r.DeleteAll(q))
}

func (r Repo) MustDeleteAll(queries ...query.Builder) {
	must(r.DeleteAll(queries...))
}

func (r Repo) deleteAll(q query.Query) error {
	return r.adapter.Delete(q, r.logger...)
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
	assocField := schema.InferAssociationField(schemaType, path[len(path)-1])

	addrs, ids := collectPreloadTarget(preload, assocField.ReferenceIndex)
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
	err := r.All(result.Interface(), where.In(assocField.ForeignColumn, ids...))
	if err != nil {
		return err
	}

	// map results.
	result = result.Elem()
	for i := 0; i < result.Len(); i++ {
		curr := result.Index(i)
		id := getPreloadID(curr.FieldByIndex(assocField.ForeignIndex))

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
