package land

type UpdateQuery interface {
}

type updateQueryBuilder struct {
	entity *entity
}

func createUpdateQuery(entity *entity) *updateQueryBuilder {
	return &updateQueryBuilder{
		entity: entity,
	}
}
