// Package postgres wraps postgres (pq) driver as an adapter for grimoire.
//
// Usage:
//	// open postgres connection.
//	adapter, err := postgres.Open("postgres://postgres@localhost/grimoire_test?sslmode=disable")
//	if err != nil {
//		panic(err)
//	}
//	defer adapter.Close()
//
//	// initialize grimoire's repo.
//	repo := grimoire.New(adapter)
package postgres

import (
	db "database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	"github.com/Fs02/grimoire/errors"
	"github.com/Fs02/grimoire/internal"
	"github.com/lib/pq"
)

// Adapter definition for postgrees database.
type Adapter struct {
	*sql.Adapter
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Open postgrees connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error

	adapter := &Adapter{
		Adapter: &sql.Adapter{
			Config: &sql.Config{
				Placeholder:         "$",
				EscapeChar:          "\"",
				Ordinal:             true,
				InsertDefaultValues: true,
				ErrorFunc:           errorFunc,
			},
		},
	}
	adapter.DB, err = db.Open("postgres", dsn)

	return adapter, err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}, loggers ...grimoire.Logger) (interface{}, error) {
	statement, args := sql.NewBuilder(adapter.Config).Returning("id").Insert(query.Collection, changes)

	var result struct {
		ID int64
	}

	_, err := adapter.Query(&result, statement, args, loggers...)
	return result.ID, err
}

// InsertAll inserts multiple records to database and returns its ids.
func (adapter *Adapter) InsertAll(query grimoire.Query, fields []string, allchanges []map[string]interface{}, loggers ...grimoire.Logger) ([]interface{}, error) {
	statement, args := sql.NewBuilder(adapter.Config).Returning("id").InsertAll(query.Collection, fields, allchanges)

	var result []struct {
		ID int64
	}

	_, err := adapter.Query(&result, statement, args, loggers...)

	ids := make([]interface{}, 0, len(result))
	for _, r := range result {
		ids = append(ids, r.ID)
	}

	return ids, err
}

// Begin begins a new transaction.
func (adapter *Adapter) Begin() (grimoire.Adapter, error) {
	Tx, err := adapter.DB.Begin()

	return &Adapter{
		&sql.Adapter{
			Config: adapter.Config,
			Tx:     Tx,
		},
	}, err
}

func errorFunc(err error) error {
	if err == nil {
		return nil
	}

	if e, ok := err.(*pq.Error); ok {
		switch e.Code {
		case "23505":
			return errors.New(e.Message, internal.ExtractString(e.Message, "constraint \"", "\""), errors.UniqueConstraint)
		case "23503":
			return errors.New(e.Message, internal.ExtractString(e.Message, "constraint \"", "\""), errors.ForeignKeyConstraint)
		case "23514":
			return errors.New(e.Message, internal.ExtractString(e.Message, "constraint \"", "\""), errors.CheckConstraint)
		}
	}

	return err
}
