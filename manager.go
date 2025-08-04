package housekeeper

import (
	"context"
	"io"
	"reflect"
)

type Manager struct {
	// Value cache
	cache map[string]reflect.Value
	// Provider registry
	providers map[string]reflect.Value
	// Configurable options
	options options
}

// Get resolves value from manager.
//
// The ptrptr should be a pointer to pointer to target value.
func (m *Manager) Get(ptrptr any) (err error) {
	return m.GetWithContext(context.TODO(), ptrptr)
}

// GetWithContext resolves value from manager with context.
//
// The ptrptr should be a pointer to pointer to target value.
func (m *Manager) GetWithContext(ctx context.Context, ptrptr any) (err error) {
	pv := reflect.ValueOf(ptrptr)
	pt := pv.Type()
	if err = assertPtrToPtrType(pt); err != nil {
		return
	}
	v, err := m.resolveValue(reflect.ValueOf(ctx), pt.Elem())
	if err == nil {
		pv.Elem().Set(v)
	}
	return
}

// Put puts value to manager.
// The ptr should be a pointer to value.
func (m *Manager) Put(ptr any) (err error) {
	pt := reflect.TypeOf(ptr)
	if err = assertPtrType(pt); err != nil {
		return
	}
	cacheKey := getTypeName(pt)
	m.cache[cacheKey] = reflect.ValueOf(ptr)
	return
}

// MustPut puts several values to manager, invalid values will be skipped, and
// all errors will be ignored.
func (m *Manager) MustPut(ptrs ...any) {
	for _, ptr := range ptrs {
		m.Put(ptr)
	}
}

// Close closes the manager, and all managed values.
func (m *Manager) Close() (err error) {
	for _, rv := range m.cache {
		v := rv.Interface()
		if c, ok := v.(io.Closer); ok {
			c.Close()
		}
	}
	clear(m.cache)
	return
}

// GetFrom gets value from manager with generic support.
//
// Deprecated: Use Get/GetWithContext instead.
func GetFrom[V any](m *Manager) (value *V, err error) {
	err = m.Get(&value)
	return
}

// New creates a value manager.
func New(opts ...Option) *Manager {
	// Create manager
	return &Manager{
		// Initialize cache
		cache: make(map[string]reflect.Value),
		// Initialize provider registry
		providers: make(map[string]reflect.Value),
		// Save options
		options: mergeOptions(opts...),
	}
}
