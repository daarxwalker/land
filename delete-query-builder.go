package land

type DeleteQuery interface {
}

type deleteQueryBuilder struct {
	entity *entity
}

func createDeleteQuery(entity *entity) *deleteQueryBuilder {
	return &deleteQueryBuilder{
		entity: entity,
	}
}
