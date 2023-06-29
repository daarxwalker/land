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
	Avg() ColumnQuery
	Count() ColumnQuery
	Sum() ColumnQuery
	ArrayAgg() ColumnQuery
	StringAgg(columns ...ColumnQuery) ColumnQuery
	Min() ColumnQuery
	Max() ColumnQuery
	Subquery(query SelectQuery) ColumnQuery
	Separator(separator string) ColumnQuery

	getPtr() *columnQueryBuilder
}

type columnQueryBuilder struct {
	*queryBuilder
	entity        *entity
	alias         string
	name          string
	subquery      string
	subqueryExist bool
	use           bool
	aggregate     string
	separator     string
	columns       []*columnQueryBuilder
}

const (
	aggregateAvg       = "avg"
	aggregateArrayAgg  = "array-agg"
	aggregateCount     = "count"
	aggregateSum       = "sum"
	aggregateMax       = "max"
	aggregateMin       = "min"
	aggregateStringAgg = "string-agg"
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

func (q *columnQueryBuilder) Subquery(query SelectQuery) ColumnQuery {
	q.subquery = query.GetSQL()
	q.subqueryExist = len(q.subquery) > 0
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
	case aggregateSum:
		col = fmt.Sprintf("SUM(%s)", col)
	case aggregateArrayAgg:
		col = fmt.Sprintf("ARRAY_AGG(%s)", col)
	case aggregateStringAgg:
		col = q.createStringAggWrapper(col)
	}
	return col
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
