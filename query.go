package grimoire

import (
	"reflect"
	"strings"
	"time"

	"github.com/Fs02/go-paranoid"
	"github.com/Fs02/grimoire/c"
	"github.com/Fs02/grimoire/changeset"
	"github.com/Fs02/grimoire/errors"
)

type Query struct {
	repo            *Repo
	Collection      string
	Fields          []string
	AsDistinct      bool
	JoinClause      []c.Join
	Condition       c.Condition
	GroupFields     []string
	HavingCondition c.Condition
	OrderClause     []c.Order
	OffsetResult    int
	LimitResult     int
	Changes         map[string]interface{}
}

// Select filter fields to be selected from database.
func (query Query) Select(fields ...string) Query {
	query.Fields = fields
	return query
}

// Distinct add distinct option to select query.
func (query Query) Distinct() Query {
	query.AsDistinct = true
	return query
}

func (query Query) Join(collection string, condition ...c.Condition) Query {
	return query.JoinWith("JOIN", collection, condition...)
}

func (query Query) JoinWith(mode string, collection string, condition ...c.Condition) Query {
	if len(condition) == 0 {
		query.JoinClause = append(query.JoinClause, c.Join{
			Mode:       mode,
			Collection: collection,
			Condition: c.And(c.Eq(
				c.I(query.Collection+"."+strings.TrimSuffix(collection, "s")+"_id"),
				c.I(collection+".id"),
			)),
		})
	} else {
		query.JoinClause = append(query.JoinClause, c.Join{
			Mode:       mode,
			Collection: collection,
			Condition:  c.And(condition...),
		})
	}

	return query
}

// Where expressions are used to filter the result set. If there is more than one where expression, they are combined with an and operator
func (query Query) Where(condition ...c.Condition) Query {
	query.Condition = query.Condition.And(condition...)
	return query
}

// OrWhere behaves exactly the same as where except it combines with any previous expression by using an OR
func (query Query) OrWhere(condition ...c.Condition) Query {
	query.Condition = query.Condition.Or(c.And(condition...))
	return query
}

func (query Query) Group(fields ...string) Query {
	query.GroupFields = fields
	return query
}

func (query Query) Having(condition ...c.Condition) Query {
	query.HavingCondition = query.HavingCondition.And(condition...)
	return query
}

func (query Query) OrHaving(condition ...c.Condition) Query {
	query.HavingCondition = query.HavingCondition.Or(c.And(condition...))
	return query
}

// Order the result returned by database.
func (query Query) Order(order ...c.Order) Query {
	query.OrderClause = append(query.OrderClause, order...)
	return query
}

// Offset the result returned by database.
func (query Query) Offset(offset int) Query {
	query.OffsetResult = offset
	return query
}

// Limit result returned by database.
func (query Query) Limit(limit int) Query {
	query.LimitResult = limit
	return query
}

// Find adds where id=? into query.
// This is short cut for Where(Eq(I("id"), 1))
func (query Query) Find(id interface{}) Query {
	return query.Where(c.Eq(c.I(query.Collection+".id"), id))
}

// Set value for insert or update operation that will replace changeset value.
func (query Query) Set(field string, value interface{}) Query {
	if query.Changes == nil {
		query.Changes = make(map[string]interface{})
	}

	query.Changes[field] = value
	return query
}

// One retrieves one result that match the query.
// If no result found, it'll return not found error.
func (query Query) One(record interface{}) error {
	query.LimitResult = 1
	count, err := query.repo.adapter.All(query, record)

	if err != nil {
		return errors.Wrap(err)
	} else if count == 0 {
		return errors.NotFoundError("no result found")
	} else {
		return nil
	}
}

// MustOne retrieves one result that match the query.
// If no result found, it'll panic.
func (query Query) MustOne(record interface{}) {
	paranoid.Panic(query.One(record))
}

// All retrieves all results that match the query.
func (query Query) All(record interface{}) error {
	_, err := query.repo.adapter.All(query, record)
	return err
}

// MustAll retrieves all results that match the query.
// It'll panic if any error eccured.
func (query Query) MustAll(record interface{}) {
	paranoid.Panic(query.All(record))
}

