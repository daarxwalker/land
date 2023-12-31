package land

import (
	"context"
	"fmt"
	"strings"
)

type SelectQuery interface {
	With(name string) WithQuery
	Context(context context.Context) SelectQuery
	Column(name ...string) ColumnQuery
	Columns(columns ...string) ColumnsQuery
	Where(entity ...Entity) ConditionQuery
	Having(entity ...Entity) ConditionQuery
	Join(entity ...Entity) JoinQuery
	Fulltext(value string) SelectQuery
	Group(columns ...string) GroupQuery
	Order(orders ...OrderParam) OrderQuery
	Offset(offset int) SelectQuery
	Distinct() SelectQuery
	Single() SelectQuery
	Limit(limit int) SelectQuery
	Exists() bool
	All() SelectQuery
	Param(param Param) SelectQuery
	GetSQL() string
	GetResult(value any)
	Exec()
	
	getPtr() *selectQueryBuilder
}

type selectQueryBuilder struct {
	*queryBuilder
	context       context.Context
	entity        *entity
	columns       []*columnsQueryBuilder
	singleColumns []*columnQueryBuilder
	joins         []*joinQueryBuilder
	wheres        []*conditionQueryBuilder
	havings       []*conditionQueryBuilder
	orders        []*orderQueryBuilder
	groups        []*groupQueryBuilder
	withs         []*withQueryBuilder
	param         Param
	distinct      bool
}

func createSelectQuery(entity *entity) *selectQueryBuilder {
	q := &selectQueryBuilder{
		queryBuilder:  createQueryBuilder().setQueryType(Select),
		context:       context.Background(),
		entity:        entity,
		columns:       make([]*columnsQueryBuilder, 0),
		singleColumns: make([]*columnQueryBuilder, 0),
		joins:         make([]*joinQueryBuilder, 0),
		wheres:        make([]*conditionQueryBuilder, 0),
		havings:       make([]*conditionQueryBuilder, 0),
		orders:        make([]*orderQueryBuilder, 0),
		groups:        make([]*groupQueryBuilder, 0),
		withs:         make([]*withQueryBuilder, 0),
		param: Param{
			Limit: DefaultLimit,
		},
		distinct: false,
	}
	return q
}

func (q *selectQueryBuilder) With(name string) WithQuery {
	w := createWithQuery(name)
	q.withs = append(q.withs, w.getPtr())
	return w
}

func (q *selectQueryBuilder) Context(context context.Context) SelectQuery {
	q.context = context
	return q
}

func (q *selectQueryBuilder) Distinct() SelectQuery {
	q.distinct = true
	return q
}

func (q *selectQueryBuilder) Column(name ...string) ColumnQuery {
	columnName := ""
	if len(name) > 0 {
		columnName = name[0]
	}
	c := createColumnQuery(q.entity, columnName)
	q.singleColumns = append(q.singleColumns, c)
	return c
}

func (q *selectQueryBuilder) Columns(columns ...string) ColumnsQuery {
	c := createColumnsQuery(q.entity, columns...)
	q.columns = append(q.columns, c)
	return c
}

func (q *selectQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Select).exec()
}

func (q *selectQueryBuilder) GetResult(dest any) {
	createQueryManager(
		q.entity, q.context,
	).setQuery(q.GetSQL()).setQueryType(Select).setDest(dest).getResult()
}

func (q *selectQueryBuilder) Exists() bool {
	var result bool
	createQueryManager(q.entity, q.context).
		setQuery(fmt.Sprintf("SELECT EXISTS(%s);", strings.TrimSuffix(q.GetSQL(), q.getQueryDivider()))).
		setQueryType(Select).
		setDest(&result).
		getResult()
	return result
}

func (q *selectQueryBuilder) GetSQL() string {
	return q.createQueryString() + q.getQueryDivider()
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

func (q *selectQueryBuilder) Group(columns ...string) GroupQuery {
	group := createGroupQuery(q.entity, columns...)
	q.groups = append(q.groups, group)
	return group
}

func (q *selectQueryBuilder) Order(orders ...OrderParam) OrderQuery {
	order := createOrderQuery(q.entity, q.columns, q.singleColumns, orders...)
	q.orders = append(q.orders, order)
	return order
}

func (q *selectQueryBuilder) Where(entity ...Entity) ConditionQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	where := createConditionQuery(e)
	q.wheres = append(q.wheres, where)
	return where
}

func (q *selectQueryBuilder) Having(entity ...Entity) ConditionQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	having := createConditionQuery(e)
	q.havings = append(q.havings, having)
	return having
}

func (q *selectQueryBuilder) Fulltext(value string) SelectQuery {
	q.param.Fulltext = value
	return q
}

func (q *selectQueryBuilder) Offset(offset int) SelectQuery {
	q.param.Offset = offset
	return q
}

func (q *selectQueryBuilder) Single() SelectQuery {
	q.param.Limit = 1
	return q
}

func (q *selectQueryBuilder) Limit(limit int) SelectQuery {
	q.param.Limit = limit
	return q
}

func (q *selectQueryBuilder) Param(param Param) SelectQuery {
	for i := range param.Order {
		param.Order[i].Dynamic = true
	}
	q.param = param
	if len(param.Order) > 0 {
		order := createOrderQuery(q.entity, q.columns, q.singleColumns, param.Order...)
		q.orders = append(q.orders, order)
	}
	return q
}

func (q *selectQueryBuilder) All() SelectQuery {
	q.param.All = true
	return q
}

func (q *selectQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	if len(q.withs) > 0 {
		result = append(result, q.createWithsPart())
	}
	result = append(result, "SELECT")
	if q.distinct {
		result = append(result, "DISTINCT")
	}
	result = append(result, strings.Join(q.createColumnsPart(), q.getColumnsDivider()))
	result = append(result, "FROM")
	result = append(result, q.createFromPart()...)
	result = append(result, q.createJoinsPart()...)
	result = append(result, q.createWheresPart()...)
	result = append(result, q.createGroupsPart()...)
	result = append(result, q.createHavingsPart()...)
	result = append(result, q.createOrdersPart()...)
	result = append(result, q.createLimit()...)
	result = append(result, q.createOffset()...)
	return strings.Join(result, " ")
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

func (q *selectQueryBuilder) createWithsPart() string {
	result := make([]string, 0)
	for _, with := range q.withs {
		result = append(result, with.createQueryString())
	}
	if len(result) > 0 {
		return "WITH" + " " + strings.Join(result, ",")
	}
	return ""
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
	q.wheres = append(q.wheres, createConditionQuery(q.entity).Column(Vectors).fulltext(q.param.Fulltext))
}

func (q *selectQueryBuilder) createWheresPart() []string {
	result := make([]string, 0)
	i := 0
	q.createFulltextConditions()
	for _, where := range q.wheres {
		if where.excludeFromZeroLevel || !where.use {
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
		i++
	}
	return result
}

func (q *selectQueryBuilder) createHavingsPart() []string {
	result := make([]string, 0)
	for i, having := range q.havings {
		if having.excludeFromZeroLevel {
			continue
		}
		condition := make([]string, 0)
		if i == 0 {
			condition = append(condition, "HAVING")
		}
		if i > 0 {
			condition = append(condition, "AND")
		}
		condition = append(condition, having.createQueryString())
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
	for _, order := range q.orders {
		orders = append(orders, order.createQueryString())
	}
	if len(orders) > 0 {
		result = append(result, "ORDER BY")
		result = append(result, strings.Join(orders, q.getColumnsDivider()))
	}
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

func (q *selectQueryBuilder) getPtr() *selectQueryBuilder {
	return q
}
