package grimoire

// type Error struct {
// 	Type string
// 	Field string
// 	Message string
// }

type Changeset struct {
	Collection string
	Entity     interface{}
	Changes    map[string]interface{}
	// errors *Error
}

// cast convert struct changeset and apply parameters as changes
func Cast(entity interface{}, params map[string]interface{}, fields ...string) *Changeset {
	// TODO: extract collection name from entity

	ch := &Changeset{
		Entity:  entity,
		Changes: make(map[string]interface{}),
	}

	for _, fi := range fields {
		if val, exist := params[fi]; exist {
			// TODO: convert fi to snake case
			ch.Changes[fi] = val
		}
	}

	return ch
}

// // change convert a struct into a changeset
// func Change(entity) *Changeset {

// }

// // Decode convert struct changeset and apply parameters as changes
// func Decode(entity interface{}, json []byte, fields string...) *Changeset {

// }

// // cast convert struct changeset and apply parameters as changes
// func CastAssociation(entity interface{}, params Param, fields string...) *Changeset {

// }

// func ValidateAcceptance() *Changeset {

// }

// func ValidateChange() *Changeset {

// }

// func ValidateConfirmation() *Changeset {

// }

// func ValidateExclusion() *Changeset {

// }

// func ValidateFormat() *Changeset {

// }

// func ValidateInclusion() *Changeset {

// }

// func ValidateLength() *Changeset {

// }

// func ValidateNumber() *Changeset {

// }

// func ValidateLength() *Changeset {

// }
