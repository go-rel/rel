package changeset

func Cast(params map[string]interface{}, fields []string, opts ...Option) *Changeset {
	ch := &Changeset{}
	ch.changes = make(map[string]interface{})

	for _, f := range fields {
		val, exist := params[f]
		if exist {
			ch.changes[f] = val
		}
	}

	return ch
}
