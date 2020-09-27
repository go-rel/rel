package rel

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"strings"
)

// Repository defines sets of available database operations.
type Repository interface {
	Adapter(ctx context.Context) Adapter
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
	Insert(ctx context.Context, record interface{}, mutators ...Mutator) error
	MustInsert(ctx context.Context, record interface{}, mutators ...Mutator)
	InsertAll(ctx context.Context, records interface{}) error
	MustInsertAll(ctx context.Context, records interface{})
	Update(ctx context.Context, record interface{}, mutators ...Mutator) error
	MustUpdate(ctx context.Context, record interface{}, mutators ...Mutator)
	UpdateAll(ctx context.Context, query Query, mutates ...Mutate) error
	MustUpdateAll(ctx context.Context, query Query, mutates ...Mutate)
	Delete(ctx context.Context, record interface{}, options ...Cascade) error
	MustDelete(ctx context.Context, record interface{}, options ...Cascade)
	DeleteAll(ctx context.Context, query Query) error
	MustDeleteAll(ctx context.Context, query Query)
	Preload(ctx context.Context, records interface{}, field string, queriers ...Querier) error
	MustPreload(ctx context.Context, records interface{}, field string, queriers ...Querier)
	Transaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type repository struct {
	rootAdapter  Adapter
	instrumenter Instrumenter
}

func (r repository) Adapter(ctx context.Context) Adapter {
	return fetchContext(ctx, r.rootAdapter).adapter
}

func (r *repository) Instrumentation(instrumenter Instrumenter) {
	r.instrumenter = instrumenter
	r.rootAdapter.Instrumentation(instrumenter)
}

// Ping database.
func (r *repository) Ping(ctx context.Context) error {
	return r.rootAdapter.Ping(ctx)
}

// Iterate through a collection of records from database in batches.
// This function returns iterator that can be used to loop all records.
// Limit, Offset and Sort query is automatically ignored.
func (r repository) Iterate(ctx context.Context, query Query, options ...IteratorOption) Iterator {
	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return newIterator(cw.ctx, cw.adapter, query, options)
}

// Aggregate calculate aggregate over the given field.
// Supported aggregate: count, sum, avg, max, min.
// Any select, group, offset, limit and sort query will be ignored automatically.
// If complex aggregation is needed, consider using All instead,
func (r repository) Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-aggregate", "aggregating records")
	defer finish(nil)

	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return r.aggregate(cw, query, aggregate, field)
}

func (r repository) aggregate(cw contextWrapper, query Query, aggregate string, field string) (int, error) {
	query.GroupQuery = GroupQuery{}
	query.LimitQuery = 0
	query.OffsetQuery = 0
	query.SortQuery = nil

	return cw.adapter.Aggregate(cw.ctx, query, aggregate, field)
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
	finish := r.instrumenter.Observe(ctx, "rel-count", "aggregating records")
	defer finish(nil)

	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return r.aggregate(cw, Build(collection, queriers...), "count", "*")
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
	finish := r.instrumenter.Observe(ctx, "rel-find", "finding a record")
	defer finish(nil)

	var (
		cw    = fetchContext(ctx, r.rootAdapter)
		doc   = NewDocument(record)
		query = Build(doc.Table(), queriers...)
	)

	return r.find(cw, doc, query)
}

// MustFind a record that match the query.
// If no result found, it'll panic.
func (r repository) MustFind(ctx context.Context, record interface{}, queriers ...Querier) {
	must(r.Find(ctx, record, queriers...))
}

func (r repository) find(cw contextWrapper, doc *Document, query Query) error {
	query = r.withDefaultScope(doc.data, query)
	cur, err := cw.adapter.Query(cw.ctx, query.Limit(1))
	if err != nil {
		return err
	}

	finish := r.instrumenter.Observe(cw.ctx, "rel-scan-one", "scanning a record")
	defer finish(nil)

	return scanOne(cur, doc)
}

