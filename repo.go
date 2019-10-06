package rel

import (
	"reflect"
	"runtime"
	"strings"
	"time"
)

var (
	now = time.Now
)

// Repo defines rel repository.
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
// Any select, group, offset, limit and sort query will be ignored automatically.
// If complex aggregation is needed, consider using All instead,
func (r Repo) Aggregate(query Query, aggregate string, field string) (int, error) {
	query.GroupQuery = GroupQuery{}
	query.LimitQuery = 0
	query.OffsetQuery = 0
	query.SortQuery = nil

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
func (r Repo) One(record interface{}, queriers ...Querier) error {
	var (
		doc   = newDocument(record)
		query = BuildQuery(doc.Table(), queriers...)
	)

	return r.one(doc, query)
}

// MustOne retrieves one result that match the query.
// If no result found, it'll panic.
func (r Repo) MustOne(record interface{}, queriers ...Querier) {
	must(r.One(record, queriers...))
}

func (r Repo) one(doc *document, query Query) error {
	cur, err := r.adapter.Query(query.Limit(1), r.logger...)
	if err != nil {
		return err
	}

	return scanOne(cur, doc)
}

// All retrieves all results that match the query.
func (r Repo) All(records interface{}, queriers ...Querier) error {
	var (
		col   = newCollection(records)
		query = BuildQuery(col.Table(), queriers...)
	)

	col.Reset()

	return r.all(col, query)
}

// MustAll retrieves all results that match the query.
// It'll panic if any error eccured.
func (r Repo) MustAll(records interface{}, queriers ...Querier) {
	must(r.All(records, queriers...))
}

func (r Repo) all(col *collection, query Query) error {
	cur, err := r.adapter.Query(query, r.logger...)
	if err != nil {
		return err
	}

	return scanMany(cur, col)
}

// Insert an record to database.
// TODO: insert all (multiple changes as multiple records)
func (r Repo) Insert(record interface{}, changers ...Changer) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	if record == nil {
		return nil
	}

	var (
		changes Changes
		doc     = newDocument(record)
	)

	if len(changers) == 0 {
		changes = BuildChanges(newStructset(doc))
	} else {
		changes = BuildChanges(changers...)
	}

	if len(changes.Assoc) > 0 {
		return r.Transaction(func(r Repo) error {
			return r.insert(doc, changes)
		})
	}

	return r.insert(doc, changes)
}

func (r Repo) insert(doc *document, changes Changes) error {
	var (
		t        = now()
		pField   = doc.PrimaryField()
		queriers = BuildQuery(doc.Table())
	)

	if err := r.saveBelongsTo(doc, &changes); err != nil {
		return err
	}

	r.putTimestamp(doc, &changes, "created_at", t)
	r.putTimestamp(doc, &changes, "updated_at", t)

	id, err := r.Adapter().Insert(queriers, changes, r.logger...)
	if err != nil {
		return err
	}

	// fetch record
	if err := r.one(doc, queriers.Where(Eq(pField, id))); err != nil {
		return err
	}

	if err := r.saveHasOne(doc, &changes); err != nil {
		return err
	}

	if err := r.saveHasMany(doc, &changes, true); err != nil {
		return err
	}

	return nil
}

// MustInsert an record to database.
// It'll panic if any error occurred.
func (r Repo) MustInsert(record interface{}, changers ...Changer) {
	must(r.Insert(record, changers...))
}

func (r Repo) InsertAll(records interface{}, changes ...Changes) error {
	if records == nil {
		return nil
	}

	var (
		col = newCollection(records)
	)

	if len(changes) == 0 {
		changes = make([]Changes, col.Len())
		for i := range changes {
			changes[i] = BuildChanges(newStructset(col.Get(i)))
		}
	}

	col.Reset()

	return r.insertAll(col, changes)
}

func (r Repo) MustInsertAll(records interface{}, changes ...Changes) {
	must(r.InsertAll(records, changes...))
}

