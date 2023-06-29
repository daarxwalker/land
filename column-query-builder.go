package land

import (
	"fmt"
	"strings"
)

type ColumnQuery interface {
	Entity(entity Entity) ColumnQuery
	Alias(alias string) ColumnQuery
	Use(use bool) ColumnQuery
	Count() ColumnQuery
}

type columnQueryBuilder struct {
	*queryBuilder
	entity *entity
	alias  string
	name   string
	use    bool
	count  bool
}

func createColumnQuery(entity *entity, name string) *columnQueryBuilder {
	return &columnQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Column),
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

func (q *columnQueryBuilder) Count() ColumnQuery {
	q.count = true
	return q
}

func (q *columnQueryBuilder) getQueryString() string {
	result := make([]string, 0)
	colSql := make([]string, 0)
	if len(q.entity.alias) > 0 {
		colSql = append(colSql, q.escape(q.entity.alias)+q.getCoupler())
	}
	colSql = append(colSql, q.escape(q.name))
	col := strings.Join(colSql, "")
	if q.count {
		col = fmt.Sprintf("COUNT(%s)", col)
	}
	result = append(result, col)
	if len(q.alias) > 0 {
		result = append(result, "AS", q.escape(q.alias))
	}
	return strings.Join(result, " ")
}
