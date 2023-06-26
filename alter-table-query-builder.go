package land

type AlterTableQuery interface {
}

type alterTableQueryBuilder struct {
	entity *entity
}

func createAlterTableQuery(entity *entity) *alterTableQueryBuilder {
	return &alterTableQueryBuilder{
		entity: entity,
	}
}
