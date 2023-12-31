package land

import (
	"context"
	"fmt"
	"reflect"
	"strings"
)

type CreateTableQuery interface {
	GetSQL() string
	Exec()
	IfNotExists() CreateTableQuery
}

type createTableQueryBuilder struct {
	*queryBuilder
	entity      *entity
	context     context.Context
	ifNotExists bool
}

func createCreateTableQuery(entity *entity) *createTableQueryBuilder {
	return &createTableQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(CreateTable),
		context:      context.Background(),
		entity:       entity,
	}
}

func (q *createTableQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *createTableQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(CreateTable).exec()
}

func (q *createTableQueryBuilder) IfNotExists() CreateTableQuery {
	q.ifNotExists = true
	return q
}

func (q *createTableQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "CREATE TABLE")
	if q.ifNotExists {
		result = append(result, "IF NOT EXISTS")
	}
	result = append(result, q.escape(q.entity.name))
	result = append(result, fmt.Sprintf("(%s)", q.createStructurePart()))
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *createTableQueryBuilder) createStructurePart() string {
	result := make([]string, 0)
	for _, c := range q.entity.columns {
		colSql := make([]string, 0)
		colSql = append(colSql, q.escape(c.name), q.createDataType(c))
		if c.options.PK {
			colSql = append(colSql, "PRIMARY KEY")
		}
		if c.options.NotNull {
			colSql = append(colSql, "NOT NULL")
		}
		if c.options.Unique {
			colSql = append(colSql, "UNIQUE")
		}
		if c.options.Default != nil {
			colSql = append(colSql, "DEFAULT", q.createValue(c, reflect.ValueOf(c.options.Default)))
		}
		if (c.options.Reference.Self && len(c.options.Reference.Column) > 0) || (c.options.Reference.Entity != nil && len(c.options.Reference.Column) > 0) {
			colSql = append(
				colSql, "REFERENCES", q.getReferenceName(c)+fmt.Sprintf("(%s)", q.escape(c.options.Reference.Column)),
			)
		}
		result = append(result, strings.Join(colSql, " "))
	}
	return strings.Join(result, q.getColumnsDivider())
}

func (q *createTableQueryBuilder) getReferenceName(c *column) string {
	if c.options.Reference.Self && c.options.Reference.Entity == nil {
		return q.escape(q.entity.name)
	}
	return q.escape(c.options.Reference.Entity.getPtr().name)
}
