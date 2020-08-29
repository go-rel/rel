package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTable(t *testing.T) {
	var table Table

	t.Run("Column", func(t *testing.T) {
		table.Column("column", String, Comment("column"))
		assert.Equal(t, Column{
			Name:    "column",
			Type:    String,
			Comment: "column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Bool", func(t *testing.T) {
		table.Bool("boolean", Comment("boolean"))
		assert.Equal(t, Column{
			Name:    "boolean",
			Type:    Bool,
			Comment: "boolean",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Int", func(t *testing.T) {
		table.Int("integer", Comment("integer"))
		assert.Equal(t, Column{
			Name:    "integer",
			Type:    Int,
			Comment: "integer",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("BigInt", func(t *testing.T) {
		table.BigInt("bigint", Comment("bigint"))
		assert.Equal(t, Column{
			Name:    "bigint",
			Type:    BigInt,
			Comment: "bigint",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Float", func(t *testing.T) {
		table.Float("float", Comment("float"))
		assert.Equal(t, Column{
			Name:    "float",
			Type:    Float,
			Comment: "float",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Decimal", func(t *testing.T) {
		table.Decimal("decimal", Comment("decimal"))
		assert.Equal(t, Column{
			Name:    "decimal",
			Type:    Decimal,
			Comment: "decimal",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("String", func(t *testing.T) {
		table.String("string", Comment("string"))
		assert.Equal(t, Column{
			Name:    "string",
			Type:    String,
			Comment: "string",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Text", func(t *testing.T) {
		table.Text("text", Comment("text"))
		assert.Equal(t, Column{
			Name:    "text",
			Type:    Text,
			Comment: "text",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Binary", func(t *testing.T) {
		table.Binary("binary", Comment("binary"))
		assert.Equal(t, Column{
			Name:    "binary",
			Type:    Binary,
			Comment: "binary",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Date", func(t *testing.T) {
		table.Date("date", Comment("date"))
		assert.Equal(t, Column{
			Name:    "date",
			Type:    Date,
			Comment: "date",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DateTime", func(t *testing.T) {
		table.DateTime("datetime", Comment("datetime"))
		assert.Equal(t, Column{
			Name:    "datetime",
			Type:    DateTime,
			Comment: "datetime",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Time", func(t *testing.T) {
		table.Time("time", Comment("time"))
		assert.Equal(t, Column{
			Name:    "time",
			Type:    Time,
			Comment: "time",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Timestamp", func(t *testing.T) {
		table.Timestamp("timestamp", Comment("timestamp"))
		assert.Equal(t, Column{
			Name:    "timestamp",
			Type:    Timestamp,
			Comment: "timestamp",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Index", func(t *testing.T) {
		table.Index([]string{"id"}, PrimaryKey, Comment("primary key"))
		assert.Equal(t, Index{
			Columns: []string{"id"},
			Type:    PrimaryKey,
			Comment: "primary key",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Unique", func(t *testing.T) {
		table.Unique([]string{"username"}, Comment("unique"))
		assert.Equal(t, Index{
			Columns: []string{"username"},
			Type:    UniqueIndex,
			Comment: "unique",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("PrimaryKey", func(t *testing.T) {
		table.PrimaryKey([]string{"id"}, Comment("primary key"))
		assert.Equal(t, Index{
			Columns: []string{"id"},
			Type:    PrimaryKey,
			Comment: "primary key",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("ForeignKey", func(t *testing.T) {
		table.ForeignKey("user_id", "users", "id", Comment("foreign key"))
		assert.Equal(t, Index{
			Columns: []string{"user_id"},
			Type:    ForeignKey,
			Reference: ForeignKeyReference{
				Table:   "users",
				Columns: []string{"id"},
			},
			Comment: "foreign key",
		}, table.Definitions[len(table.Definitions)-1])
	})
}

func TestAlterTable(t *testing.T) {
	var table AlterTable

	t.Run("RenameColumn", func(t *testing.T) {
		table.RenameColumn("column", "new_column")
		assert.Equal(t, Column{
			Op:      SchemaRename,
			Name:    "column",
			NewName: "new_column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("AlterColumn", func(t *testing.T) {
		table.AlterColumn("column", Bool, Comment("column"))
		assert.Equal(t, Column{
			Op:      SchemaAlter,
			Name:    "column",
			Type:    Bool,
			Comment: "column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DropColumn", func(t *testing.T) {
		table.DropColumn("column")
		assert.Equal(t, Column{
			Op:   SchemaDrop,
			Name: "column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("RenameIndex", func(t *testing.T) {
		table.RenameIndex("index", "new_index")
		assert.Equal(t, Index{
			Op:      SchemaRename,
			Name:    "index",
			NewName: "new_index",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DropIndex", func(t *testing.T) {
		table.DropIndex("index")
		assert.Equal(t, Index{
			Op:   SchemaDrop,
			Name: "index",
		}, table.Definitions[len(table.Definitions)-1])
	})
}

func TestCreateTable(t *testing.T) {
	var (
		options = []TableOption{
			Comment("comment"),
			Options("options"),
		}
		table = createTable("table", options)
	)

	assert.Equal(t, Table{
		Name:    "table",
		Comment: "comment",
		Options: "options",
	}, table)
}
