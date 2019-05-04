package grimoire

type tableName interface {
	TableName() string
}

func getTableName(record interface{}) string {
	if tn, ok := record.(tableName); ok {
		return tn.TableName()
	}

	return "users" // TODO
}
