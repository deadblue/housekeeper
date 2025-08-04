package housekeeper

import (
	"fmt"
	"log"
	"reflect"
)

const (
	_TagAutowire = "autowire"
)

func (m *Manager) wireStructFields(
	ctxVal reflect.Value,
	st reflect.Type,
	sv reflect.Value,
	stack ...string,
) (err error) {
	typeName := stack[0]
	for index := range st.NumField() {
		ft := st.Field(index)
		// Skip non-autowire field
		if _, found := ft.Tag.Lookup(_TagAutowire); !found {
			continue
		}
		fv := sv.Field(index)
		if !fv.CanSet() {
			log.Printf("Skip unexported field: %s.%s", typeName, ft.Name)
			continue
		}
		// Assign value to field
		if val, err := m.resolveValue(ctxVal, ft.Type, stack...); err == nil {
			fv.Set(val)
		} else {
			// Breaks for loop and returns error
			return fmt.Errorf(
				"resolve field \"%s.%s\" value failed: %w",
				typeName, ft.Name, err,
			)
		}
	}
	return nil
}
