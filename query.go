package grimoire

import (
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
	query.JoinClause = append(query.JoinClause, c.Join{
		Mode:       mode,
		Collection: collection,
		Condition:  c.And(condition...),
	})

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
	query.Condition = query.Condition.And(c.Eq(c.I("id"), id))
	return query
}

func (query Query) One(doc interface{}) error {
	query.LimitResult = 1
	qs, args := query.repo.adapter.Find(query)
	return errors.Wrap(query.repo.adapter.Query(doc, qs, args))
}

func (query Query) MustOne(doc interface{}) {
	paranoid.Panic(query.One(doc))
}

func (query Query) All(doc interface{}) error {
	qs, args := query.repo.adapter.Find(query)
	err := query.repo.adapter.Query(doc, qs, args)

	// ignore not found error
	if e, ok := err.(errors.Error); ok && e.NotFoundError() {
		return nil
	}

	return errors.Wrap(err)
}

func (query Query) MustAll(doc interface{}) {
	paranoid.Panic(query.All(doc))
}

func (query Query) Insert(doc interface{}, ch changeset.Changeset) error {
	qs, args := query.repo.adapter.Insert(query, ch)
	id, _, err := query.repo.adapter.Exec(qs, args)
	if err != nil {
		return errors.Wrap(err)
	}

	return errors.Wrap(query.Find(id).One(doc))
}

func (query Query) MustInsert(doc interface{}, ch changeset.Changeset) {
	paranoid.Panic(query.Insert(doc, ch))
}

func (query Query) Update(doc interface{}, ch changeset.Changeset) error {
	qs, args := query.repo.adapter.Update(query, ch)
	_, _, err := query.repo.adapter.Exec(qs, args)
	if err != nil {
		return errors.Wrap(err)
	}

	return errors.Wrap(query.All(doc))
}

func (query Query) MustUpdate(doc interface{}, ch changeset.Changeset) {
	paranoid.Panic(query.Update(doc, ch))
}

func (query Query) Delete() error {
	qs, args := query.repo.adapter.Delete(query)
	_, _, err := query.repo.adapter.Exec(qs, args)
	return errors.Wrap(err)
}

func (query Query) MustDelete() {
	paranoid.Panic(query.Delete())
}

// func (query Query) Replace(doc interface{}) error {
// 	return nil
// }
