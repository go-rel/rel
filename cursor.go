package grimoire

import (
	"github.com/Fs02/grimoire/errors"
)

type Cursor interface {
	Close() error
	Fields() ([]string, error)
	Next() bool
	Scan(...interface{}) error
	NopScanner() interface{}
}

func scanOne(cur Cursor, collec Collection) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	if !cur.Next() {
		return errors.NewUnexpected("TODO: no result")
	}

	var (
		scanners = collec.Add().Scanners(fields)
	)

	return cur.Scan(scanners...)
}

func scanMany(cur Cursor, collec Collection) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	for cur.Next() {
		var (
			doc      = collec.Add()
			scanners = doc.Scanners(fields)
		)

		if err := cur.Scan(scanners...); err != nil {
			return err
		}
	}

	return nil
}

func scanMulti(cur Cursor, keyField string, collecs map[interface{}][]Collection) error {
	defer cur.Close()

	fields, err := cur.Fields()
	if err != nil {
		return err
	}

	var (
		keyScanners = make([]interface{}, len(fields))
		keyIndex    = -1
	)

	for i, field := range fields {
		if keyField == field {
			keyIndex = i
		}

		keyScanners[i] = cur.NopScanner()
	}

	if keyIndex < 0 {
		panic("grimoire: TODO")
	}

	// get the first key
	for k := range collecs {
		keyScanners[keyIndex] = k
		break
	}

	// scan the result
	for cur.Next() {
		// scan key
		if err := cur.Scan(keyScanners...); err != nil {
			return err
		}

		for _, collec := range collecs[keyScanners[keyIndex]] {
			var (
				doc = collec.Add()
			)

			if err := cur.Scan(doc.Scanners(fields)); err != nil {
				return err
			}
		}
	}

	return nil
}
