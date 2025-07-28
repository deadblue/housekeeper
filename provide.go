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
	// Check provider type
	if rt.Kind() != reflect.Func {
		return errProviderShouldBeFunction
	}
	// Check provider function signature
	if rt.NumOut() == 0 {
		return errProviderShouldReturn
	}
	valueType := rt.Out(0)
	if assertPtrType(valueType) != nil {
		return errProviderShouldReturnPtr
	}
	// Register provider
	providerKey := getTypeName(valueType)
	m.providers[providerKey] = rv
	return
}

func (m *Manager) provideValue(pt reflect.Type) (pv reflect.Value, err error) {
	// Search provider
	providerKey := getTypeName(pt)
	provider, found := m.providers[providerKey]
	if !found {
		// Simply uses new for target type
		pv = reflect.New(pt.Elem())
		return
	}
	// Prepare provider arguments
	provType := provider.Type()
	argCount := provType.NumIn()
	args := make([]reflect.Value, argCount)
	for i := range argCount {
		// TODO: Handle circular-reference
		args[i], err = m.getValue(provType.In(i))
		if err != nil {
			err = errors.Join(
				fmt.Errorf("can not prepare argument #%d for provider: %s", i, getTypeName(pt)),
				err,
			)
			return
		}
	}
	// Call provider
	results := provider.Call(args)
	if err = findError(results); err == nil {
		pv = results[0]
	}
	return
}
