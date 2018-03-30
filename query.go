package grimoire

import (
	"strings"

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

func (query Query) Select(fields ...string) Query {
	query.Fields = fields
	return query
}

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

func (query Query) Order(order ...c.Order) Query {
	query.OrderClause = append(query.OrderClause, order...)
	return query
}

func (query Query) Offset(offset int) Query {
	query.OffsetResult = offset
	return query
}

func (query Query) Limit(limit int) Query {
	query.LimitResult = limit
	return query
}

func (query Query) Find(id interface{}) Query {
	return query.Where(c.Eq(c.I("id"), id))
}

func (query Query) Set(field string, value interface{}) Query {
	if query.Changes == nil {
		query.Changes = make(map[string]interface{})
	}

	query.Changes[field] = value
	return query
}

func (query Query) One(doc interface{}) error {
	query.LimitResult = 1
	qs, args := query.repo.adapter.Find(query)
	count, err := query.repo.adapter.Query(doc, qs, args)

	if err != nil {
		return errors.Wrap(err)
	} else if count == 0 {
		return errors.NotFoundError("no result found")
	} else {
		return nil
	}
}

func (query Query) MustOne(doc interface{}) {
	paranoid.Panic(query.One(doc))
}

func (query Query) All(doc interface{}) error {
	qs, args := query.repo.adapter.Find(query)
	_, err := query.repo.adapter.Query(doc, qs, args)
	return errors.Wrap(err)
}

func (query Query) MustAll(doc interface{}) {
	paranoid.Panic(query.All(doc))
}

func (query Query) Insert(doc interface{}, chs ...*changeset.Changeset) error {
	var ids []interface{}

	if len(chs) > 0 {
		for _, ch := range chs {
			changes := make(map[string]interface{})

			for k, v := range ch.Changes() {
				changes[k] = v
			}

			// apply query changes
			for k, v := range query.Changes {
				changes[k] = v
			}

			qs, args := query.repo.adapter.Insert(query, changes)
			id, _, err := query.repo.adapter.Exec(qs, args)
			if err != nil {
				return errors.Wrap(err)
			}

			ids = append(ids, id)
		}
	} else {
		qs, args := query.repo.adapter.Insert(query, query.Changes)
		id, _, err := query.repo.adapter.Exec(qs, args)
		if err != nil {
			return errors.Wrap(err)
		}

		ids = append(ids, id)
	}

	if doc == nil || len(ids) == 0 {
		return nil
	} else if len(ids) == 1 {
		return errors.Wrap(query.Find(ids[0]).One(doc))
	} else {
		return errors.Wrap(query.Where(c.In(c.I("id"), ids...)).All(doc))
	}
}

func (query Query) MustInsert(doc interface{}, chs ...*changeset.Changeset) {
	paranoid.Panic(query.Insert(doc, chs...))
}

func (query Query) Update(doc interface{}, chs ...*changeset.Changeset) error {
	changes := make(map[string]interface{})

	// only take the first changeset if any
	if len(chs) != 0 {
		for k, v := range chs[0].Changes() {
			changes[k] = v
		}
	}

	// apply query changes
	for k, v := range query.Changes {
		changes[k] = v
	}

	// nothing to update
	if len(changes) == 0 {
		return nil
	}

	// perform update
	qs, args := query.repo.adapter.Update(query, changes)
	_, _, err := query.repo.adapter.Exec(qs, args)
	if err != nil {
		return errors.Wrap(err)
	}

	// should not fetch updated record(s) if not necessery
	if doc != nil {
		return errors.Wrap(query.All(doc))
	}

	return nil
}

func (query Query) MustUpdate(doc interface{}, chs ...*changeset.Changeset) {
	paranoid.Panic(query.Update(doc, chs...))
}

func (query Query) Delete() error {
	qs, args := query.repo.adapter.Delete(query)
	_, _, err := query.repo.adapter.Exec(qs, args)
	return errors.Wrap(err)
}

func (query Query) MustDelete() {
	paranoid.Panic(query.Delete())
}
