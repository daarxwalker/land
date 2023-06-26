package land

type DropTableQuery interface {
}

type dropTableQueryBuilder struct {
	entity *entity
}

func createDropTableQuery(entity *entity) *dropTableQueryBuilder {
	return &dropTableQueryBuilder{
		entity: entity,
	}
}
