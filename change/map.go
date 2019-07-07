package change

type Map map[string]interface{}

func (m Map) Build(changes *Changes) {
	for field, value := range m {
		changes.Set(Set(field, value))
	}
}