// FindAll records that match the query.
func (r repository) FindAll(ctx context.Context, records interface{}, queriers ...Querier) error {
	finish := r.instrumenter.Observe(ctx, "rel-find-all", "finding all records")
	defer finish(nil)

	var (
		cw    = fetchContext(ctx, r.rootAdapter)
		col   = NewCollection(records)
		query = Build(col.Table(), queriers...)
	)

	col.Reset()

	return r.findAll(cw, col, query)
}

// MustFindAll records that match the query.
// It'll panic if any error eccured.
func (r repository) MustFindAll(ctx context.Context, records interface{}, queriers ...Querier) {
	must(r.FindAll(ctx, records, queriers...))
}

func (r repository) findAll(cw contextWrapper, col *Collection, query Query) error {
	query = r.withDefaultScope(col.data, query)
	cur, err := cw.adapter.Query(cw.ctx, query)
	if err != nil {
		return err
	}

	finish := r.instrumenter.Observe(cw.ctx, "rel-scan-all", "scanning all records")
	defer finish(nil)

	return scanAll(cur, col)
}

// FindAndCountAll is convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
// Limit and Offset property will be ignored when performing count query.
func (r repository) FindAndCountAll(ctx context.Context, records interface{}, queriers ...Querier) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-find-and-count-all", "finding all records")
	defer finish(nil)

	var (
		cw    = fetchContext(ctx, r.rootAdapter)
		col   = NewCollection(records)
		query = Build(col.Table(), queriers...)
	)

	col.Reset()

	if err := r.findAll(cw, col, query); err != nil {
		return 0, err
	}

	return r.aggregate(cw, query, "count", "*")
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
func (r repository) Insert(ctx context.Context, record interface{}, mutators ...Mutator) error {
	finish := r.instrumenter.Observe(ctx, "rel-insert", "inserting a record")
	defer finish(nil)

	if record == nil {
		return nil
	}

	var (
		cw       = fetchContext(ctx, r.rootAdapter)
		doc      = NewDocument(record)
		mutation = Apply(doc, mutators...)
	)

	if !mutation.IsAssocEmpty() && mutation.Cascade == true {
		return r.transaction(cw, func(cw contextWrapper) error {
			return r.insert(cw, doc, mutation)
		})
	}

	return r.insert(cw, doc, mutation)
}

func (r repository) insert(cw contextWrapper, doc *Document, mutation Mutation) error {
	var (
		pField   string
		pFields  = doc.PrimaryFields()
		queriers = Build(doc.Table())
	)

	if mutation.Cascade {
		if err := r.saveBelongsTo(cw, doc, &mutation); err != nil {
			return err
		}
	}

	if len(pFields) == 1 {
		pField = pFields[0]
	}

	pValue, err := cw.adapter.Insert(cw.ctx, queriers, pField, mutation.Mutates)
	if err != nil {
		return mutation.ErrorFunc.transform(err)
	}

	// update primary value
	if pField != "" {
		doc.SetValue(pField, pValue)
	}

	if mutation.Reload {
		var (
			filter = filterDocument(doc)
		)

		// fetch record
		if err := r.find(cw, doc, queriers.Where(filter)); err != nil {
			return err
		}
	}

	if mutation.Cascade {
		if err := r.saveHasOne(cw, doc, &mutation); err != nil {
			return err
		}

		if err := r.saveHasMany(cw, doc, &mutation, true); err != nil {
			return err
		}
	}

	return nil
}

// MustInsert an record to database.
// It'll panic if any error occurred.
func (r repository) MustInsert(ctx context.Context, record interface{}, mutators ...Mutator) {
	must(r.Insert(ctx, record, mutators...))
}

func (r repository) InsertAll(ctx context.Context, records interface{}) error {
	finish := r.instrumenter.Observe(ctx, "rel-insert-all", "inserting multiple records")
	defer finish(nil)

	if records == nil {
		return nil
	}

	var (
		cw   = fetchContext(ctx, r.rootAdapter)
		col  = NewCollection(records)
		muts = make([]Mutation, col.Len())
	)

	for i := range muts {
		doc := col.Get(i)
		muts[i] = Apply(doc, newStructset(doc, false))
	}

	return r.insertAll(cw, col, muts)
}

