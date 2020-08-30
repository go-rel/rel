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
	NewName   string
	Reference ForeignKeyReference
	Options   string
}

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

func renameKey(name string, newName string, options []KeyOption) Key {
	key := Key{
		Op:      SchemaRename,
		Name:    name,
		NewName: newName,
	}

	applyKeyOptions(&key, options)
	return key
}

func dropKey(name string, options []KeyOption) Key {
	key := Key{
		Op:   SchemaDrop,
		Name: name,
	}

	applyKeyOptions(&key, options)
	return key
}

// KeyOption interface.
// Available options are: Comment, Options.
type KeyOption interface {
	applyKey(key *Key)
}

func applyKeyOptions(key *Key, options []KeyOption) {
	for i := range options {
		options[i].applyKey(key)
	}
}

// OnDelete option for foreign key.
type OnDelete string

func (od OnDelete) applyKey(key *Key) {
	key.Reference.OnDelete = string(od)
}

// OnUpdate option for foreign key.
type OnUpdate string

func (ou OnUpdate) applyKey(key *Key) {
	key.Reference.OnUpdate = string(ou)
}