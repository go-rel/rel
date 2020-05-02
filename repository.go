package rel

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"strings"
)

// Repository defines sets of available database operations.
// TODO: support update all.
type Repository interface {
	Adapter() Adapter
	Instrumentation(instrumenter Instrumenter)
	Ping(ctx context.Context) error
	Iterate(ctx context.Context, query Query, option ...IteratorOption) Iterator
	Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error)
	MustAggregate(ctx context.Context, query Query, aggregate string, field string) int
	Count(ctx context.Context, collection string, queriers ...Querier) (int, error)
	MustCount(ctx context.Context, collection string, queriers ...Querier) int
	Find(ctx context.Context, record interface{}, queriers ...Querier) error
	MustFind(ctx context.Context, record interface{}, queriers ...Querier)
	FindAll(ctx context.Context, records interface{}, queriers ...Querier) error
	MustFindAll(ctx context.Context, records interface{}, queriers ...Querier)
	FindAndCountAll(ctx context.Context, records interface{}, queriers ...Querier) (int, error)
	MustFindAndCountAll(ctx context.Context, records interface{}, queriers ...Querier) int
	Insert(ctx context.Context, record interface{}, modifiers ...Modifier) error
	MustInsert(ctx context.Context, record interface{}, modifiers ...Modifier)
	InsertAll(ctx context.Context, records interface{}) error
	MustInsertAll(ctx context.Context, records interface{})
	Update(ctx context.Context, record interface{}, modifiers ...Modifier) error
	MustUpdate(ctx context.Context, record interface{}, modifiers ...Modifier)
	Delete(ctx context.Context, record interface{}) error
	MustDelete(ctx context.Context, record interface{})
	DeleteAll(ctx context.Context, query Query) error
	MustDeleteAll(ctx context.Context, query Query)
	Preload(ctx context.Context, records interface{}, field string, queriers ...Querier) error
	MustPreload(ctx context.Context, records interface{}, field string, queriers ...Querier)
	Transaction(ctx context.Context, fn func(Repository) error) error
}

type repository struct {
	adapter       Adapter
	instrumenter  Instrumenter
	inTransaction bool
}

func (r repository) Adapter() Adapter {
	return r.adapter
}

func (r *repository) Instrumentation(instrumenter Instrumenter) {
	r.instrumenter = instrumenter
	r.adapter.Instrumentation(instrumenter)
}

// Instrument call instrumenter, if no instrumenter is set, this will be a no op.
func (r *repository) instrument(ctx context.Context, op string, message string) func(err error) {
	if r.instrumenter != nil {
		return r.instrumenter(ctx, op, message)
	}

	return func(err error) {}
}

// Ping database.
func (r *repository) Ping(ctx context.Context) error {
	return r.adapter.Ping(ctx)
}

// Iterate through a collection of records from database in batches.
// This function returns iterator that can be used to loop all records.
// Limit, Offset and Sort query is automatically ignored.
func (r repository) Iterate(ctx context.Context, query Query, options ...IteratorOption) Iterator {
	return newIterator(ctx, r.adapter, query, options)
}

// Aggregate calculate aggregate over the given field.
// Supported aggregate: count, sum, avg, max, min.
// Any select, group, offset, limit and sort query will be ignored automatically.
// If complex aggregation is needed, consider using All instead,
func (r repository) Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	finish := r.instrument(ctx, "rel-aggregate", "aggregating records")
	defer finish(nil)

	return r.aggregate(ctx, query, aggregate, field)
}

func (r repository) aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	query.GroupQuery = GroupQuery{}
	query.LimitQuery = 0
	query.OffsetQuery = 0
	query.SortQuery = nil

	return r.adapter.Aggregate(ctx, query, aggregate, field)
}

// MustAggregate calculate aggregate over the given field.
// It'll panic if any error eccured.
func (r repository) MustAggregate(ctx context.Context, query Query, aggregate string, field string) int {
	result, err := r.Aggregate(ctx, query, aggregate, field)
	must(err)
	return result
}

// Count retrieves count of results that match the query.
func (r repository) Count(ctx context.Context, collection string, queriers ...Querier) (int, error) {
	finish := r.instrument(ctx, "rel-count", "aggregating records")
	defer finish(nil)

	return r.aggregate(ctx, Build(collection, queriers...), "count", "*")
}

