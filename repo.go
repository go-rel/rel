package grimoire

import (
	"reflect"
	"strings"

	"github.com/Fs02/grimoire/errors"
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
// Supported aggregate: count, sum, avg, max, min.
func (r Repo) Aggregate(query Query, aggregate string, field string) (int, error) {
	return r.adapter.Aggregate(query, aggregate, field, r.logger...)
}

// MustAggregate calculate aggregate over the given field.
// It'll panic if any error eccured.
func (r Repo) MustAggregate(query Query, aggregate string, field string) int {
	result, err := r.Aggregate(query, aggregate, field)
	must(err)
	return result
}

// Count retrieves count of results that match the query.
func (r Repo) Count(collection string, queriers ...Querier) (int, error) {
	return r.Aggregate(BuildQuery(collection, queriers...), "count", "*")
}

// MustCount retrieves count of results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustCount(collection string, queriers ...Querier) int {
	count, err := r.Count(collection, queriers...)
	must(err)
	return count
}

// One retrieves one result that match the query.
// If no result found, it'll return not found error.
func (r Repo) One(entity interface{}, queriers ...Querier) error {
	var (
		doc   = newDocument(entity)
		query = BuildQuery(doc.Table(), queriers...).Limit(1)
	)

	cur, err := r.adapter.Query(query, r.logger...)
	if err != nil {
		return err
	}

	return scanOne(cur, doc)
}

// MustOne retrieves one result that match the query.
// If no result found, it'll panic.
func (r Repo) MustOne(entity interface{}, queriers ...Querier) {
	must(r.One(entity, queriers...))
}

// All retrieves all results that match the query.
func (r Repo) All(entities interface{}, queriers ...Querier) error {
	var (
		col   = newCollection(entities)
		query = BuildQuery(col.Table(), queriers...)
	)

	cur, err := r.adapter.Query(query, r.logger...)
	if err != nil {
		return err
	}

	return scanMany(cur, col)
}

// MustAll retrieves all results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustAll(entities interface{}, queriers ...Querier) {
	must(r.All(entities, queriers...))
}

// Insert a record to database.
// TODO: insert all (multiple changes as multiple records)
func (r Repo) Insert(record interface{}, changers ...Changer) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil || len(changers) == 0 {
		return nil
	}

	// TODO: transform changeset error
	return transformError(r.insert(record, BuildChanges(changers...)))
}

func (r Repo) insert(record interface{}, changes Changes) error {
	var (
		doc      = newDocument(record)
		pField   = doc.PrimaryField()
		queriers = BuildQuery(doc.Table())
	)

	if err := r.upsertBelongsTo(doc, &changes); err != nil {
		return err
	}

	// TODO: put timestamp (updated_at, created_at) if those fields exist.
	id, err := r.Adapter().Insert(queriers, changes, r.logger...)
	if err != nil {
		return err
	}

	// fetch record
	if err := r.One(record, Eq(pField, id)); err != nil {
		return err
	}

	if err := r.upsertHasOne(doc, &changes, id); err != nil {
		return err
	}

	if err := r.upsertHasMany(doc, &changes, id, true); err != nil {
		return err
	}

	return nil
}

// MustInsert a record to database.
// It'll panic if any error occurred.
func (r Repo) MustInsert(record interface{}, changers ...Changer) {
	must(r.Insert(record, changers...))
}

func (r Repo) InsertAll(record interface{}, changes []Changes) error {
	return transformError(r.insertAll(record, changes))
}

func (r Repo) MustInsertAll(record interface{}, changes []Changes) {
	must(r.InsertAll(record, changes))
}

// TODO: support assocs
func (r Repo) insertAll(record interface{}, changes []Changes) error {
	if len(changes) == 0 {
		return nil
	}

	var (
		col      = newCollection(record)
		pField   = col.PrimaryField()
		queriers = BuildQuery(col.Table())
		fields   = make([]string, 0, len(changes[0].Fields))
		fieldMap = make(map[string]struct{}, len(changes[0].Fields))
	)

	for i := range changes {
		for _, ch := range changes[i].Changes {
			if _, exist := fieldMap[ch.Field]; !exist {
				fieldMap[ch.Field] = struct{}{}
				fields = append(fields, ch.Field)
			}
		}
	}

	ids, err := r.adapter.InsertAll(queriers, fields, changes, r.logger...)
	if err != nil {
		return err
	}

	cur, err := r.adapter.Query(queriers.Where(In(pField, ids...)), r.logger...)
	if err != nil {
		return err
	}

	return scanMany(cur, col)
}

