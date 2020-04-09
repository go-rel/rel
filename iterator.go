package rel

import (
	"context"
	"io"
)

// Iterator alllows iterating through all record in database in batch.
type Iterator interface {
	io.Closer
	Next(record interface{}) error
}

// IteratorOption is used to configure iteration behaviour, such as batch size, start id and finish id.
type IteratorOption interface {
	apply(*iterator)
}

// BatchSize specifies the size of iterator batch. Defaults to 1000.
type BatchSize int

func (bs BatchSize) apply(i *iterator) {
	i.batchSize = int(bs)
}

// Start specfies the primary value to start from (inclusive).
type Start int

func (s Start) apply(i *iterator) {
	i.start = int(s)
}

// Finish specfies the primary value to finish at (inclusive).
type Finish int

func (f Finish) apply(i *iterator) {
	i.finish = int(f)
}

type iterator struct {
	ctx       context.Context
	start     interface{}
	finish    interface{}
	batchSize int
	current   int
	query     Query
	adapter   Adapter
	cursor    Cursor
	fields    []string
	closed    bool
}

func (i *iterator) Close() error {
	if !i.closed && i.cursor != nil {
		i.closed = true
		return i.cursor.Close()
	}

	return nil
}

func (i *iterator) Next(record interface{}) error {
	if i.current%i.batchSize == 0 {
		if err := i.fetch(i.ctx, record); err != nil {
			return err
		}
	}

	if !i.cursor.Next() {
		return io.EOF
	}

	var (
		doc      = NewDocument(record)
		scanners = doc.Scanners(i.fields)
	)

	i.current++
	return i.cursor.Scan(scanners...)
}

func (i *iterator) fetch(ctx context.Context, record interface{}) error {
	if i.current == 0 {
		i.init(record)
	} else {
		i.cursor.Close()
	}

	i.query = i.query.Limit(i.batchSize).Offset(i.current)

	cursor, err := i.adapter.Query(ctx, i.query)
	if err != nil {
		return err
	}

	fields, err := cursor.Fields()
	if err != nil {
		return err
	}

	i.cursor = cursor
	i.fields = fields

	return nil
}

func (i *iterator) init(record interface{}) {
	var (
		doc = NewDocument(record)
	)

	if i.query.Table == "" {
		i.query.Table = doc.Table()
	}

	if i.start != nil {
		i.query = i.query.Where(Gte(doc.PrimaryField(), i.start))
	}

	if i.finish != nil {
		i.query = i.query.Where(Lte(doc.PrimaryField(), i.finish))
	}

	i.query = i.query.SortAsc(doc.PrimaryField())
}

func newIterator(ctx context.Context, adapter Adapter, query Query, options []IteratorOption) Iterator {
	it := &iterator{
		ctx:       ctx,
		batchSize: 1000,
		query:     query,
		adapter:   adapter,
	}

	for i := range options {
		options[i].apply(it)
	}

	return it
}