// MustCount retrieves count of results that match the query.
// It'll panic if any error eccured.
func (r repository) MustCount(ctx context.Context, collection string, queriers ...Querier) int {
	count, err := r.Count(ctx, collection, queriers...)
	must(err)
	return count
}

// Find a record that match the query.
// If no result found, it'll return not found error.
func (r repository) Find(ctx context.Context, record interface{}, queriers ...Querier) error {
	finish := r.instrument(ctx, "rel-find", "finding a record")
	defer finish(nil)

	var (
		doc   = NewDocument(record)
		query = Build(doc.Table(), queriers...)
	)

	return r.find(ctx, doc, query)
}

// MustFind a record that match the query.
// If no result found, it'll panic.
func (r repository) MustFind(ctx context.Context, record interface{}, queriers ...Querier) {
	must(r.Find(ctx, record, queriers...))
}

func (r repository) find(ctx context.Context, doc *Document, query Query) error {
	query = r.withDefaultScope(doc.data, query)
	cur, err := r.adapter.Query(ctx, query.Limit(1))
	if err != nil {
		return err
	}

	finish := r.instrument(ctx, "rel-scan-one", "scanning a record")
	if err := scanOne(cur, doc); err != nil {
		finish(err)
		return err
	}
	finish(nil)

	// init dirty.
	if doc.Dirty() != nil {
		doc.Dirty().init(doc)
	}

	return nil
}

// FindAll records that match the query.
func (r repository) FindAll(ctx context.Context, records interface{}, queriers ...Querier) error {
	finish := r.instrument(ctx, "rel-find-all", "finding all records")
	defer finish(nil)

	var (
		col   = NewCollection(records)
		query = Build(col.Table(), queriers...)
	)

	col.Reset()

	return r.findAll(ctx, col, query)
}

// MustFindAll records that match the query.
// It'll panic if any error eccured.
func (r repository) MustFindAll(ctx context.Context, records interface{}, queriers ...Querier) {
	must(r.FindAll(ctx, records, queriers...))
}

func (r repository) findAll(ctx context.Context, col *Collection, query Query) error {
	query = r.withDefaultScope(col.data, query)
	cur, err := r.adapter.Query(ctx, query)
	if err != nil {
		return err
	}

	finish := r.instrument(ctx, "rel-scan-all", "scanning all records")
	if err := scanAll(cur, col); err != nil {
		finish(err)
		return err
	}
	finish(nil)

	// init dirties.
	for i := 0; i < col.Len(); i++ {
		if doc := col.Get(i); doc.Dirty() != nil {
			defer doc.Dirty().init(doc)
		}
	}

	return nil
}

// FindAndCountAll is convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
// Limit and Offset property will be ignored when performing count query.
func (r repository) FindAndCountAll(ctx context.Context, records interface{}, queriers ...Querier) (int, error) {
	finish := r.instrument(ctx, "rel-find-and-count-all", "finding all records")
	defer finish(nil)

	var (
		col   = NewCollection(records)
		query = Build(col.Table(), queriers...)
	)

	col.Reset()

	if err := r.findAll(ctx, col, query); err != nil {
		return 0, err
	}

	return r.aggregate(ctx, query, "count", "*")
}

// MustFindAndCountAll is convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
// Limit and Offset property will be ignored when performing count query.
// It'll panic if any error eccured.
func (r repository) MustFindAndCountAll(ctx context.Context, records interface{}, queriers ...Querier) int {
	count, err := r.FindAndCountAll(ctx, records, queriers...)
	must(err)

	return count
}

// Insert an record to database.
func (r repository) Insert(ctx context.Context, record interface{}, modifiers ...Modifier) error {
	finish := r.instrument(ctx, "rel-insert", "inserting a record")
	defer finish(nil)

	if record == nil {
		return nil
	}

	var (
		modification Modification
		doc          = NewDocument(record)
	)

	if len(modifiers) == 0 {
		modification = Apply(doc, newStructset(doc, false))
	} else {
		modification = Apply(doc, modifiers...)
	}

	if len(modification.Assoc) > 0 {
		return r.Transaction(ctx, func(r Repository) error {
			return r.(*repository).insert(ctx, doc, modification)
		})
	}

	return r.insert(ctx, doc, modification)
}