// TODO: support assocs
func (r Repo) insertAll(col *collection, changes []Changes) error {
	if len(changes) == 0 {
		return nil
	}

	var (
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

	return r.all(col, queriers.Where(In(pField, ids...)))
}

// Update an record in database.
// It'll panic if any error occurred.
// not supported:
// - update has many (will be replaced by default)
// - replace has one or has many - may cause duplicate record, update instead
func (r Repo) Update(record interface{}, changers ...Changer) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: support not returning via changeset table inference
	// TODO: make sure primary id not changed
	if record == nil {
		return nil
	}

	var (
		changes Changes
		doc     = newDocument(record)
		pField  = doc.PrimaryField()
		pValue  = doc.PrimaryValue()
	)

	if len(changers) == 0 {
		changes = BuildChanges(newStructset(doc))
	} else {
		changes = BuildChanges(changers...)
	}

	if len(changes.Assoc) > 0 {
		return r.Transaction(func(r Repo) error {
			return r.update(doc, changes, Eq(pField, pValue))
		})
	}

	return r.update(doc, changes, Eq(pField, pValue))
}

func (r Repo) update(doc *document, changes Changes, filter FilterQuery) error {
	if err := r.saveBelongsTo(doc, &changes); err != nil {
		return err
	}

	if !changes.Empty() {
		var (
			queriers = BuildQuery(doc.Table(), filter)
		)

		r.putTimestamp(doc, &changes, "updated_at", now())

		if err := r.adapter.Update(queriers, changes, r.logger...); err != nil {
			return err
		}

		if err := r.one(doc, queriers); err != nil {
			return err
		}
	}

	if err := r.saveHasOne(doc, &changes); err != nil {
		return err
	}

	if err := r.saveHasMany(doc, &changes, false); err != nil {
		return err
	}

	return nil
}

// MustUpdate an record in database.
// It'll panic if any error occurred.
func (r Repo) MustUpdate(record interface{}, changers ...Changer) {
	must(r.Update(record, changers...))
}