func (r repository) MustInsertAll(ctx context.Context, records interface{}) {
	must(r.InsertAll(ctx, records))
}

// TODO: support assocs
func (r repository) insertAll(cw contextWrapper, col *Collection, mutation []Mutation) error {
	if len(mutation) == 0 {
		return nil
	}

	var (
		pField      string
		pFields     = col.PrimaryFields()
		queriers    = Build(col.Table())
		fields      = make([]string, 0, len(mutation[0].Mutates))
		fieldMap    = make(map[string]struct{}, len(mutation[0].Mutates))
		bulkMutates = make([]map[string]Mutate, len(mutation))
	)

	// TODO: baypassable if it's predictable.
	for i := range mutation {
		for field := range mutation[i].Mutates {
			if _, exist := fieldMap[field]; !exist {
				fieldMap[field] = struct{}{}
				fields = append(fields, field)
			}
		}
		bulkMutates[i] = mutation[i].Mutates
	}

	if len(pFields) == 1 {
		pField = pFields[0]
	}

	ids, err := cw.adapter.InsertAll(cw.ctx, queriers, pField, fields, bulkMutates)
	if err != nil {
		return mutation[0].ErrorFunc.transform(err)
	}

	// apply ids
	if pField != "" {
		for i, id := range ids {
			col.Get(i).SetValue(pField, id)
		}
	}

	return nil
}

// Update an record in database.
// It'll panic if any error occurred.
func (r repository) Update(ctx context.Context, record interface{}, mutators ...Mutator) error {
	finish := r.instrumenter.Observe(ctx, "rel-update", "updating a record")
	defer finish(nil)

	if record == nil {
		return nil
	}

	var (
		cw       = fetchContext(ctx, r.rootAdapter)
		doc      = NewDocument(record)
		filter   = filterDocument(doc)
		mutation = Apply(doc, mutators...)
	)

	if !mutation.IsAssocEmpty() && mutation.Cascade == true {
		return r.transaction(cw, func(cw contextWrapper) error {
			return r.update(cw, doc, mutation, filter)
		})
	}

	return r.update(cw, doc, mutation, filter)
}

func (r repository) update(cw contextWrapper, doc *Document, mutation Mutation, filter FilterQuery) error {
	if mutation.Cascade {
		if err := r.saveBelongsTo(cw, doc, &mutation); err != nil {
			return err
		}
	}

	if !mutation.IsMutatesEmpty() {
		var (
			query = r.withDefaultScope(doc.data, Build(doc.Table(), filter, mutation.Unscoped))
		)

		if updatedCount, err := cw.adapter.Update(cw.ctx, query, mutation.Mutates); err != nil {
			return mutation.ErrorFunc.transform(err)
		} else if updatedCount == 0 {
			return NotFoundError{}
		}

		if mutation.Reload {
			if err := r.find(cw, doc, query); err != nil {
				return err
			}
		}
	}

	if mutation.Cascade {
		if err := r.saveHasOne(cw, doc, &mutation); err != nil {
			return err
		}

		if err := r.saveHasMany(cw, doc, &mutation, false); err != nil {
			return err
		}
	}

	return nil
}

// MustUpdate an record in database.
// It'll panic if any error occurred.
func (r repository) MustUpdate(ctx context.Context, record interface{}, mutators ...Mutator) {
	must(r.Update(ctx, record, mutators...))
}

// TODO: support deletion
func (r repository) saveBelongsTo(cw contextWrapper, doc *Document, mutation *Mutation) error {
	for _, field := range doc.BelongsTo() {
		assocMuts, changed := mutation.Assoc[field]
		if !changed || len(assocMuts.Mutations) == 0 {
			continue
		}

		var (
			assoc            = doc.Association(field)
			assocDoc, loaded = assoc.Document()
			assocMut         = assocMuts.Mutations[0]
		)

		if loaded {
			filter, err := filterBelongsTo(assoc)
			if err != nil {
				return err
			}

			if err := r.update(cw, assocDoc, assocMut, filter); err != nil {
				return err
			}
		} else {
			if err := r.insert(cw, assocDoc, assocMut); err != nil {
				return err
			}

			var (
				rField = assoc.ReferenceField()
				fValue = assoc.ForeignValue()
			)

			mutation.Add(Set(rField, fValue))
			doc.SetValue(rField, fValue)
		}
	}

	return nil
}

