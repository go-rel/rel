package rel

import (
	"reflect"
	"sync"

	"github.com/serenize/snaker"
)

var associationMetaCache sync.Map

type associationKey struct {
	rt reflect.Type
	// string repr of index, because []int is not hashable
	index string
}

// AssociationType defines the type of association in database.
type AssociationType uint8

const (
	// BelongsTo association.
	BelongsTo = iota
	// HasOne association.
	HasOne
	// HasMany association.
	HasMany
)

type cachedAssociationMeta struct {
	typ            AssociationType
	targetIndex    []int
	referenceField string
	referenceIndex []int
	foreignField   string
	foreignIndex   []int
	through        string
	autoload       bool
	autosave       bool
}

type AssociationMeta struct {
	rt reflect.Type
	cachedAssociationMeta
}

// Type of association.
func (am AssociationMeta) Type() AssociationType {
	return am.typ
}

// ReferenceField of the association.
func (am AssociationMeta) ReferenceField() string {
	return am.referenceField
}

// ForeignField of the association.
func (am AssociationMeta) ForeignField() string {
	return am.foreignField
}

// Through return intermediary association.
func (am AssociationMeta) Through() string {
	return am.through
}

// Autoload assoc setting when parent is loaded.
func (am AssociationMeta) Autoload() bool {
	return am.autoload
}

// Autosave setting when parent is created/updated/deleted.
func (am AssociationMeta) Autosave() bool {
	return am.autosave
}

// Document returns association target document meta.
func (am AssociationMeta) DocumentMeta() DocumentMeta {
	var (
		rt = am.rt.FieldByIndex(am.targetIndex).Type
	)

	if rt.Kind() == reflect.Ptr {
		rt = rt.Elem()
	}

	if rt.Kind() == reflect.Slice {
		rt = rt.Elem()
	}

	return getDocumentMeta(rt, false)
}

func getAssociationMeta(rt reflect.Type, index []int) AssociationMeta {
	var (
		key = associationKey{
			rt:    rt,
			index: encodeIndices(index),
		}
	)

	if val, cached := associationMetaCache.Load(key); cached {
		return AssociationMeta{
			rt:                    rt,
			cachedAssociationMeta: val.(cachedAssociationMeta),
		}
	}

	var (
		sf        = rt.FieldByIndex(index)
		ft        = sf.Type
		ref       = sf.Tag.Get("ref")
		fk        = sf.Tag.Get("fk")
		fName, _  = fieldName(sf)
		assocMeta = cachedAssociationMeta{
			targetIndex: index,
			through:     sf.Tag.Get("through"),
			autoload:    sf.Tag.Get("auto") == "true" || sf.Tag.Get("autoload") == "true",
			autosave:    sf.Tag.Get("auto") == "true" || sf.Tag.Get("autosave") == "true",
		}
	)

	if assocMeta.autosave && assocMeta.through != "" {
		panic("rel: autosave is not supported for has one/has many through association")
	}

	for ft.Kind() == reflect.Ptr || ft.Kind() == reflect.Slice {
		ft = ft.Elem()
	}

	var (
		refDocMeta = getDocumentMeta(rt, true)
		fkDocMeta  = getDocumentMeta(ft, true)
	)

	// Try to guess ref and fk if not defined.
	if ref == "" || fk == "" {
		// TODO: replace "id" with inferred primary field
		if assocMeta.through != "" {
			ref = "id"
			fk = "id"
		} else if _, isBelongsTo := refDocMeta.index[fName+"_id"]; isBelongsTo {
			ref = fName + "_id"
			fk = "id"
		} else {
			ref = "id"
			fk = snaker.CamelToSnake(rt.Name()) + "_id"
		}
	}

	if id, exist := refDocMeta.index[ref]; !exist {
		panic("rel: references (" + ref + ") field not found ")
	} else {
		assocMeta.referenceIndex = id
		assocMeta.referenceField = ref
	}

	if id, exist := fkDocMeta.index[fk]; !exist {
		panic("rel: foreign_key (" + fk + ") field not found")
	} else {
		assocMeta.foreignIndex = id
		assocMeta.foreignField = fk
	}

	// guess assoc type
	if sf.Type.Kind() == reflect.Slice || (sf.Type.Kind() == reflect.Ptr && sf.Type.Elem().Kind() == reflect.Slice) {
		assocMeta.typ = HasMany
	} else {
		if len(assocMeta.referenceField) > len(assocMeta.foreignField) {
			assocMeta.typ = BelongsTo
		} else {
			assocMeta.typ = HasOne
		}
	}

	associationMetaCache.Store(key, assocMeta)

	return AssociationMeta{
		rt:                    rt,
		cachedAssociationMeta: assocMeta,
	}
}
