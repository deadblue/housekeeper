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

// resolveValue searches value from cache, or makes value and puts it to cache
// when absent.
//
// The input |pt| should be a pointer type.
func (m *Manager) resolveValue(
	ctxVal reflect.Value,
	pt reflect.Type,
	stack ...string,
) (pv reflect.Value, err error) {
	if err = assertPtrType(pt); err != nil {
		return
	}
	cacheKey := getTypeName(pt)
	var found bool
	if pv, found = m.cache[cacheKey]; !found {
		if pv, err = m.makeValue(ctxVal, pt, stack...); err == nil {
			m.cache[cacheKey] = pv
		}
	}
	return
}

// makeValue makes a value follow these steps:
//
//   - Provide value.
//   - Wire fields that has autoware tag.
//   - Call Init method on value.
func (m *Manager) makeValue(
	ctxVal reflect.Value,
	pt reflect.Type,
	stack ...string,
) (pv reflect.Value, err error) {
	// Circular-reference checking when make value
	typeName := getTypeName(pt)
	if slices.Contains(stack, typeName) {
		err = errCircularReference
		return
	}
	nextStack := append([]string{typeName}, stack...)

	// Provide value
	if pv, err = m.provideValue(ctxVal, pt, nextStack...); err != nil {
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
		err = m.wireStructFields(ctxVal, elemType, pv.Elem(), nextStack...)
		if err != nil {
			return
		}
	}
	// Call Init method on pv when present
	err = m.initValue(ctxVal, pt, pv, nextStack...)
	return
}

func (m *Manager) initValue(
	ctxVal reflect.Value,
	pt reflect.Type,
	pv reflect.Value,
	stack ...string,
) (err error) {
	// Find init method
	im, found := pt.MethodByName(m.options.InitMethodName)
	if !found {
		return
	}
	// Prepare method parameters
	paramCount := im.Type.NumIn()
	params := make([]reflect.Value, paramCount)
	for i := range paramCount {
		// Set receiver
		if i == 0 {
			params[i] = pv
			continue
		}
		// Prepare other parameters
		if paramType := im.Type.In(i); isContextType(paramType) {
			params[i] = ctxVal
		} else {
			params[i], err = m.resolveValue(ctxVal, paramType, stack...)
			if err != nil {
				return fmt.Errorf(
					"resolve method \"%s.%s\" argument #%d failed: %w",
					stack[0], im.Name, i, err,
				)
			}
		}
	}
	// Call init method
	return findError(im.Func.Call(params))
}
