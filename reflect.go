package land

import "reflect"

type ref struct {
	v    reflect.Value
	t    reflect.Type
	kind reflect.Kind
}
