package grimoire

type Changer interface {
	Build(changes *Changes)
}

func BuildChanges(changers ...Changer) Changes {
	changes := Changes{
		Fields:       make(map[string]int),
		Changes:      make([]Change, 0, len(changers)),
		Assoc:        make(map[string]int),
		AssocChanges: make([][]Changes, 0),
	}

	for i := range changers {
		changers[i].Build(&changes)
	}

	return changes
}
