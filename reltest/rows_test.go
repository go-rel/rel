package reltest

import (
	"database/sql/driver"
	"errors"
	"strconv"
	"testing"

	"github.com/Fs02/rel"

	"github.com/stretchr/testify/assert"
)

type mockScannerValuer struct {
	value    int
	scanErr  error
	valueErr error
}

func (msv *mockScannerValuer) Scan(src interface{}) error {
	msv.value = src.(int)
	return msv.scanErr
}

func (msv mockScannerValuer) Value() (driver.Value, error) {
	return msv.value, msv.valueErr
}

type rowsTestRecord struct {
	ID    int
	Value mockScannerValuer
}

func TestRows_Scan(t *testing.T) {
	id := 1
	tests := []struct {
		src interface{}
		dst interface{}
	}{
		{
			src: &Author{ID: 1, Name: "Del Piero"},
			dst: &Author{},
		},
		{
			src: &Rating{ID: 2, Score: 100, BookID: 5},
			dst: &Rating{},
		},
		{
			src: &Poster{ID: 3, Image: "image.png", BookID: 6},
			dst: &Poster{},
		},
		{
			src: &Book{ID: 4, Title: "REL for dummies", AuthorID: &id, Views: 1000},
			dst: &Book{},
		},
	}

	for i := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				rows      = newRows(tests[i].src)
				fields, _ = rows.Fields()
				scanners  = rel.NewDocument(tests[i].dst).Scanners(fields)
			)

			assert.NotNil(t, fields)
			assert.True(t, rows.Next())
			assert.Nil(t, rows.NopScanner())
			assert.Nil(t, rows.Scan(scanners...))

			assert.Equal(t, tests[i].src, tests[i].dst)

			assert.False(t, rows.Next())
			assert.Nil(t, rows.Close())
		})
	}
}

func TestRows_Scan_collectiion(t *testing.T) {
	id := 1
	tests := []struct {
		src interface{}
		dst interface{}
	}{
		{
			src: &[]Author{{ID: 1, Name: "Del Piero"}},
			dst: &[]Author{},
		},
		{
			src: &[]Rating{{ID: 2, Score: 100, BookID: 5}},
			dst: &[]Rating{},
		},
		{
			src: &[]Poster{{ID: 3, Image: "image.png", BookID: 6}},
			dst: &[]Poster{},
		},
		{
			src: &[]Book{{ID: 4, Title: "REL for dummies", AuthorID: &id, Views: 1000}},
			dst: &[]Book{},
		},
	}

	for i := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			var (
				rows      = newRows(tests[i].src)
				fields, _ = rows.Fields()
				scanners  = rel.NewCollection(tests[i].dst).Add().Scanners(fields)
			)

			assert.NotNil(t, fields)
			assert.True(t, rows.Next())
			assert.Nil(t, rows.NopScanner())
			assert.Nil(t, rows.Scan(scanners...))

			assert.Equal(t, tests[i].src, tests[i].dst)

			assert.False(t, rows.Next())
			assert.Nil(t, rows.Close())
		})
	}
}

func TestRows_Scan_scanner(t *testing.T) {
	var (
		src       = struct{ Value int }{Value: 1}
		dst       = rowsTestRecord{}
		rows      = newRows(src)
		fields, _ = rows.Fields()
		scanners  = rel.NewDocument(&dst).Scanners(fields)
	)

	assert.NotNil(t, fields)
	assert.True(t, rows.Next())
	assert.Nil(t, rows.NopScanner())
	assert.Nil(t, rows.Scan(scanners...))

	assert.Equal(t, src.Value, dst.Value.value)
	assert.Nil(t, rows.Close())
}

func TestRows_Scan_scannerError(t *testing.T) {
	var (
		src       = struct{ Value int }{Value: 1}
		dst       = rowsTestRecord{Value: mockScannerValuer{scanErr: errors.New("scan error")}}
		rows      = newRows(src)
		fields, _ = rows.Fields()
		scanners  = rel.NewDocument(&dst).Scanners(fields)
	)

	assert.NotNil(t, fields)
	assert.True(t, rows.Next())
	assert.Nil(t, rows.NopScanner())
	assert.Equal(t, errors.New("scan error"), rows.Scan(scanners...))
	assert.Nil(t, rows.Close())
}

func TestRows_Scan_scannerValuer(t *testing.T) {
	var (
		src       = rowsTestRecord{ID: 1, Value: mockScannerValuer{value: 2}}
		dst       = rowsTestRecord{}
		rows      = newRows(src)
		fields, _ = rows.Fields()
		scanners  = rel.NewDocument(&dst).Scanners(fields)
	)

	assert.NotNil(t, fields)
	assert.True(t, rows.Next())
	assert.Nil(t, rows.NopScanner())
	assert.Nil(t, rows.Scan(scanners...))

	assert.Equal(t, src, dst)
	assert.Nil(t, rows.Close())
}

func TestRows_Scan_scannerValuerError(t *testing.T) {
	var (
		src       = rowsTestRecord{ID: 1, Value: mockScannerValuer{value: 2, valueErr: errors.New("value error")}}
		dst       = rowsTestRecord{}
		rows      = newRows(src)
		fields, _ = rows.Fields()
		scanners  = rel.NewDocument(&dst).Scanners(fields)
	)

	assert.NotNil(t, fields)
	assert.True(t, rows.Next())
	assert.Nil(t, rows.NopScanner())
	assert.Equal(t, errors.New("value error"), rows.Scan(scanners...))
	assert.Nil(t, rows.Close())
}

func TestRows_Scan_notAssignable(t *testing.T) {
	var (
		src       = struct{ Value int }{Value: 1}
		dst       = struct{ Value *string }{}
		rows      = newRows(src)
		fields, _ = rows.Fields()
		scanners  = rel.NewDocument(&dst).Scanners(fields)
	)

	assert.NotNil(t, fields)
	assert.True(t, rows.Next())
	assert.Nil(t, rows.NopScanner())
	assert.Equal(t, errors.New("reltest: cannot assign value from type *int to *string"), rows.Scan(scanners...))
	assert.Nil(t, rows.Close())
}

func TestRows_Fields_emptyRows(t *testing.T) {
	var (
		src         = []struct{}{}
		rows        = newRows(src)
		fields, err = rows.Fields()
	)

	assert.Nil(t, fields)
	assert.Nil(t, err)
}
