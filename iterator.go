package rel

import (
	"context"
	"fmt"
	"io"
)

// Iterator allows iterating through all record in database in batch.
type Iterator interface {
	io.Closer
	Next(record interface{}) error
}

// IteratorOption is used to configure iteration behaviour, such as batch size, start id and finish id.
type IteratorOption interface {
	apply(*iterator)
}

type batchSize int

func (bs batchSize) apply(i *iterator) {
	i.batchSize = int(bs)
}

// String representation.
func (bs batchSize) String() string {
	return fmt.Sprintf("rel.BatchSize(%d)", bs)
}

// BatchSize specifies the size of iterator batch. Defaults to 1000.
func BatchSize(size int) IteratorOption {
	return batchSize(size)
}

type start []interface{}

func (s start) apply(i *iterator) {
	i.start = s
}

// String representation.
func (s start) String() string {
	return fmt.Sprintf("rel.Start(%s)", fmtifaces(s))
}

// Start specifies the primary value to start from (inclusive).
func Start(id ...interface{}) IteratorOption {
	return start(id)
}

type finish []interface{}

func (f finish) apply(i *iterator) {
	i.finish = f
}

// String representation.
func (f finish) String() string {
	return fmt.Sprintf("rel.Finish(%s)", fmtifaces(f))
}

// Finish specifies the primary value to finish at (inclusive).
func Finish(id ...interface{}) IteratorOption {
	return finish(id)
}

type iterator struct {
	ctx       context.Context
	start     []interface{}
	finish    []interface{}
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

	if len(i.start) > 0 {
		i.query = i.query.Where(filterDocumentPrimary(doc.PrimaryFields(), i.start, FilterGteOp))
	}

	if len(i.finish) > 0 {
		i.query = i.query.Where(filterDocumentPrimary(doc.PrimaryFields(), i.finish, FilterLteOp))
	}

	i.query = i.query.SortAsc(doc.PrimaryFields()...)
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