// TODO: suppprt deletion
func (r repository) saveHasOne(cw contextWrapper, doc *Document, mutation *Mutation) error {
	for _, field := range doc.HasOne() {
		assocMuts, changed := mutation.Assoc[field]
		if !changed || len(assocMuts.Mutations) == 0 {
			continue
		}

		var (
			assoc            = doc.Association(field)
			assocDoc, loaded = assoc.Document()
			assocMut         = assocMuts.Mutations[0]
		)

		if loaded {
			filter, err := filterHasOne(assoc, assocDoc)
			if err != nil {
				return err
			}

			if err := r.update(cw, assocDoc, assocMut, filter); err != nil {
				return err
			}
		} else {
			var (
				fField = assoc.ForeignField()
				rValue = assoc.ReferenceValue()
			)

			assocMut.Add(Set(fField, rValue))
			assocDoc.SetValue(fField, rValue)

			if err := r.insert(cw, assocDoc, assocMut); err != nil {
				return err
			}
		}
	}

	return nil
}

// saveHasMany expects has many mutation to be ordered the same as the recrods in collection.
func (r repository) saveHasMany(cw contextWrapper, doc *Document, mutation *Mutation, insertion bool) error {
	for _, field := range doc.HasMany() {
		assocMuts, changed := mutation.Assoc[field]
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
			muts       = assocMuts.Mutations
			deletedIDs = assocMuts.DeletedIDs
		)

		// this shouldn't happen unless there's bug in the mutator.
		if len(muts) != col.Len() {
			panic("rel: invalid mutator")
		}

		if !insertion {
			var (
				filter = Eq(fField, rValue)
			)

			if deletedIDs == nil {
				// if it's nil, then clear old association (used by structset).
				if _, err := r.deleteAll(cw, col.data.flag, Build(table, filter)); err != nil {
					return err
				}
			} else if len(deletedIDs) > 0 {
				if _, err := r.deleteAll(cw, col.data.flag, Build(table, filter.AndIn(pField, deletedIDs...))); err != nil {
					return err
				}
			}
		}

		// update and filter for bulk insertion.
		updateCount := 0
		for i := range muts {
			var (
				assocDoc = col.Get(i)
			)

			// When deleted IDs is nil, it's assumed that association will be replaced.
			// hence any update request is ignored here.
			if deletedIDs != nil && !isZero(assocDoc.PrimaryValue()) {
				var (
					fValue, _ = assocDoc.Value(fField)
					filter    = filterDocument(assocDoc).AndEq(fField, rValue)
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
					muts[i], muts[updateCount] = muts[updateCount], muts[i]
				}

				if err := r.update(cw, assocDoc, muts[updateCount], filter); err != nil {
					return err
				}

				updateCount++
			} else {
				muts[i].Add(Set(fField, rValue))
				assocDoc.SetValue(fField, rValue)
			}
		}

		if len(muts)-updateCount > 0 {
			var (
				insertMuts = muts
				insertCol  = col
			)

			if updateCount > 0 {
				insertMuts = muts[updateCount:]
				insertCol = col.Slice(updateCount, len(muts))
			}

			if err := r.insertAll(cw, insertCol, insertMuts); err != nil {
				return err
			}
		}

	}

	return nil
}

func (r repository) UpdateAll(ctx context.Context, query Query, mutates ...Mutate) error {
	finish := r.instrumenter.Observe(ctx, "rel-update-all", "updating multiple records")
	defer finish(nil)

	var (
		err  error
		cw   = fetchContext(ctx, r.rootAdapter)
		muts = make(map[string]Mutate, len(mutates))
	)

	for _, mut := range mutates {
		muts[mut.Field] = mut
	}

	if len(muts) > 0 {
		_, err = cw.adapter.Update(cw.ctx, query, muts)
	}

	return err
}

