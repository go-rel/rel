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
	SetLogger(logger ...Logger)
	Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error)
	MustAggregate(ctx context.Context, query Query, aggregate string, field string) int
	Count(ctx context.Context, collection string, queriers ...Querier) (int, error)
	MustCount(ctx context.Context, collection string, queriers ...Querier) int
	Find(ctx context.Context, record interface{}, queriers ...Querier) error
	MustFind(ctx context.Context, record interface{}, queriers ...Querier)
	FindAll(ctx context.Context, records interface{}, queriers ...Querier) error
	MustFindAll(ctx context.Context, records interface{}, queriers ...Querier)
	Insert(ctx context.Context, record interface{}, modifiers ...Modifier) error
	MustInsert(ctx context.Context, record interface{}, modifiers ...Modifier)
	InsertAll(ctx context.Context, records interface{}) error
	MustInsertAll(ctx context.Context, records interface{})
	Update(ctx context.Context, record interface{}, modifiers ...Modifier) error
	MustUpdate(ctx context.Context, record interface{}, modifiers ...Modifier)
	Delete(ctx context.Context, record interface{}) error
	MustDelete(ctx context.Context, record interface{})
	DeleteAll(ctx context.Context, queriers ...Querier) error
	MustDeleteAll(ctx context.Context, queriers ...Querier)
	Preload(ctx context.Context, records interface{}, field string, queriers ...Querier) error
	MustPreload(ctx context.Context, records interface{}, field string, queriers ...Querier)
	Transaction(ctx context.Context, fn func(Repository) error) error
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
func (r repository) Aggregate(ctx context.Context, query Query, aggregate string, field string) (int, error) {
	query.GroupQuery = GroupQuery{}
	query.LimitQuery = 0
	query.OffsetQuery = 0
	query.SortQuery = nil

	return r.adapter.Aggregate(ctx, query, aggregate, field, r.logger...)
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
	return r.Aggregate(ctx, Build(collection, queriers...), "count", "*")
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
	cur, err := r.adapter.Query(ctx, query.Limit(1), r.logger...)
	if err != nil {
		return err
	}

	return scanOne(cur, doc)
}

// FindAll records that match the query.
func (r repository) FindAll(ctx context.Context, records interface{}, queriers ...Querier) error {
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
	cur, err := r.adapter.Query(ctx, query, r.logger...)
	if err != nil {
		return err
	}

	return scanMany(cur, col)
}

// Insert an record to database.
func (r repository) Insert(ctx context.Context, record interface{}, modifiers ...Modifier) error {
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

	pValue, err := r.Adapter().Insert(ctx, queriers, modification.Modifies, r.logger...)
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

	ids, err := r.adapter.InsertAll(ctx, queriers, fields, bulkModifies, r.logger...)
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
			query             = r.withDefaultScope(doc.data, Build(doc.Table(), filter))
			updatedCount, err = r.adapter.Update(ctx, query, modification.Modifies, r.logger...)
		)

		if err != nil {
			return err
		}

		if updatedCount == 0 {
			return NotFoundError{}
		}

		if modification.Reload {
			if err := r.find(ctx, doc, query.Unscoped()); err != nil {
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
		deletedCount, err = r.adapter.Update(ctx, query, modifies, r.logger...)
	} else {
		deletedCount, err = r.adapter.Delete(ctx, query, r.logger...)
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

func (r repository) DeleteAll(ctx context.Context, queriers ...Querier) error {
	var (
		q = Build("", queriers...)
	)

	return r.deleteAll(ctx, Invalid, q)
}

func (r repository) MustDeleteAll(ctx context.Context, queriers ...Querier) {
	must(r.DeleteAll(ctx, queriers...))
}

func (r repository) deleteAll(ctx context.Context, flag DocumentFlag, query Query) error {
	var (
		err error
	)

	if flag.Is(HasDeletedAt) {
		modifies := map[string]Modify{"deleted_at": Set("deleted_at", nil)}
		_, err = r.adapter.Update(ctx, query, modifies, r.logger...)
	} else {
		_, err = r.adapter.Delete(ctx, query, r.logger...)
	}

	return err
}

// Preload loads association with given query.
func (r repository) Preload(ctx context.Context, records interface{}, field string, queriers ...Querier) error {
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
		targets, table, keyField, keyType, ddata = r.mapPreloadTargets(sl, path)
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
		cur, err = r.adapter.Query(ctx, r.withDefaultScope(ddata, query), r.logger...)
	)

	if err != nil {
		return err
	}

	return scanMulti(cur, keyField, keyType, targets)
}

// MustPreload loads association with given query.
// It'll panic if any error occurred.
func (r repository) MustPreload(ctx context.Context, records interface{}, field string, queriers ...Querier) {
	must(r.Preload(ctx, records, field, queriers...))
}

func (r repository) mapPreloadTargets(sl slice, path []string) (map[interface{}][]slice, string, string, reflect.Type, documentData) {
	type frame struct {
		index int
		doc   *Document
	}

	var (
		table     string
		keyField  string
		keyType   reflect.Type
		ddata     documentData
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

	return mapTarget, table, keyField, keyType, ddata
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
	adp, err := r.adapter.Begin(ctx)
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
	return &repository{
		adapter: adapter,
		logger:  []Logger{DefaultLogger},
	}
}
