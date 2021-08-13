package reltest

import (
	"reflect"

	"github.com/go-rel/rel"
)

type any struct{}

func (any) String() string {
	return "(Any)"
}

var Any interface{} = any{}

func matchQuery(a rel.Query, b rel.Query) bool {
	return matchTable(a.Table, b.Table) &&
		matchSelectQuery(a.SelectQuery, b.SelectQuery) &&
		matchJoinQuery(a.JoinQuery, b.JoinQuery) &&
		matchFilterQuery(a.WhereQuery, b.WhereQuery) &&
		matchGroupQuery(a.GroupQuery, b.GroupQuery) &&
		matchSortQuery(a.SortQuery, b.SortQuery) &&
		a.OffsetQuery == b.OffsetQuery &&
		a.LimitQuery == b.LimitQuery &&
		a.LockQuery == b.LockQuery &&
		matchSQLQuery(a.SQLQuery, b.SQLQuery) &&
		a.UnscopedQuery == b.UnscopedQuery &&
		a.ReloadQuery == b.ReloadQuery &&
		a.CascadeQuery == b.CascadeQuery &&
		reflect.DeepEqual(a.PreloadQuery, b.PreloadQuery)

}

func matchTable(a string, b string) bool {
	return a == "" || b == "" || a == b
}

func matchSelectQuery(a rel.SelectQuery, b rel.SelectQuery) bool {
	return a.OnlyDistinct == b.OnlyDistinct && reflect.DeepEqual(a.Fields, b.Fields)
}

func matchJoinQuery(a []rel.JoinQuery, b []rel.JoinQuery) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i].Mode != b[i].Mode || a[i].Table != b[i].Table || a[i].From != b[i].From || a[i].To != b[i].To || reflect.DeepEqual(a[i].Arguments, b[i].Arguments) {
			return false
		}
	}

	return true
}

func matchFilterQuery(a rel.FilterQuery, b rel.FilterQuery) bool {
	if a.Type != b.Type || a.Field != b.Field || (a.Value != b.Value && a.Value != Any) || len(a.Inner) != len(b.Inner) {
		return false
	}

	switch v := a.Value.(type) {
	case rel.SubQuery:
		if bSubQuery, _ := b.Value.(rel.SubQuery); v.Prefix != bSubQuery.Prefix || !matchQuery(v.Query, bSubQuery.Query) {
			return false
		}
	case rel.Query:
		if bQuery, ok := b.Value.(rel.Query); !ok || !matchQuery(v, bQuery) {
			return false
		}
	default:
		if a.Value != b.Value && a.Value != Any {
			return false
		}
	}

	for i := range a.Inner {
		if !matchFilterQuery(a.Inner[i], b.Inner[i]) {
			return false
		}
	}

	return true
}

func matchGroupQuery(a rel.GroupQuery, b rel.GroupQuery) bool {
	return reflect.DeepEqual(a.Fields, b.Fields) && matchFilterQuery(a.Filter, b.Filter)
}

func matchSortQuery(a []rel.SortQuery, b []rel.SortQuery) bool {
	return reflect.DeepEqual(a, b)
}

func matchSQLQuery(a rel.SQLQuery, b rel.SQLQuery) bool {
	if a.Statement != b.Statement && len(a.Values) != len(b.Values) {
		return false
	}

	for i := range a.Values {
		if a.Values[i] != b.Values[i] && a.Values[i] != Any {
			return false
		}
	}

	return true
}

func matchMutation(a rel.Mutation, b rel.Mutation) bool {
	if len(a.Mutates) != len(b.Mutates) || len(a.Assoc) != len(b.Assoc) || a.Unscoped != b.Unscoped || a.Reload != b.Reload || a.Cascade != b.Cascade {
		return false
	}

	for i := range a.Mutates {
		if !matchMutate(a.Mutates[i], b.Mutates[i]) {
			return false
		}
	}

	for i := range a.Assoc {
		if !matchAssocMutation(a.Assoc[i], b.Assoc[i]) {
			return false
		}
	}

	return true
}

func matchMutate(a rel.Mutate, b rel.Mutate) bool {
	return a.Type == b.Type && a.Field == b.Field && (a.Value == b.Value || a.Value == Any)
}

func matchMutates(a []rel.Mutate, b []rel.Mutate) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if !matchMutate(a[i], b[i]) {
			return false
		}
	}

	return true
}

func matchAssocMutation(a rel.AssocMutation, b rel.AssocMutation) bool {
	if len(a.Mutations) != len(b.Mutations) || reflect.DeepEqual(a.DeletedIDs, b.DeletedIDs) {
		return false
	}

	for i := range a.Mutations {
		if !matchMutation(a.Mutations[i], b.Mutations[i]) {
			return false
		}
	}

	return true
}
