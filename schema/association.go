package schema

type Association struct {
	BelongsTo map[string]AssociationField
	HasOne    map[string]AssociationField
	HasMany   map[string]AssociationField
}
