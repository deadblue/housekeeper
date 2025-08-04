package housekeeper

import (
	"errors"
	"fmt"
	"reflect"
)

var (
	errProviderShouldBeFunction = errors.New("provider should be a function")
	errProviderShouldReturn     = errors.New("provider should return values")
	errProviderShouldReturnPtr  = errors.New("the first returning value of provider should be pointer type")
)

// Provide registers type provider to manager.
func (m *Manager) Provide(provider any) (err error) {
	// Validate provider
	rv := reflect.ValueOf(provider)
	rt := rv.Type()
	provName := getTypeName(rt)
	// Check provider type
	if rt.Kind() != reflect.Func {
		return fmt.Errorf("invalid provider \"%s\": %w", provName, errProviderShouldBeFunction)
	}
	// Check provider function signature
	if rt.NumOut() == 0 {
		return fmt.Errorf("invalid provider \"%s\": %w", provName, errProviderShouldReturn)
	}
	valueType := rt.Out(0)
	if assertPtrType(valueType) != nil {
		return fmt.Errorf("invalid provider \"%s\": %w", provName, errProviderShouldReturnPtr)
	}
	// Register provider
	typeName := getTypeName(valueType)
	m.providers[typeName] = rv
	return
}

// MustProvide registers several type providers to manager.
// Invalid providers will be skipped, and all errors will be ingored.
//
// Use this method ONLY when you are sure all providers are valid.
func (m *Manager) MustProvide(providers ...any) {
	for _, provider := range providers {
		m.Provide(provider)
	}
}

func (m *Manager) provideValue(
	ctxVal reflect.Value,
	pt reflect.Type,
	stack ...string,
) (pv reflect.Value, err error) {
	// Search provider
	typeName := stack[0]
	provider, found := m.providers[typeName]
	if !found {
		// Simply uses new for target type
		pv = reflect.New(pt.Elem())
		return
	}
	// Prepare provider arguments
	provType := provider.Type()
	paramCount := provType.NumIn()
	params := make([]reflect.Value, paramCount)
	for i := range paramCount {
		if paramType := provType.In(i); isContextType(paramType) {
			params[i] = ctxVal
		} else {
			params[i], err = m.resolveValue(ctxVal, paramType, stack...)
			if err != nil {
				err = fmt.Errorf(
					"resolve provider %s argument #%d failed: %w",
					getTypeName(provType), i, err,
				)
				return
			}
		}
	}
	// Call provider
	results := provider.Call(params)
	if err = findError(results); err == nil {
		pv = results[0]
	}
	return
}
