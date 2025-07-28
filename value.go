package housekeeper

import (
	"reflect"
)

const (
	_MethodInit = "Init"
)

// getValue tries to get value from cache, or make value and put it to cache
// when absent.
//
// The input |pt| should be a pointer type.
func (m *Manager) getValue(pt reflect.Type) (pv reflect.Value, err error) {
	if err = assertPtrType(pt); err != nil {
		return
	}
	cacheKey := getTypeName(pt)
	var found bool
	if pv, found = m.cache[cacheKey]; !found {
		if pv, err = m.makeValue(pt); err == nil {
			m.cache[cacheKey] = pv
		}
	}
	return
}

func (m *Manager) makeValue(pt reflect.Type) (pv reflect.Value, err error) {
	if pv, err = m.provideValue(pt); err != nil {
		return
	}
	// Call Init method on pv when present
	if err = m.initValue(pt, pv); err != nil {
		pv = reflect.Zero(pt)
	}
	return
}

func (m *Manager) initValue(pt reflect.Type, pv reflect.Value) (err error) {
	// Find init method
	im, found := pt.MethodByName(_MethodInit)
	if !found {
		return
	}
	// Prepare method argument
	ft := im.Func.Type()
	argCount := ft.NumIn()
	args := make([]reflect.Value, argCount)
	args[0] = pv
	for i := 1; i < argCount; i++ {
		// TODO: Handle circular-reference
		args[i], err = m.getValue(ft.In(i))
		if err != nil {
			return
		}
	}
	// Call init method
	return findError(im.Func.Call(args))
}
