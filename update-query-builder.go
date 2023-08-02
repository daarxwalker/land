package land

import (
	"context"
	"reflect"
	"slices"
	"strings"
	
	"github.com/iancoleman/strcase"
)

type UpdateQuery interface {
	SetColumns(columns ...string) UpdateQuery
	SetValues(value any) UpdateQuery
	GetSQL() string
	GetResult(dest any)
	Exec()
	SetVectors(values ...any) UpdateQuery
	Return(columns ...string) UpdateQuery
	Where(entity ...Entity) ConditionQuery
}

type updateQueryBuilder struct {
	*queryBuilder
	entity   *entity
	context  context.Context
	data     ref
	wheres   []*conditionQueryBuilder
	vectors  string
	returns  []string
	columns  []string
	isReturn bool
}

var (
	updateQueryReservedColumns = []string{UpdatedAt, Vectors}
)

func createUpdateQuery(entity *entity) *updateQueryBuilder {
	return &updateQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Update),
		context:      context.Background(),
		entity:       entity,
		wheres:       make([]*conditionQueryBuilder, 0),
		returns:      make([]string, 0),
	}
}

func (q *updateQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *updateQueryBuilder) GetResult(dest any) {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Update).setDest(dest).getResult()
}

func (q *updateQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Update).exec()
}

func (q *updateQueryBuilder) SetValues(data any) UpdateQuery {
	q.data.t = reflect.TypeOf(data)
	q.data.v = reflect.ValueOf(data)
	q.data.kind = q.data.v.Kind()
	if q.data.v.Kind() == reflect.Ptr {
		q.data.v = q.data.v.Elem()
	}
	return q
}

func (q *updateQueryBuilder) SetColumns(columns ...string) UpdateQuery {
	q.columns = columns
	return q
}

func (q *updateQueryBuilder) Return(columns ...string) UpdateQuery {
	q.returns = append(q.returns, columns...)
	q.isReturn = true
	return q
}

func (q *updateQueryBuilder) SetVectors(values ...any) UpdateQuery {
	q.vectors = createTSVectors(values...)
	return q
}

func (q *updateQueryBuilder) Where(entity ...Entity) ConditionQuery {
	e := q.entity
	if len(entity) > 0 {
		e = entity[0].getPtr()
	}
	where := createConditionQuery(e)
	q.wheres = append(q.wheres, where)
	return where
}

func (q *updateQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "UPDATE", q.escape(q.entity.name), "AS", q.escape(q.entity.alias))
	result = append(result, "SET", q.createSetsPart())
	result = append(result, q.createWheresPart()...)
	result = append(result, q.createReturnPart()...)
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *updateQueryBuilder) createSetsPart() string {
	result := make([]string, 0)
	for _, c := range q.entity.columns {
		if len(q.columns) > 0 && !slices.Contains(q.columns, c.name) {
			continue
		}
		if c.name == Id || c.name == CreatedAt || !q.data.v.IsValid() || (c.name == Vectors && len(q.vectors) == 0) {
			continue
		}
		setSql := make([]string, 0)
		setSql = append(setSql, q.escape(c.name), "=")
		if slices.Contains(updateQueryReservedColumns, c.name) {
			switch c.name {
			case UpdatedAt:
				setSql = append(setSql, CurrentTimestamp)
			case Vectors:
				setSql = append(setSql, q.vectors)
			}
			result = append(result, strings.Join(setSql, " "))
			continue
		}
		var field reflect.Value
		switch q.data.kind {
		case reflect.Map:
			field = q.getMapValue(q.data.v, c.name)
		case reflect.Struct:
			field = q.data.v.FieldByName(strcase.ToCamel(c.name))
		}
		if !field.IsValid() {
			continue
		}
		if field.IsZero() {
			setSql = append(setSql, q.createZeroValue(c, field))
		}
		if !field.IsZero() {
			setSql = append(setSql, q.createValue(c, field))
		}
		result = append(result, strings.Join(setSql, " "))
	}
	return strings.Join(result, q.getColumnsDivider())
}

func (q *updateQueryBuilder) createZeroValue(c *column, field reflect.Value) string {
	if !c.options.NotNull {
		return "NULL"
	}
	switch c.dataType {
	case Boolean:
		return "false"
	default:
		return q.createValue(c, field)
	}
}

func (q *updateQueryBuilder) createReturnPart() []string {
	result := make([]string, 0)
	if !q.isReturn {
		return result
	}
	result = append(result, "RETURNING")
	if len(q.returns) == 0 {
		result = append(result, "*")
		return result
	}
	returnCols := make([]string, len(q.returns))
	for i, r := range q.returns {
		returnCols[i] = q.escape(r)
	}
	result = append(result, strings.Join(returnCols, q.getColumnsDivider()))
	return result
}

func (q *updateQueryBuilder) createWheresPart() []string {
	result := make([]string, 0)
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
