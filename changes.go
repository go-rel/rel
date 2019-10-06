package rel

type Changer interface {
	Build(changes *Changes)
}

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

// TODO: assoc changes
// Use Assoc fields in Changes?
// Table name not stored here, but handled by repo logic.
// TODO: handle deleteion
//	- Answer: Changes should be forward only operation, no delete change is supported (use changeset instead).
// Implement iterator to be used by adapter api?
// Not safe to be used multiple time. some operation my alter changes data.
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

func (c Changes) Empty() bool {
	return len(c.changes) == 0
}

func (c Changes) Count() int {
	return len(c.changes)
}

func (c Changes) AssocCount() int {
	return len(c.assocChanges)
}

func (c Changes) All() []Change {
	return c.changes
}

func (c Changes) Get(field string) (Change, bool) {
	if index, ok := c.fields[field]; ok {
		return c.changes[index], true
	}

	return Change{}, false
}

func (c *Changes) Set(ch Change) {
	if index, exist := c.fields[ch.Field]; exist {
		c.changes[index] = ch
	} else {
		c.fields[ch.Field] = len(c.changes)
		c.changes = append(c.changes, ch)
	}
}

func (c Changes) GetValue(field string) (interface{}, bool) {
	var (
		ch, ok = c.Get(field)
	)

	return ch.Value, ok
}

func (c *Changes) SetValue(field string, value interface{}) {
	c.Set(Set(field, value))
}

func (c Changes) GetAssoc(field string) (AssocChanges, bool) {
	if index, ok := c.assoc[field]; ok {
		return c.assocChanges[index], true
	}

	return AssocChanges{}, false
}

func (c *Changes) SetAssoc(field string, chs ...Changes) {
	if index, exist := c.assoc[field]; exist {
		c.assocChanges[index].Changes = chs
	} else {
		c.appendAssoc(field, AssocChanges{
			Changes: chs,
		})
	}
}

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

type ChangeOp int

const (
	ChangeSetOp ChangeOp = iota
	ChangeIncOp
	ChangeDecOp
	ChangeFragmentOp
)

type Change struct {
	Type  ChangeOp
	Field string
	Value interface{}
}

func (c Change) Build(changes *Changes) {
	changes.Set(c)
}

func Set(field string, value interface{}) Change {
	return Change{
		Type:  ChangeSetOp,
		Field: field,
		Value: value,
	}
}

func Inc(field string) Change {
	return IncBy(field, 1)
}

func IncBy(field string, n int) Change {
	return Change{
		Type:  ChangeIncOp,
		Field: field,
		Value: n,
	}
}

func Dec(field string) Change {
	return DecBy(field, 1)
}

func DecBy(field string, n int) Change {
	return Change{
		Type:  ChangeDecOp,
		Field: field,
		Value: n,
	}
}

func ChangeFragment(raw string, args ...interface{}) Change {
	return Change{
		Type:  ChangeFragmentOp,
		Field: raw,
		Value: args,
	}
}
