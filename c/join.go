package c

// Join defines join information in query.
type Join struct {
	Mode       string
	Collection string
	Condition  Condition
}
