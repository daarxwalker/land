package land

import (
	"context"
	"strings"
)

type WithQuery interface {
	Query(query SelectQuery) WithQuery
	Name(name string) WithQuery
	Alias(alias string) WithQuery

	getPtr() *withQueryBuilder
}

type withQueryBuilder struct {
	*queryBuilder
	context context.Context
	query   *selectQueryBuilder
	name    string
	alias   string
}

func createWithQuery(name string) *withQueryBuilder {
	return &withQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Update),
		context:      context.Background(),
		name:         name,
	}
}

func (q *withQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, q.name, "AS")
	result = append(result, "("+strings.TrimSuffix(q.query.createQueryString(), q.getQueryDivider())+")")
	return strings.Join(result, " ")
}

func (q *withQueryBuilder) Alias(alias string) WithQuery {
	q.alias = alias
	return q
}

func (q *withQueryBuilder) Name(name string) WithQuery {
	q.name = name
	return q
}

func (q *withQueryBuilder) Query(query SelectQuery) WithQuery {
	q.query = query.getPtr()
	return q
}

func (q *withQueryBuilder) getPtr() *withQueryBuilder {
	return q
}
