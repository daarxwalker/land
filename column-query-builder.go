package land

import (
	"fmt"
	"strings"
)

type ColumnQuery interface {
	Empty() ColumnQuery
	Entity(entity Entity) ColumnQuery
	Alias(alias string) ColumnQuery
	Use(use bool) ColumnQuery
	Not() ColumnQuery
	Avg() ColumnQuery
	Count() ColumnQuery
	Length() ColumnQuery
	Sum() ColumnQuery
	ArrayAgg() ColumnQuery
	StringAgg(columns ...ColumnQuery) ColumnQuery
	Min() ColumnQuery
	Max() ColumnQuery
	Subquery(query SelectQuery) ColumnQuery
	Separator(separator string) ColumnQuery
	Webalize() ColumnQuery
	Greater(value any) ColumnQuery
	GreaterEqual(value any) ColumnQuery
	Less(value any) ColumnQuery
	LessEqual(value any) ColumnQuery
	Equal(value any) ColumnQuery
	
	getPtr() *columnQueryBuilder
}

type columnQueryBuilder struct {
	*queryBuilder
	entity            *entity
	alias             string
	name              string
	subquery          string
	subqueryExist     bool
	webalize          bool
	use               bool
	negation          bool
	aggregate         string
	compareExpression string
	compareValue      any
	separator         string
	columns           []*columnQueryBuilder
}

const (
	aggregateAvg       = "avg"
	aggregateArrayAgg  = "array-agg"
	aggregateCount     = "count"
	aggregateLength    = "length"
	aggregateSum       = "sum"
	aggregateMax       = "max"
	aggregateMin       = "min"
	aggregateStringAgg = "string-agg"
)

const (
	compareGreater      = "greater"
	compareGreaterEqual = "greater-equal"
	compareLess         = "less"
	compareLessEqual    = "less-equal"
	compareEqual        = "equal"
)

const (
	operatorDoublePipe string = "||"
)

func createColumnQuery(entity *entity, name string) *columnQueryBuilder {
	return &columnQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Column),
		entity:       entity,
		use:          true,
		name:         name,
	}
}

func (q *columnQueryBuilder) Empty() ColumnQuery {
	q.entity = nil
	return q
}

func (q *columnQueryBuilder) Entity(entity Entity) ColumnQuery {
	q.entity = entity.getPtr()
	return q
}

func (q *columnQueryBuilder) Alias(alias string) ColumnQuery {
	q.alias = alias
	return q
}

func (q *columnQueryBuilder) Use(use bool) ColumnQuery {
	q.use = use
	return q
}

func (q *columnQueryBuilder) Not() ColumnQuery {
	q.negation = true
	return q
}

func (q *columnQueryBuilder) Subquery(query SelectQuery) ColumnQuery {
	q.subquery = strings.TrimSuffix(query.GetSQL(), q.getQueryDivider())
	q.subqueryExist = len(q.subquery) > 0
	return q
}

func (q *columnQueryBuilder) Greater(value any) ColumnQuery {
	q.compareExpression = compareGreater
	q.compareValue = value
	return q
}

func (q *columnQueryBuilder) GreaterEqual(value any) ColumnQuery {
	q.compareExpression = compareGreaterEqual
	q.compareValue = value
	return q
}

func (q *columnQueryBuilder) Less(value any) ColumnQuery {
	q.compareExpression = compareLess
	q.compareValue = value
	return q
}
func (q *columnQueryBuilder) LessEqual(value any) ColumnQuery {
	q.compareExpression = compareLessEqual
	q.compareValue = value
	return q
}
func (q *columnQueryBuilder) Equal(value any) ColumnQuery {
	q.compareExpression = compareEqual
	q.compareValue = value
	return q
}

func (q *columnQueryBuilder) Avg() ColumnQuery {
	q.aggregate = aggregateAvg
	return q
}

func (q *columnQueryBuilder) ArrayAgg() ColumnQuery {
	q.aggregate = aggregateArrayAgg
	return q
}

func (q *columnQueryBuilder) StringAgg(columns ...ColumnQuery) ColumnQuery {
	q.aggregate = aggregateStringAgg
	for _, c := range columns {
		q.columns = append(q.columns, c.getPtr())
	}
	return q
}

func (q *columnQueryBuilder) Count() ColumnQuery {
	q.aggregate = aggregateCount
	return q
}

