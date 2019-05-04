package internal


type CustomSchema struct {
	UUID string
}

func (c CustomSchema) TableName() string {
	return "users"
}

func (c CustomSchema) PrimaryKey() (string, interface{}) {
	return "uuid", c.UUID
}
