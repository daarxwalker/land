package land

import (
	"reflect"
	"strings"
)

type ConditionQuery interface {
	And(queries ...ConditionQuery) ConditionQuery
	Column(column string) ConditionQuery
	Contains(value any) ConditionQuery
	Equal(value any) ConditionQuery
	Like(value any) ConditionQuery
	Not() ConditionQuery
	Null() ConditionQuery
	Or(queries ...ConditionQuery) ConditionQuery
	Use(use bool) ConditionQuery

	fulltext(value string) *conditionQueryBuilder
	getPtr() *conditionQueryBuilder
}

type conditionQueryBuilder struct {
	*queryBuilder
	entity               *entity
	orQueries            []*conditionQueryBuilder
	andQueries           []*conditionQueryBuilder
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

func createConditionQuery(entity *entity) *conditionQueryBuilder {
	return &conditionQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Where),
		entity:       entity,
		use:          true,
		orQueries:    make([]*conditionQueryBuilder, 0),
		andQueries:   make([]*conditionQueryBuilder, 0),
	}
}

func (q *conditionQueryBuilder) And(queries ...ConditionQuery) ConditionQuery {
	for _, query := range queries {
		qq := query.getPtr()
		qq.excludeFromZeroLevel = true
		q.andQueries = append(q.andQueries, qq)
	}
	return q
}

func (q *conditionQueryBuilder) Column(column string) ConditionQuery {
	q.column = column
	return q
}

func (q *conditionQueryBuilder) Contains(value any) ConditionQuery {
	q.whereType = whereContains
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}

func (q *conditionQueryBuilder) Equal(value any) ConditionQuery {
	q.whereType = whereEqual
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}

func (q *conditionQueryBuilder) Like(value any) ConditionQuery {
	q.whereType = whereLike
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}

func (q *conditionQueryBuilder) Not() ConditionQuery {
	q.negation = true
	return q
}

func (q *conditionQueryBuilder) Null() ConditionQuery {
	q.whereType = whereNull
	return q
}

func (q *conditionQueryBuilder) Or(queries ...ConditionQuery) ConditionQuery {
	for _, query := range queries {
		qq := query.getPtr()
		qq.excludeFromZeroLevel = true
		q.orQueries = append(q.orQueries, qq)
	}
	return q
}

func (q *conditionQueryBuilder) Use(use bool) ConditionQuery {
	q.use = use
	return q
}

func (q *conditionQueryBuilder) createQueryString() string {
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

func (q *conditionQueryBuilder) getOperator() string {
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

func (q *conditionQueryBuilder) getValue() string {
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

func (q *conditionQueryBuilder) getColumn() *column {
	for _, c := range q.entity.columns {
		if c.name == q.column {
			return c
		}
	}
	return nil
}

func (q *conditionQueryBuilder) getPtr() *conditionQueryBuilder {
	return q
}

func (q *conditionQueryBuilder) fulltext(value string) *conditionQueryBuilder {
	q.whereType = whereFulltext
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	return q
}
