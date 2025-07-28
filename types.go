package housekeeper

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errInvalidPtrType = errors.New("value should be a pointer")

	errInvalidPtrToPtrType = errors.New("value should be a pointer to pointer")
)

func assertPtrType(t reflect.Type) error {
	if t.Kind() != reflect.Pointer ||
		t.Elem().Kind() == reflect.Pointer {
		return errInvalidPtrType
	}
	return nil
}

func assertPtrToPtrType(r reflect.Type) error {
	if r.Kind() != reflect.Pointer ||
		r.Elem().Kind() != reflect.Pointer ||
		r.Elem().Elem().Kind() == reflect.Pointer {
		return errInvalidPtrToPtrType
	}
	return nil
}

func getTypeName(t reflect.Type) string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	return fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
}