// Update a record in database.
// It'll panic if any error occurred.
func (r Repo) Update(record interface{}, changers ...Changer) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil || len(changers) == 0 {
		return nil
	}

	var (
		doc     = newDocument(record)
		pField  = doc.PrimaryField()
		pValue  = doc.PrimaryValue()
		changes = BuildChanges(changers...)
	)

	return r.update(record, changes, Eq(pField, pValue))
}

func (r Repo) update(record interface{}, changes Changes, filter FilterQuery) error {
	if changes.Empty() {
		return nil
	}

	var (
		doc      = newDocument(record)
		queriers = BuildQuery(doc.Table(), filter)
	)

	// TODO: update timestamp (updated_at) from form

	// perform update
	err := r.adapter.Update(queriers, changes, r.logger...)
	if err != nil {
		// TODO: changeset error
		return transformError(err)
	}

	return r.One(record, queriers)
}

// MustUpdate a record in database.
// It'll panic if any error occurred.
func (r Repo) MustUpdate(record interface{}, changers ...Changer) {
	must(r.Update(record, changers...))
}

func (r Repo) upsertBelongsTo(doc Document, changes *Changes) error {
	for _, field := range doc.BelongsTo() {
		allAssocChanges, changed := changes.GetAssoc(field)
		if !changed || len(allAssocChanges) == 0 {
			continue
		}

		var (
			assocChanges   = allAssocChanges[0]
			assoc          = doc.Association(field)
			fValue         = assoc.ForeignValue()
			target, loaded = assoc.Target()
			doc            = target.(Document)
		)

		if loaded {
			var (
				pField = doc.PrimaryField()
				pValue = doc.PrimaryValue()
			)

			if pch, exist := assocChanges.Get(pField); exist && pch.Value != pValue {
				panic("cannot update assoc: inconsistent primary value")
			}

			var (
				filter = Eq(assoc.ForeignField(), fValue)
			)

			if err := r.update(doc, assocChanges, filter); err != nil {
				return err
			}
		} else {
			if err := r.insert(doc, assocChanges); err != nil {
				return err
			}

			changes.SetValue(assoc.ReferenceField(), assoc.ForeignValue())
		}
	}

	return nil
}

