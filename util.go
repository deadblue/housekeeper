package housekeeper

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errInvalidValueForGet = errors.New("value should be a pointer to pointer")

	errInvalidValue = errors.New("value should be a pointer")

	errorType = reflect.TypeFor[error]()
)

func checkTypeForGet(r reflect.Type) error {
	if r.Kind() != reflect.Pointer ||
		r.Elem().Kind() != reflect.Pointer ||
		r.Elem().Elem().Kind() == reflect.Pointer {
		return errInvalidValueForGet
	}
	return nil
}

func checkType(t reflect.Type) error {
	if t.Kind() != reflect.Pointer ||
		t.Elem().Kind() == reflect.Pointer {
		return errInvalidValue
	}
	return nil
}

func getCacheKey(t reflect.Type) string {
	t = t.Elem()
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}

func findError(vals []reflect.Value) error {
	for _, val := range vals {
		if val.Type().AssignableTo(errorType) {
			if val.IsNil() {
				return nil
			} else {
				return val.Interface().(error)
			}
		}
	}
	return nil
}
