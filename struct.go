package grimoire

type documentChanger struct {
	doc Document
}

func (dc documentChanger) Build(changes *Changes) {
	var (
		fields      = dc.doc.Fields()
		values      = dc.doc.Values()
		pField      = dc.doc.PrimaryField()
		belongsTo   = dc.doc.BelongsTo()
		hasOne      = dc.doc.HasOne()
		hasMany     = dc.doc.HasMany()
		assocFields = make(map[string]struct{}, len(belongsTo)+len(hasOne)+len(hasMany))
	)

	for _, field := range belongsTo {
		assocFields[field] = struct{}{}
		dc.buildAssoc(field, changes)
	}

	for _, field := range hasOne {
		assocFields[field] = struct{}{}
		dc.buildAssoc(field, changes)
	}

	for _, field := range hasMany {
		assocFields[field] = struct{}{}
		dc.buildAssocMany(field, changes)
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

func (dc documentChanger) buildAssoc(field string, changes *Changes) {
	var (
		assoc = dc.doc.Association(field)
	)

	if col, loaded := assoc.Target(); loaded {
		var (
			ch = BuildChanges(documentChanger{doc: col.Get(0)})
		)

		changes.SetAssoc(field, ch)
	}
}

func (dc documentChanger) buildAssocMany(field string, changes *Changes) {
	var (
		assoc = dc.doc.Association(field)
	)

	if col, loaded := assoc.Target(); loaded {
		var (
			chs = make([]Changes, col.Len())
		)

		for i := range chs {
			chs[i] = BuildChanges(changeDoc(col.Get(i)))
		}

		changes.SetAssoc(field, chs...)
	}
}

func changeDoc(doc Document) Changer {
	return documentChanger{
		doc: doc,
	}
}

func Struct(record interface{}) Changer {
	return changeDoc(newDocument(record))
}
