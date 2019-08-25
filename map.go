package grimoire

type Map map[string]interface{}

func (m Map) Build(changes *Changes) {
	for field, value := range m {
		switch v := value.(type) {
		case Map:
			changes.SetAssoc(field, BuildChanges(v))
		case []Map:
			var (
				chs = make([]Changes, len(v))
			)

			for i := range v {
				chs[i] = BuildChanges(v[i])
			}
			changes.SetAssoc(field, chs...)
		default:
			changes.SetValue(field, v)
		}
	}
}
