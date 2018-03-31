package changeset

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/Fs02/grimoire/errors"
)

var CastAssocErrorMessage = "{field} is invalid"

type changefunc func(interface{}, map[string]interface{}) *Changeset

func CastAssoc(ch *Changeset, field string, fn changefunc, opts ...Option) {
	options := Options{
		Message: CastAssocErrorMessage,
	}
	options.Apply(opts)

	par, pexist := ch.params[field]
	typ, texist := ch.types[field]
	valid := true

	if pexist && texist {
		if typ.Kind() == reflect.Struct {
			valid = castOne(ch, field, typ, par, fn)
		} else if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Struct {
			valid = castMany(ch, field, typ, par, fn)
		}
	}

	if !valid {
		msg := strings.Replace(options.Message, "{field}", field, 1)
		AddError(ch, field, msg)
	}
}

func castOne(ch *Changeset, field string, typ reflect.Type, par interface{}, fn changefunc) bool {
	mpar, ok := par.(map[string]interface{})
	if !ok {
		return false
	}

	var innerch *Changeset

	if val, exist := ch.values[field]; exist && val != nil {
		innerch = fn(val, mpar)
	} else {
		innerch = fn(reflect.Zero(typ).Interface(), mpar)
	}

	ch.changes[field] = innerch

	// add errors to main errors
	for _, err := range innerch.errors {
		e := err.(errors.Error)
		AddError(ch, field+"."+e.Field, e.Message)
	}

	return true
}

func castMany(ch *Changeset, field string, typ reflect.Type, par interface{}, fn changefunc) bool {
	spar, ok := par.([]map[string]interface{})
	if !ok {
		return false
	}

	chs := make([]*Changeset, 0, len(spar))
	entity := reflect.Zero(typ.Elem()).Interface()

	for i, par := range spar {
		innerch := fn(entity, par)
		chs = append(chs, innerch)

		// add errors to main errors
		for _, err := range innerch.errors {
			e := err.(errors.Error)
			AddError(ch, field+"["+strconv.Itoa(i)+"]."+e.Field, e.Message)
		}
	}

	ch.changes[field] = chs
	return true
}
