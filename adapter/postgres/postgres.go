package postgres

import (
	db "database/sql"

	"github.com/Fs02/grimoire"
	"github.com/Fs02/grimoire/adapter/sql"
	_ "github.com/lib/pq"
)

// Adapter definition for mysql database.
type Adapter struct {
	*sql.Adapter
}

var _ grimoire.Adapter = (*Adapter)(nil)

// Open mysql connection using dsn.
func Open(dsn string) (*Adapter, error) {
	var err error
	adapter := &Adapter{
		&sql.Adapter{
			Placeholder: "$",
			IsOrdinal:   true,
			ErrorFunc:   errorFunc,
		},
	}

	adapter.DB, err = db.Open("postgres", dsn)
	return adapter, err
}

// Insert inserts a record to database and returns its id.
func (adapter *Adapter) Insert(query grimoire.Query, changes map[string]interface{}, logger grimoire.Logger) (interface{}, error) {
	statement, args := sql.NewBuilder(adapter.Placeholder, adapter.IsOrdinal).
		Returning("id").
		Insert(query.Collection, changes)

	var result struct {
		ID int64
	}

	_, err := adapter.Query(&result, statement, args, logger)
	return result.ID, err
}

// InsertAll inserts all record to database and returns its ids.
func (adapter *Adapter) InsertAll(query grimoire.Query, fields []string, allchanges []map[string]interface{}, logger grimoire.Logger) ([]interface{}, error) {
	statement, args := sql.NewBuilder(adapter.Placeholder, adapter.IsOrdinal).Returning("id").InsertAll(query.Collection, fields, allchanges)

	var result []struct {
		ID int64
	}

	_, err := adapter.Query(&result, statement, args, logger)

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
			Placeholder:   adapter.Placeholder,
			IsOrdinal:     adapter.IsOrdinal,
			IncrementFunc: adapter.IncrementFunc,
			ErrorFunc:     adapter.ErrorFunc,
			Tx:            Tx,
		},
	}, err
}

func errorFunc(err error) error {
	if err == nil {
		return nil
		// } else if e, ok := err.(*mysql.MySQLError); ok && e.Number == 1062 {
		// 	return errors.DuplicateError(e.Message, "")
	}

	return err
}
