package land

import (
	"fmt"
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
	Subquery(subquery SelectQuery) ConditionQuery
	Use(use bool) ConditionQuery
	Webalize() ConditionQuery
	
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
	subquery             string
	valueRef             ref
	excludeFromZeroLevel bool
	use                  bool
	negation             bool
	webalize             bool
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
	q.createValueRef(value)
	return q
}

func (q *conditionQueryBuilder) Equal(value any) ConditionQuery {
	q.whereType = whereEqual
	q.createValueRef(value)
	return q
}

func (q *conditionQueryBuilder) Like(value any) ConditionQuery {
	q.whereType = whereLike
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	q.valueRef.kind = q.valueRef.v.Kind()
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

func (q *conditionQueryBuilder) Subquery(query SelectQuery) ConditionQuery {
	q.subquery = strings.TrimSuffix(query.GetSQL(), q.getQueryDivider())
	return q
}

func (q *conditionQueryBuilder) Use(use bool) ConditionQuery {
	q.use = use
	return q
}

func (q *conditionQueryBuilder) Webalize() ConditionQuery {
	q.webalize = true
	return q
}

func (q *conditionQueryBuilder) createValueRef(value any) {
	q.valueRef.t = reflect.TypeOf(value)
	q.valueRef.v = reflect.ValueOf(value)
	q.valueRef.kind = q.valueRef.v.Kind()
	q.valueRef.safe = q.valueRef.t == reflect.TypeOf(Safe{})
}

func (q *conditionQueryBuilder) createQueryString() string {
	shouldBeGrouped := len(q.orQueries) > 0 || len(q.andQueries) > 0
	subqueryExist := len(q.subquery) > 0
	result := make([]string, 0)
	if subqueryExist {
		result = append(result, fmt.Sprintf("(%s)", q.subquery))
	}
	if !shouldBeGrouped {
		if !subqueryExist {
			columnSql := make([]string, 0)
			if len(q.entity.alias) > 0 {
				columnSql = append(columnSql, q.escape(q.entity.alias)+q.getCoupler())
			}
			columnSql = append(columnSql, q.escape(q.column))
			col := strings.Join(columnSql, "")
			if q.webalize {
				col = webalize(col)
			}
			result = append(result, col)
		}
		result = append(result, q.getOperator())
		value := q.getValue()
		if len(value) > 0 {
			result = append(result, value)
		}
	}
	for i, item := range q.andQueries {
		if (i != 0 && shouldBeGrouped) || !shouldBeGrouped {
			result = append(result, "AND")
		}
		result = append(result, item.createQueryString())
	}
	for i, item := range q.orQueries {
		if (i != 0 && shouldBeGrouped) || !shouldBeGrouped {
			result = append(result, "OR")
		}
		result = append(result, item.createQueryString())
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
	if column == nil && len(q.subquery) == 0 {
		return ""
	}
	if column == nil && len(q.subquery) > 0 {
		return q.createValueWithUnknownColumn(q.valueRef)
	}
	if q.valueRef.safe {
		return fmt.Sprintf("%s", q.valueRef.v.Interface().(Safe).Value)
	}
	value := q.createValue(column, q.valueRef.v)
	if q.webalize {
		value = webalize(value)
	}
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
