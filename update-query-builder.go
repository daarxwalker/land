package land

import (
	"reflect"
	"strings"

	"github.com/iancoleman/strcase"
)

type UpdateQuery interface {
	SetData(value any) UpdateQuery
	GetSQL() string
	SetTSVectors(values ...any) UpdateQuery
	Return(columns ...string) UpdateQuery
}

type updateQueryBuilder struct {
	*queryBuilder
	entity   *entity
	data     ref
	vectors  string
	returns  []string
	isReturn bool
}

func createUpdateQuery(entity *entity) *updateQueryBuilder {
	return &updateQueryBuilder{
		queryBuilder: createQueryBuilder(),
		entity:       entity,
		returns:      make([]string, 0),
	}
}

func (q *updateQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *updateQueryBuilder) SetData(data any) UpdateQuery {
	q.data.t = reflect.TypeOf(data)
	q.data.v = reflect.ValueOf(data)
	if q.data.v.Kind() == reflect.Ptr {
		q.data.v = q.data.v.Elem()
	}
	return q
}

func (q *updateQueryBuilder) Return(columns ...string) UpdateQuery {
	q.returns = append(q.returns, columns...)
	q.isReturn = true
	return q
}

func (q *updateQueryBuilder) SetTSVectors(values ...any) UpdateQuery {
	q.vectors = createTSVectors(values...)
	return q
}

func (q *updateQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "UPDATE", q.escape(q.entity.name))
	result = append(result, "SET", q.createSetsPart())
	result = append(result, q.createReturnPart()...)
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *updateQueryBuilder) createSetsPart() string {
	result := make([]string, 0)
	for _, c := range q.entity.columns {
		if c.name == Id || !q.data.v.IsValid() {
			continue
		}
		setSql := make([]string, 0)
		setSql = append(setSql, q.escape(c.name), "=")
		field := q.data.v.FieldByName(strcase.ToCamel(c.name))
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

func (q *updateQueryBuilder) getValueOfInvalidField(c *column) string {
	switch c.name {
	case CreatedAt:
		return reflect.ValueOf(c.options.Default).String()
	case UpdatedAt:
		return reflect.ValueOf(c.options.Default).String()
	case Vectors:
		return q.vectors
	default:
		return ""
	}
}
