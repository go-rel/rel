package reltest

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"reflect"

	"github.com/Fs02/rel"
)

type data interface {
	Len() int
	Get(index int) *rel.Document
}

type rows struct {
	data    data
	current int
}

func (r *rows) Close() error {
	return nil
}

func (r *rows) Fields() ([]string, error) {
	if r.data.Len() > 0 {
		return r.data.Get(0).Fields(), nil
	}

	return nil, nil
}

func (r *rows) Next() bool {
	r.current++
	return r.current <= r.data.Len()
}

func (r *rows) Scan(dsts ...interface{}) error {
	var (
		doc    = r.data.Get(r.current - 1)
		fields = doc.Fields()
	)

	for i := range dsts {
		var (
			dst    = dsts[i]
			src, _ = doc.Value(fields[i])
		)

		if scanner, ok := dst.(sql.Scanner); ok {
			// TODO: convert value to basic type before passing to scanner when it's not coming from valuer.
			if valuer, ok := src.(driver.Valuer); ok {
				value, err := valuer.Value()
				if err != nil {
					return err
				}

				src = value
			}

			if err := scanner.Scan(src); err != nil {
				return err
			}
		} else {
			var (
				sv = reflect.ValueOf(src)
				dv = reflect.ValueOf(dst).Elem()
			)

			if dv.Kind() == reflect.Ptr && sv.Kind() != reflect.Ptr {
				nsv := reflect.New(sv.Type())
				nsv.Elem().Set(sv)
				sv = nsv
			}

			// TODO: convert value.
			if !sv.Type().AssignableTo(dv.Type()) {
				return errors.New("reltest: cannot assign " + fields[i] + " from type " + sv.Type().String() + " to " + dv.Type().String())
			}

			dv.Set(sv)
		}
	}

	return nil
}

func (r *rows) NopScanner() interface{} {
	return nil
}

func newRows(records interface{}) rel.Cursor {
	var (
		data data
		rt   = reflect.TypeOf(records)
	)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		data = rel.NewCollection(records, true)
	} else {
		data = rel.NewDocument(records, true)
	}

	return &rows{
		data: data,
	}
}
