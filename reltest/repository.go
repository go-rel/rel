package reltest

import (
	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

func must(err error) {
	if err != nil {
		panic(err)
	}
}

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

var _ rel.Repository = (*Repository)(nil)

// Adapter provides a mock function with given fields:
func (r *Repository) Adapter() rel.Adapter {
	return nil
}

// SetLogger provides a mock function with given fields: logger
func (r *Repository) SetLogger(logger ...rel.Logger) {
}

// Aggregate provides a mock function with given fields: query, aggregate, field
func (r *Repository) Aggregate(query rel.Query, aggregate string, field string) (int, error) {
	ret := r.Called(query, aggregate, field)

	var r0 int
	if rf, ok := ret.Get(0).(func(rel.Query, string, string) int); ok {
		r0 = rf(query, aggregate, field)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(rel.Query, string, string) error); ok {
		r1 = rf(query, aggregate, field)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MustAggregate provides a mock function with given fields: query, aggregate, field
func (r *Repository) MustAggregate(query rel.Query, aggregate string, field string) int {
	result, err := r.Aggregate(query, aggregate, field)
	must(err)
	return result
}

// Count provides a mock function with given fields: collection, queriers
func (r *Repository) Count(collection string, queriers ...rel.Querier) (int, error) {
	args := make([]interface{}, len(queriers)+1)
	args[0] = collection

	for i := range queriers {
		args[i+1] = queriers[i]
	}

	ret := r.Called(args...)

	var r0 int
	if rf, ok := ret.Get(0).(func(string, ...rel.Querier) int); ok {
		r0 = rf(collection, queriers...)
	} else {
		r0 = ret.Get(0).(int)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, ...rel.Querier) error); ok {
		r1 = rf(collection, queriers...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MustCount provides a mock function with given fields: collection, queriers
func (r *Repository) MustCount(collection string, queriers ...rel.Querier) int {
	count, err := r.Count(collection, queriers...)
	must(err)
	return count
}

// Find provides a mock function with given fields: record, queriers
func (r *Repository) Find(record interface{}, queriers ...rel.Querier) error {
	return r.findCalled("Find", record, queriers)
}

// MustFind provides a mock function with given fields: record, queriers
func (r *Repository) MustFind(record interface{}, queriers ...rel.Querier) {
	must(r.Find(record, queriers...))
}

// ExpectFind apply mocks and expectations for Find
func (r *Repository) ExpectFind(record interface{}, queriers ...rel.Querier) *ExpectFind {
	return NewExpectFind(r, queriers)
}

// FindAll provides a mock function with given fields: records, queriers
func (r *Repository) FindAll(records interface{}, queriers ...rel.Querier) error {
	return r.findCalled("FindAll", records, queriers)
}

func (r *Repository) findCalled(method string, record interface{}, queriers []rel.Querier) error {
	args := make([]interface{}, len(queriers)+1)
	args[0] = record

	for i := range queriers {
		args[i+1] = queriers[i]
	}

	ret := r.MethodCalled(method, args...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...rel.Querier) error); ok {
		r0 = rf(record, queriers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ExpectFindAll apply mocks and expectations for FindAll
func (r *Repository) ExpectFindAll(record interface{}, queriers ...rel.Querier) *ExpectFindAll {
	return NewExpectFindAll(r, queriers)
}

// MustFindAll provides a mock function with given fields: records, queriers
func (r *Repository) MustFindAll(records interface{}, queriers ...rel.Querier) {
	must(r.FindAll(records, queriers...))
}

// Insert provides a mock function with given fields: record, changers
func (r *Repository) Insert(record interface{}, changers ...rel.Changer) error {
	return r.modifyCalled("Insert", record, changers)
}

// MustInsert provides a mock function with given fields: record, changers
func (r *Repository) MustInsert(record interface{}, changers ...rel.Changer) {
	must(r.Insert(record, changers...))
}

// InsertAll provides a mock function with given fields: records, changes
func (r *Repository) InsertAll(records interface{}, changes ...rel.Changes) error {
	args := make([]interface{}, len(changes)+1)
	args[0] = records

	for i := range changes {
		args[i+1] = changes[i]
	}

	ret := r.Called(args...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...rel.Changes) error); ok {
		r0 = rf(records, changes...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MustInsertAll provides a mock function with given fields: records, changes
func (r *Repository) MustInsertAll(records interface{}, changes ...rel.Changes) {
	must(r.InsertAll(records, changes...))
}

// Update provides a mock function with given fields: record, changers
func (r *Repository) Update(record interface{}, changers ...rel.Changer) error {
	return r.modifyCalled("Update", record, changers)
}

// MustUpdate provides a mock function with given fields: record, changers
func (r *Repository) MustUpdate(record interface{}, changers ...rel.Changer) {
	must(r.Update(record, changers...))
}

// Save provides a mock function with given fields: record, changers
func (r *Repository) Save(record interface{}, changers ...rel.Changer) error {
	return r.modifyCalled("Save", record, changers)
}

// MustSave provides a mock function with given fields: record, changers
func (r *Repository) MustSave(record interface{}, changers ...rel.Changer) {
	must(r.Save(record, changers...))
}

func (r *Repository) modifyCalled(method string, record interface{}, changers []rel.Changer) error {
	args := make([]interface{}, len(changers)+1)
	args[0] = record

	for i := range changers {
		args[i+1] = changers[i]
	}

	ret := r.MethodCalled(method, args...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, ...rel.Changer) error); ok {
		r0 = rf(record, changers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Delete provides a mock function with given fields: record
func (r *Repository) Delete(record interface{}) error {
	ret := r.Called(record)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(record)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MustDelete provides a mock function with given fields: record
func (r *Repository) MustDelete(record interface{}) {
	must(r.Delete(record))
}

// DeleteAll provides a mock function with given fields: queriers
func (r *Repository) DeleteAll(queriers ...rel.Querier) error {
	args := make([]interface{}, len(queriers))
	for i := range queriers {
		args[i] = queriers[i]
	}

	ret := r.Called(args...)

	var r0 error
	if rf, ok := ret.Get(0).(func(...rel.Querier) error); ok {
		r0 = rf(queriers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MustDeleteAll provides a mock function with given fields: queriers
func (r *Repository) MustDeleteAll(queriers ...rel.Querier) {
	must(r.DeleteAll(queriers...))
}

// Preload provides a mock function with given fields: records, field, queriers
func (r *Repository) Preload(records interface{}, field string, queriers ...rel.Querier) error {
	args := make([]interface{}, len(queriers)+2)
	args[0] = records
	args[1] = field

	for i := range queriers {
		args[i+2] = queriers[i]
	}

	ret := r.Called(args...)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}, string, ...rel.Querier) error); ok {
		r0 = rf(records, field, queriers...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MustPreload provides a mock function with given fields: records, field, queriers
func (r *Repository) MustPreload(records interface{}, field string, queriers ...rel.Querier) {
	must(r.Preload(records, field, queriers...))
}

// Transaction provides a mock function with given fields: fn
func (r *Repository) Transaction(fn func(rel.Repository) error) error {
	ret := r.Called(fn)

	var r0 error
	if rf, ok := ret.Get(0).(func(func(rel.Repository) error) error); ok {
		r0 = rf(fn)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
