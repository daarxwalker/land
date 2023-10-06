package land

import (
	"context"
	"strings"
)

type TruncateQuery interface {
	GetSQL() string
	Exec()
	RestartIdentity() TruncateQuery
	Cascade() TruncateQuery
}

type truncateQueryBuilder struct {
	*queryBuilder
	entity          *entity
	context         context.Context
	restartIdentity bool
	cascade         bool
}

func createTruncateQuery(entity *entity) *truncateQueryBuilder {
	return &truncateQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Truncate),
		context:      context.Background(),
		entity:       entity,
	}
}

func (q *truncateQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *truncateQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Truncate).exec()
}

func (q *truncateQueryBuilder) RestartIdentity() TruncateQuery {
	q.restartIdentity = true
	return q
}

func (q *truncateQueryBuilder) Cascade() TruncateQuery {
	q.cascade = true
	return q
}

func (q *truncateQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, Truncate)
	result = append(result, q.escape(q.entity.name))
	if q.restartIdentity {
		result = append(result, "RESTART IDENTITY")
	}
	if q.cascade {
		result = append(result, "CASCADE")
	}
	return strings.Join(result, " ") + q.getQueryDivider()
}
