package land

import (
	"database/sql"
	"fmt"

	"github.com/iancoleman/strcase"
)

type Entity interface {
	ErrorManager
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
	Truncate() TruncateQuery
	Column(name string) Safe

	getPtr() *entity
}

type entity struct {
	*errorManager
	errorHandler *errorHandler
	land         *land
	alias        string
	name         string
	columns      []*column
	fulltext     []*entity
}

func createEntity(land *land, name string) *entity {
	e := &entity{
		errorManager: createErrorManager(),
		errorHandler: createErrorHandler(land),
		land:         land,
		name:         name,
		columns:      make([]*column, 0),
		fulltext:     make([]*entity, 0),
	}
	e.createIdColumn()
	return e
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

func (e *entity) Truncate() TruncateQuery {
	return createTruncateQuery(e)
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
	e.columns = append(
		e.columns, &column{name: Vectors, dataType: TsVector, options: ColOpts{NotNull: true, Default: "", Exclude: true}},
	)
	return e
}

func (e *entity) SetCreatedAt() Entity {
	e.columns = append(
		e.columns,
		&column{name: CreatedAt, dataType: e.getDateDataType(), options: ColOpts{NotNull: true, Default: CurrentTimestamp}},
	)
	return e
}

func (e *entity) SetUpdatedAt() Entity {
	e.columns = append(
		e.columns,
		&column{name: UpdatedAt, dataType: e.getDateDataType(), options: ColOpts{NotNull: true, Default: CurrentTimestamp}},
	)
	return e
}

func (e *entity) Column(name string) Safe {
	return Safe{Value: fmt.Sprintf(`"%s"."%s"`, e.alias, strcase.ToSnake(name))}
}

func (e *entity) getDateDataType() string {
	if e.land.config.Timezone {
		return TimestampWithZone
	}
	return Timestamp
}

func (e *entity) getIdDataType() string {
	if e.land.db == nil {
		return Serial
	}
	switch e.land.db.connector.dbtype {
	case Postgres:
		return Serial
	default:
		return Serial
	}
}

func (e *entity) createIdColumn() Entity {
	e.columns = append(
		e.columns, &column{name: Id, dataType: e.getIdDataType(), options: ColOpts{PK: true, NotNull: true, Unique: true}},
	)
	return e
}

func (e *entity) connection() *sql.DB {
	return e.land.db.connection
}

func (e *entity) getPtr() *entity {
	return e
}
