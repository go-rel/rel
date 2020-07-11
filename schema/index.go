package schema

// IndexType definition.
type IndexType string

const (
	// Simple IndexType.
	Simple IndexType = "index"
	// UniqueIndex IndexType.
	UniqueIndex IndexType = "unique"
	// PrimaryKey IndexType.
	PrimaryKey IndexType = "primary key"
	// ForeignKey IndexType.
	ForeignKey IndexType = "foreign key"
)

// Index definition.
type Index struct {
	Op        Op
	Name      string
	Type      IndexType
	Columns   []string
	Reference Reference
	NewName   string
	Comment   string
	Options   string
}

// Reference definition.
type Reference struct {
	Table    string
	Column   string
	OnDelete string
	OnUpdate string
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

func addForeignKey(column string, refTable string, refColumn string, options []IndexOption) Index {
	index := Index{
		Op:      Add,
		Columns: []string{column},
		Reference: Reference{
			Table:  refTable,
			Column: refColumn,
		},
		Type: ForeignKey,
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
