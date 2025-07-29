package housekeeper

import (
	"errors"
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
	if pv, err = m.provideValue(pt, nextStack...); err != nil {
		return
	}
	// Autowire struct fields
	if elemType := pt.Elem(); elemType.Kind() == reflect.Struct {
		if err = m.wireStructFields(elemType, pv.Elem(), nextStack...); err != nil {
			pv = reflect.Zero(pt)
			return
		}
	}
	// Call Init method on pv when present
	if err = m.initValue(pt, pv, nextStack...); err != nil {
		pv = reflect.Zero(pt)
	}
	return
}

func (m *Manager) initValue(pt reflect.Type, pv reflect.Value, stack ...string) (err error) {
	// Find init method
	im, found := pt.MethodByName(m.options.InitMethodName)
	if !found {
		return
	}
	// Prepare method argument
	ft := im.Func.Type()
	argCount := ft.NumIn()
	args := make([]reflect.Value, argCount)
	args[0] = pv
	for i := 1; i < argCount; i++ {
		args[i], err = m.getValue(ft.In(i), stack...)
		if err != nil {
			return
		}
	}
	// Call init method
	return findError(im.Func.Call(args))
}
