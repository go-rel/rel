package change

type Changes struct {
	Fields  map[string]int
	Changes []Change
}

func (c Changes) Empty() bool {
	return len(c.Changes) == 0
}

func (c Changes) Get(field string) (Change, bool) {
	if index, ok := c.Fields[field]; ok {
		return c.Changes[index], true
	}

	return Change{}, false
}

func (c *Changes) Apply(builders ...Builder) {
	for i := range builders {
		builders[i].Build(c)
	}
}