func (r repository) insert(ctx context.Context, doc *Document, modification Modification) error {
	var (
		pField   = doc.PrimaryField()
		queriers = Build(doc.Table())
	)

	if err := r.saveBelongsTo(ctx, doc, &modification); err != nil {
		return err
	}

	pValue, err := r.Adapter().Insert(ctx, queriers, modification.Modifies)
	if err != nil {
		return err
	}

	if modification.Reload {
		// fetch record
		if err := r.find(ctx, doc, queriers.Where(Eq(pField, pValue))); err != nil {
			return err
		}
	} else {
		// update primary value
		doc.SetValue(pField, pValue)
	}

	if err := r.saveHasOne(ctx, doc, &modification); err != nil {
		return err
	}

	if err := r.saveHasMany(ctx, doc, &modification, true); err != nil {
		return err
	}

	return nil
}

// MustInsert an record to database.
// It'll panic if any error occurred.
func (r repository) MustInsert(ctx context.Context, record interface{}, modifiers ...Modifier) {
	must(r.Insert(ctx, record, modifiers...))
}

func (r repository) InsertAll(ctx context.Context, records interface{}) error {
	finish := r.instrument(ctx, "rel-insert-all", "inserting multiple records")
	defer finish(nil)

	if records == nil {
		return nil
	}

	var (
		col  = NewCollection(records)
		mods = make([]Modification, col.Len())
	)

	for i := range mods {
		doc := col.Get(i)
		mods[i] = Apply(doc, newStructset(doc, false))
	}

	return r.insertAll(ctx, col, mods)
}

func (r repository) MustInsertAll(ctx context.Context, records interface{}) {
	must(r.InsertAll(ctx, records))
}

// TODO: support assocs
func (r repository) insertAll(ctx context.Context, col *Collection, modification []Modification) error {
	if len(modification) == 0 {
		return nil
	}

	var (
		pField       = col.PrimaryField()
		queriers     = Build(col.Table())
		fields       = make([]string, 0, len(modification[0].Modifies))
		fieldMap     = make(map[string]struct{}, len(modification[0].Modifies))
		bulkModifies = make([]map[string]Modify, len(modification))
	)

	// TODO: baypassable if it's predictable.
	for i := range modification {
		for field := range modification[i].Modifies {
			if _, exist := fieldMap[field]; !exist {
				fieldMap[field] = struct{}{}
				fields = append(fields, field)
			}
		}
		bulkModifies[i] = modification[i].Modifies
	}

	ids, err := r.adapter.InsertAll(ctx, queriers, fields, bulkModifies)
	if err != nil {
		return err
	}

	// apply ids
	for i, id := range ids {
		col.Get(i).SetValue(pField, id)
	}

	return nil
}

// Update an record in database.
// It'll panic if any error occurred.
// not supported:
// - update has many (will be replaced by default)
// - replacing has one or belongs to assoc may cause duplicate record, please ensure database level unique constraint enabled.
func (r repository) Update(ctx context.Context, record interface{}, modifiers ...Modifier) error {
	finish := r.instrument(ctx, "rel-update", "updating a record")
	defer finish(nil)

	if record == nil {
		return nil
	}

	var (
		modification Modification
		doc          = NewDocument(record)
		pField       = doc.PrimaryField()
		pValue       = doc.PrimaryValue()
	)

	if len(modifiers) == 0 {
		modification = Apply(doc, newStructset(doc, false))
	} else {
		modification = Apply(doc, modifiers...)
	}

	if len(modification.Assoc) > 0 {
		return r.Transaction(ctx, func(r Repository) error {
			return r.(*repository).update(ctx, doc, modification, Eq(pField, pValue))
		})
	}

	return r.update(ctx, doc, modification, Eq(pField, pValue))
}

