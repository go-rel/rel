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
	Name    string
	Type    IndexType
	Columns []string
	NewName string
	Options string
}

func createIndex(columns []string, typ IndexType, options []IndexOption) Index {
	index := Index{
		Op:      SchemaCreate,
		Columns: columns,
		Type:    typ,
	}

	applyIndexOptions(&index, options)
	return index
}

func renameIndex(name string, newName string, options []IndexOption) Index {
	index := Index{
		Op:      SchemaRename,
		Name:    name,
		NewName: newName,
	}

	applyIndexOptions(&index, options)
	return index
}

func dropIndex(name string, options []IndexOption) Index {
	index := Index{
		Op:   SchemaDrop,
		Name: name,
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
