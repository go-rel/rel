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
			Name:   "column",
			Type:   String,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("ID", func(t *testing.T) {
		table.ID("id")
		assert.Equal(t, Column{
			Name:    "id",
			Type:    ID,
			Primary: true,
			Constr:  ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("BigID", func(t *testing.T) {
		table.BigID("big_id")
		assert.Equal(t, Column{
			Name:    "big_id",
			Type:    BigID,
			Primary: true,
			Constr:  ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("IDNotPrimaryKey", func(t *testing.T) {
		table.ID("id", Primary(false))
		assert.Equal(t, Column{
			Name:   "id",
			Type:   ID,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Bool", func(t *testing.T) {
		table.Bool("boolean")
		assert.Equal(t, Column{
			Name:   "boolean",
			Type:   Bool,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("SmallInt", func(t *testing.T) {
		table.SmallInt("smallint")
		assert.Equal(t, Column{
			Name:   "smallint",
			Type:   SmallInt,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Int", func(t *testing.T) {
		table.Int("integer")
		assert.Equal(t, Column{
			Name:   "integer",
			Type:   Int,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("BigInt", func(t *testing.T) {
		table.BigInt("bigint")
		assert.Equal(t, Column{
			Name:   "bigint",
			Type:   BigInt,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Float", func(t *testing.T) {
		table.Float("float")
		assert.Equal(t, Column{
			Name:   "float",
			Type:   Float,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Decimal", func(t *testing.T) {
		table.Decimal("decimal")
		assert.Equal(t, Column{
			Name:   "decimal",
			Type:   Decimal,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("String", func(t *testing.T) {
		table.String("string")
		assert.Equal(t, Column{
			Name:   "string",
			Type:   String,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Text", func(t *testing.T) {
		table.Text("text")
		assert.Equal(t, Column{
			Name:   "text",
			Type:   Text,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("JSON", func(t *testing.T) {
		table.JSON("json")
		assert.Equal(t, Column{
			Name:   "json",
			Type:   JSON,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Date", func(t *testing.T) {
		table.Date("date")
		assert.Equal(t, Column{
			Name:   "date",
			Type:   Date,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DateTime", func(t *testing.T) {
		table.DateTime("datetime")
		assert.Equal(t, Column{
			Name:   "datetime",
			Type:   DateTime,
			Constr: ColumnConstraintType,
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Time", func(t *testing.T) {
		table.Time("time")
		assert.Equal(t, Column{
			Name:   "time",
			Type:   Time,
			Constr: ColumnConstraintType,
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

func TestTable_Description(t *testing.T) {
	assert.Equal(t, "create table tests", Table{Name: "tests"}.description())
}

func TestTable_InternalMigration(t *testing.T) {
	assert.NotPanics(t, func() { Table{}.internalMigration() })
}