func (r repository) update(ctx context.Context, doc *Document, modification Modification, filter FilterQuery) error {
	if err := r.saveBelongsTo(ctx, doc, &modification); err != nil {
		return err
	}

	if len(modification.Modifies) != 0 {
		var (
			query             = r.withDefaultScope(doc.data, Build(doc.Table(), filter, modification.Unscoped))
			updatedCount, err = r.adapter.Update(ctx, query, modification.Modifies)
		)

		if err != nil {
			return err
		}

		if updatedCount == 0 {
			return NotFoundError{}
		}

		if modification.Reload {
			if err := r.find(ctx, doc, query); err != nil {
				return err
			}
		}
	}

	if err := r.saveHasOne(ctx, doc, &modification); err != nil {
		return err
	}

	if err := r.saveHasMany(ctx, doc, &modification, false); err != nil {
		return err
	}

	return nil
}

// MustUpdate an record in database.
// It'll panic if any error occurred.
func (r repository) MustUpdate(ctx context.Context, record interface{}, modifiers ...Modifier) {
	must(r.Update(ctx, record, modifiers...))
}

// TODO: support deletion
func (r repository) saveBelongsTo(ctx context.Context, doc *Document, modification *Modification) error {
	for _, field := range doc.BelongsTo() {
		assocMods, changed := modification.Assoc[field]
		if !changed || len(assocMods.Modifications) == 0 {
			continue
		}

		var (
			assoc            = doc.Association(field)
			assocDoc, loaded = assoc.Document()
			assocMod         = assocMods.Modifications[0]
		)

		if loaded {
			var (
				fValue = assoc.ForeignValue()
			)

			if assoc.ReferenceValue() != fValue {
				return ConstraintError{
					Key:  assoc.ReferenceField(),
					Type: ForeignKeyConstraint,
					Err:  errors.New("rel: inconsistent belongs to ref and fk"),
				}
			}

			var (
				filter = Eq(assoc.ForeignField(), fValue)
			)

			if err := r.update(ctx, assocDoc, assocMod, filter); err != nil {
				return err
			}
		} else {
			if err := r.insert(ctx, assocDoc, assocMod); err != nil {
				return err
			}

			var (
				rField = assoc.ReferenceField()
				fValue = assoc.ForeignValue()
			)

			modification.Add(Set(rField, fValue))
			doc.SetValue(rField, fValue)
		}
	}

	return nil
}

// TODO: suppprt deletion
func (r repository) saveHasOne(ctx context.Context, doc *Document, modification *Modification) error {
	for _, field := range doc.HasOne() {
		assocMods, changed := modification.Assoc[field]
		if !changed || len(assocMods.Modifications) == 0 {
			continue
		}

		var (
			assoc            = doc.Association(field)
			fField           = assoc.ForeignField()
			rValue           = assoc.ReferenceValue()
			assocDoc, loaded = assoc.Document()
			pField           = assocDoc.PrimaryField()
			pValue           = assocDoc.PrimaryValue()
			assocMod         = assocMods.Modifications[0]
		)

		if loaded {
			if rValue != assoc.ForeignValue() {
				return ConstraintError{
					Key:  fField,
					Type: ForeignKeyConstraint,
					Err:  errors.New("rel: inconsistent has one ref and fk"),
				}
			}

			var (
				filter = Eq(pField, pValue).AndEq(fField, rValue)
			)

			if err := r.update(ctx, assocDoc, assocMod, filter); err != nil {
				return err
			}
		} else {
			assocMod.Add(Set(fField, rValue))

			if err := r.insert(ctx, assocDoc, assocMod); err != nil {
				return err
			}
		}

		assocDoc.SetValue(fField, rValue)
	}

	return nil
}

