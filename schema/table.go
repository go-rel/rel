package schema

// Table definition.
type Table struct {
	Name    string
	Columns []Column
	Indices []Index
}

// Column defines a column with name and type.
func (t *Table) Column(name string, typ ColumnType, options ...Option) {
	column := Column{Name: name, Type: typ}
	t.Columns = append(t.Columns, column)
}

// Index defines an index for columns.
func (t *Table) Index(columns []string, options ...Option) {
	index := Index{Columns: columns}
	t.Indices = append(t.Indices, index)
}

// Boolean defines a column with name and Boolean type.
func (t *Table) Boolean(name string, options ...Option) {
	t.Column(name, Boolean, options...)
}

// Integer defines a column with name and Integer type.
func (t *Table) Integer(name string, options ...Option) {
	t.Column(name, Integer, options...)
}

// BigInt defines a column with name and BigInt type.
func (t *Table) BigInt(name string, options ...Option) {
	t.Column(name, BigInt, options...)
}

// Float defines a column with name and Float type.
func (t *Table) Float(name string, options ...Option) {
	t.Column(name, Float, options...)
}

// Decimal defines a column with name and Decimal type.
func (t *Table) Decimal(name string, options ...Option) {
	t.Column(name, Decimal, options...)
}

// String defines a column with name and String type.
func (t *Table) String(name string, options ...Option) {
	t.Column(name, String, options...)
}

// Text defines a column with name and Text type.
func (t *Table) Text(name string, options ...Option) {
	t.Column(name, Text, options...)
}

// Date defines a column with name and Date type.
func (t *Table) Date(name string, options ...Option) {
	t.Column(name, Date, options...)
}

// DateTime defines a column with name and DateTime type.
func (t *Table) DateTime(name string, options ...Option) {
	t.Column(name, DateTime, options...)
}

// Time defines a column with name and Time type.
func (t *Table) Time(name string, options ...Option) {
	t.Column(name, Time, options...)
}

// Timestamp defines a column with name and Timestamp type.
func (t *Table) Timestamp(name string, options ...Option) {
	t.Column(name, Timestamp, options...)
}

func (t Table) migrate() {}
