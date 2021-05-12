package rel

import (
	"context"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIterator(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		query   = From("users")
		cur1    = createCursor(5)
		cur2    = createCursor(5)
		cur3    = createCursor(3)
		options = []IteratorOption{BatchSize(5)}
		it      = newIterator(context.TODO(), adapter, query, options)
	)

	query = query.From("users").SortAsc("id").Limit(5)
	adapter.On("Query", query).Return(cur1, nil).Once()
	adapter.On("Query", query.Offset(5)).Return(cur2, nil).Once()
	adapter.On("Query", query.Offset(10)).Return(cur3, nil).Once()

	recordsCount := 0
	for {
		if err := it.Next(&user); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
		}

		assert.NotEqual(t, 0, user.ID)
		recordsCount++
	}
	it.Close()

	assert.Equal(t, 13, recordsCount)

	// the last next is not called because it's already refetched.
	// call here to make expectation pass.
	cur1.Next()
	cur2.Next()

	adapter.AssertExpectations(t)
	cur1.AssertExpectations(t)
	cur2.AssertExpectations(t)
	cur3.AssertExpectations(t)
}

func TestIterator_setTableName(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		query   = Query{}
		cur     = createCursor(1)
		it      = newIterator(context.TODO(), adapter, query, nil)
	)

	adapter.On("Query", query.From("users").SortAsc("id").Limit(1000)).Return(cur, nil).Once()

	for {
		if err := it.Next(&user); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
		}

		assert.NotEqual(t, 0, user.ID)
	}
	it.Close()

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestIterator_setStartAndFinishID(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		query   = From("users")
		cur     = createCursor(1)
		options = []IteratorOption{Start(10), Finish(20)}
		it      = newIterator(context.TODO(), adapter, query, options)
	)

	adapter.On("Query", query.Where(Gte("id", 10).AndLte("id", 20)).SortAsc("id").Limit(1000)).Return(cur, nil).Once()

	for {
		if err := it.Next(&user); err == io.EOF {
			break
		} else {
			assert.Nil(t, err)
		}

		assert.NotEqual(t, 0, user.ID)
	}
	it.Close()

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestIterator_cursorFieldsError(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		query   = From("users")
		cur     = &testCursor{}
		it      = newIterator(context.TODO(), adapter, query, nil)
		err     = errors.New("cursor error")
	)

	adapter.On("Query", query.SortAsc("id").Limit(1000)).Return(cur, nil).Once()
	cur.On("Fields").Return([]string{}, err)

	defer it.Close()
	for {
		if err := it.Next(&user); err == io.EOF {
			break
		} else {
			assert.Equal(t, err, err)
		}

		assert.Equal(t, 0, user.ID)
		break
	}

	adapter.AssertExpectations(t)
	cur.AssertExpectations(t)
}

func TestIterator_queryError(t *testing.T) {
	var (
		user    User
		adapter = &testAdapter{}
		query   = From("users")
		cur     = &testCursor{}
		it      = newIterator(context.TODO(), adapter, query, nil)
		err     = errors.New("query error")
	)

	adapter.On("Query", query.SortAsc("id").Limit(1000)).Return(cur, err).Once()

	defer it.Close()
	for {
		if err := it.Next(&user); err == io.EOF {
			break
		} else {
			assert.Equal(t, err, err)
		}

		assert.Equal(t, 0, user.ID)
		break
	}

	adapter.AssertExpectations(t)
}
