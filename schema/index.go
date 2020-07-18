package schema

// IndexType definition.
type IndexType string

const (
	// SimpleIndex IndexType.
	SimpleIndex IndexType = "INDEX"
	// UniqueIndex IndexType.
	UniqueIndex IndexType = "UNIQUE INDEX"
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
	Op        Op
	Name      string
	Type      IndexType
	Columns   []string
	NewName   string
	Reference ForeignKeyReference
	Comment   string
	Options   string
}

func addIndex(columns []string, typ IndexType, options []IndexOption) Index {
	index := Index{
		Op:      Add,
		Columns: columns,
		Type:    typ,
	}

	applyIndexOptions(&index, options)
	return index
}

// TODO: support multi columns
func addForeignKey(column string, refTable string, refColumn string, options []IndexOption) Index {
	index := Index{
		Op:      Add,
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
		Op:      Rename,
		Name:    name,
		NewName: newName,
	}

	applyIndexOptions(&index, options)
	return index
}

func dropIndex(name string, options []IndexOption) Index {
	index := Index{
		Op:   Drop,
		Name: name,
	}

	applyIndexOptions(&index, options)
	return index
}
