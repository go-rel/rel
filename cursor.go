package rel

import (
	"reflect"
)

// Cursor is interface to work with database result (used by adapter).
type Cursor interface {
	Close() error
	Fields() ([]string, error)
	Next() bool
	Scan(...any) error
	NopScanner() any // TODO: conflict with manual scanners interface
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

func scanMulti(cur Cursor, keyField string, keyType reflect.Type, cols map[any][]slice) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	keyFound := false
	for _, field := range fields {
		if keyField == field {
			keyFound = true
		}
	}

	if !keyFound && fields != nil {
		panic("rel: primary key row does not exists")
	}

	var doc *Document
	for k := range cols {
		for _, col := range cols[k] {
			doc = col.NewDocument()
			break
		}
		break
	}

	// scan the result
	for cur.Next() {
		// scan key
		if err := cur.Scan(doc.Scanners(fields)...); err != nil {
			return err
		}

		key, found := doc.Value(keyField)
		mustTrue(found, "rel: key field not found")

		for _, col := range cols[key] {
			col.Append(doc)
		}

		// create new doc for next scan
		doc = doc.NewDocument()
	}

	return nil
}
