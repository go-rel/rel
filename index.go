package rel

// IndexType definition.
type IndexType string

const (
	// SimpleIndex IndexType.
	SimpleIndex IndexType = "INDEX"
	// UniqueIndex IndexType.
	UniqueIndex IndexType = "UNIQUE"
)

// Index definition.
type Index struct {
	Op      SchemaOp
	Table   string
	Name    string
	Type    IndexType
	Columns []string
	Options string
}

func (Index) internalMigration() {}

func createIndex(table string, columns []string, typ IndexType, options []IndexOption) Index {
	index := Index{
		Op:      SchemaCreate,
		Table:   table,
		Columns: columns,
		Type:    typ,
	}

	applyIndexOptions(&index, options)
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

func (n Name) applyIndex(index *Index) {
	index.Name = string(n)
}

func (n Name) applyKey(key *Key) {
	key.Name = string(n)
}
