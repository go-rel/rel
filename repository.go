package rel

import (
	"reflect"
	"runtime"
	"strings"
)

// Repository defines sets of available database operations.
// TODO: InsertAll only accepts records not changes, this eliminates the needs of exposing changes and making the api more consistent.
// TODO: support update all.
type Repository interface {
	Adapter() Adapter
	SetLogger(logger ...Logger)
	Aggregate(query Query, aggregate string, field string) (int, error)
	MustAggregate(query Query, aggregate string, field string) int
	Count(collection string, queriers ...Querier) (int, error)
	MustCount(collection string, queriers ...Querier) int
	Find(record interface{}, queriers ...Querier) error
	MustFind(record interface{}, queriers ...Querier)
	FindAll(records interface{}, queriers ...Querier) error
	MustFindAll(records interface{}, queriers ...Querier)
	Insert(record interface{}, changers ...Changer) error
	MustInsert(record interface{}, changers ...Changer)
	InsertAll(records interface{}, changes ...Changes) error
	MustInsertAll(records interface{}, changes ...Changes)
	Update(record interface{}, changers ...Changer) error
	MustUpdate(record interface{}, changers ...Changer)
	Delete(record interface{}) error
	MustDelete(record interface{})
	DeleteAll(queriers ...Querier) error
	MustDeleteAll(queriers ...Querier)
	Preload(records interface{}, field string, queriers ...Querier) error
	MustPreload(records interface{}, field string, queriers ...Querier)
	Transaction(fn func(Repository) error) error
}

type repository struct {
	adapter       Adapter
	logger        []Logger
	inTransaction bool
}

func (r repository) Adapter() Adapter {
	return r.adapter
}

func (r *repository) SetLogger(logger ...Logger) {
	r.logger = logger
}

// Aggregate calculate aggregate over the given field.
// Supported aggregate: count, sum, avg, max, min.
// Any select, group, offset, limit and sort query will be ignored automatically.
// If complex aggregation is needed, consider using All instead,
func (r repository) Aggregate(query Query, aggregate string, field string) (int, error) {
	query.GroupQuery = GroupQuery{}
	query.LimitQuery = 0
	query.OffsetQuery = 0
	query.SortQuery = nil

	return r.adapter.Aggregate(query, aggregate, field, r.logger...)
}

// MustAggregate calculate aggregate over the given field.
// It'll panic if any error eccured.
func (r repository) MustAggregate(query Query, aggregate string, field string) int {
	result, err := r.Aggregate(query, aggregate, field)
	must(err)
	return result
}

// Count retrieves count of results that match the query.
func (r repository) Count(collection string, queriers ...Querier) (int, error) {
	return r.Aggregate(BuildQuery(collection, queriers...), "count", "*")
}

// MustCount retrieves count of results that match the query.
// It'll panic if any error eccured.
func (r repository) MustCount(collection string, queriers ...Querier) int {
	count, err := r.Count(collection, queriers...)
	must(err)
	return count
}

// Find a record that match the query.
// If no result found, it'll return not found error.
func (r repository) Find(record interface{}, queriers ...Querier) error {
	var (
		doc   = NewDocument(record)
		query = BuildQuery(doc.Table(), queriers...)
	)

	return r.find(doc, query)
}

// MustFind a record that match the query.
// If no result found, it'll panic.
func (r repository) MustFind(record interface{}, queriers ...Querier) {
	must(r.Find(record, queriers...))
}

func (r repository) find(doc *Document, query Query) error {
	cur, err := r.adapter.Query(query.Limit(1), r.logger...)
	if err != nil {
		return err
	}

	return scanOne(cur, doc)
}

// FindAll records that match the query.
func (r repository) FindAll(records interface{}, queriers ...Querier) error {
	var (
		col   = NewCollection(records)
		query = BuildQuery(col.Table(), queriers...)
	)

	col.Reset()

	return r.findAll(col, query)
}

// MustFindAll records that match the query.
// It'll panic if any error eccured.
func (r repository) MustFindAll(records interface{}, queriers ...Querier) {
	must(r.FindAll(records, queriers...))
}

func (r repository) findAll(col *Collection, query Query) error {
	cur, err := r.adapter.Query(query, r.logger...)
	if err != nil {
		return err
	}

	return scanMany(cur, col)
}

// Insert an record to database.
func (r repository) Insert(record interface{}, changers ...Changer) error {
	// TODO: perform reference check on library level for record instead of adapter level
	if record == nil {
		return nil
	}

	var (
		err     error
		changes Changes
		doc     = NewDocument(record)
	)

	if len(changers) == 0 {
		changes, err = ApplyChanges(doc, newStructset(doc))
	} else {
		changes, err = ApplyChanges(doc, changers...)
	}

	if err != nil {
		return err
	}

	if changes.AssocCount() > 0 {
		return r.Transaction(func(r Repository) error {
			return r.(*repository).insert(doc, changes)
		})
	}

	return r.insert(doc, changes)
}

