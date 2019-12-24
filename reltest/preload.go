package reltest

import (
	"reflect"
	"strings"

	"github.com/Fs02/rel"
	"github.com/stretchr/testify/mock"
)

// Preload asserts and simulate preload function for test.
type Preload struct {
	*Expect
}

type slice interface {
	ReflectValue() reflect.Value
	Reset()
	Get(index int) *rel.Document
	Len() int
}

// Result sets the result of Preload query.
func (p *Preload) Result(records interface{}) {
	p.Run(func(args mock.Arguments) {
		var (
			target = p.asSlice(args[0])
			result = p.asSlice(records)
			path   = strings.Split(args[1].(string), ".")
		)

		p.preload(target, result, path)
	})
}

// For match expect calls for given record.
func (p *Preload) For(record interface{}) *Preload {
	p.Arguments[0] = record
	return p
}

// ForType match expect calls for given type.
// Type must include package name, example: `model.User`.
func (p *Preload) ForType(typ string) *Preload {
	return p.For(mock.AnythingOfType("*" + strings.TrimPrefix(typ, "*")))
}

func (p Preload) asSlice(v interface{}) slice {
	var (
		sl slice
		rt = reflect.TypeOf(v)
	)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		sl = rel.NewCollection(v)
	} else {
		sl = rel.NewDocument(v)
	}

	return sl
}

func (p Preload) preload(target slice, result slice, path []string) {
	type frame struct {
		index int
		doc   *rel.Document
	}

	var (
		mapResult map[interface{}]reflect.Value
		stack     = make([]frame, target.Len())
	)

	// init stack
	for i := 0; i < len(stack); i++ {
		stack[i] = frame{index: 0, doc: target.Get(i)}
	}

	for len(stack) > 0 {
		var (
			n       = len(stack) - 1
			top     = stack[n]
			assocs  = top.doc.Association(path[top.index])
			hasMany = assocs.Type() == rel.HasMany
		)

		stack = stack[:n]

		if top.index == len(path)-1 {
			var (
				curr   slice
				rValue = assocs.ReferenceValue()
				fField = assocs.ForeignField()
			)

			if rValue == nil {
				continue
			}

			if hasMany {
				curr, _ = assocs.Collection()
			} else {
				curr, _ = assocs.Document()
			}

			curr.Reset()

			if mapResult == nil {
				mapResult = p.mapResult(result, fField, hasMany)
			}

			if rv, ok := mapResult[rValue]; ok {
				curr.ReflectValue().Set(rv)
			}
		} else {
			if assocs.Type() == rel.HasMany {
				var (
					col, loaded = assocs.Collection()
				)

				if !loaded {
					continue
				}

				stack = append(stack, make([]frame, col.Len())...)
				for i := 0; i < col.Len(); i++ {
					stack[n+i] = frame{
						index: top.index + 1,
						doc:   col.Get(i),
					}
				}
			} else {
				if doc, loaded := assocs.Document(); loaded {
					stack = append(stack, frame{
						index: top.index + 1,
						doc:   doc,
					})
				}
			}
		}
	}
}

func (p Preload) mapResult(result slice, fField string, hasMany bool) map[interface{}]reflect.Value {
	var (
		mapResult = make(map[interface{}]reflect.Value)
	)

	for i := 0; i < result.Len(); i++ {
		var (
			doc       = result.Get(i)
			rv        = doc.ReflectValue()
			fValue, _ = doc.Value(fField)
		)

		if hasMany {
			if _, ok := mapResult[fValue]; !ok {
				mapResult[fValue] = reflect.MakeSlice(reflect.SliceOf(rv.Type()), 0, 0)
			}

			mapResult[fValue] = reflect.Append(mapResult[fValue], rv)
		} else {
			mapResult[fValue] = rv
		}
	}

	return mapResult
}

// ExpectPreload to be called with given field and queries.
func ExpectPreload(r *Repository, field string, queriers []rel.Querier) *Preload {
	return &Preload{
		Expect: newExpect(r, "Preload",
			[]interface{}{mock.Anything, field, queriers},
			[]interface{}{nil},
		),
	}
}
