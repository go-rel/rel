package grimoire

type Lock string

func (l Lock) Build(query *Query) {
	query.LockClause = l
}

func ForUpdate() Lock {
	return "FOR UPDATE"
}
