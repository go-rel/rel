package rel

// KeyType definition.
type KeyType string

const (
	// PrimaryKey KeyType.
	PrimaryKey KeyType = "PRIMARY KEY"
	// ForeignKey KeyType.
	ForeignKey KeyType = "FOREIGN KEY"
	// UniqueKey KeyType.
	UniqueKey = "UNIQUE"
)

// ForeignKeyReference definition.
type ForeignKeyReference struct {
	Table    string
	Columns  []string
	OnDelete string
	OnUpdate string
}

// Key definition.
type Key struct {
	Op        SchemaOp
	Name      string
	Type      KeyType
	Columns   []string
	Rename    string
	Reference ForeignKeyReference
	Options   string
}

func (Key) internalTableDefinition() {}

func createKeys(columns []string, typ KeyType, options []KeyOption) Key {
	key := Key{
		Op:      SchemaCreate,
		Columns: columns,
		Type:    typ,
	}

	applyKeyOptions(&key, options)
	return key
}

func createPrimaryKeys(columns []string, options []KeyOption) Key {
	return createKeys(columns, PrimaryKey, options)
}

func createForeignKey(column string, refTable string, refColumn string, options []KeyOption) Key {
	key := Key{
		Op:      SchemaCreate,
		Type:    ForeignKey,
		Columns: []string{column},
		Reference: ForeignKeyReference{
			Table:   refTable,
			Columns: []string{refColumn},
		},
	}

	applyKeyOptions(&key, options)
	return key
}

// TODO: Rename and Drop, PR welcomed.
