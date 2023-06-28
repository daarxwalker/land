package land

import (
	"slices"
	"strings"

	"github.com/iancoleman/strcase"
)

type OrderQuery interface {
	Asc(column string) OrderQuery
	Desc(column string) OrderQuery
	Slice(orders []OrderParam) OrderQuery
}

type orderQueryBuilder struct {
	*queryBuilder
	entity        *entity
	columns       []*columnsQueryBuilder
	singleColumns []*columnQueryBuilder
	orders        []OrderParam
}

const (
	orderAsc  = "ASC"
	orderDesc = "DESC"
)

func createOrderQuery(entity *entity, columns []*columnsQueryBuilder, singleColumns []*columnQueryBuilder) *orderQueryBuilder {
	return &orderQueryBuilder{
		queryBuilder:  createQueryBuilder().setQueryType(Order),
		entity:        entity,
		columns:       columns,
		singleColumns: singleColumns,
		orders:        make([]OrderParam, 0),
	}
}

func (q *orderQueryBuilder) Asc(column string) OrderQuery {
	q.orders = append(q.orders, OrderParam{Key: strcase.ToSnake(column), Direction: orderAsc})
	return q
}

func (q *orderQueryBuilder) Desc(column string) OrderQuery {
	q.orders = append(q.orders, OrderParam{Key: strcase.ToSnake(column), Direction: orderDesc})
	return q
}

func (q *orderQueryBuilder) Slice(orders []OrderParam) OrderQuery {
	for _, o := range orders {
		o.Key = strcase.ToSnake(o.Key)
		o.Direction = strings.ToUpper(o.Direction)
		q.orders = append(q.orders, o)
	}
	return q
}

func (q *orderQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	for _, order := range q.orders {
		if order.Dynamic {
			for _, c := range q.columns {
				if !slices.Contains(c.columns, order.Key) {
					continue
				}
				result = append(result, q.createColumnSql(c.entity, order))
			}
			for _, c := range q.singleColumns {
				if c.name != order.Key {
					continue
				}
				result = append(result, q.createColumnSql(c.entity, order))
			}
		}
		if !order.Dynamic {
			result = append(result, q.createColumnSql(q.entity, order))
		}
	}
	return strings.Join(result, q.getColumnsDivider())
}

func (q *orderQueryBuilder) createColumnSql(entity *entity, order OrderParam) string {
	result := make([]string, 0)
	if len(entity.alias) > 0 {
		result = append(result, q.escape(entity.alias), q.getCoupler())
	}
	result = append(result, q.escape(order.Key))
	result = append(result, " "+order.Direction)
	return strings.Join(result, "")
}
