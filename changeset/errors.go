package changeset

type Error struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message"`
}

func (e Error) Error() string {
	return e.Message
}

type Errors []Error

func (es Errors) Error() string {
	var messages string

	for i, e := range es {
		messages += e.Error()

		if i < len(es)-1 {
			messages += ", "
		}
	}

	return messages
}
