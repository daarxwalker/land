package land

import (
	"strings"
)

type ColumnQuery interface {
	Entity(entity Entity) ColumnQuery
	Alias(alias string) ColumnQuery
	Use(use bool) ColumnQuery
}

type columnQueryBuilder struct {
	*queryBuilder
	entity *entity
	alias  string
	name   string
	use    bool
}

func createColumnQuery(entity *entity, name string) *columnQueryBuilder {
	return &columnQueryBuilder{
		queryBuilder: createQueryBuilder(),
		entity:       entity,
		use:          true,
		name:         name,
	}
}

func (q *columnQueryBuilder) Entity(entity Entity) ColumnQuery {
	q.entity = entity.getPtr()
	return q
}

func (q *columnQueryBuilder) Alias(alias string) ColumnQuery {
	q.alias = alias
	return q
}

func (q *columnQueryBuilder) Use(use bool) ColumnQuery {
	q.use = use
	return q
}

func (q *columnQueryBuilder) getQueryString() string {
	result := make([]string, 0)
	colSql := make([]string, 0)
	if len(q.entity.alias) > 0 {
		colSql = append(colSql, q.escape(q.entity.alias)+q.getCoupler())
	}
	colSql = append(colSql, q.escape(q.name))
	result = append(result, strings.Join(colSql, ""))
	if len(q.alias) > 0 {
		result = append(result, "AS", q.escape(q.name))
	}
	return strings.Join(result, " ")
}