// Insert records to database.
func (query Query) Insert(record interface{}, chs ...*changeset.Changeset) error {
	var ids []interface{}

	if len(chs) == 1 {
		// single insert
		ch := chs[0]
		changes := make(map[string]interface{})
		cloneChangeset(changes, ch.Changes())
		putTimestamp(changes, "created_at", ch.Types())
		putTimestamp(changes, "updated_at", ch.Types())
		cloneQuery(changes, query.Changes)

		id, err := query.repo.adapter.Insert(query, changes)
		if err != nil {
			return errors.Wrap(err)
		}
		ids = append(ids, id)
	} else if len(chs) > 1 {
		// multiple insert
		fields := getFields(query, chs)

		allchanges := make([]map[string]interface{}, len(chs))
		for i, ch := range chs {
			changes := make(map[string]interface{})
			cloneChangeset(changes, ch.Changes())
			putTimestamp(changes, "created_at", ch.Types())
			putTimestamp(changes, "updated_at", ch.Types())
			cloneQuery(changes, query.Changes)

			allchanges[i] = changes
		}

		var err error
		ids, err = query.repo.adapter.InsertAll(query, fields, allchanges)
		if err != nil {
			return errors.Wrap(err)
		}
	} else if len(query.Changes) > 0 {
		// set only
		id, err := query.repo.adapter.Insert(query, query.Changes)
		if err != nil {
			return errors.Wrap(err)
		}
		ids = append(ids, id)
	}

	if record == nil || len(ids) == 0 {
		return nil
	} else if len(ids) == 1 {
		return errors.Wrap(query.Find(ids[0]).One(record))
	} else {
		return errors.Wrap(query.Where(c.In(c.I("id"), ids...)).All(record))
	}
}

// MustInsert records to database.
// It'll panic if any error occurred.
func (query Query) MustInsert(record interface{}, chs ...*changeset.Changeset) {
	paranoid.Panic(query.Insert(record, chs...))
}

// Update records in database.
// It'll panic if any error occurred.
func (query Query) Update(record interface{}, chs ...*changeset.Changeset) error {
	changes := make(map[string]interface{})

	// only take the first changeset if any
	if len(chs) != 0 {
		cloneChangeset(changes, chs[0].Changes())
		putTimestamp(changes, "updated_at", chs[0].Types())
	}

	cloneQuery(changes, query.Changes)

	// nothing to update
	if len(changes) == 0 {
		return nil
	}

	// perform update
	err := query.repo.adapter.Update(query, changes)
	if err != nil {
		return errors.Wrap(err)
	}

	// should not fetch updated record(s) if not necessery
	if record != nil {
		return errors.Wrap(query.All(record))
	}

	return nil
}

// MustUpdate records in database.
// It'll panic if any error occurred.
func (query Query) MustUpdate(record interface{}, chs ...*changeset.Changeset) {
	paranoid.Panic(query.Update(record, chs...))
}

// Delete deletes all results that match the query.
func (query Query) Delete() error {
	return errors.Wrap(query.repo.adapter.Delete(query))
}

// MustDelete deletes all results that match the query.
// It'll panic if any error eccured.
func (query Query) MustDelete() {
	paranoid.Panic(query.Delete())
}

func cloneChangeset(out map[string]interface{}, changes map[string]interface{}) {
	for k, v := range changes {
		// filter out changeset
		if _, ok := v.(*changeset.Changeset); ok {
			continue
		}

		if _, ok := v.([]*changeset.Changeset); ok {
			continue
		}

		out[k] = v
	}
}

func cloneQuery(out map[string]interface{}, changes map[string]interface{}) {
	for k, v := range changes {
		out[k] = v
	}
}

func putTimestamp(out map[string]interface{}, field string, types map[string]reflect.Type) {
	if typ, ok := types[field]; ok && typ == reflect.TypeOf(time.Time{}) {
		out[field] = time.Now().Round(time.Second)
	}
}

func getFields(query Query, chs []*changeset.Changeset) []string {
	fields := make([]string, 0, len(chs[0].Types()))

	for f := range chs[0].Types() {
		if f == "created_at" || f == "updated_at" {
			fields = append(fields, f)
			continue
		}

		if _, exist := query.Changes[f]; exist {
			fields = append(fields, f)
		}

		for _, ch := range chs {
			if _, exist := ch.Changes()[f]; exist {
				fields = append(fields, f)
				break
			}
		}
	}

	return fields
}
