package rel

// ForceCascade mutation.
// It forces cascade save of the specified has many associations even if
// the collection is empty or not loaded. This is useful to delete all
// associated records by setting the collection to nil or an empty slice.
//
// Example:
//
//	repo.Update(ctx, &cart, ForceCascade{"items"})
//
// Without this option, a nil or empty collection is treated as "not loaded"
// and the associated records are left untouched.
type ForceCascade []string

// Apply mutation.
func (fc ForceCascade) Apply(doc *Document, mutation *Mutation) {
	mutation.ForceCascade = append(mutation.ForceCascade, fc...)
}

// has returns true if the given association field is forced to cascade.
func (fc ForceCascade) has(field string) bool {
	for i := range fc {
		if fc[i] == field {
			return true
		}
	}

	return false
}