// saveHasMany expects has many modification to be ordered the same as the recrods in collection.
func (r repository) saveHasMany(ctx context.Context, doc *Document, modification *Modification, insertion bool) error {
	for _, field := range doc.HasMany() {
		assocMods, changed := modification.Assoc[field]
		if !changed {
			continue
		}

		var (
			assoc      = doc.Association(field)
			col, _     = assoc.Collection()
			table      = col.Table()
			pField     = col.PrimaryField()
			fField     = assoc.ForeignField()
			rValue     = assoc.ReferenceValue()
			mods       = assocMods.Modifications
			deletedIDs = assocMods.DeletedIDs
		)

		// this shouldn't happen unless there's bug in the modifier.
		if len(mods) != col.Len() {
			panic("rel: invalid modifier")
		}

		if !insertion {
			var (
				filter = Eq(fField, rValue)
			)

			if deletedIDs == nil {
				// if it's nil, then clear old association (used by structset).
				if err := r.deleteAll(ctx, col.data.flag, Build(table, filter)); err != nil {
					return err
				}
			} else if len(deletedIDs) > 0 {
				if err := r.deleteAll(ctx, col.data.flag, Build(table, filter.AndIn(pField, deletedIDs...))); err != nil {
					return err
				}
			}
		}

		// update and filter for bulk insertion.
		updateCount := 0
		for i := range mods {
			var (
				assocDoc = col.Get(i)
				pValue   = assocDoc.PrimaryValue()
			)

			if !isZero(pValue) {
				var (
					fValue, _ = assocDoc.Value(fField)
					filter    = Eq(pField, pValue).AndEq(fField, rValue)
				)

				if rValue != fValue {
					return ConstraintError{
						Key:  fField,
						Type: ForeignKeyConstraint,
						Err:  errors.New("rel: inconsistent has many ref and fk"),
					}
				}

				if updateCount < i {
					col.Swap(updateCount, i)
					mods[i], mods[updateCount] = mods[updateCount], mods[i]
				}

				if err := r.update(ctx, assocDoc, mods[i], filter); err != nil {
					return err
				}

				updateCount++
			} else {
				mods[i].Add(Set(fField, rValue))
				assocDoc.SetValue(fField, rValue)
			}
		}

		if len(mods)-updateCount > 0 {
			var (
				insertMods = mods
				insertCol  = col
			)

			if updateCount > 0 {
				insertMods = mods[updateCount:]
				insertCol = col.Slice(updateCount, len(mods))
			}

			if err := r.insertAll(ctx, insertCol, insertMods); err != nil {
				return err
			}
		}

	}

	return nil
}

// Delete single entry.
func (r repository) Delete(ctx context.Context, record interface{}) error {
	finish := r.instrument(ctx, "rel-delete", "deleting a record")
	defer finish(nil)

	var (
		err          error
		deletedCount int
		doc          = NewDocument(record)
		table        = doc.Table()
		pField       = doc.PrimaryField()
		pValue       = doc.PrimaryValue()
		query        = Build(table, Eq(pField, pValue))
	)

	if doc.Flag(HasDeletedAt) {
		modifies := map[string]Modify{"deleted_at": Set("deleted_at", now())}
		deletedCount, err = r.adapter.Update(ctx, query, modifies)
	} else {
		deletedCount, err = r.adapter.Delete(ctx, query)
	}

	if err == nil && deletedCount == 0 {
		return NotFoundError{}
	}

	return err
}

// MustDelete single entry.
// It'll panic if any error eccured.
func (r repository) MustDelete(ctx context.Context, record interface{}) {
	must(r.Delete(ctx, record))
}

func (r repository) DeleteAll(ctx context.Context, query Query) error {
	finish := r.instrument(ctx, "rel-delete-all", "deleting multiple records")
	defer finish(nil)

	return r.deleteAll(ctx, Invalid, query)
}

func (r repository) MustDeleteAll(ctx context.Context, query Query) {
	must(r.DeleteAll(ctx, query))
}

func (r repository) deleteAll(ctx context.Context, flag DocumentFlag, query Query) error {
	var (
		err error
	)

	if flag.Is(HasDeletedAt) {
		modifies := map[string]Modify{"deleted_at": Set("deleted_at", nil)}
		_, err = r.adapter.Update(ctx, query, modifies)
	} else {
		_, err = r.adapter.Delete(ctx, query)
	}

	return err
}

// Preload loads association with given query.
func (r repository) Preload(ctx context.Context, records interface{}, field string, queriers ...Querier) error {
	finish := r.instrument(ctx, "rel-preload", "preloading associations")
	defer finish(nil)

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
		targets, table, keyField, keyType, ddata, dirties = r.mapPreloadTargets(sl, path)
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
		query    = Build(table, append(queriers, In(keyField, ids...))...)
		cur, err = r.adapter.Query(ctx, r.withDefaultScope(ddata, query))
	)

	if err != nil {
		return err
	}

	scanFinish := r.instrument(ctx, "rel-scan-multi", "scanning all records to multiple targets")
	if err := scanMulti(cur, keyField, keyType, targets); err != nil {
		scanFinish(err)
		return err
	}
	scanFinish(nil)

	// init dirties.
	for i := range dirties {
		dirties[i].initAssoc()
	}

	return nil
}

