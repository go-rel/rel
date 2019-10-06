package rel

// Changer is interface for a record changer.
type Changer interface {
	Build(changes *Changes)
}

// BuildChanges using given changers.
func BuildChanges(changers ...Changer) Changes {
	changes := Changes{
		fields: make(map[string]int),
		assoc:  make(map[string]int),
	}

	for i := range changers {
		changers[i].Build(&changes)
	}

	return changes
}

// Changes represents value to be inserted or updated to database.
// It's not safe to be used multiple time. some operation my alter changes data.
type Changes struct {
	fields       map[string]int // TODO: not copy friendly
	changes      []Change
	assoc        map[string]int
	assocChanges []AssocChanges
}

type AssocChanges struct {
	Changes []Changes
	// if nil, has many associations will be cleared.
	StaleIDs []interface{}
}

// Empty returns true if changes is empty.
func (c Changes) Empty() bool {
	return len(c.changes) == 0
}

// Count returns count of changes.
func (c Changes) Count() int {
	return len(c.changes)
}

// AssocCount returns count of associations being changes.
func (c Changes) AssocCount() int {
	return len(c.assocChanges)
}

// All return array of change.
func (c Changes) All() []Change {
	return c.changes
}

// Get a change by field name.
func (c Changes) Get(field string) (Change, bool) {
	if index, ok := c.fields[field]; ok {
		return c.changes[index], true
	}

	return Change{}, false
}

// Set a change directly, will existing value replace if it's already exists.
func (c *Changes) Set(ch Change) {
	if index, exist := c.fields[ch.Field]; exist {
		c.changes[index] = ch
	} else {
		c.fields[ch.Field] = len(c.changes)
		c.changes = append(c.changes, ch)
	}
}

// GetValue of change by field name.
func (c Changes) GetValue(field string) (interface{}, bool) {
	var (
		ch, ok = c.Get(field)
	)

	return ch.Value, ok
}

// SetValue using field name and changed value.
func (c *Changes) SetValue(field string, value interface{}) {
	c.Set(Set(field, value))
}

// GetAssoc by field name.
func (c Changes) GetAssoc(field string) (AssocChanges, bool) {
	if index, ok := c.assoc[field]; ok {
		return c.assocChanges[index], true
	}

	return AssocChanges{}, false
}

// SetAssoc by field name.
func (c *Changes) SetAssoc(field string, chs ...Changes) {
	if index, exist := c.assoc[field]; exist {
		c.assocChanges[index].Changes = chs
	} else {
		c.appendAssoc(field, AssocChanges{
			Changes: chs,
		})
	}
}

// SetStaleAssoc sets list of ids marked as stale.
func (c *Changes) SetStaleAssoc(field string, ids []interface{}) {
	if index, exist := c.assoc[field]; exist {
		c.assocChanges[index].StaleIDs = ids
	} else {
		c.appendAssoc(field, AssocChanges{
			StaleIDs: ids,
		})
	}
}

func (c *Changes) appendAssoc(field string, ac AssocChanges) {
	c.assoc[field] = len(c.assocChanges)
	c.assocChanges = append(c.assocChanges, ac)
}

// ChangeOp represents type of change operation.
type ChangeOp int

const (
	ChangeInvalidOp ChangeOp = iota
	ChangeSetOp
	ChangeIncOp
	ChangeDecOp
	ChangeFragmentOp
)

// Change defines information of a change instruction.
type Change struct {
	Type  ChangeOp
	Field string
	Value interface{}
}

// Build changes using this change.
func (c Change) Build(changes *Changes) {
	changes.Set(c)
}

// Set create a change using set operation.
func Set(field string, value interface{}) Change {
	return Change{
		Type:  ChangeSetOp,
		Field: field,
		Value: value,
	}
}

// Inc create a change using increment operation.
func Inc(field string) Change {
	return IncBy(field, 1)
}

// IncBy create a change using increment operation with custom increment value.
func IncBy(field string, n int) Change {
	return Change{
		Type:  ChangeIncOp,
		Field: field,
		Value: n,
	}
}

// Dec create a change using deccrement operation.
func Dec(field string) Change {
	return DecBy(field, 1)
}

// DecBy create a change using decrement operation with custom decrement value.
func DecBy(field string, n int) Change {
	return Change{
		Type:  ChangeDecOp,
		Field: field,
		Value: n,
	}
}

// ChangeFragment create a change operation using random fragment operation.
// Only available for Update.
func ChangeFragment(raw string, args ...interface{}) Change {
	return Change{
		Type:  ChangeFragmentOp,
		Field: raw,
		Value: args,
	}
}
