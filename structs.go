package housekeeper

import (
	"errors"
	"fmt"
	"log"
	"reflect"
)

const (
	_TagAutowire = "autowire"
)

func (m *Manager) structAutowireFields(st reflect.Type, sv reflect.Value) (err error) {
	typeName := fmt.Sprintf("%s/%s", st.PkgPath(), st.Name())
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
		// TODO: Circular-reference checking
		if val, err := m.getValue(ft.Type); err == nil {
			fv.Set(val)
		} else {
			// Breaks for loop and returns error
			return errors.Join(
				fmt.Errorf("can not get value for field: %s.%s", typeName, ft.Name),
				err,
			)
		}
	}
	return nil
}
