package grimoire

type Structset struct {
	doc *Document
}

func (s Structset) Build(changes *Changes) {
	var (
		fields      = s.doc.Fields()
		values      = s.doc.Values()
		pField      = s.doc.PrimaryField()
		belongsTo   = s.doc.BelongsTo()
		hasOne      = s.doc.HasOne()
		hasMany     = s.doc.HasMany()
		assocFields = make(map[string]struct{}, len(belongsTo)+len(hasOne)+len(hasMany))
	)

	for _, field := range belongsTo {
		assocFields[field] = struct{}{}
		s.buildAssoc(field, changes)
	}

	for _, field := range hasOne {
		assocFields[field] = struct{}{}
		s.buildAssoc(field, changes)
	}

	for _, field := range hasMany {
		assocFields[field] = struct{}{}
		s.buildAssocMany(field, changes)
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

func (s Structset) buildAssoc(field string, changes *Changes) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			doc, _ = assoc.Document()
			ch     = BuildChanges(Structset{doc: doc})
		)

		changes.SetAssoc(field, ch)
	}
}

func (s Structset) buildAssocMany(field string, changes *Changes) {
	var (
		assoc = s.doc.Association(field)
	)

	if !assoc.IsZero() {
		var (
			col, _ = assoc.Collection()
			chs    = make([]Changes, col.Len())
		)

		for i := range chs {
			chs[i] = BuildChanges(newStructset(col.Get(i)))
		}

		changes.SetAssoc(field, chs...)
	}
}

func newStructset(doc *Document) Structset {
	return Structset{
		doc: doc,
	}
}

func NewStructset(record interface{}) Structset {
	return newStructset(newDocument(record))
}
