package land

import (
	"context"
	"strings"
)

type DeleteQuery interface {
	GetSQL() string
	Exec()
	GetResult(dest any)
	Return(columns ...string) DeleteQuery
	Where(entity ...Entity) ConditionQuery
}

type deleteQueryBuilder struct {
	*queryBuilder
	entity   *entity
	context  context.Context
	wheres   []*conditionQueryBuilder
	returns  []string
	isReturn bool
}

func createDeleteQuery(entity *entity) *deleteQueryBuilder {
	return &deleteQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Delete),
		entity:       entity,
		context:      context.Background(),
		wheres:       make([]*conditionQueryBuilder, 0),
		returns:      make([]string, 0),
	}
}

func (q *deleteQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *deleteQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Delete).exec()
}

func (q *deleteQueryBuilder) GetResult(dest any) {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Select).setDest(dest).getResult()
}

func (q *deleteQueryBuilder) Return(columns ...string) DeleteQuery {
	q.returns = append(q.returns, columns...)
	q.isReturn = true
	return q
}

func (q *deleteQueryBuilder) Where(entity ...Entity) ConditionQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	where := createConditionQuery(e)
	q.wheres = append(q.wheres, where)
	return where
}

func (q *deleteQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "DELETE FROM", q.escape(q.entity.name), "AS", q.escape(q.entity.alias))
	result = append(result, q.createWheresPart()...)
	result = append(result, q.createReturnPart()...)
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *deleteQueryBuilder) createColumnsPart() string {
	result := make([]string, 0)
	for _, c := range q.entity.columns {
		if c.name == Id {
			continue
		}
		result = append(result, q.escape(c.name))
	}
	return strings.Join(result, q.getColumnsDivider())
}

func (q *deleteQueryBuilder) createWheresPart() []string {
	result := make([]string, 0)
	for i, where := range q.wheres {
		if where.excludeFromZeroLevel {
			continue
		}
		condition := make([]string, 0)
		if i == 0 {
			condition = append(condition, "WHERE")
		}
		if i > 0 {
			condition = append(condition, "AND")
		}
		condition = append(condition, where.createQueryString())
		result = append(result, strings.Join(condition, " "))
	}
	return result
}

func (q *deleteQueryBuilder) createReturnPart() []string {
	result := make([]string, 0)
	if !q.isReturn {
		return result
	}
	result = append(result, "RETURNING")
	if len(q.returns) == 0 {
		result = append(result, "*")
		return result
	}
	returnCols := make([]string, len(q.returns))
	for i, r := range q.returns {
		returnCols[i] = q.escape(r)
	}
	result = append(result, strings.Join(returnCols, q.getColumnsDivider()))
	return result
}
