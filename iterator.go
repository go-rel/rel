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

type iterator struct {
	batch   int
	current int
	query   Query
	adapter Adapter
	cursor  Cursor
	fields  []string
}

func (i iterator) Close() error {
	return i.cursor.Close()
}

func (i *iterator) Next(ctx context.Context, record interface{}) error {
	if i.current%i.batch == 0 {
		if err := i.fetch(ctx); err != nil {
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

	return i.cursor.Scan(scanners...)
}

func (i *iterator) fetch(ctx context.Context) error {
	var (
		query = i.query.Limit(i.batch).Offset(i.current)
	)

	cursor, err := i.adapter.Query(ctx, query)
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
