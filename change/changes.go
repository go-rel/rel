package change

type Changes struct {
	Fields  map[string]int
	Changes []Change
}

func (c Changes) Empty() bool {
	return len(c.Changes) == 0
}

func (c *Changes) Apply(builders ...Builder) {
	for i := range builders {
		builders[i].Build(c)
	}
}
