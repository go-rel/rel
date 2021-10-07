package rel

// Index definition.
type Index struct {
	Op       SchemaOp
	Table    string
	Name     string
	Unique   bool
	Columns  []string
	Optional bool
	Filter   FilterQuery
	Options  string
}

func (i Index) description() string {
	return i.Op.String() + " index " + i.Name + " on " + i.Table
}

func (Index) internalMigration() {}

func createIndex(table string, name string, columns []string, options []IndexOption) Index {
	index := Index{
		Op:      SchemaCreate,
		Table:   table,
		Name:    name,
		Columns: columns,
	}

	applyIndexOptions(&index, options)
	return index
}

func createUniqueIndex(table string, name string, columns []string, options []IndexOption) Index {
	index := createIndex(table, name, columns, options)
	index.Unique = true
	return index
}

func dropIndex(table string, name string, options []IndexOption) Index {
	index := Index{
		Op:    SchemaDrop,
		Table: table,
		Name:  name,
	}

	applyIndexOptions(&index, options)
	return index
}

// IndexOption interface.
// Available options are: Comment, Options.
type IndexOption interface {
	applyIndex(index *Index)
}

func applyIndexOptions(index *Index, options []IndexOption) {
	for i := range options {
		options[i].applyIndex(index)
	}
}

// Name option for defining custom index name.
type Name string

func (n Name) applyKey(key *Key) {
	key.Name = string(n)
}