func (r repository) MustUpdateAll(ctx context.Context, query Query, mutates ...Mutate) {
	must(r.UpdateAll(ctx, query, mutates...))
}

// Delete single entry.
func (r repository) Delete(ctx context.Context, record interface{}, options ...Cascade) error {
	finish := r.instrumenter.Observe(ctx, "rel-delete", "deleting a record")
	defer finish(nil)

	var (
		cw      = fetchContext(ctx, r.rootAdapter)
		doc     = NewDocument(record)
		cascade = Cascade(false)
	)

	if len(options) > 0 {
		cascade = options[0]
	}

	if cascade {
		return r.transaction(cw, func(cw contextWrapper) error {
			return r.delete(cw, doc, filterDocument(doc), cascade)
		})
	}

	return r.delete(cw, doc, filterDocument(doc), cascade)
}

func (r repository) delete(cw contextWrapper, doc *Document, filter FilterQuery, cascade Cascade) error {
	var (
		table = doc.Table()
		query = Build(table, filter)
	)

	if cascade {
		if err := r.deleteHasOne(cw, doc, cascade); err != nil {
			return err
		}

		if err := r.deleteHasMany(cw, doc); err != nil {
			return err
		}
	}

	deletedCount, err := r.deleteAll(cw, doc.data.flag, query)
	if err == nil && deletedCount == 0 {
		err = NotFoundError{}
	}

	if err == nil && cascade {
		if err := r.deleteBelongsTo(cw, doc, cascade); err != nil {
			return err
		}
	}

	return err
}