// TODO: support deletion
func (r Repo) saveBelongsTo(doc *document, changes *Changes) error {
	for _, field := range doc.BelongsTo() {
		ac, changed := changes.GetAssoc(field)
		if !changed || len(ac.Changes) == 0 {
			continue
		}

		var (
			assocChanges = ac.Changes[0]
			assoc        = doc.Association(field)
			fValue       = assoc.ForeignValue()
			doc, loaded  = assoc.Document()
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

// TODO: suppprt deletion
func (r Repo) saveHasOne(doc *document, changes *Changes) error {
	for _, field := range doc.HasOne() {
		ac, changed := changes.GetAssoc(field)
		if !changed || len(ac.Changes) == 0 {
			continue
		}

		var (
			assocChanges = ac.Changes[0]
			assoc        = doc.Association(field)
			fField       = assoc.ForeignField()
			rValue       = assoc.ReferenceValue()
			doc, loaded  = assoc.Document()
			pField       = doc.PrimaryField()
			pValue       = doc.PrimaryValue()
		)

		if loaded {
			if pch, exist := assocChanges.Get(pField); exist && pch.Value != pValue {
				panic("cannot update assoc: inconsistent primary key")
			}

			var (
				filter = Eq(pField, pValue).AndEq(fField, rValue)
			)

			if err := r.update(doc, assocChanges, filter); err != nil {
				return err
			}
		} else {
			assocChanges.SetValue(fField, rValue)

			if err := r.insert(doc, assocChanges); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r Repo) saveHasMany(doc *document, changes *Changes, insertion bool) error {
	for _, field := range doc.HasMany() {
		ac, changed := changes.GetAssoc(field)
		if !changed {
			continue
		}

		var (
			assoc       = doc.Association(field)
			col, loaded = assoc.Collection()
			table       = col.Table()
			pField      = col.PrimaryField()
			fField      = assoc.ForeignField()
			rValue      = assoc.ReferenceValue()
		)

		col.Reset()

		if !insertion {
			if !loaded {
				panic("rel: association must be loaded to update")
			}

			var (
				filter = Eq(fField, rValue)
			)

			// if deleted ids is specified, then only delete those.
			// if it's nill, then clear old association (used by structset).
			if len(ac.StaleIDs) > 0 {
				if err := r.deleteAll(BuildQuery(table, filter.AndIn(pField, ac.StaleIDs...))); err != nil {
					return err
				}
			} else if ac.StaleIDs == nil {
				if err := r.deleteAll(BuildQuery(table, filter)); err != nil {
					return err
				}
			}
		}

		// update and filter for bulk insertion in place
		// TODO: load updated result once
		n := 0
		for _, ch := range ac.Changes {
			if pChange, changed := ch.Get(pField); changed {
				var (
					filter = Eq(pField, pChange.Value).AndEq(fField, rValue)
				)

				if err := r.update(col.Add(), ch, filter); err != nil {
					return err
				}
			} else {
				ch.SetValue(fField, rValue)
				ac.Changes[n] = ch
				n++
			}
		}
		ac.Changes = ac.Changes[:n]

		if err := r.insertAll(col, ac.Changes); err != nil {
			return err
		}
	}

	return nil
}

func (r Repo) Save(record interface{}, changers ...Changer) error {
	if record == nil {
		return nil
	}

	var (
		doc = newDocument(record)
	)

	if len(changers) == 0 {
		changers = []Changer{newStructset(doc)}
	}

	return r.save(doc, BuildChanges(changers...))
}

func (r Repo) save(doc *document, changes Changes) error {
	var (
		pField = doc.PrimaryField()
		pValue = doc.PrimaryValue()
	)

	if isZero(pValue) {
		return r.insert(doc, changes)
	}

	return r.update(doc, changes, Eq(pField, pValue))
}

// Delete single entry.
func (r Repo) Delete(record interface{}) error {
	var (
		doc    = newDocument(record)
		table  = doc.Table()
		pField = doc.PrimaryField()
		pValue = doc.PrimaryValue()
		q      = BuildQuery(table, Eq(pField, pValue))
	)

	return r.adapter.Delete(q, r.logger...)
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

	return r.deleteAll(q)
}

func (r Repo) MustDeleteAll(queriers ...Querier) {
	must(r.DeleteAll(queriers...))
}

func (r Repo) deleteAll(q Query) error {
	return r.adapter.Delete(q, r.logger...)
}

// Preload loads association with given query.
func (r Repo) Preload(records interface{}, field string, queriers ...Querier) error {
	var (
		sl   slice
		path = strings.Split(field, ".")
		rt   = reflect.TypeOf(records)
	)

	if rt.Kind() != reflect.Ptr {
		panic("rel: record parameter must be a pointer.")
	}

	rt = rt.Elem()
	if rt.Kind() == reflect.Slice {
		sl = newCollection(records)
	} else {
		sl = newDocument(records)
	}

	var (
		targets, table, keyField, keyType = r.mapPreloadTargets(sl, path)
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

	var (
		query    = BuildQuery(table, append(queriers, In(keyField, ids...))...)
		cur, err = r.adapter.Query(query, r.logger...)
	)

	if err != nil {
		return err
	}

	return scanMulti(cur, keyField, keyType, targets)
}

func (r Repo) mapPreloadTargets(sl slice, path []string) (map[interface{}][]slice, string, string, reflect.Type) {
	type frame struct {
		index int
		doc   *document
	}

	var (
		table     string
		keyField  string
		keyType   reflect.Type
		mapTarget = make(map[interface{}][]slice)
		stack     = make([]frame, sl.Len())
	)

	// init stack
	for i := 0; i < len(stack); i++ {
		stack[i] = frame{index: 0, doc: sl.Get(i)}
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
				target slice
				ref    = assocs.ReferenceValue()
			)

			if ref == nil {
				continue
			}

			if assocs.Type() == HasMany {
				target, _ = assocs.Collection()
			} else {
				target, _ = assocs.Document()
			}

			target.Reset()
			mapTarget[ref] = append(mapTarget[ref], target)

			if table == "" {
				table = target.Table()
				keyField = assocs.ForeignField()
				keyType = reflect.TypeOf(ref)
			}
		} else {
			if assocs.Type() == HasMany {
				var (
					col, loaded = assocs.Collection()
				)

				if !loaded {
					continue
				}

				stack = append(stack, make([]frame, col.Len())...)
				for i := 0; i < col.Len(); i++ {
					stack[n+i] = frame{
						index: top.index + 1,
						doc:   col.Get(i),
					}
				}
			} else {
				if doc, loaded := assocs.Document(); loaded {
					stack = append(stack, frame{
						index: top.index + 1,
						doc:   doc,
					})
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

				switch e := p.(type) {
				case runtime.Error:
					panic(e)
				case error:
					err = e
				default:
					panic(e)
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

func (r Repo) putTimestamp(doc *document, changes *Changes, field string, t time.Time) {
	if typ, ok := doc.Type(field); ok && typ == rtTime {
		changes.SetValue(field, t)
	}
}
