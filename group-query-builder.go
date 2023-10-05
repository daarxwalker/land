package land

import (
	"strings"
)

type GroupQuery interface {
	Entity(entity Entity) GroupQuery
}

type groupQueryBuilder struct {
	*queryBuilder
	entity  *entity
	columns []string
}

func createGroupQuery(entity *entity, columns ...string) *groupQueryBuilder {
	return &groupQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Group),
		entity:       entity,
		columns:      columns,
	}
}

func (q *groupQueryBuilder) Entity(entity Entity) GroupQuery {
	q.entity = entity.getPtr()
	return q
}

func (q *groupQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	for _, col := range q.columns {
		colSql := make([]string, 0)
		if len(q.entity.alias) > 0 {
			colSql = append(colSql, q.escape(q.entity.alias), q.getCoupler())
		}
		colSql = append(colSql, q.escape(col))
		result = append(result, strings.Join(colSql, ""))
	}
	return strings.Join(result, q.getColumnsDivider())
}
