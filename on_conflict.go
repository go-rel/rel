package rel

// OnConflict mutation.
type OnConflict struct {
	Keys         []string
	Ignore       bool
	Replace      bool
	Fragment     string
	FragmentArgs []interface{}
}

// Apply mutation.
func (ocm OnConflict) Apply(doc *Document, mutation *Mutation) {
	if ocm.Keys == nil && ocm.Fragment == "" {
		ocm.Keys = doc.PrimaryFields()
	}

	mutation.OnConflict = ocm
}

// OnConflictIgnore insertion when conflict happens.
func OnConflictIgnore() OnConflict {
	return OnConflict{Ignore: true}
}

// OnConflictKeyIgnore insertion when conflict happens on specific keys.
//
// Specifying key is not supported by all database and may be ignored.
func OnConflictKeyIgnore(key string) OnConflict {
	return OnConflictKeysIgnore([]string{key})
}

// OnConflictKeysIgnore insertion when conflict happens on specific keys.
//
// Specifying key is not supported by all database and may be ignored.
func OnConflictKeysIgnore(keys []string) OnConflict {
	return OnConflict{Keys: keys, Ignore: true}
}

// OnConflictReplace insertion when conflict happens.
func OnConflictReplace() OnConflict {
	return OnConflict{Replace: true}
}

// OnConflictKeyReplace insertion when conflict happens on specific keys.
//
// Specifying key is not supported by all database and may be ignored.
func OnConflictKeyReplace(key string) OnConflict {
	return OnConflictKeysReplace([]string{key})
}

// OnConflictKeysReplace insertion when conflict happens on specific keys.
//
// Specifying key is not supported by all database and may be ignored.
func OnConflictKeysReplace(keys []string) OnConflict {
	return OnConflict{Keys: keys, Replace: true}
}

// OnConflictFragment allows to write custom sql for on conflict.
//
// This will add custom sql after ON CONFLICT, example: ON CONFLICT [FRAGMENT]
func OnConflictFragment(sql string, args ...interface{}) OnConflict {
	return OnConflict{Fragment: sql, FragmentArgs: args}
}
