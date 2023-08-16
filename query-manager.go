package land

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
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

func (m *queryManager) getRowData(row *sql.Rows, columnsTypes []*sql.ColumnType) (reflect.Value, error) {
	result := m.createRowDataModel()
	model := make([]any, len(columnsTypes))
	columns := make([]string, len(columnsTypes))
	for i, ct := range columnsTypes {
		columns[i] = strcase.ToCamel(ct.Name())
		switch strings.ToLower(ct.DatabaseTypeName()) {
		case Varchar, Char, Text:
			model[i] = &sql.NullString{}
		case Int4, Int8:
			model[i] = &sql.NullInt64{}
		case Float4, Float8:
			model[i] = &sql.NullFloat64{}
		case Bool, Boolean:
			model[i] = &sql.NullBool{}
		case Byte, Bytea:
			model[i] = &sql.NullByte{}
		case Timestamp, TimestampWithZone:
			model[i] = &sql.NullTime{}
		}
	}
	if err := row.Scan(model...); err != nil {
		return result, err
	}
	for i, c := range columns {
		switch modelValue := model[i].(type) {
		case *sql.NullFloat64:
			m.setValue(result, c, modelValue.Float64)
		case *sql.NullInt64:
			m.setValue(result, c, int(modelValue.Int64))
		case *sql.NullString:
			m.setValue(result, c, modelValue.String)
		case *sql.NullBool:
			m.setValue(result, c, modelValue.Bool)
		case *sql.NullByte:
			m.setValue(result, c, modelValue.Byte)
		case *sql.NullTime:
			m.setValue(result, c, modelValue.Time)
		}
	}
	return result, nil
}

func (m *queryManager) setValue(model reflect.Value, key string, value any) {
	v := reflect.ValueOf(value)
	if m.isDestMap() || m.isDestSliceOfMaps() {
		model.SetMapIndex(reflect.ValueOf(key), v)
		return
	}
	if model.Kind() != reflect.Struct {
		model.Set(v)
		return
	}
	m.setValueToStruct(model, key, v)
}

func (m *queryManager) setValueToStruct(model reflect.Value, key string, value reflect.Value) {
	f := model.FieldByName(key)
	if !f.IsValid() || !value.IsValid() {
		return
	}
	f.Set(value)
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
		rowData, err := m.getRowData(rows, columnsTypes)
		if err != nil {
			panic(
				Error{
					Error: err,
					Query: m.query,
				},
			)
		}
		m.setRowDataToResult(rowData)
	}
}

func (m *queryManager) createRowDataModel() reflect.Value {
	var result reflect.Value
	if !m.isDestSlice() {
		result = m.destRef.v
		return result
	}
	switch m.destRef.t.Elem().Kind() {
	case reflect.Map:
		result = reflect.MakeMapWithSize(m.destRef.t.Elem(), 0)
	default:
		result = reflect.New(m.destRef.t.Elem()).Elem()
	}
	return result
}

func (m *queryManager) setRowDataToResult(rowData reflect.Value) {
	if m.isDestSlice() {
		m.destRef.v.Set(reflect.Append(m.destRef.v, rowData))
		return
	}
	m.destRef.v.Set(rowData)
}

func (m *queryManager) isDestMap() bool {
	return m.resultType == reflect.Map.String()
}

func (m *queryManager) isDestSlice() bool {
	return m.resultType == reflect.Slice.String()
}

func (m *queryManager) isDestSliceOfMaps() bool {
	return m.isDestSlice() && m.destRef.t.Elem().Kind() == reflect.Map
}

func (m *queryManager) isDestSliceOfStructs() bool {
	return m.isDestSlice() && m.destRef.t.Elem().Kind() == reflect.Struct
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
