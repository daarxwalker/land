package land

import (
	"fmt"
	"reflect"
)

type queryBuilder struct {
}

func createQueryBuilder() *queryBuilder {
	return &queryBuilder{}
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

func (q *queryBuilder) createValue(column *column, value reflect.Value) string {
	if column == nil {
		return ""
	}
	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}
	switch column.dataType {
	case TsVector:
		return fmt.Sprintf(`%s`, createTSQuery(value.String()))
	case Varchar:
		return fmt.Sprintf(`'%s'`, value.String())
	case Text:
		return fmt.Sprintf(`'%s'`, value.String())
	case Char:
		return fmt.Sprintf(`'%s'`, value.String())
	case Int:
		return fmt.Sprintf(`%d`, value.Int())
	case BigInt:
		return fmt.Sprintf(`%d`, value.Int())
	case Float:
		return fmt.Sprintf(`%f`, value.Float())
	case Boolean:
		return fmt.Sprintf(`%t`, value.Bool())
	case Timestamp:
		return fmt.Sprintf(`%v`, value.Interface())
	case TimestampWithZone:
		return fmt.Sprintf(`%v`, value.Interface())
	default:
		return ""
	}
}
