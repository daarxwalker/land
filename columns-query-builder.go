package land

import (
	"strings"
)

type ColumnsQuery interface {
	Entity(entity Entity) ColumnsQuery
	Use(use bool) ColumnsQuery
}

type columnsQueryBuilder struct {
	*queryBuilder
	entity  *entity
	columns []string
	use     bool
}

func createColumnsQuery(entity *entity, columns ...string) *columnsQueryBuilder {
	c := &columnsQueryBuilder{
		queryBuilder: createQueryBuilder(),
		entity:       entity,
		columns:      make([]string, 0),
		use:          true,
	}
	c.columns = append(c.columns, columns...)
	return c
}

func (q *columnsQueryBuilder) Entity(entity Entity) ColumnsQuery {
	q.entity = entity.getPtr()
	return q
}

func (q *columnsQueryBuilder) Use(use bool) ColumnsQuery {
	q.use = use
	return q
}

func (q *columnsQueryBuilder) getQueryStringSlice() []string {
	result := make([]string, 0)
	for _, c := range q.columns {
		colSql := make([]string, 0)
		if len(q.entity.alias) > 0 {
			colSql = append(colSql, q.escape(q.entity.alias), q.getCoupler())
		}
		colSql = append(colSql, q.escape(c))
		result = append(result, strings.Join(colSql, ""))
	}
	return result
}
