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
	if ocm.Keys == nil {
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

// OnConflictFragment insertion when conflict happens.
func OnConflictFragment(sql string, args ...interface{}) OnConflict {
	return OnConflict{Fragment: sql, FragmentArgs: args}
}

// OnConflictKeyFragment insertion when conflict happens on specific keys.
//
// Specifying key is not supported by all database and may be ignored.
func OnConflictKeyFragment(key string, sql string, args ...interface{}) OnConflict {
	return OnConflictKeysFragment([]string{key}, sql, args...)
}

// OnConflictKeysFragment insertion when conflict happens on specific keys.
//
// Specifying key is not supported by all database and may be ignored.
func OnConflictKeysFragment(keys []string, sql string, args ...interface{}) OnConflict {
	return OnConflict{Keys: keys, Fragment: sql, FragmentArgs: args}
}
