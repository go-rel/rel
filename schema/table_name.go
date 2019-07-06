package schema

import (
	"sync"

	"github.com/azer/snakecase"
	"github.com/jinzhu/inflection"
)

var tableNamesCache sync.Map

type tableName interface {
	TableName() string
}

// InferTableName from struct definition, fallback to reflection is not defined.
func InferTableName(record interface{}) string {
	if tn, ok := record.(tableName); ok {
		return tn.TableName()
	}

	rt, _ := reflectInternalType(record)

	// check for cache
	if name, cached := tableNamesCache.Load(rt); cached {
		return name.(string)
	}

	name := inflection.Plural(rt.Name())
	name = snakecase.SnakeCase(name)

	tableNamesCache.Store(rt, name)

	return name
}