func (r repository) deleteBelongsTo(cw contextWrapper, doc *Document, cascade Cascade) error {
	for _, field := range doc.BelongsTo() {
		var (
			assoc            = doc.Association(field)
			assocDoc, loaded = assoc.Document()
		)

		if loaded {
			filter, err := filterBelongsTo(assoc)
			if err != nil {
				return err
			}

			if err := r.delete(cw, assocDoc, filter, cascade); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r repository) deleteHasOne(cw contextWrapper, doc *Document, cascade Cascade) error {
	for _, field := range doc.HasOne() {
		var (
			assoc            = doc.Association(field)
			assocDoc, loaded = assoc.Document()
		)

		if loaded {
			filter, err := filterHasOne(assoc, assocDoc)
			if err != nil {
				return err
			}

			if err := r.delete(cw, assocDoc, filter, cascade); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r repository) deleteHasMany(cw contextWrapper, doc *Document) error {
	for _, field := range doc.HasMany() {
		var (
			assoc       = doc.Association(field)
			col, loaded = assoc.Collection()
		)

		if loaded {
			var (
				table  = col.Table()
				fField = assoc.ForeignField()
				rValue = assoc.ReferenceValue()
				filter = Eq(fField, rValue).And(filterCollection(col))
			)

			if _, err := r.deleteAll(cw, col.data.flag, Build(table, filter)); err != nil {
				return err
			}
		}
	}

	return nil
}

// MustDelete single entry.
// It'll panic if any error eccured.
func (r repository) MustDelete(ctx context.Context, record interface{}, options ...Cascade) {
	must(r.Delete(ctx, record, options...))
}

// DeleteAll records athat matches query.
func (r repository) DeleteAll(ctx context.Context, query Query) error {
	finish := r.instrumenter.Observe(ctx, "rel-delete-all", "deleting multiple records")
	defer finish(nil)

	var (
		cw     = fetchContext(ctx, r.rootAdapter)
		_, err = r.deleteAll(cw, Invalid, query)
	)

	return err
}

// MustDeleteAll records athat matches query.
// It'll panic if any error eccured.
func (r repository) MustDeleteAll(ctx context.Context, query Query) {
	must(r.DeleteAll(ctx, query))
}

func (r repository) deleteAll(cw contextWrapper, flag DocumentFlag, query Query) (int, error) {
	if flag.Is(HasDeletedAt) {
		mutates := map[string]Mutate{"deleted_at": Set("deleted_at", now())}
		return cw.adapter.Update(cw.ctx, query, mutates)
	}

	return cw.adapter.Delete(cw.ctx, query)
}

// Preload loads association with given query.
// If association is already loaded, this will do nothing.
// To force preloading even though association is already loaeded, add `Reload(true)` as query.
func (r repository) Preload(ctx context.Context, records interface{}, field string, queriers ...Querier) error {
	finish := r.instrumenter.Observe(ctx, "rel-preload", "preloading associations")
	defer finish(nil)

	var (
		sl   slice
		cw   = fetchContext(ctx, r.rootAdapter)
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
		targets, table, keyField, keyType, ddata, loaded = r.mapPreloadTargets(sl, path)
		ids                                              = r.targetIDs(targets)
		query                                            = Build(table, append(queriers, In(keyField, ids...))...)
	)

	if len(targets) == 0 || loaded && !bool(query.ReloadQuery) {
		return nil
	}

	var (
		cur, err = cw.adapter.Query(cw.ctx, r.withDefaultScope(ddata, query))
	)

	if err != nil {
		return err
	}

	scanFinish := r.instrumenter.Observe(ctx, "rel-scan-multi", "scanning all records to multiple targets")
	defer scanFinish(nil)

	return scanMulti(cur, keyField, keyType, targets)
}

// MustPreload loads association with given query.
// It'll panic if any error occurred.
func (r repository) MustPreload(ctx context.Context, records interface{}, field string, queriers ...Querier) {
	must(r.Preload(ctx, records, field, queriers...))
}

func (r repository) mapPreloadTargets(sl slice, path []string) (map[interface{}][]slice, string, string, reflect.Type, documentData, bool) {
	type frame struct {
		index int
		doc   *Document
	}

	var (
		table     string
		keyField  string
		keyType   reflect.Type
		ddata     documentData
		loaded    = true
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
				target       slice
				targetLoaded bool
				ref          = assocs.ReferenceValue()
			)

			if ref == nil {
				continue
			}

			if assocs.Type() == HasMany {
				target, targetLoaded = assocs.Collection()
			} else {
				target, targetLoaded = assocs.Document()
			}

			target.Reset()
			mapTarget[ref] = append(mapTarget[ref], target)
			loaded = loaded && targetLoaded

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

	return mapTarget, table, keyField, keyType, ddata, loaded
}

func (r repository) targetIDs(targets map[interface{}][]slice) []interface{} {
	var (
		ids = make([]interface{}, len(targets))
		i   = 0
	)

	for key := range targets {
		ids[i] = key
		i++
	}

	return ids
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
func (r repository) Transaction(ctx context.Context, fn func(ctx context.Context) error) error {
	finish := r.instrumenter.Observe(ctx, "rel-transaction", "transaction")
	defer finish(nil)

	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return r.transaction(cw, func(cw contextWrapper) error {
		return fn(cw.ctx)
	})
}

func (r repository) transaction(cw contextWrapper, fn func(cw contextWrapper) error) error {
	adp, err := cw.adapter.Begin(cw.ctx)
	if err != nil {
		return err
	}

	// wrap trx adapter to new context.
	cw = wrapContext(cw.ctx, adp)

	func() {
		defer func() {
			if p := recover(); p != nil {
				_ = cw.adapter.Rollback(cw.ctx)

				switch e := p.(type) {
				case runtime.Error:
					panic(e)
				case error:
					err = e
				default:
					panic(e)
				}
			} else if err != nil {
				_ = cw.adapter.Rollback(cw.ctx)
			} else {
				err = cw.adapter.Commit(cw.ctx)
			}
		}()

		err = fn(cw)
	}()

	return err
}

// New create new repo using adapter.
func New(adapter Adapter) Repository {
	repo := &repository{
		rootAdapter:  adapter,
		instrumenter: DefaultLogger,
	}

	repo.Instrumentation(DefaultLogger)

	return repo
}
