package land

import (
	"context"
	"reflect"
	"strings"
	
	"github.com/iancoleman/strcase"
)

type InsertQuery interface {
	SetValues(value any) InsertQuery
	CustomId() InsertQuery
	CustomTimestamp() InsertQuery
	GetSQL() string
	Exec()
	GetResult(dest any)
	SetVectors(values ...any) InsertQuery
	Return(columns ...string) InsertQuery
}

type insertQueryBuilder struct {
	*queryBuilder
	entity          *entity
	context         context.Context
	data            ref
	vectors         string
	returns         []string
	isReturn        bool
	customId        bool
	customTimestamp bool
}

func createInsertQuery(entity *entity) *insertQueryBuilder {
	return &insertQueryBuilder{
		queryBuilder: createQueryBuilder().setQueryType(Insert),
		entity:       entity,
		context:      context.Background(),
		returns:      make([]string, 0),
	}
}

func (q *insertQueryBuilder) GetSQL() string {
	return q.createQueryString()
}

func (q *insertQueryBuilder) CustomId() InsertQuery {
	q.customId = true
	return q
}

func (q *insertQueryBuilder) CustomTimestamp() InsertQuery {
	q.customTimestamp = true
	return q
}

func (q *insertQueryBuilder) Exec() {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Insert).exec()
}

func (q *insertQueryBuilder) GetResult(dest any) {
	createQueryManager(q.entity, q.context).setQuery(q.GetSQL()).setQueryType(Insert).setDest(dest).getResult()
}

func (q *insertQueryBuilder) SetValues(data any) InsertQuery {
	q.data.t = reflect.TypeOf(data)
	q.data.v = reflect.ValueOf(data)
	q.data.kind = q.data.v.Kind()
	if q.data.v.Kind() == reflect.Ptr {
		q.data.v = q.data.v.Elem()
	}
	return q
}

func (q *insertQueryBuilder) Return(columns ...string) InsertQuery {
	q.returns = append(q.returns, columns...)
	q.isReturn = true
	return q
}

func (q *insertQueryBuilder) SetVectors(values ...any) InsertQuery {
	q.vectors = createTSVectors(values...)
	return q
}

func (q *insertQueryBuilder) createQueryString() string {
	result := make([]string, 0)
	result = append(result, "INSERT", "INTO", q.escape(q.entity.name))
	result = append(result, "("+q.createColumnsPart()+")")
	result = append(result, "VALUES")
	result = append(result, "("+q.createValuesPart()+")")
	result = append(result, q.createReturnPart()...)
	return strings.Join(result, " ") + q.getQueryDivider()
}

func (q *insertQueryBuilder) createColumnsPart() string {
	result := make([]string, 0)
	for _, c := range q.entity.columns {
		if !q.customId && c.name == Id {
			continue
		}
		result = append(result, q.escape(c.name))
	}
	return strings.Join(result, q.getColumnsDivider())
}

func (q *insertQueryBuilder) createValuesPart() string {
	result := make([]string, 0)
	for _, c := range q.entity.columns {
		if !q.data.v.IsValid() {
			continue
		}
		if !q.customId && c.name == Id {
			continue
		}
		if !q.customTimestamp && c.name == CreatedAt {
			result = append(result, CurrentTimestamp)
			continue
		}
		if !q.customTimestamp && c.name == UpdatedAt {
			result = append(result, CurrentTimestamp)
			continue
		}
		if c.name == Vectors {
			if len(q.vectors) == 0 {
				result = append(result, createTSVectors(""))
			}
			if len(q.vectors) > 0 {
				result = append(result, q.vectors)
			}
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
		if field.IsZero() && !c.options.NotNull {
			result = append(result, "NULL")
			continue
		}
		if field.IsZero() && c.options.NotNull {
			result = append(result, q.createValue(c, q.createDefaultValue(c, reflect.ValueOf(c.options.Default))))
			continue
		}
		result = append(result, q.createValue(c, field))
	}
	return strings.Join(result, q.getColumnsDivider())
}

func (q *insertQueryBuilder) createReturnPart() []string {
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
