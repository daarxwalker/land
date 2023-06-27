package land

import (
	"reflect"
	"strings"
)

type WhereQuery interface {
	And(queries ...WhereQuery) WhereQuery
	Column(column string) WhereQuery
	Contains(value any) WhereQuery
	Equal(value any) WhereQuery
	Like(value any) WhereQuery
	Not() WhereQuery
	Null() WhereQuery
	Or(queries ...WhereQuery) WhereQuery
	Use(use bool) WhereQuery

	fulltext(value string) *whereQueryBuilder
	getPtr() *whereQueryBuilder
}

type whereQueryBuilder struct {
	*queryBuilder
	entity               *entity
	orQueries            []*whereQueryBuilder
	andQueries           []*whereQueryBuilder
	whereType            string
	column               string
	valueRef             ref
	excludeFromZeroLevel bool
	use                  bool
	negation             bool
}

const (
	whereContains = "contains"
	whereEqual    = "equal"
	whereFulltext = "fulltext"
	whereLike     = "like"
	whereNull     = "null"
)

func createWhereQuery(entity *entity) *whereQueryBuilder {
	return &whereQueryBuilder{
		queryBuilder: createQueryBuilder(),
		entity:       entity,
		use:          true,
		orQueries:    make([]*whereQueryBuilder, 0),
		andQueries:   make([]*whereQueryBuilder, 0),
	}
}

func (q *whereQueryBuilder) And(queries ...WhereQuery) WhereQuery {
	for _, query := range queries {
		qq := query.getPtr()
		qq.excludeFromZeroLevel = true
		q.andQueries = append(q.andQueries, qq)
	}
	return q
}

func (q *whereQueryBuilder) Column(column string) WhereQuery {
	q.column = column
	return q
}

func (q *whereQueryBuilder) Contains(value any) WhereQuery {
	q.whereType = whereContains
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}

func (q *whereQueryBuilder) Equal(value any) WhereQuery {
	q.whereType = whereEqual
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}

func (q *whereQueryBuilder) Like(value any) WhereQuery {
	q.whereType = whereLike
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}

func (q *whereQueryBuilder) Not() WhereQuery {
	q.negation = true
	return q
}

func (q *whereQueryBuilder) Null() WhereQuery {
	q.whereType = whereNull
	return q
}

func (q *whereQueryBuilder) Or(queries ...WhereQuery) WhereQuery {
	for _, query := range queries {
		qq := query.getPtr()
		qq.excludeFromZeroLevel = true
		q.orQueries = append(q.orQueries, qq)
	}
	return q
}

func (q *whereQueryBuilder) Use(use bool) WhereQuery {
	q.use = use
	return q
}

func (q *whereQueryBuilder) createQueryString() string {
	shouldBeGrouped := len(q.orQueries) > 0 || len(q.andQueries) > 0
	result := make([]string, 0)
	columnSql := make([]string, 0)
	if len(q.entity.alias) > 0 {
		columnSql = append(columnSql, q.escape(q.entity.alias)+q.getCoupler())
	}
	columnSql = append(columnSql, q.escape(q.column))
	result = append(result, strings.Join(columnSql, ""))
	result = append(result, q.getOperator())
	result = append(result, q.getValue())
	for _, item := range q.andQueries {
		result = append(result, "AND", item.createQueryString())
	}
	for _, item := range q.orQueries {
		result = append(result, "OR", item.createQueryString())
	}
	resultStr := strings.Join(result, " ")
	if shouldBeGrouped {
		return "(" + resultStr + ")"
	}
	return resultStr
}

func (q *whereQueryBuilder) getOperator() string {
	switch q.whereType {
	case whereContains:
		if q.negation {
			return "NOT IN"
		}
		return "IN"
	case whereEqual:
		if q.negation {
			return "!="
		}
		return "="
	case whereFulltext:
		return "@@"
	case whereLike:
		if q.negation {
			return "NOT LIKE"
		}
		return "LIKE"
	case whereNull:
		if q.negation {
			return "IS NOT NULL"
		}
		return "IS NULL"
	default:
		return ""
	}
}

func (q *whereQueryBuilder) getValue() string {
	column := q.getColumn()
	if column == nil {
		return ""
	}
	value := q.createValue(column, q.valueRef.v)
	if q.whereType == whereContains {
		return "(" + value + ")"
	}
	return value
}

func (q *whereQueryBuilder) getColumn() *column {
	for _, c := range q.entity.columns {
		if c.name == q.column {
			return c
		}
	}
	return nil
}

func (q *whereQueryBuilder) getPtr() *whereQueryBuilder {
	return q
}

func (q *whereQueryBuilder) fulltext(value string) *whereQueryBuilder {
	q.whereType = whereFulltext
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}
