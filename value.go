package housekeeper

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
)

var (
	errCircularReference = errors.New("circular reference")
)

// getValue tries to get value from cache, or make value and put it to cache
// when absent.
//
// The input |pt| should be a pointer type.
func (m *Manager) getValue(pt reflect.Type, stack ...string) (pv reflect.Value, err error) {
	if err = assertPtrType(pt); err != nil {
		return
	}
	cacheKey := getTypeName(pt)
	var found bool
	if pv, found = m.cache[cacheKey]; !found {
		if pv, err = m.makeValue(pt, stack...); err == nil {
			m.cache[cacheKey] = pv
		}
	}
	return
}

func (m *Manager) makeValue(pt reflect.Type, stack ...string) (pv reflect.Value, err error) {
	// Circular-reference checking when make value
	typeName := getTypeName(pt)
	if slices.Contains(stack, typeName) {
		err = errCircularReference
		return
	}
	nextStack := append([]string{typeName}, stack...)

	// Provide value
	if pv, err = m.provideValue(pt, nextStack...); err != nil {
		return
	}
	// Reset pv when error is not nil
	defer func() {
		if err != nil {
			pv = reflect.Zero(pt)
		}
	}()
	// Autowire struct fields
	if elemType := pt.Elem(); elemType.Kind() == reflect.Struct {
		if err = m.wireStructFields(elemType, pv.Elem(), nextStack...); err != nil {
			return
		}
	}
	// Call Init method on pv when present
	err = m.initValue(pt, pv, nextStack...)
	return
}

func (m *Manager) initValue(pt reflect.Type, pv reflect.Value, stack ...string) (err error) {
	// Find init method
	im, found := pt.MethodByName(m.options.InitMethodName)
	if !found {
		return
	}
	// Prepare method argument
	argCount := im.Type.NumIn()
	args := make([]reflect.Value, argCount)
	for i := range argCount {
		if i == 0 {
			args[i] = pv
		} else {
			args[i], err = m.getValue(im.Type.In(i), stack...)
			if err != nil {
				return fmt.Errorf(
					"resolve method \"%s.%s\" argument #%d failed: %w",
					stack[0], im.Name, i, err,
				)
			}
		}
	}
	// Call init method
	return findError(im.Func.Call(args))
}
