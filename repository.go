package rel

import (
	"context"
	"errors"
	"reflect"
	"runtime"
	"strings"
)

// Repository for interacting with database.
type Repository interface {
	// Adapter used in this repository.
	Adapter(ctx context.Context) Adapter

	// Instrumentation defines callback to be used as instrumenter.
	Instrumentation(instrumenter Instrumenter)

	// Ping database.
	Ping(ctx context.Context) error

	// Iterate through a collection of entities from database in batches.
	// This function returns iterator that can be used to loop all entities.
	// Limit, Offset and Sort query is automatically ignored.
	Iterate(ctx context.Context, query Query, option ...IteratorOption) Iterator

	// Aggregate over the given field.
	// Supported aggregate: count, sum, avg, max, min.
	// Any select, group, offset, limit and sort query will be ignored automatically.
	// If complex aggregation is needed, consider using All instead.
	Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error)

	// MustAggregate over the given field.
	// Supported aggregate: count, sum, avg, max, min.
	// Any select, group, offset, limit and sort query will be ignored automatically.
	// If complex aggregation is needed, consider using All instead.
	// It'll panic if any error eccured.
	MustAggregate(ctx context.Context, query Query, aggregate string, field string) int

	// Count entities that match the query.
	Count(ctx context.Context, collection string, queriers ...Querier) (int, error)

	// MustCount entities that match the query.
	// It'll panic if any error eccured.
	MustCount(ctx context.Context, collection string, queriers ...Querier) int

	// Find a entity that match the query.
	// If no result found, it'll return not found error.
	Find(ctx context.Context, entity any, queriers ...Querier) error

	// MustFind a entity that match the query.
	// If no result found, it'll panic.
	MustFind(ctx context.Context, entity any, queriers ...Querier)

	// FindAll entities that match the query.
	FindAll(ctx context.Context, entities any, queriers ...Querier) error

	// MustFindAll entities that match the query.
	// It'll panic if any error eccured.
	MustFindAll(ctx context.Context, entities any, queriers ...Querier)

	// FindAndCountAll entities that match the query.
	// This is a convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
	// Limit and Offset property will be ignored when performing count query.
	FindAndCountAll(ctx context.Context, entities any, queriers ...Querier) (int, error)

	// MustFindAndCountAll entities that match the query.
	// This is a convenient method that combines FindAll and Count. It's useful when dealing with queries related to pagination.
	// Limit and Offset property will be ignored when performing count query.
	// It'll panic if any error eccured.
	MustFindAndCountAll(ctx context.Context, entities any, queriers ...Querier) int

	// Insert a entity to database.
	Insert(ctx context.Context, entity any, mutators ...Mutator) error

	// MustInsert an entity to database.
	// It'll panic if any error occurred.
	MustInsert(ctx context.Context, entity any, mutators ...Mutator)

	// InsertAll entities.
	// Does not supports application cascade insert.
	InsertAll(ctx context.Context, entities any, mutators ...Mutator) error

	// MustInsertAll entities.
	// It'll panic if any error occurred.
	// Does not supports application cascade insert.
	MustInsertAll(ctx context.Context, entities any, mutators ...Mutator)

	// Update a entity in database.
	// It'll panic if any error occurred.
	Update(ctx context.Context, entity any, mutators ...Mutator) error

	// MustUpdate a entity in database.
	// It'll panic if any error occurred.
	MustUpdate(ctx context.Context, entity any, mutators ...Mutator)

	// UpdateAny entities tha match the query.
	// Returns number of updated entities and error.
	UpdateAny(ctx context.Context, query Query, mutates ...Mutate) (int, error)

	// MustUpdateAny entities that match the query.
	// It'll panic if any error occurred.
	// Returns number of updated entities.
	MustUpdateAny(ctx context.Context, query Query, mutates ...Mutate) int

	// Delete a entity.
	Delete(ctx context.Context, entity any, mutators ...Mutator) error

	// MustDelete a entity.
	// It'll panic if any error eccured.
	MustDelete(ctx context.Context, entity any, mutators ...Mutator)

	// DeleteAll entities.
	// Does not supports application cascade delete.
	DeleteAll(ctx context.Context, entities any) error

	// MustDeleteAll entities.
	// It'll panic if any error occurred.
	// Does not supports application cascade delete.
	MustDeleteAll(ctx context.Context, entities any)

	// DeleteAny entities that match the query.
	// Returns number of deleted entities and error.
	DeleteAny(ctx context.Context, query Query) (int, error)

	// MustDeleteAny entities that match the query.
	// It'll panic if any error eccured.
	// Returns number of updated entities.
	MustDeleteAny(ctx context.Context, query Query) int

	// Preload association with given query.
	// This function can accepts either a struct or a slice of structs.
	// If association is already loaded, this will do nothing.
	// To force preloading even though association is already loaeded, add `Reload(true)` as query.
	Preload(ctx context.Context, entities any, field string, queriers ...Querier) error

	// MustPreload association with given query.
	// This function can accepts either a struct or a slice of structs.
	// It'll panic if any error occurred.
	MustPreload(ctx context.Context, entities any, field string, queriers ...Querier)

	// Exec raw statement.
	// Returns last inserted id, rows affected and error.
	Exec(ctx context.Context, statement string, args ...any) (int, int, error)

	// MustExec raw statement.
	// Returns last inserted id, rows affected and error.
	MustExec(ctx context.Context, statement string, args ...any) (int, int)

	// Transaction performs transaction with given function argument.
	// Transaction scope/connection is automatically passed using context.
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

