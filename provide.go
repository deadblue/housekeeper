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
//
// # Specification:
//
//   - Provider should be function.
//   - Provider function must have at least one result.
//   - The first result should be a pointer type, it will be treated as target value type.
//   - The second or later result can be an error type, it will be treated as providing error.
//   - When the function has parameters, all parameter types should be pointer type.
//   - The function can not have a variadic parameter.
//
// # Example:
//
// Valid providers:
//   - func () *ResultType
//   - func (arg1 *Arg1Type, arg2 *Arg2Type) (*ResultType, error)
//
// Invalid providers:
//   - func (arg1 *Arg1Type) (error, *ResultType)
//   - func (arg1 string) *ResultType
//   - func (args ...int) (*ResultType, error)
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

func (m *Manager) provideValue(pt reflect.Type, stack ...string) (pv reflect.Value, err error) {
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
	argCount := provType.NumIn()
	args := make([]reflect.Value, argCount)
	for i := range argCount {
		args[i], err = m.getValue(provType.In(i), stack...)
		if err != nil {
			err = fmt.Errorf(
				"resolve provider %s argument #%d failed: %w",
				getTypeName(provType), i, err,
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