func (r Repo) upsertHasOne(doc Document, changes *Changes, id interface{}) error {
	for _, field := range doc.HasOne() {
		allAssocChanges, changed := changes.GetAssoc(field)
		if !changed || len(allAssocChanges) == 0 {
			continue
		}

		var (
			assocChanges   = allAssocChanges[0]
			assoc          = doc.Association(field)
			fField         = assoc.ForeignField()
			rValue         = assoc.ReferenceValue()
			target, loaded = assoc.Target()
			doc            = target.(Document)
			pField         = doc.PrimaryField()
			pValue         = doc.PrimaryValue()
		)

		if loaded {
			if pch, exist := assocChanges.Get(pField); exist && pch.Value != pValue {
				panic("cannot update assoc: inconsistent primary key")
			}

			var (
				filter = Eq(pField, pValue).AndEq(fField, rValue)
			)

			if err := r.update(target, assocChanges, filter); err != nil {
				return err
			}
		} else {
			assocChanges.SetValue(fField, rValue)

			if err := r.insert(target, assocChanges); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r Repo) upsertHasMany(doc Document, changes *Changes, id interface{}, insertion bool) error {
	for _, field := range doc.HasMany() {
		changes, changed := changes.GetAssoc(field)
		if !changed {
			continue
		}

		var (
			assoc          = doc.Association(field)
			target, loaded = assoc.Target()
			table          = target.Table()
			fField         = assoc.ForeignField()
			rValue         = assoc.ReferenceValue()
		)

		if !insertion {
			if !loaded {
				panic("grimoire: association must be loaded to update")
			}

			var (
				pField  = target.PrimaryField()
				pValues = target.PrimaryValue().([]interface{})
			)

			if len(pValues) > 0 {
				var (
					filter = Eq(fField, rValue).AndIn(pField, pValues...)
				)

				if err := r.deleteAll(BuildQuery(table, filter)); err != nil {
					return err
				}
			}
		}

		// set assocs
		for i := range changes {
			changes[i].SetValue(fField, rValue)
		}

		if err := r.insertAll(target, changes); err != nil {
			return err
		}

	}

	return nil
}

// Delete single entry.
func (r Repo) Delete(entity interface{}) error {
	var (
		doc    = newDocument(entity)
		table  = doc.Table()
		pField = doc.PrimaryField()
		pValue = doc.PrimaryValue()
		q      = BuildQuery(table, Eq(pField, pValue))
	)

	return transformError(r.adapter.Delete(q, r.logger...))
}

// MustDelete single entry.
// It'll panic if any error eccured.
func (r Repo) MustDelete(record interface{}) {
	must(r.Delete(record))
}

func (r Repo) DeleteAll(queriers ...Querier) error {
	var (
		q = BuildQuery("", queriers...)
	)

	return transformError(r.deleteAll(q))
}

func (r Repo) MustDeleteAll(queriers ...Querier) {
	must(r.DeleteAll(queriers...))
}

func (r Repo) deleteAll(q Query) error {
	return r.adapter.Delete(q, r.logger...)
}

// Preload loads association with given query.
func (r Repo) Preload(entities interface{}, field string, queriers ...Querier) error {
	var (
		col  Collection
		path = strings.Split(field, ".")
		rt   = reflect.TypeOf(entities)
	)

	if rt.Kind() != reflect.Ptr {
		panic("grimoire: record parameter must be a pointer.")
	}

	rt = rt.Elem()
	if rt.Kind() == reflect.Slice {
		col = newCollection(entities)
	} else {
		col = newDocument(entities)
	}

	var (
		targets, table, keyField, keyType = r.mapPreloadTargets(col, path)
	)

	if len(targets) == 0 {
		return nil
	}

	var (
		ids = make([]interface{}, len(targets))
		i   = 0
	)

	for key := range targets {
		ids[i] = key
		i++
	}

	cur, err := r.adapter.Query(BuildQuery(table, In(keyField, ids...)), r.logger...)
	if err != nil {
		return err
	}

	return scanMulti(cur, keyField, keyType, targets)
}

func (r Repo) mapPreloadTargets(col Collection, path []string) (map[interface{}][]Collection, string, string, reflect.Type) {
	type frame struct {
		index int
		doc   Document
	}

	var (
		table     string
		keyField  string
		keyType   reflect.Type
		mapTarget = make(map[interface{}][]Collection)
		stack     = make([]frame, col.Len())
	)

	// init stack
	for i := 0; i < len(stack); i++ {
		stack[i] = frame{index: 0, doc: col.Get(i)}
	}

	for len(stack) > 0 {
		var (
			n      = len(stack) - 1
			top    = stack[n]
			assocs = top.doc.Association(path[top.index])
		)

		stack = stack[:n]

		if top.index == len(path)-1 {
			var (
				ref = assocs.ReferenceValue()
			)

			if ref == nil {
				continue
			}

			var (
				target, _ = assocs.Target()
			)

			target.Reset()
			mapTarget[ref] = append(mapTarget[ref], target)

			if table == "" {
				table = target.Table()
				keyField = assocs.ForeignField()
				keyType = reflect.TypeOf(ref)
			}
		} else {
			var (
				target, loaded = assocs.Target()
			)

			if !loaded {
				continue
			}

			stack = append(stack, make([]frame, target.Len())...)
			for i := 0; i < target.Len(); i++ {
				stack[n+i] = frame{
					index: top.index + 1,
					doc:   target.Get(i),
				}
			}
		}

	}

	return mapTarget, table, keyField, keyType
}

// MustPreload loads association with given query.
// It'll panic if any error occurred.
func (r Repo) MustPreload(record interface{}, field string, queriers ...Querier) {
	must(r.Preload(record, field, queriers...))
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
