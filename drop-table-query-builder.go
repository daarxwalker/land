package land

import (
	"strings"
)

type DropTableQuery interface {
	Cascade() DropTableQuery
	GetSQL() string
	IfExists() DropTableQuery
}

type dropTableQueryBuilder struct {
	*queryBuilder
	entity   *entity
	ifExists bool
	cascade  bool
}

func createDropTableQuery(entity *entity) *dropTableQueryBuilder {
	return &dropTableQueryBuilder{
		queryBuilder: createQueryBuilder(),
		entity:       entity,
	}
}

func (q *dropTableQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *dropTableQueryBuilder) IfExists() DropTableQuery {
	q.ifExists = true
	return q
}

func (q *dropTableQueryBuilder) Cascade() DropTableQuery {
	q.cascade = true
	return q
}

func (q *dropTableQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "DROP TABLE")
	if q.ifExists {
		result = append(result, "IF EXISTS")
	}
	result = append(result, q.escape(q.entity.name))
	if q.cascade {
		result = append(result, "CASCADE")
	}
	return strings.Join(result, " ") + q.getQueryDivider()
}
