package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	var table Table

	t.Run("Column", func(t *testing.T) {
		table.Column("column", String)
		assert.Equal(t, Column{
			Name: "column",
			Type: String,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Bool", func(t *testing.T) {
		table.Bool("boolean")
		assert.Equal(t, Column{
			Name: "boolean",
			Type: Bool,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Int", func(t *testing.T) {
		table.Int("integer")
		assert.Equal(t, Column{
			Name: "integer",
			Type: Int,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("BigInt", func(t *testing.T) {
		table.BigInt("bigint")
		assert.Equal(t, Column{
			Name: "bigint",
			Type: BigInt,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Float", func(t *testing.T) {
		table.Float("float")
		assert.Equal(t, Column{
			Name: "float",
			Type: Float,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Decimal", func(t *testing.T) {
		table.Decimal("decimal")
		assert.Equal(t, Column{
			Name: "decimal",
			Type: Decimal,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("String", func(t *testing.T) {
		table.String("string")
		assert.Equal(t, Column{
			Name: "string",
			Type: String,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Text", func(t *testing.T) {
		table.Text("text")
		assert.Equal(t, Column{
			Name: "text",
			Type: Text,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Date", func(t *testing.T) {
		table.Date("date")
		assert.Equal(t, Column{
			Name: "date",
			Type: Date,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DateTime", func(t *testing.T) {
		table.DateTime("datetime")
		assert.Equal(t, Column{
			Name: "datetime",
			Type: DateTime,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Time", func(t *testing.T) {
		table.Time("time")
		assert.Equal(t, Column{
			Name: "time",
			Type: Time,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Timestamp", func(t *testing.T) {
		table.Timestamp("timestamp")
		assert.Equal(t, Column{
			Name: "timestamp",
			Type: Timestamp,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("PrimaryKey", func(t *testing.T) {
		table.PrimaryKey("id")
		assert.Equal(t, Key{
			Columns: []string{"id"},
			Type:    PrimaryKey,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("ForeignKey", func(t *testing.T) {
		table.ForeignKey("user_id", "users", "id")
		assert.Equal(t, Key{
			Columns: []string{"user_id"},
			Type:    ForeignKey,
			Reference: ForeignKeyReference{
				Table:   "users",
				Columns: []string{"id"},
			},
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Unique", func(t *testing.T) {
		table.Unique([]string{"username"})
		assert.Equal(t, Key{
			Columns: []string{"username"},
			Type:    UniqueKey,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Fragment", func(t *testing.T) {
		table.Fragment("SQL")
		assert.Equal(t, Raw("SQL"), table.Definitions[len(table.Definitions)-1])
	})
}

func TestAlterTable(t *testing.T) {
	var table AlterTable

	t.Run("RenameColumn", func(t *testing.T) {
		table.RenameColumn("column", "new_column")
		assert.Equal(t, Column{
			Op:     SchemaRename,
			Name:   "column",
			Rename: "new_column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DropColumn", func(t *testing.T) {
		table.DropColumn("column")
		assert.Equal(t, Column{
			Op:   SchemaDrop,
			Name: "column",
		}, table.Definitions[len(table.Definitions)-1])
	})
}

func TestCreateTable(t *testing.T) {
	var (
		options = []TableOption{
			Options("options"),
			Optional(true),
		}
		table = createTable("table", options)
	)

	assert.Equal(t, Table{
		Name:     "table",
		Optional: true,
		Options:  "options",
	}, table)
}

func TestTable_description(t *testing.T) {
	assert.Equal(t, "create table tests", Table{Name: "tests"}.description())
}