func (q *columnQueryBuilder) Length() ColumnQuery {
	q.aggregate = aggregateLength
	return q
}

func (q *columnQueryBuilder) Sum() ColumnQuery {
	q.aggregate = aggregateSum
	return q
}

func (q *columnQueryBuilder) Min() ColumnQuery {
	q.aggregate = aggregateMin
	return q
}

func (q *columnQueryBuilder) Max() ColumnQuery {
	q.aggregate = aggregateMax
	return q
}

func (q *columnQueryBuilder) Separator(separator string) ColumnQuery {
	q.separator = separator
	return q
}

func (q *columnQueryBuilder) Webalize() ColumnQuery {
	q.webalize = true
	q.alias = q.name
	return q
}

func (q *columnQueryBuilder) getQueryString() string {
	if q.shouldUseEmpty() {
		return q.createEmptyString()
	}
	result := make([]string, 0)
	colSql := make([]string, 0)
	if len(q.entity.alias) > 0 {
		colSql = append(colSql, q.escape(q.entity.alias)+q.getCoupler())
	}
	if q.subqueryExist {
		colSql = append(colSql, fmt.Sprintf("(%s)", q.subquery))
	}
	if !q.subqueryExist {
		colSql = append(colSql, q.escape(q.name))
	}
	col := strings.Join(colSql, "")
	if len(q.aggregate) > 0 {
		col = q.createAggregateWrapper(col)
	}
	if len(q.compareExpression) > 0 && q.compareValue != nil {
		col = q.createCompare(col)
	}
	if q.webalize {
		col = webalize(col)
	}
	result = append(result, col)
	if len(q.alias) > 0 {
		result = append(result, "AS", q.escape(q.alias))
	}
	return strings.Join(result, " ")
}

func (q *columnQueryBuilder) getPtr() *columnQueryBuilder {
	return q
}

func (q *columnQueryBuilder) shouldUseEmpty() bool {
	return q.entity == nil && len(q.separator) > 0
}

func (q *columnQueryBuilder) createEmptyString() string {
	return fmt.Sprintf("'%s'", q.separator)
}

func (q *columnQueryBuilder) createAggregateWrapper(col string) string {
	switch q.aggregate {
	case aggregateAvg:
		col = fmt.Sprintf("AVG(%s)", col)
	case aggregateMin:
		col = fmt.Sprintf("MIN(%s)", col)
	case aggregateMax:
		col = fmt.Sprintf("MAX(%s)", col)
	case aggregateCount:
		col = fmt.Sprintf("COUNT(%s)", col)
	case aggregateLength:
		col = fmt.Sprintf("LENGTH(%s)", col)
	case aggregateSum:
		col = fmt.Sprintf("SUM(%s)", col)
	case aggregateArrayAgg:
		col = fmt.Sprintf("ARRAY_AGG(%s)", col)
	case aggregateStringAgg:
		col = q.createStringAggWrapper(col)
	}
	return col
}

func (q *columnQueryBuilder) createCompare(column string) string {
	var value string
	switch q.compareValue.(type) {
	case string:
		value = fmt.Sprintf("'%v'", value)
	default:
		value = fmt.Sprintf("%v", value)
	}
	switch q.compareExpression {
	case compareGreater:
		return fmt.Sprintf("%s > %s", column, value)
	case compareGreaterEqual:
		return fmt.Sprintf("%s >= %s", column, value)
	case compareLess:
		return fmt.Sprintf("%s < %s", column, value)
	case compareLessEqual:
		return fmt.Sprintf("%s <= %s", column, value)
	case compareEqual:
		if q.negation {
			return fmt.Sprintf("%s != %s", column, value)
		}
		return fmt.Sprintf("%s = %s", column, value)
	default:
		return ""
	}
}

func (q *columnQueryBuilder) createStringAggWrapper(col string) string {
	stringAggCols := make([]string, 0)
	stringAggCols = append(stringAggCols, col)
	for _, c := range q.columns {
		stringAggCols = append(stringAggCols, c.getQueryString())
	}
	return fmt.Sprintf(
		"STRING_AGG(%s, '%s')",
		strings.Join(
			stringAggCols,
			fmt.Sprintf(" %s ", operatorDoublePipe),
		),
		q.separator,
	)
}