func (r *repository) Ping(ctx context.Context) error {
	return r.rootAdapter.Ping(ctx)
}

func (r repository) Iterate(ctx context.Context, query Query, options ...IteratorOption) Iterator {
	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return newIterator(cw.ctx, cw.adapter, query, options)
}

func (r repository) Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-aggregate", "aggregating entities")
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

func (r repository) MustAggregate(ctx context.Context, query Query, aggregate string, field string) int {
	result, err := r.Aggregate(ctx, query, aggregate, field)
	must(err)
	return result
}

func (r repository) Count(ctx context.Context, collection string, queriers ...Querier) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-count", "aggregating entities")
	defer finish(nil)

	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return r.aggregate(cw, Build(collection, queriers...), "count", "*")
}

func (r repository) MustCount(ctx context.Context, collection string, queriers ...Querier) int {
	count, err := r.Count(ctx, collection, queriers...)
	must(err)
	return count
}

func (r repository) Find(ctx context.Context, entity any, queriers ...Querier) error {
	finish := r.instrumenter.Observe(ctx, "rel-find", "finding a entity")
	defer finish(nil)

	var (
		cw    = fetchContext(ctx, r.rootAdapter)
		doc   = NewDocument(entity)
		query = Build(doc.Table(), queriers...).Populate(doc.Meta())
	)

	return r.find(cw, doc, query)
}

func (r repository) MustFind(ctx context.Context, entity any, queriers ...Querier) {
	must(r.Find(ctx, entity, queriers...))
}

func (r repository) find(cw contextWrapper, doc *Document, query Query) error {
	query = r.withDefaultScope(doc.meta, query, true)
	cur, err := cw.adapter.Query(cw.ctx, query.Limit(1))
	if err != nil {
		return err
	}

	finish := r.instrumenter.Observe(cw.ctx, "rel-scan-one", "scanning a entity")
	if err := scanOne(cur, doc); err != nil {
		finish(err)
		return err
	}
	finish(nil)

	for i := range query.PreloadQuery {
		if err := r.preload(cw, doc, query.PreloadQuery[i], nil); err != nil {
			return err
		}
	}

	return nil
}

func (r repository) FindAll(ctx context.Context, entities any, queriers ...Querier) error {
	finish := r.instrumenter.Observe(ctx, "rel-find-all", "finding all entities")
	defer finish(nil)

	var (
		cw    = fetchContext(ctx, r.rootAdapter)
		col   = NewCollection(entities)
		query = Build(col.Table(), queriers...).Populate(col.Meta())
	)

	col.Reset()

	return r.findAll(cw, col, query)
}

func (r repository) MustFindAll(ctx context.Context, entities any, queriers ...Querier) {
	must(r.FindAll(ctx, entities, queriers...))
}

func (r repository) findAll(cw contextWrapper, col *Collection, query Query) error {
	query = r.withDefaultScope(col.meta, query, true)
	cur, err := cw.adapter.Query(cw.ctx, query)
	if err != nil {
		return err
	}

	finish := r.instrumenter.Observe(cw.ctx, "rel-scan-all", "scanning all entities")
	if err := scanAll(cur, col); err != nil {
		finish(err)
		return err
	}
	finish(nil)

	for i := range query.PreloadQuery {
		if err := r.preload(cw, col, query.PreloadQuery[i], nil); err != nil {
			return err
		}
	}

	return nil
}

