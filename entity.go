package land

import "database/sql"

type Entity interface {
	SetAlias(alias string) Entity
	SetColumn(name, dataType string, options ...ColOpts) Entity
	SetFulltext(entities ...Entity) Entity
	SetCreatedAt() Entity
	SetUpdatedAt() Entity
	Select() SelectQuery
	Insert() InsertQuery
	Update() UpdateQuery
	Delete() DeleteQuery
	CreateTable() CreateTableQuery
	AlterTable() AlterTableQuery
	DropTable() DropTableQuery

	getPtr() *entity
}

type entity struct {
	entityManager *entityManager
	error         *errorManager
	alias         string
	name          string
	columns       []*column
	fulltext      []*entity
}

func createEntity(entityManager *entityManager, name string) *entity {
	return &entity{
		entityManager: entityManager,
		error:         createErrorManager(),
		name:          name,
		columns:       make([]*column, 0),
		fulltext:      make([]*entity, 0),
	}
}

func (e *entity) Select() SelectQuery {
	return createSelectQuery(e)
}

func (e *entity) Insert() InsertQuery {
	return createInsertQuery(e)
}

func (e *entity) Update() UpdateQuery {
	return createUpdateQuery(e)
}

func (e *entity) Delete() DeleteQuery {
	return createDeleteQuery(e)
}

func (e *entity) CreateTable() CreateTableQuery {
	return createCreateTableQuery(e)
}

func (e *entity) AlterTable() AlterTableQuery {
	return createAlterTableQuery(e)
}

func (e *entity) DropTable() DropTableQuery {
	return createDropTableQuery(e)
}

func (e *entity) SetAlias(alias string) Entity {
	e.alias = alias
	return e
}

func (e *entity) SetColumn(name, dataType string, options ...ColOpts) Entity {
	opts := ColOpts{}
	if len(options) > 0 {
		opts = options[0]
	}
	c := createColumn(name, dataType, opts)
	e.columns = append(e.columns, c)
	return e
}

func (e *entity) SetFulltext(entities ...Entity) Entity {
	e.fulltext = append(e.fulltext, e)
	for _, ent := range entities {
		e.fulltext = append(e.fulltext, ent.getPtr())
	}
	e.columns = append(e.columns, &column{name: Vectors, dataType: TsVector, options: ColOpts{NotNull: true, Default: "", Exclude: true}})
	return e
}

func (e *entity) SetCreatedAt() Entity {
	e.columns = append(e.columns, &column{name: CreatedAt, dataType: e.getDateDataType(), options: ColOpts{NotNull: true, Default: CurrentTimestamp}})
	return e
}

func (e *entity) SetUpdatedAt() Entity {
	e.columns = append(e.columns, &column{name: UpdatedAt, dataType: e.getDateDataType(), options: ColOpts{NotNull: true, Default: CurrentTimestamp}})
	return e
}

func (e *entity) getDateDataType() string {
	if e.entityManager.land.config.Timezone {
		return TimestampWithZone
	}
	return Timestamp
}

func (e *entity) getIdDataType() string {
	switch e.entityManager.land.config.DatabaseType {
	case Postgres:
		return Serial
	default:
		return Serial
	}
}

func (e *entity) setIdColumn() Entity {
	e.columns = append(e.columns, &column{name: Id, dataType: e.getIdDataType(), options: ColOpts{PK: true, NotNull: true, Unique: true}})
	return e
}

func (e *entity) connection() *sql.DB {
	return e.entityManager.land.db.connection
}

func (e *entity) getPtr() *entity {
	return e
}
