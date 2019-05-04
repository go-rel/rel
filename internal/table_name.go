package internal

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

	typ := reflectInternalType(record)
	if name, cached := tableNamesCache.Load(typ); cached {
		return name.(string)
	}

	name := inflection.Plural(typ.Name())
	name = snakecase.SnakeCase(name)
	tableNamesCache.Store(typ, name)

	return name
}
