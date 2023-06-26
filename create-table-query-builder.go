package land

type CreateTableQuery interface {
}

type createTableQueryBuilder struct {
	entity *entity
}

func createCreateTableQuery(entity *entity) *createTableQueryBuilder {
	return &createTableQueryBuilder{
		entity: entity,
	}
}
