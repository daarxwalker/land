package land

import (
	"fmt"
	"reflect"
	"strings"
)

type CreateTableQuery interface {
	GetSQL() string
	IfNotExists() CreateTableQuery
}

type createTableQueryBuilder struct {
	*queryBuilder
	entity      *entity
	ifNotExists bool
}

func createCreateTableQuery(entity *entity) *createTableQueryBuilder {
	return &createTableQueryBuilder{
		queryBuilder: createQueryBuilder(),
		entity:       entity,
	}
}

func (q *createTableQueryBuilder) GetSQL() string {
	return q.createQueryString()
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
		if c.options.NotNull {
			colSql = append(colSql, "NOT NULL")
		}
		if c.options.Unique {
			colSql = append(colSql, "UNIQUE")
		}
		if c.options.Default != nil {
			colSql = append(colSql, "DEFAULT", q.createValue(c, reflect.ValueOf(c.options.Default)))
		}
		if c.options.Reference.Entity != nil && len(c.options.Reference.Column) > 0 {
			colSql = append(colSql, "REFERENCES", q.escape(c.options.Reference.Entity.getPtr().name)+fmt.Sprintf("(%s)", q.escape(c.options.Reference.Column)))
		}
		result = append(result, strings.Join(colSql, " "))
	}
	return strings.Join(result, q.getColumnsDivider())
}