func (r repository) FindAndCountAll(ctx context.Context, entities any, queriers ...Querier) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-find-and-count-all", "finding all entities")
	defer finish(nil)

	var (
		cw    = fetchContext(ctx, r.rootAdapter)
		col   = NewCollection(entities)
		query = Build(col.Table(), queriers...).Populate(col.Meta())
	)

	col.Reset()

	if err := r.findAll(cw, col, query); err != nil {
		return 0, err
	}

	return r.aggregate(cw, r.withDefaultScope(col.meta, query, false), "count", "*")
}

func (r repository) MustFindAndCountAll(ctx context.Context, entities any, queriers ...Querier) int {
	count, err := r.FindAndCountAll(ctx, entities, queriers...)
	must(err)

	return count
}

func (r repository) Insert(ctx context.Context, entity any, mutators ...Mutator) error {
	finish := r.instrumenter.Observe(ctx, "rel-insert", "inserting a entity")
	defer finish(nil)

	if entity == nil {
		return nil
	}

	var (
		cw       = fetchContext(ctx, r.rootAdapter)
		doc      = NewDocument(entity)
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

	pValue, err := cw.adapter.Insert(cw.ctx, queriers, pField, mutation.Mutates, mutation.OnConflict)
	if err != nil {
		return mutation.ErrorFunc.transform(err)
	}

	// update primary value
	if pField != "" {
		doc.SetValue(pField, pValue)
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

func (r repository) MustInsert(ctx context.Context, entity any, mutators ...Mutator) {
	must(r.Insert(ctx, entity, mutators...))
}

func (r repository) InsertAll(ctx context.Context, entities any, mutators ...Mutator) error {
	finish := r.instrumenter.Observe(ctx, "rel-insert-all", "inserting multiple entities")
	defer finish(nil)

	if entities == nil {
		return nil
	}

	var (
		cw   = fetchContext(ctx, r.rootAdapter)
		col  = NewCollection(entities)
		muts = make([]Mutation, col.Len())
	)

	for i := range muts {
		doc := col.Get(i)
		if i == 0 {
			// only need to apply options from first one
			muts[i] = Apply(doc, mutators...)
		} else {
			muts[i] = Apply(doc)
		}
	}

	return r.insertAll(cw, col, muts)
}

func (r repository) MustInsertAll(ctx context.Context, entities any, mutators ...Mutator) {
	must(r.InsertAll(ctx, entities, mutators...))
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
		onConflict  = mutation[0].OnConflict
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

	ids, err := cw.adapter.InsertAll(cw.ctx, queriers, pField, fields, bulkMutates, onConflict)
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

func (r repository) Update(ctx context.Context, entity any, mutators ...Mutator) error {
	finish := r.instrumenter.Observe(ctx, "rel-update", "updating a entity")
	defer finish(nil)

	if entity == nil {
		return nil
	}

	var (
		cw       = fetchContext(ctx, r.rootAdapter)
		doc      = NewDocument(entity)
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

func (r repository) lockVersion(doc Document, unscoped Unscoped) (int, bool) {
	if unscoped {
		return 0, false
	}
	if doc.Flag(HasVersioning) {
		versionRaw, _ := doc.Value("lock_version")
		version, _ := versionRaw.(int)
		return version, true
	}
	return 0, false
}

func (r repository) update(cw contextWrapper, doc *Document, mutation Mutation, filter FilterQuery) error {
	if mutation.Cascade {
		if err := r.saveBelongsTo(cw, doc, &mutation); err != nil {
			return err
		}
	}

	if !mutation.IsMutatesEmpty() {
		if err := r.applyMutates(cw, doc, mutation, filter); err != nil {
			return err
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

func (r repository) applyMutates(cw contextWrapper, doc *Document, mutation Mutation, filter FilterQuery) (dbErr error) {
	var (
		baseQueries = []Querier{filter, mutation.Unscoped, mutation.Cascade}
		queries     = baseQueries
	)

	if version, ok := r.lockVersion(*doc, mutation.Unscoped); ok {
		Set("lock_version", version+1).Apply(doc, &mutation)
		queries = append(queries, lockVersion(version))
		defer func() {
			if dbErr != nil {
				doc.SetValue("lock_version", version)
			}
		}()
	}

	var (
		pField string
		query  = r.withDefaultScope(doc.meta, Build(doc.Table(), queries...).Populate(doc.Meta()), false)
	)

	if len(doc.meta.primaryField) == 1 {
		pField = doc.PrimaryField()
	}

	if updatedCount, err := cw.adapter.Update(cw.ctx, query, pField, mutation.Mutates); err != nil {
		return mutation.ErrorFunc.transform(err)
	} else if updatedCount == 0 {
		return NotFoundError{}
	}

	if mutation.Reload {
		baseQuery := r.withDefaultScope(doc.meta, Build(doc.Table(), baseQueries...).Populate(doc.Meta()), false)
		if err := r.find(cw, doc, baseQuery.UsePrimary()); err != nil {
			return err
		}
	}

	return nil
}

func (r repository) MustUpdate(ctx context.Context, entity any, mutators ...Mutator) {
	must(r.Update(ctx, entity, mutators...))
}

// TODO: support deletion
func (r repository) saveBelongsTo(cw contextWrapper, doc *Document, mutation *Mutation) error {
	for _, field := range doc.BelongsTo() {
		var (
			assoc              = doc.Association(field)
			assocMuts, changed = mutation.Assoc[field]
		)

		if !assoc.Autosave() || !changed || len(assocMuts.Mutations) == 0 {
			continue
		}

		var (
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
		var (
			assoc              = doc.Association(field)
			assocMuts, changed = mutation.Assoc[field]
		)

		if !assoc.Autosave() || !changed || len(assocMuts.Mutations) == 0 {
			continue
		}

		var (
			assocDoc, loaded = assoc.Document()
			assocMut         = assocMuts.Mutations[0]
		)

		if loaded && (assoc.ForeignField() == "" || !isZero(assoc.ForeignValue())) {
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
		var (
			assoc              = doc.Association(field)
			assocMuts, changed = mutation.Assoc[field]
		)

		if !assoc.Autosave() || !changed {
			continue
		}

		var (
			col, _     = assoc.Collection()
			table      = col.Table()
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
				if _, err := r.deleteAny(cw, col.meta.flag, Build(table, filter).Populate(col.Meta())); err != nil {
					return err
				}
			} else if len(deletedIDs) > 0 {
				filter = filter.AndIn(col.PrimaryField(), deletedIDs...)
				if _, err := r.deleteAny(cw, col.meta.flag, Build(table, filter).Populate(col.Meta())); err != nil {
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
			var fValue, _ = assocDoc.Value(fField)
			if deletedIDs != nil && !isZero(assocDoc.PrimaryValue()) && !isZero(fValue) {
				var (
					filter = filterDocument(assocDoc).AndEq(fField, rValue)
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

func (r repository) UpdateAny(ctx context.Context, query Query, mutates ...Mutate) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-update-any", "updating multiple entities")
	defer finish(nil)

	var (
		err          error
		updatedCount int
		cw           = fetchContext(ctx, r.rootAdapter)
		muts         = make(map[string]Mutate, len(mutates))
	)

	for _, mut := range mutates {
		muts[mut.Field] = mut
	}

	if len(muts) > 0 {
		updatedCount, err = cw.adapter.Update(cw.ctx, query, "", muts)
	}

	return updatedCount, err
}

func (r repository) MustUpdateAny(ctx context.Context, query Query, mutates ...Mutate) int {
	updatedCount, err := r.UpdateAny(ctx, query, mutates...)
	must(err)
	return updatedCount
}

func (r repository) Delete(ctx context.Context, entity any, mutators ...Mutator) error {
	finish := r.instrumenter.Observe(ctx, "rel-delete", "deleting a entity")
	defer finish(nil)

	var (
		cw       = fetchContext(ctx, r.rootAdapter)
		doc      = NewDocument(entity)
		mutation = applyMutators(nil, false, false, mutators...)
	)

	if mutation.Cascade {
		return r.transaction(cw, func(cw contextWrapper) error {
			return r.delete(cw, doc, filterDocument(doc), mutation)
		})
	}

	return r.delete(cw, doc, filterDocument(doc), mutation)
}

func (r repository) delete(cw contextWrapper, doc *Document, filter FilterQuery, mutation Mutation) error {
	var filters []Querier = []Querier{filter, mutation.Unscoped}

	if version, ok := r.lockVersion(*doc, mutation.Unscoped); ok {
		filters = append(filters, lockVersion(version))
	}

	var (
		table = doc.Table()
		query = Build(table, filters...).Populate(doc.Meta())
	)

	if mutation.Cascade {
		if err := r.deleteHasOne(cw, doc, true); err != nil {
			return err
		}

		if err := r.deleteHasMany(cw, doc); err != nil {
			return err
		}
	}

	deletedCount, err := r.deleteAny(cw, doc.meta.flag, query)
	if err == nil && deletedCount == 0 {
		err = NotFoundError{}
	}

	if err == nil && mutation.Cascade {
		if err := r.deleteBelongsTo(cw, doc, true); err != nil {
			return err
		}
	}

	return err
}

func (r repository) deleteBelongsTo(cw contextWrapper, doc *Document, cascade Cascade) error {
	for _, field := range doc.BelongsTo() {
		var (
			assoc = doc.Association(field)
		)

		if !assoc.Autosave() {
			continue
		}

		if assocDoc, loaded := assoc.Document(); loaded {
			filter, err := filterBelongsTo(assoc)
			if err != nil {
				return err
			}

			if err := r.delete(cw, assocDoc, filter, Mutation{Cascade: cascade}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r repository) deleteHasOne(cw contextWrapper, doc *Document, cascade Cascade) error {
	for _, field := range doc.HasOne() {
		var (
			assoc = doc.Association(field)
		)

		if !assoc.Autosave() {
			continue
		}

		if assocDoc, loaded := assoc.Document(); loaded {
			filter, err := filterHasOne(assoc, assocDoc)
			if err != nil {
				return err
			}

			if err := r.delete(cw, assocDoc, filter, Mutation{Cascade: cascade}); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r repository) deleteHasMany(cw contextWrapper, doc *Document) error {
	for _, field := range doc.HasMany() {
		var (
			assoc = doc.Association(field)
		)

		if !assoc.Autosave() {
			continue
		}

		if col, loaded := assoc.Collection(); loaded && col.Len() != 0 {
			var (
				table  = col.Table()
				fField = assoc.ForeignField()
				rValue = assoc.ReferenceValue()
				filter = Eq(fField, rValue).And(filterCollection(col))
			)

			if _, err := r.deleteAny(cw, col.meta.flag, Build(table, filter).Populate(doc.Meta())); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r repository) MustDelete(ctx context.Context, entity any, mutators ...Mutator) {
	must(r.Delete(ctx, entity, mutators...))
}

func (r repository) DeleteAll(ctx context.Context, entities any) error {
	finish := r.instrumenter.Observe(ctx, "rel-delete-all", "deleting entities")
	defer finish(nil)

	var (
		cw  = fetchContext(ctx, r.rootAdapter)
		col = NewCollection(entities)
	)

	if col.Len() == 0 {
		return nil
	}

	var (
		query  = Build(col.Table(), filterCollection(col)).Populate(col.Meta())
		_, err = r.deleteAny(cw, col.meta.flag, query)
	)

	return err
}

func (r repository) MustDeleteAll(ctx context.Context, entities any) {
	must(r.DeleteAll(ctx, entities))
}

func (r repository) DeleteAny(ctx context.Context, query Query) (int, error) {
	finish := r.instrumenter.Observe(ctx, "rel-delete-any", "deleting multiple entities")
	defer finish(nil)

	var (
		cw = fetchContext(ctx, r.rootAdapter)
	)

	return r.deleteAny(cw, Invalid, query)
}

func (r repository) MustDeleteAny(ctx context.Context, query Query) int {
	deletedCount, err := r.DeleteAny(ctx, query)
	must(err)
	return deletedCount
}

func (r repository) deleteAny(cw contextWrapper, flag DocumentFlag, query Query) (int, error) {
	hasDeletedAt := flag.Is(HasDeletedAt)
	hasDeleted := flag.Is(HasDeleted)
	mutates := make(map[string]Mutate, 1)
	if hasDeletedAt {
		mutates["deleted_at"] = Set("deleted_at", Now())
	}
	if hasDeleted {
		mutates["deleted"] = Set("deleted", true)
		if flag.Is(HasUpdatedAt) && !hasDeletedAt {
			mutates["updated_at"] = Set("updated_at", Now())
		}
	}
	if hasDeletedAt || hasDeleted {
		if flag.Is(HasVersioning) {
			mutates["lock_version"] = Inc("lock_version")
		}
		return cw.adapter.Update(cw.ctx, query, "", mutates)
	}

	return cw.adapter.Delete(cw.ctx, query)
}

func (r repository) Preload(ctx context.Context, entities any, field string, queriers ...Querier) error {
	finish := r.instrumenter.Observe(ctx, "rel-preload", "preloading associations")
	defer finish(nil)

	var (
		sl slice
		cw = fetchContext(ctx, r.rootAdapter)
		rt = reflect.TypeOf(entities)
	)

	if rt.Kind() != reflect.Ptr {
		panic("rel: entity parameter must be a pointer.")
	}

	rt = rt.Elem()
	if rt.Kind() == reflect.Slice {
		sl = NewCollection(entities)
	} else {
		sl = NewDocument(entities)
	}

	return r.preload(cw, sl, field, queriers)
}

func (r repository) preload(cw contextWrapper, entities slice, field string, queriers []Querier) error {
	var (
		path                                             = strings.Split(field, ".")
		targets, table, keyField, keyType, ddata, loaded = r.mapPreloadTargets(entities, path)
		ids                                              = r.targetIDs(targets)
		inClauseLength                                   = 999
	)

	// Create separate queries if the amount of ids is more than inClauseLength.
	for {
		if len(ids) == 0 {
			break
		}

		// necessary check to avoid slicing beyond
		// slice capacity
		if len(ids) < inClauseLength {
			inClauseLength = len(ids)
		}

		idsChunk := ids[0:inClauseLength]
		ids = ids[inClauseLength:]

		query := Build(table, append(queriers, In(keyField, idsChunk...))...).Populate(entities.Meta())
		if len(targets) == 0 || loaded && !bool(query.ReloadQuery) {
			return nil
		}

		var (
			cur, err = cw.adapter.Query(cw.ctx, r.withDefaultScope(ddata, query, false))
		)

		if err != nil {
			return err
		}

		scanFinish := r.instrumenter.Observe(cw.ctx, "rel-scan-multi", "scanning all entities to multiple targets")
		// Note: Calling scanMulti multiple times with the same targets works
		// only if the cursor of each execution only contains a new set of keys.
		// That is here the case as each select is with a unique set of ids.
		err = scanMulti(cur, keyField, keyType, targets)
		scanFinish(err)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r repository) MustPreload(ctx context.Context, entities any, field string, queriers ...Querier) {
	must(r.Preload(ctx, entities, field, queriers...))
}

func (r repository) mapPreloadTargets(sl slice, path []string) (map[any][]slice, string, string, reflect.Type, DocumentMeta, bool) {
	type frame struct {
		index int
		doc   *Document
	}

	var (
		table     string
		keyField  string
		keyType   reflect.Type
		meta      DocumentMeta
		loaded    = true
		mapTarget = make(map[any][]slice)
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
				target, targetLoaded = assocs.LazyDocument()
			}

			target.Reset()
			mapTarget[ref] = append(mapTarget[ref], target)
			loaded = loaded && targetLoaded

			if table == "" {
				table = target.Table()
				keyField = assocs.ForeignField()
				keyType = reflect.TypeOf(ref)

				if doc, ok := target.(*Document); ok {
					meta = doc.meta
				}

				if col, ok := target.(*Collection); ok {
					meta = col.meta
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
				if doc, loaded := assocs.LazyDocument(); loaded {
					stack = append(stack, frame{
						index: top.index + 1,
						doc:   doc,
					})
				}
			}
		}

	}

	return mapTarget, table, keyField, keyType, meta, loaded
}

func (r repository) targetIDs(targets map[any][]slice) []any {
	var (
		ids = make([]any, len(targets))
		i   = 0
	)

	for key := range targets {
		ids[i] = key
		i++
	}

	return ids
}

func (r repository) withDefaultScope(meta DocumentMeta, query Query, preload bool) Query {
	if query.UnscopedQuery {
		return query
	}

	if meta.flag.Is(HasDeleted) {
		query = query.Where(Eq("deleted", false))
	} else if meta.flag.Is(HasDeletedAt) {
		query = query.Where(Nil("deleted_at"))
	}

	if preload && bool(query.CascadeQuery) {
		query.PreloadQuery = append(meta.preload, query.PreloadQuery...)
	}

	return query
}

// Exec raw statement.
// Returns last inserted id, rows affected and error.
func (r repository) Exec(ctx context.Context, stmt string, args ...any) (int, int, error) {
	lastInsertedId, rowsAffected, err := r.Adapter(ctx).Exec(ctx, stmt, args)
	return int(lastInsertedId), int(rowsAffected), err
}

// MustExec raw statement.
// Returns last inserted id, rows affected and error.
func (r repository) MustExec(ctx context.Context, stmt string, args ...any) (int, int) {
	lastInsertedId, rowsAffected, err := r.Exec(ctx, stmt, args...)
	must(err)
	return lastInsertedId, rowsAffected
}

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
