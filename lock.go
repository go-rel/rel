package rel

type Lock string

func (l Lock) Build(query *Query) {
	query.LockQuery = l
}

func ForUpdate() Lock {
	return "FOR UPDATE"
}