func (r repository) insert(doc *Document, changes Changes) error {
	var (
		pField   = doc.PrimaryField()
		queriers = BuildQuery(doc.Table())
	)

	if err := r.saveBelongsTo(doc, &changes); err != nil {
		return err
	}

	id, err := r.Adapter().Insert(queriers, changes, r.logger...)
	if err != nil {
		return err
	}

	// fetch record
	if err := r.find(doc, queriers.Where(Eq(pField, id))); err != nil {
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
func (r repository) MustInsert(record interface{}, changers ...Changer) {
	must(r.Insert(record, changers...))
}

func (r repository) InsertAll(records interface{}, changes ...Changes) error {
	if records == nil {
		return nil
	}

	var (
		err error
		col = NewCollection(records)
	)

	if len(changes) == 0 {
		changes = make([]Changes, col.Len())
		for i := range changes {
			changes[i], err = ApplyChanges(nil, newStructset(col.Get(i)))
			if err != nil {
				return err
			}
		}
	}

	col.Reset()

	return r.insertAll(col, changes)
}

func (r repository) MustInsertAll(records interface{}, changes ...Changes) {
	must(r.InsertAll(records, changes...))
}

// TODO: support assocs
func (r repository) insertAll(col *Collection, changes []Changes) error {
	if len(changes) == 0 {
		return nil
	}

	var (
		pField   = col.PrimaryField()
		queriers = BuildQuery(col.Table())
		fields   = make([]string, 0, changes[0].Count())
		fieldMap = make(map[string]struct{}, changes[0].Count())
	)

	for i := range changes {
		for _, ch := range changes[i].All() {
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

	return r.findAll(col, queriers.Where(In(pField, ids...)))
}

// Update an record in database.
// It'll panic if any error occurred.
// not supported:
// - update has many (will be replaced by default)
// - replace has one or has many - may cause duplicate record, update instead
func (r repository) Update(record interface{}, changers ...Changer) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: make sure primary id not changed
	if record == nil {
		return nil
	}

	var (
		err     error
		changes Changes
		doc     = NewDocument(record)
		pField  = doc.PrimaryField()
		pValue  = doc.PrimaryValue()
	)

	if len(changers) == 0 {
		changes, err = ApplyChanges(doc, newStructset(doc))
	} else {
		changes, err = ApplyChanges(doc, changers...)
	}

	if err != nil {
		return err
	}

	if len(changes.assoc) > 0 {
		return r.Transaction(func(r Repository) error {
			return r.(*repository).update(doc, changes, Eq(pField, pValue))
		})
	}

	return r.update(doc, changes, Eq(pField, pValue))
}

func (r repository) update(doc *Document, changes Changes, filter FilterQuery) error {
	if err := r.saveBelongsTo(doc, &changes); err != nil {
		return err
	}

	if !changes.Empty() {
		var (
			queriers = BuildQuery(doc.Table(), filter)
		)

		if err := r.adapter.Update(queriers, changes, r.logger...); err != nil {
			return err
		}

		if err := r.find(doc, queriers); err != nil {
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
func (r repository) MustUpdate(record interface{}, changers ...Changer) {
	must(r.Update(record, changers...))
}

// TODO: support deletion
func (r repository) saveBelongsTo(doc *Document, changes *Changes) error {
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
func (r repository) saveHasOne(doc *Document, changes *Changes) error {
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

func (r repository) saveHasMany(doc *Document, changes *Changes, insertion bool) error {
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

// Delete single entry.
func (r repository) Delete(record interface{}) error {
	var (
		doc    = NewDocument(record)
		table  = doc.Table()
		pField = doc.PrimaryField()
		pValue = doc.PrimaryValue()
		q      = BuildQuery(table, Eq(pField, pValue))
	)

	return r.adapter.Delete(q, r.logger...)
}

// MustDelete single entry.
// It'll panic if any error eccured.
func (r repository) MustDelete(record interface{}) {
	must(r.Delete(record))
}

func (r repository) DeleteAll(queriers ...Querier) error {
	var (
		q = BuildQuery("", queriers...)
	)

	return r.deleteAll(q)
}

func (r repository) MustDeleteAll(queriers ...Querier) {
	must(r.DeleteAll(queriers...))
}

func (r repository) deleteAll(q Query) error {
	return r.adapter.Delete(q, r.logger...)
}

// Preload loads association with given query.
func (r repository) Preload(records interface{}, field string, queriers ...Querier) error {
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
		sl = NewCollection(records)
	} else {
		sl = NewDocument(records)
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

// MustPreload loads association with given query.
// It'll panic if any error occurred.
func (r repository) MustPreload(records interface{}, field string, queriers ...Querier) {
	must(r.Preload(records, field, queriers...))
}

func (r repository) mapPreloadTargets(sl slice, path []string) (map[interface{}][]slice, string, string, reflect.Type) {
	type frame struct {
		index int
		doc   *Document
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

// Transaction performs transaction with given function argument.
func (r repository) Transaction(fn func(Repository) error) error {
	adp, err := r.adapter.Begin()
	if err != nil {
		return err
	}

	txRepo := &repository{
		adapter:       adp,
		logger:        []Logger{DefaultLogger},
		inTransaction: true,
	}

	func() {
		defer func() {
			if p := recover(); p != nil {
				_ = txRepo.adapter.Rollback()

				switch e := p.(type) {
				case runtime.Error:
					panic(e)
				case error:
					err = e
				default:
					panic(e)
				}
			} else if err != nil {
				_ = txRepo.adapter.Rollback()
			} else {
				err = txRepo.adapter.Commit()
			}
		}()

		err = fn(txRepo)
	}()

	return err
}

// New create new repo using adapter.
func New(adapter Adapter) Repository {
	return &repository{
		adapter: adapter,
		logger:  []Logger{DefaultLogger},
	}
}
