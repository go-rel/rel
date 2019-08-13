package grimoire

import (
	"reflect"

	"github.com/Fs02/grimoire/errors"
)

type Cursor interface {
	Close() error
	Fields() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	NopScanner() interface{} // TODO: conflict with manual scanners interface
}

func scanOne(cur Cursor, col Collection) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	if !cur.Next() {
		return errors.New("no result found", "", errors.NotFound)
	}

	var (
		scanners = col.Add().Scanners(fields)
	)

	return cur.Scan(scanners...)
}

func scanMany(cur Cursor, col Collection) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	for cur.Next() {
		var (
			doc      = col.Add()
			scanners = doc.Scanners(fields)
		)

		if err := cur.Scan(scanners...); err != nil {
			return err
		}
	}

	return nil
}

func scanMulti(cur Cursor, keyField string, keyType reflect.Type, cols map[interface{}][]Collection) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	var (
		found       = false
		keyValue    = reflect.New(keyType)
		nopScanner  = cur.NopScanner()
		keyScanners = make([]interface{}, len(fields))
	)

	for i, field := range fields {
		if keyField == field {
			found = true
			keyScanners[i] = keyValue.Interface()
		} else {
			keyScanners[i] = nopScanner
		}
	}

	if !found {
		panic("grimoire: TODO")
	}

	// scan the result
	for cur.Next() {
		// scan key
		if err := cur.Scan(keyScanners...); err != nil {
			return err
		}

		var (
			key = reflect.Indirect(keyValue).Interface()
		)

		for _, col := range cols[key] {
			var (
				doc      = col.Add()
				scanners = doc.Scanners(fields)
			)

			if err := cur.Scan(scanners...); err != nil {
				return err
			}
		}
	}

	return nil
}
