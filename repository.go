package rel

import (
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
	Aggregate(query Query, aggregate string, field string) (int, error)
	MustAggregate(query Query, aggregate string, field string) int
	Count(collection string, queriers ...Querier) (int, error)
	MustCount(collection string, queriers ...Querier) int
	Find(record interface{}, queriers ...Querier) error
	MustFind(record interface{}, queriers ...Querier)
	FindAll(records interface{}, queriers ...Querier) error
	MustFindAll(records interface{}, queriers ...Querier)
	Insert(record interface{}, modifiers ...Modifier) error
	MustInsert(record interface{}, modifiers ...Modifier)
	InsertAll(records interface{}) error
	MustInsertAll(records interface{})
	Update(record interface{}, modifiers ...Modifier) error
	MustUpdate(record interface{}, modifiers ...Modifier)
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
	return r.Aggregate(Build(collection, queriers...), "count", "*")
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
		query = Build(doc.Table(), queriers...)
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
		query = Build(col.Table(), queriers...)
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
func (r repository) Insert(record interface{}, modifiers ...Modifier) error {
	// TODO: perform reference check on library level for record instead of adapter level
	if record == nil {
		return nil
	}

	var (
		modification Modification
		doc          = NewDocument(record)
	)

	if len(modifiers) == 0 {
		modification = Apply(doc, newStructset(doc))
	} else {
		modification = Apply(doc, modifiers...)
	}

	if len(modification.Assoc) > 0 {
		return r.Transaction(func(r Repository) error {
			return r.(*repository).insert(doc, modification)
		})
	}

	return r.insert(doc, modification)
}

func (r repository) insert(doc *Document, modification Modification) error {
	var (
		pField   = doc.PrimaryField()
		queriers = Build(doc.Table())
	)

	if err := r.saveBelongsTo(doc, &modification); err != nil {
		return err
	}

	pValue, err := r.Adapter().Insert(queriers, modification.Modifies, r.logger...)
	if err != nil {
		return err
	}

	if modification.Reload {
		// fetch record
		if err := r.find(doc, queriers.Where(Eq(pField, pValue))); err != nil {
			return err
		}
	} else {
		// update primary value
		doc.SetValue(pField, pValue)
	}

	if err := r.saveHasOne(doc, &modification); err != nil {
		return err
	}

	if err := r.saveHasMany(doc, &modification, true); err != nil {
		return err
	}

	return nil
}

// MustInsert an record to database.
// It'll panic if any error occurred.
func (r repository) MustInsert(record interface{}, modifiers ...Modifier) {
	must(r.Insert(record, modifiers...))
}

func (r repository) InsertAll(records interface{}) error {
	if records == nil {
		return nil
	}

	var (
		col  = NewCollection(records)
		mods = make([]Modification, col.Len())
	)

	for i := range mods {
		doc := col.Get(i)
		mods[i] = Apply(doc, newStructset(doc))
	}

	return r.insertAll(col, mods)
}

func (r repository) MustInsertAll(records interface{}) {
	must(r.InsertAll(records))
}

// TODO: support assocs
func (r repository) insertAll(col *Collection, modification []Modification) error {
	if len(modification) == 0 {
		return nil
	}

	var (
		pField   = col.PrimaryField()
		queriers = Build(col.Table())
		fields   = make([]string, 0, len(modification[0].Modifies))
		fieldMap = make(map[string]struct{}, len(modification[0].Modifies))
		modifies = make([]map[string]Modify, len(modification))
	)

	for i := range modification {
		for field := range modification[i].Modifies {
			if _, exist := fieldMap[field]; !exist {
				fieldMap[field] = struct{}{}
				fields = append(fields, field)
			}
		}
		modifies[i] = modification[i].Modifies
	}

	ids, err := r.adapter.InsertAll(queriers, fields, modifies, r.logger...)
	if err != nil {
		return err
	}

	// TODO: reload
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
func (r repository) Update(record interface{}, modifiers ...Modifier) error {
	// TODO: perform reference check on library level for record instead of adapter level
	// TODO: make sure primary id not changed
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
		modification = Apply(doc, newStructset(doc))
	} else {
		modification = Apply(doc, modifiers...)
	}

	if len(modification.Assoc) > 0 {
		return r.Transaction(func(r Repository) error {
			return r.(*repository).update(doc, modification, Eq(pField, pValue))
		})
	}

	return r.update(doc, modification, Eq(pField, pValue))
}

func (r repository) update(doc *Document, modification Modification, filter FilterQuery) error {
	if err := r.saveBelongsTo(doc, &modification); err != nil {
		return err
	}

	if len(modification.Modifies) != 0 {
		var (
			queriers = Build(doc.Table(), filter)
		)

		if err := r.adapter.Update(queriers, modification.Modifies, r.logger...); err != nil {
			return err
		}

		if modification.Reload {
			if err := r.find(doc, queriers); err != nil {
				return err
			}
		}
	}

	if err := r.saveHasOne(doc, &modification); err != nil {
		return err
	}

	if err := r.saveHasMany(doc, &modification, false); err != nil {
		return err
	}

	return nil
}

// MustUpdate an record in database.
// It'll panic if any error occurred.
func (r repository) MustUpdate(record interface{}, modifiers ...Modifier) {
	must(r.Update(record, modifiers...))
}

// TODO: support deletion
func (r repository) saveBelongsTo(doc *Document, modification *Modification) error {
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

			if err := r.update(assocDoc, assocMod, filter); err != nil {
				return err
			}
		} else {
			if err := r.insert(assocDoc, assocMod); err != nil {
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
func (r repository) saveHasOne(doc *Document, modification *Modification) error {
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

			if err := r.update(assocDoc, assocMod, filter); err != nil {
				return err
			}
		} else {
			assocMod.Add(Set(fField, rValue))

			if err := r.insert(assocDoc, assocMod); err != nil {
				return err
			}
		}

		assocDoc.SetValue(fField, rValue)
	}

	return nil
}

// saveHasMany expects has many modification to be ordered the same as the recrods in collection.
func (r repository) saveHasMany(doc *Document, modification *Modification, insertion bool) error {
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
				if err := r.deleteAll(Build(table, filter)); err != nil {
					return err
				}
			} else if len(deletedIDs) > 0 {
				if err := r.deleteAll(Build(table, filter.AndIn(pField, deletedIDs...))); err != nil {
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

				if err := r.update(assocDoc, mods[i], filter); err != nil {
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

			if err := r.insertAll(insertCol, insertMods); err != nil {
				return err
			}
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
		q      = Build(table, Eq(pField, pValue))
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
		q = Build("", queriers...)
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
		query    = Build(table, append(queriers, In(keyField, ids...))...)
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
