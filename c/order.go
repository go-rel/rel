package c

type Order struct {
	Field string
	Order int
}

func Asc(field string) Order {
	return Order{
		Field: field,
		Order: 1,
	}
}

func Desc(field string) Order {
	return Order{
		Field: field,
		Order: -1,
	}
}

func (order Order) Asc() bool {
	return order.Order >= 0
}

func (order Order) Desc() bool {
	return order.Order < 0
}
