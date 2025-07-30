package housekeeper

import (
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

// Get returns value from manager.
// The ptrptr should be a pointer to pointer to target value type.
func (m *Manager) Get(ptrptr any) (err error) {
	pv := reflect.ValueOf(ptrptr)
	pt := pv.Type()
	if err = assertPtrToPtrType(pt); err != nil {
		return
	}
	v, err := m.getValue(pt.Elem())
	if err == nil {
		pv.Elem().Set(v)
	}
	return
}

// Put puts value to manager, ptr should be a pointer to value.
// Caller can manually put some values into manager, such as configuration, etc.
func (m *Manager) Put(ptr any) (err error) {
	pt := reflect.TypeOf(ptr)
	if err = assertPtrType(pt); err != nil {
		return
	}
	cacheKey := getTypeName(pt)
	m.cache[cacheKey] = reflect.ValueOf(ptr)
	return
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
