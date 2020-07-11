package schema

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

	t.Run("Boolean", func(t *testing.T) {
		table.Boolean("boolean", Comment("boolean"))
		assert.Equal(t, Column{
			Name:    "boolean",
			Type:    Boolean,
			Comment: "boolean",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("Integer", func(t *testing.T) {
		table.Integer("integer", Comment("integer"))
		assert.Equal(t, Column{
			Name:    "integer",
			Type:    Integer,
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
		table.PrimaryKey("id", Comment("primary key"))
		assert.Equal(t, Index{
			Columns: []string{"id"},
			Type:    PrimaryKey,
			Comment: "primary key",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("ForeignKey", func(t *testing.T) {
		table.ForeignKey("user_id", "users", "id", Comment("foreign key"))
		assert.Equal(t, Index{
			Columns: []string{"user_id", "users", "id"},
			Type:    ForeignKey,
			Comment: "foreign key",
		}, table.Definitions[len(table.Definitions)-1])
	})
}

func TestAlterTable(t *testing.T) {
	var table AlterTable

	t.Run("RenameColumn", func(t *testing.T) {
		table.RenameColumn("column", "new_column")
		assert.Equal(t, Column{
			Op:      Rename,
			Name:    "column",
			NewName: "new_column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("AlterColumn", func(t *testing.T) {
		table.AlterColumn("column", Boolean, Comment("column"))
		assert.Equal(t, Column{
			Op:      Alter,
			Name:    "column",
			Type:    Boolean,
			Comment: "column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DropColumn", func(t *testing.T) {
		table.DropColumn("column")
		assert.Equal(t, Column{
			Op:   Drop,
			Name: "column",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("RenameIndex", func(t *testing.T) {
		table.RenameIndex("index", "new_index")
		assert.Equal(t, Index{
			Op:      Rename,
			Name:    "index",
			NewName: "new_index",
		}, table.Definitions[len(table.Definitions)-1])
	})

	t.Run("DropIndex", func(t *testing.T) {
		table.DropIndex("index")
		assert.Equal(t, Index{
			Op:   Drop,
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