// MustPreload loads association with given query.
// It'll panic if any error occurred.
func (r repository) MustPreload(ctx context.Context, records interface{}, field string, queriers ...Querier) {
	must(r.Preload(ctx, records, field, queriers...))
}

func (r repository) mapPreloadTargets(sl slice, path []string) (map[interface{}][]slice, string, string, reflect.Type, documentData, []*Dirty) {
	type frame struct {
		index int
		doc   *Document
		dirty *Dirty
	}

	var (
		table         string
		keyField      string
		keyType       reflect.Type
		ddata         documentData
		parentDirties []*Dirty
		mapTarget     = make(map[interface{}][]slice)
		stack         = make([]frame, sl.Len())
	)

	// init stack
	for i := 0; i < len(stack); i++ {
		doc := sl.Get(i)
		stack[i] = frame{index: 0, doc: doc, dirty: doc.Dirty()}
	}

	for len(stack) > 0 {
		var (
			n      = len(stack) - 1
			top    = stack[n]
			field  = path[top.index]
			assocs = top.doc.Association(field)
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

			if top.dirty != nil && top.dirty.doc != nil {
				parentDirties = append(parentDirties, top.dirty)
			}

			if table == "" {
				table = target.Table()
				keyField = assocs.ForeignField()
				keyType = reflect.TypeOf(ref)

				if doc, ok := target.(*Document); ok {
					ddata = doc.data
				}

				if col, ok := target.(*Collection); ok {
					ddata = col.data
				}
			}
		} else {
			if assocs.Type() == HasMany {
				var (
					dirties     map[interface{}]*Dirty
					col, loaded = assocs.Collection()
				)

				if !loaded {
					continue
				}

				if top.dirty != nil {
					dirties = top.dirty.assocMany[field]
				}

				stack = append(stack, make([]frame, col.Len())...)
				for i := 0; i < col.Len(); i++ {
					var (
						dirty *Dirty
						doc   = col.Get(i)
					)

					if dirties != nil {
						dirty = dirties[doc.PrimaryValue()]
					}

					stack[n+i] = frame{
						index: top.index + 1,
						doc:   doc,
						dirty: dirty,
					}
				}
			} else {
				if doc, loaded := assocs.Document(); loaded {
					var dirty *Dirty
					if top.dirty != nil {
						dirty = top.dirty.assoc[field]
					}

					stack = append(stack, frame{
						index: top.index + 1,
						doc:   doc,
						dirty: dirty,
					})
				}
			}
		}

	}

	return mapTarget, table, keyField, keyType, ddata, parentDirties
}

func (r repository) withDefaultScope(ddata documentData, query Query) Query {
	if query.UnscopedQuery {
		return query
	}

	if ddata.flag.Is(HasDeletedAt) {
		query = query.Where(Nil("deleted_at"))
	}

	return query
}

// Transaction performs transaction with given function argument.
func (r repository) Transaction(ctx context.Context, fn func(Repository) error) error {
	finish := r.instrument(ctx, "rel-transaction", "transaction")
	defer finish(nil)

	adp, err := r.adapter.Begin(ctx)
	if err != nil {
		return err
	}

	txRepo := &repository{
		adapter:       adp,
		instrumenter:  r.instrumenter,
		inTransaction: true,
	}

	func() {
		defer func() {
			if p := recover(); p != nil {
				_ = txRepo.adapter.Rollback(ctx)

				switch e := p.(type) {
				case runtime.Error:
					panic(e)
				case error:
					err = e
				default:
					panic(e)
				}
			} else if err != nil {
				_ = txRepo.adapter.Rollback(ctx)
			} else {
				err = txRepo.adapter.Commit(ctx)
			}
		}()

		err = fn(txRepo)
	}()

	return err
}

// New create new repo using adapter.
func New(adapter Adapter) Repository {
	repo := &repository{
		adapter:      adapter,
		instrumenter: DefaultLogger,
	}

	repo.Instrumentation(DefaultLogger)

	return repo
}
