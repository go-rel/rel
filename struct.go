package grimoire

type structChanger struct {
	doc Document
}

func (cs structChanger) Build(changes *Changes) {
	var (
		fields      = cs.doc.Fields()
		values      = cs.doc.Values()
		pField      = cs.doc.PrimaryField()
		belongsTo   = cs.doc.BelongsTo()
		hasOne      = cs.doc.HasOne()
		hasMany     = cs.doc.HasMany()
		assocFields = make(map[string]struct{}, len(belongsTo)+len(hasOne)+len(hasMany))
	)

	for _, field := range belongsTo {
		assocFields[field] = struct{}{}
		cs.buildAssoc(field, changes)
	}

	for _, field := range hasOne {
		assocFields[field] = struct{}{}
		cs.buildAssoc(field, changes)
	}

	for _, field := range hasMany {
		assocFields[field] = struct{}{}
		cs.buildAssocMany(field, changes)
	}

	for field, i := range fields {
		if field == pField {
			continue
		}

		if _, isAssoc := assocFields[field]; isAssoc {
			continue
		}

		var (
			value = values[i]
		)

		if !isZero(value) {
			changes.SetValue(field, value)
		}
	}
}

func (cs structChanger) buildAssoc(field string, changes *Changes) {
	var (
		assoc = cs.doc.Association(field)
	)

	if col, loaded := assoc.Target(); loaded {
		var (
			ch = BuildChanges(structChanger{doc: col.Get(0)})
		)

		changes.SetAssoc(field, ch)
	}
}

func (cs structChanger) buildAssocMany(field string, changes *Changes) {
	var (
		assoc = cs.doc.Association(field)
	)

	if col, loaded := assoc.Target(); loaded {
		var (
			chs = make([]Changes, col.Len())
		)

		for i := range chs {
			chs[i] = BuildChanges(structChanger{doc: col.Get(i)})
		}

		changes.SetAssoc(field, chs...)
	}
}

func Struct(entity interface{}) Changer {
	return structChanger{
		doc: newDocument(entity),
	}
}
