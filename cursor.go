package rel

import (
	"database/sql"
	"reflect"
)

// Cursor is interface to work with database result (used by adapter).
type Cursor interface {
	Close() error
	Fields() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	NopScanner() interface{} // TODO: conflict with manual scanners interface
}

func scanOne(cur Cursor, doc *Document) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	if !cur.Next() {
		return NotFoundError{}
	}

	var (
		scanners = doc.Scanners(fields)
	)

	return cur.Scan(scanners...)
}

func scanAll(cur Cursor, col *Collection) error {
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

func scanMulti(cur Cursor, keyField string, keyType reflect.Type, cols map[interface{}][]slice) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	var (
		found       = false
		keyValue    = reflect.New(keyType)
		keyScanners = make([]interface{}, len(fields))
	)

	for i, field := range fields {
		if keyField == field {
			found = true
			keyScanners[i] = keyValue.Interface()
		} else {
			// need to create distinct copies
			// otherwise next scan result will be corrupted
			keyScanners[i] = &sql.RawBytes{}
		}
	}

	if !found {
		panic("rel: primary key row does not exists")
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
