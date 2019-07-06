package change

type ChangeOp int

const (
	SetOp ChangeOp = iota
	IncOp
	DecOp
	FragmentOp
)

type Change struct {
	Type  ChangeOp
	Field string
	Value interface{}
}

func (c Change) Build(changes *Changes) {
	if index, exist := changes.Fields[c.Field]; exist {
		changes.Changes[index] = c
	} else {
		changes.Fields[c.Field] = len(changes.Changes)
		changes.Changes = append(changes.Changes, c)
	}
}

func Set(field string, value interface{}) Change {
	return Change{
		Type:  SetOp,
		Field: field,
		Value: value,
	}
}

func Inc(field string) Change {
	return IncBy(field, 1)
}

func IncBy(field string, n int) Change {
	return Change{
		Type:  IncOp,
		Field: field,
		Value: n,
	}
}

func Dec(field string) Change {
	return DecBy(field, 1)
}

func DecBy(field string, n int) Change {
	return Change{
		Type:  DecOp,
		Field: field,
		Value: n,
	}
}

func Fragment(raw string, args ...interface{}) Change {
	return Change{
		Type:  FragmentOp,
		Field: raw,
		Value: args,
	}
}
