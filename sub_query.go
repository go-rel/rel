package rel

// SubQuery warps a query into: Prefix (Query)
type SubQuery struct {
	Prefix string
	Query  Query
}

// All warp a query into ALL(sub-query)
//
// Some database may not support this keyword, please consult to your database documentation.
func All(sub Query) SubQuery {
	return SubQuery{
		Prefix: "ALL",
		Query:  sub,
	}
}

// Any warp a query into ANY(sub-query)
//
// Some database may not support this keyword, please consult to your database documentation.
func Any(sub Query) SubQuery {
	return SubQuery{
		Prefix: "ANY",
		Query:  sub,
	}
}
