package land

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"time"
)

type queryBuilder struct {
	queryType string
}

func createQueryBuilder() *queryBuilder {
	return &queryBuilder{}
}

func (q *queryBuilder) setQueryType(queryType string) *queryBuilder {
	q.queryType = queryType
	return q
}

func (q *queryBuilder) escape(value string) string {
	return fmt.Sprintf(`"%s"`, value)
}

func (q *queryBuilder) getColumnsDivider() string {
	return ","
}

func (q *queryBuilder) getCoupler() string {
	return "."
}

func (q *queryBuilder) getQueryDivider() string {
	return ";"
}

func (q *queryBuilder) getTimestampFormat() string {
	return "YYYY-MM-DD HH:MI:SS"
}

func (q *queryBuilder) createDataType(c *column) string {
	if c.options.Limit > 0 && slices.Contains([]string{Varchar, Char}, c.dataType) {
		return strings.ToUpper(c.dataType) + fmt.Sprintf("(%d)", c.options.Limit)
	}
	return strings.ToUpper(c.dataType)
}

func (q *queryBuilder) getMapValue(mapValue reflect.Value, key string) reflect.Value {
	return mapValue.MapIndex(reflect.ValueOf(key))
}

func (q *queryBuilder) createValue(column *column, value reflect.Value) string {
	if column == nil {
		return ""
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	if value.Kind() == reflect.Slice {
		return q.createSliceValue(value)
	}
	return q.createAnyValue(column, value)
}

func (q *queryBuilder) createSliceValue(value reflect.Value) string {
	sliceItems := make([]string, value.Len())
	sliceType := value.Type().Elem().Kind()
	for i := 0; i < value.Len(); i++ {
		switch sliceType {
		case reflect.String:
			sliceItems[i] = fmt.Sprintf("'%s'", value.Index(i).String())
		case reflect.Int:
			sliceItems[i] = fmt.Sprintf("%d", value.Index(i).Int())
		}
	}
	return strings.Join(sliceItems, ",")
}

func (q *queryBuilder) createAnyValue(column *column, value reflect.Value) string {
	if value.Kind() == reflect.Interface {
		value = reflect.ValueOf(value.Interface())
	}
	if !q.validateValueKind(column.dataType, value) {
		return ""
	}
	return q.getValueByColumnDataType(column.dataType, value)
}

func (q *queryBuilder) validateValueKind(dataType string, value reflect.Value) bool {
	kind := value.Kind()
	switch dataType {
	case TsVector:
		return kind == reflect.String
	case Varchar:
		return kind == reflect.String
	case Text:
		return kind == reflect.String
	case Char:
		return kind == reflect.String
	case Serial:
		return kind == reflect.Int
	case Int:
		return kind == reflect.Int
	case BigInt:
		return kind == reflect.Int
	case Float:
		return kind == reflect.Float32 || kind == reflect.Float64 || kind == reflect.Int
	case Bool:
		return kind == reflect.Bool
	case Boolean:
		return kind == reflect.Bool
	case Timestamp:
		return kind == reflect.String || kind == reflect.Struct
	case TimestampWithZone:
		return kind == reflect.String || kind == reflect.Struct
	default:
		return false
	}
}

func (q *queryBuilder) getValueByColumnDataType(dataType string, value reflect.Value) string {
	kind := value.Kind()
	switch dataType {
	case TsVector:
		if q.queryType == Where {
			return fmt.Sprintf(`%s`, createTSQuery(value.String()))
		}
		return fmt.Sprintf(`%s`, createTSVectors(value.String()))
	case Varchar:
		return fmt.Sprintf(`'%s'`, value.String())
	case Text:
		return fmt.Sprintf(`'%s'`, value.String())
	case Char:
		return fmt.Sprintf(`'%s'`, value.String())
	case Serial:
		return fmt.Sprintf(`%d`, value.Int())
	case Int:
		return fmt.Sprintf(`%d`, value.Int())
	case BigInt:
		return fmt.Sprintf(`%d`, value.Int())
	case Float:
		if value.Kind() == reflect.Int {
			return fmt.Sprintf(`%d`, value.Int())
		}
		return fmt.Sprintf(`%f`, value.Float())
	case Bool:
		return fmt.Sprintf(`%t`, value.Bool())
	case Boolean:
		return fmt.Sprintf(`%t`, value.Bool())
	case Timestamp:
		if kind == reflect.String {
			return fmt.Sprintf("%s", value.String())
		}
		return fmt.Sprintf("'%s'", value.Interface().(time.Time).Format(time.DateTime))
	case TimestampWithZone:
		if kind == reflect.String {
			return fmt.Sprintf("%s", value.String())
		}
		return fmt.Sprintf("'%s'", value.Interface().(time.Time).Format(time.DateTime))
	default:
		return ""
	}
}

func (q *queryBuilder) createDefaultValue(column *column, value reflect.Value) reflect.Value {
	if column.dataType == TsVector {
		if q.queryType == Where {
			return reflect.ValueOf(createTSQuery(""))
		}
		return reflect.ValueOf(createTSVectors(""))
	}
	if value.IsValid() {
		return value
	}
	if slices.Contains([]string{Timestamp, TimestampWithZone}, column.dataType) {
		return reflect.ValueOf(CurrentTimestamp)
	}
	if slices.Contains([]string{Bool, Boolean}, column.dataType) {
		return reflect.ValueOf(false)
	}
	if slices.Contains([]string{Int, BigInt, Float}, column.dataType) {
		return reflect.ValueOf(0)
	}
	return reflect.ValueOf("")
}

func (q *queryBuilder) createValueWithUnknownColumn(valueRef ref) string {
	switch valueRef.kind {
	case reflect.String:
		return fmt.Sprintf(`'%s'`, valueRef.v.String())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf(`%d`, int(valueRef.v.Int()))
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf(`%f`, valueRef.v.Float())
	case reflect.Bool:
		return fmt.Sprintf(`%t`, valueRef.v.Bool())
	default:
		return fmt.Sprintf("%v", valueRef.v.Interface())
	}
}
