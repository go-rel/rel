package rel

// IndexType definition.
type IndexType string

const (
	// SimpleIndex IndexType.
	SimpleIndex IndexType = "INDEX"
	// UniqueIndex IndexType.
	UniqueIndex IndexType = "UNIQUE"
	// PrimaryKey IndexType.
	PrimaryKey IndexType = "PRIMARY KEY"
	// ForeignKey IndexType.
	ForeignKey IndexType = "FOREIGN KEY"
)

// ForeignKeyReference definition.
type ForeignKeyReference struct {
	Table    string
	Columns  []string
	OnDelete string
	OnUpdate string
}

// Index definition.
type Index struct {
	Op        SchemaOp
	Name      string
	Type      IndexType
	Columns   []string
	NewName   string
	Reference ForeignKeyReference
	Options   string
}

func addIndex(columns []string, typ IndexType, options []IndexOption) Index {
	index := Index{
		Op:      SchemaAdd,
		Columns: columns,
		Type:    typ,
	}

	applyIndexOptions(&index, options)
	return index
}

// TODO: support multi columns
func addForeignKey(column string, refTable string, refColumn string, options []IndexOption) Index {
	index := Index{
		Op:      SchemaAdd,
		Type:    ForeignKey,
		Columns: []string{column},
		Reference: ForeignKeyReference{
			Table:   refTable,
			Columns: []string{refColumn},
		},
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
