package land

import (
	"strings"
)

type JoinQuery interface {
	Left() JoinQuery
	Right() JoinQuery
	Inner() JoinQuery
	Entity(entity Entity) JoinQuery
	Column(column string) JoinQuery
	On(entity Entity, column ...string)
}

type joinQueryBuilder struct {
	*queryBuilder
	joinType   string
	entity     *entity
	column     string
	joinEntity *entity
	joinColumn string
}

const (
	joinLeft  = "LEFT"
	joinRight = "RIGHT"
	joinInner = "INNER"
)

func createJoinQuery(entity *entity) *joinQueryBuilder {
	return &joinQueryBuilder{
		queryBuilder: createQueryBuilder(),
		joinType:     joinLeft,
		entity:       entity,
		column:       Id,
	}
}

func (q *joinQueryBuilder) Left() JoinQuery {
	q.joinType = joinLeft
	return q
}

func (q *joinQueryBuilder) Right() JoinQuery {
	q.joinType = joinRight
	return q
}

func (q *joinQueryBuilder) Inner() JoinQuery {
	q.joinType = joinInner
	return q
}

func (q *joinQueryBuilder) Entity(entity Entity) JoinQuery {
	q.entity = entity.getPtr()
	return q
}

func (q *joinQueryBuilder) Column(column string) JoinQuery {
	q.column = column
	return q
}

func (q *joinQueryBuilder) On(entity Entity, column ...string) {
	q.joinEntity = entity.getPtr()
	if len(column) > 0 {
		q.joinColumn = column[0]
	}
	if len(column) == 0 {
		q.joinColumn = Id
	}
}

func (q *joinQueryBuilder) createQueryString() []string {
	result := make([]string, 0)
	result = append(result, q.joinType, "JOIN", q.joinEntity.name)
	if len(q.joinEntity.alias) > 0 {
		result = append(result, "AS", q.joinEntity.alias)
	}
	result = append(result, "ON")
	first := make([]string, 0)
	if len(q.entity.alias) > 0 {
		first = append(first, q.escape(q.entity.alias), q.getCoupler())
	}
	first = append(first, q.escape(q.column))
	result = append(result, strings.Join(first, ""))
	result = append(result, "=")
	second := make([]string, 0)
	if len(q.joinEntity.alias) > 0 {
		second = append(second, q.escape(q.joinEntity.alias), q.getCoupler())
	}
	second = append(second, q.escape(q.joinColumn))
	result = append(result, strings.Join(second, ""))
	return result
}
