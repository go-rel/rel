package schema

// Schema definition.
type Schema struct {
	Version int
	Tables  []Table
}

// Table with name and its definition.
func (s Schema) Table(name string, fn func(t *Table), options ...Option) {
	table := Table{Name: name}
	fn(&table)
	s.Tables = append(s.Tables, table)
}

// New Schema.
func New(version int) Schema {
	return Schema{Version: version}
}
