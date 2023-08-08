package land

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
	
	"github.com/iancoleman/strcase"
)

type queryManager struct {
	entity     *entity
	context    context.Context
	queryType  string
	query      string
	dest       any
	destRef    ref
	resultType string
	errors     []error
	duration   time.Duration
}

func createQueryManager(entity *entity, context context.Context) *queryManager {
	return &queryManager{
		entity:  entity,
		context: context,
	}
}

func (m *queryManager) setQuery(query string) *queryManager {
	m.query = query
	return m
}

func (m *queryManager) setQueryType(queryType string) *queryManager {
	m.queryType = queryType
	return m
}

func (m *queryManager) setDest(dest any) *queryManager {
	t := reflect.TypeOf(dest)
	v := reflect.ValueOf(dest)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	m.dest = dest
	m.destRef = ref{t: t, v: v}
	m.resultType = v.Kind().String()
	return m
}

func (m *queryManager) getResult() {
	defer m.entity.errorHandler.recover()
	m.scan()
	m.log()
}

func (m *queryManager) exec() {
	defer m.entity.errorHandler.recover()
	_, err := m.connection().Exec(m.query)
	m.entity.errorManager.check(err, m.query)
	m.log()
}

func (m *queryManager) scan() {
	start := time.Now()
	rows, err := m.connection().QueryContext(m.context, m.query)
	m.duration = time.Now().Sub(start)
	m.entity.errorManager.check(err, m.query)
	defer func() {
		m.entity.errorManager.check(rows.Close(), m.query)
	}()
	columnsTypes, err := rows.ColumnTypes()
	for rows.Next() {
		row := make([]any, len(columnsTypes))
		for i, _ := range columnsTypes {
			row[i] = new(any)
		}
		err = rows.Scan(row...)
		rowModel := m.createResultDataModel()
		if !rowModel.IsValid() {
			continue
		}
		for i, ct := range columnsTypes {
			col := reflect.ValueOf(row[i])
			if col.Kind() == reflect.Ptr {
				col = col.Elem()
			}
			if col.Kind() == reflect.Interface {
				col = col.Elem()
			}
			m.setResultFieldValue(rowModel, ct, m.standardizeFieldValue(col))
		}
		m.fillResultWithDataModel(rowModel)
	}
}

func (m *queryManager) createResultDataModel() reflect.Value {
	var model reflect.Value
	if m.resultType == reflect.Slice.String() {
		switch m.destRef.t.Elem().Kind() {
		case reflect.Struct:
			model = reflect.New(m.destRef.t.Elem()).Elem()
		case reflect.Map:
			model = reflect.MakeMapWithSize(m.destRef.t.Elem(), 0)
		}
		return model
	}
	model = m.destRef.v
	return model
}

func (m *queryManager) fillResultWithDataModel(rowModel reflect.Value) {
	if m.resultType == reflect.Slice.String() {
		m.destRef.v.Set(reflect.Append(m.destRef.v, rowModel))
		return
	}
	m.destRef.v.Set(rowModel)
}

func (m *queryManager) setResultFieldValue(rowModel reflect.Value, ct *sql.ColumnType, value reflect.Value) {
	
	if m.isResultMap() || m.isResultSliceOfMaps() {
		rowModel.SetMapIndex(reflect.ValueOf(ct.Name()), value)
		return
	}
	if rowModel.Kind() != reflect.Struct {
		rowModel.Set(value)
		return
	}
	field := rowModel.FieldByName(strcase.ToCamel(ct.Name()))
	if !field.IsValid() || !value.IsValid() {
		return
	}
	field.Set(value)
}

func (m *queryManager) standardizeFieldValue(value reflect.Value) reflect.Value {
	kind := value.Kind()
	if slices.Contains(
		[]string{
			reflect.Int8.String(), reflect.Int16.String(), reflect.Int32.String(), reflect.Int64.String(),
		}, kind.String(),
	) {
		return reflect.ValueOf(int(value.Int()))
	}
	return value
}

func (m *queryManager) isResultMap() bool {
	return m.resultType == reflect.Map.String()
}

func (m *queryManager) isResultSliceOfMaps() bool {
	return m.resultType == reflect.Slice.String() && m.destRef.t.Elem().Kind() == reflect.Map
}

func (m *queryManager) createScanRowColumn(columnType string) any {
	typeName := strings.ToLower(columnType)
	if strings.Contains(typeName, Int) {
		return new(int)
	}
	if strings.Contains(typeName, Char) || strings.Contains(typeName, Text) {
		return new(string)
	}
	if strings.Contains(typeName, Bool) {
		return new(bool)
	}
	return new(any)
}

func (m *queryManager) connection() *sql.DB {
	return m.entity.land.db.connection
}

func (m *queryManager) log() {
	if !m.entity.land.config.Log {
		return
	}
	fmt.Println(fmt.Sprintf("%s in %v: %s", strings.ToUpper(m.queryType), m.duration, m.query))
}
