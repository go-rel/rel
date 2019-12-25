package rel

// Map can be used as changes for repository insert or update operation.
// This allows inserting or updating only on specified field.
// Insert/Update of has one or belongs to can be done using other Map as a value.
// Insert/Update of has many can be done using slice of Map as a value.
type Map map[string]interface{}

// Build changes from map.
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
