package change

type Builder interface {
	Build(changes *Changes)
}

func Build(builders ...Builder) Changes {
	changes := Changes{
		Fields:  make(map[string]int),
		Changes: make([]Change, 0, len(builders)),
	}

	for i := range builders {
		builders[i].Build(&changes)
	}

	return changes
}
