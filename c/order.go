package c

// Order defines order information of query.
type Order struct {
	Field I
	Order int
}

// Asc orders field with ascending order.
func Asc(field I) Order {
	return Order{
		Field: field,
		Order: 1,
	}
}

// Desc orders field with descending order.
func Desc(field I) Order {
	return Order{
		Field: field,
		Order: -1,
	}
}

// Asc returns true if order is ascending.
func (order Order) Asc() bool {
	return order.Order >= 0
}

// Desc returns true if order is descending.
func (order Order) Desc() bool {
	return order.Order < 0
}
