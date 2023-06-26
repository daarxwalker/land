package land

import (
	"context"
	"fmt"
	"strings"
)

type SelectQuery interface {
	Column(name string) ColumnQuery
	Columns(columns ...string) ColumnsQuery
	Where(entity ...Entity) WhereQuery
	Join(entity ...Entity) JoinQuery
	Fulltext(value string) SelectQuery
	Group(entity ...Entity) GroupQuery
	Order(entity ...Entity) OrderQuery
	Offset(offset int) SelectQuery
	Limit(limit int) SelectQuery
	All() SelectQuery
	Param(param Param) SelectQuery
	GetSQL() string
	Exec()
}

type selectQueryBuilder struct {
	*queryBuilder
	context       context.Context
	entity        *entity
	columns       []*columnsQueryBuilder
	singleColumns []*columnQueryBuilder
	joins         []*joinQueryBuilder
	wheres        []*whereQueryBuilder
	orders        []*orderQueryBuilder
	groups        []*groupQueryBuilder
	param         Param
}

func createSelectQuery(entity *entity) *selectQueryBuilder {
	return &selectQueryBuilder{
		queryBuilder:  createQueryBuilder(),
		context:       context.Background(),
		entity:        entity,
		columns:       make([]*columnsQueryBuilder, 0),
		singleColumns: make([]*columnQueryBuilder, 0),
		joins:         make([]*joinQueryBuilder, 0),
		wheres:        make([]*whereQueryBuilder, 0),
		orders:        make([]*orderQueryBuilder, 0),
		groups:        make([]*groupQueryBuilder, 0),
		param: Param{
			Limit: DefaultLimit,
		},
	}
}

func (q *selectQueryBuilder) Column(name string) ColumnQuery {
	c := createColumnQuery(q.entity, name)
	q.singleColumns = append(q.singleColumns, c)
	return c
}

func (q *selectQueryBuilder) Columns(columns ...string) ColumnsQuery {
	c := createColumnsQuery(q.entity, columns...)
	q.columns = append(q.columns, c)
	return c
}

func (q *selectQueryBuilder) Exec() {
	_, err := q.entity.connection().Exec(q.createQueryString())
	if err != nil {
		panic(
			Error{error: err},
		)
	}
}

func (q *selectQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *selectQueryBuilder) Join(entity ...Entity) JoinQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	join := createJoinQuery(e)
	q.joins = append(q.joins, join)
	return join
}

func (q *selectQueryBuilder) Group(entity ...Entity) GroupQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	group := createGroupQuery(e)
	q.groups = append(q.groups, group)
	return group
}

func (q *selectQueryBuilder) Order(entity ...Entity) OrderQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	order := createOrderQuery(e, q.columns, q.singleColumns)
	q.orders = append(q.orders, order)
	return order
}

func (q *selectQueryBuilder) Where(entity ...Entity) WhereQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	where := createWhereQuery(e)
	q.wheres = append(q.wheres, where)
	return where
}

func (q *selectQueryBuilder) Fulltext(value string) SelectQuery {
	q.param.Fulltext = value
	return q
}

func (q *selectQueryBuilder) Offset(offset int) SelectQuery {
	q.param.Offset = offset
	return q
}

func (q *selectQueryBuilder) Limit(limit int) SelectQuery {
	q.param.Limit = limit
	return q
}

func (q *selectQueryBuilder) Param(param Param) SelectQuery {
	for i := range param.Order {
		param.Order[i].dynamic = true
	}
	q.param = param
	return q
}

func (q *selectQueryBuilder) All() SelectQuery {
	q.param.All = true
	return q
}

func (q *selectQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "SELECT")
	result = append(result, strings.Join(q.createColumnsPart(), q.getColumnsDivider()))
	result = append(result, "FROM")
	result = append(result, q.createFromPart()...)
	result = append(result, q.createJoinsPart()...)
	result = append(result, q.createWheresPart()...)
	result = append(result, q.createGroupsPart()...)
	result = append(result, q.createOrdersPart()...)
	result = append(result, q.createLimit()...)
	result = append(result, q.createOffset()...)
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *selectQueryBuilder) createColumnsPart() []string {
	result := make([]string, 0)
	if len(q.columns) == 0 && len(q.singleColumns) == 0 {
		result = append(result, "*")
		return result
	}
	for _, query := range q.columns {
		if !query.use {
			continue
		}
		result = append(result, query.getQueryStringSlice()...)
	}
	for _, column := range q.singleColumns {
		if !column.use {
			continue
		}
		result = append(result, column.getQueryString())
	}
	return result
}

func (q *selectQueryBuilder) createFromPart() []string {
	result := make([]string, 0)
	result = append(result, q.escape(q.entity.name))
	if len(q.entity.alias) > 0 {
		result = append(result, "AS")
		result = append(result, q.escape(q.entity.alias))
	}
	return result
}

func (q *selectQueryBuilder) createJoinsPart() []string {
	result := make([]string, 0)
	for _, join := range q.joins {
		result = append(result, strings.Join(join.createQueryString(), " "))
	}
	return result
}

func (q *selectQueryBuilder) createFulltextConditions() {
	if len(q.param.Fulltext) == 0 {
		return
	}
	q.wheres = append(q.wheres, createWhereQuery(q.entity).Column(Vectors).fulltext(q.param.Fulltext))
	for _, join := range q.joins {
		q.wheres = append(q.wheres, createWhereQuery(join.joinEntity).Column(Vectors).fulltext(q.param.Fulltext))
	}
}

func (q *selectQueryBuilder) createWheresPart() []string {
	result := make([]string, 0)
	q.createFulltextConditions()
	for i, where := range q.wheres {
		if where.excludeFromZeroLevel {
			continue
		}
		condition := make([]string, 0)
		if i == 0 {
			condition = append(condition, "WHERE")
		}
		if i > 0 {
			condition = append(condition, "AND")
		}
		condition = append(condition, where.createQueryString())
		result = append(result, strings.Join(condition, " "))
	}
	return result
}

func (q *selectQueryBuilder) createGroupsPart() []string {
	result := make([]string, 0)
	groups := make([]string, 0)
	if len(q.groups) == 0 {
		return result
	}
	result = append(result, "GROUP BY")
	for _, group := range q.groups {
		groups = append(groups, group.createQueryString())
	}
	result = append(result, strings.Join(groups, q.getColumnsDivider()))
	return result
}

func (q *selectQueryBuilder) createOrdersPart() []string {
	result := make([]string, 0)
	orders := make([]string, 0)
	if len(q.orders) == 0 {
		return result
	}
	result = append(result, "ORDER BY")
	for _, order := range q.orders {
		orders = append(orders, order.createQueryString())
	}
	result = append(result, strings.Join(orders, q.getColumnsDivider()))
	return result
}

func (q *selectQueryBuilder) createLimit() []string {
	result := make([]string, 0)
	if q.param.All {
		return result
	}
	if q.param.Limit > 0 {
		result = append(result, fmt.Sprintf("LIMIT %d", q.param.Limit))
	}
	return result
}

func (q *selectQueryBuilder) createOffset() []string {
	result := make([]string, 0)
	if q.param.Offset > 0 {
		result = append(result, fmt.Sprintf("OFFSET %d", q.param.Offset))
	}
	return result
}
