package rel

// Map can be used as changes for repository insert or update operation.
// This allows inserting or updating only on specified field.
// Insert/Update of has one or belongs to can be done using other Map as a value.
// Insert/Update of has many can be done using slice of Map as a value.
type Map map[string]interface{}

// Apply changes.
func (m Map) Apply(doc *Document, changes *Changes) error {
	for field, value := range m {
		switch v := value.(type) {
		case Map:
			// TODO: apply assoc
			assoc, _ := ApplyChanges(nil, v)
			changes.SetAssoc(field, assoc)
		case []Map:
			var (
				chs = make([]Changes, len(v))
			)

			// TODO: apply assoc
			for i := range v {
				assoc, _ := ApplyChanges(nil, v[i])
				chs[i] = assoc
			}
			changes.SetAssoc(field, chs...)
		default:
			changes.SetValue(field, v)
		}
	}

	return nil
}
