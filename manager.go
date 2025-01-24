package housekeeper

import (
	"io"
	"reflect"
)

const (
	_MethodInit = "Init"
)

type Manager struct {
	cache map[string]reflect.Value
}

/*
Get returns value from manager.

The ptrPtr should be a pointer to pointer to target value type.

For exmaple:

	type Foo struct {}

	mgr := New()
	defer mgr.Close()

	var foo *Foo
	if err := mgr.Get(&foo); err != nil {
		log.Fatal(err)
	}
	// TODO: Works with foo
*/
func (m *Manager) Get(ptrPtr any) (err error) {
	pv := reflect.ValueOf(ptrPtr)
	pt := pv.Type()
	if err = checkTypeForGet(pt); err != nil {
		return
	}
	v, err := m.getValue(pt.Elem())
	if err == nil {
		pv.Elem().Set(v)
	}
	return
}

// Put puts value to manager.
func (m *Manager) Put(valuePtr any) (err error) {
	t := reflect.TypeOf(valuePtr)
	if err = checkType(t); err != nil {
		return
	}
	cacheKey := getCacheKey(t)
	m.cache[cacheKey] = reflect.ValueOf(valuePtr)
	return
}

// getValue gets value from cache.
//
// The input t should be a pointer type.
func (m *Manager) getValue(t reflect.Type) (v reflect.Value, err error) {
	if err = checkType(t); err != nil {
		return
	}
	cacheKey := getCacheKey(t)
	var found bool
	if v, found = m.cache[cacheKey]; !found {
		if v, err = m.newValue(t); err == nil {
			m.cache[cacheKey] = v
		}
	}
	return
}

func (m *Manager) newValue(t reflect.Type) (v reflect.Value, err error) {
	v = reflect.New(t.Elem())
	if err = m.initValue(t, v); err != nil {
		v = reflect.Zero(t)
	}
	return
}

func (m *Manager) initValue(t reflect.Type, v reflect.Value) (err error) {
	// Find init method
	im, ok := t.MethodByName(_MethodInit)
	if !ok {
		return
	}
	// Prepare method argument
	ft := im.Func.Type()
	argCount := ft.NumIn()
	args := make([]reflect.Value, argCount)
	args[0] = v
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

// Close closes manager, all values managed by manager will be closed.
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

// GetFrom gets value from manager with generic.
func GetFrom[V any](m *Manager) (value *V, err error) {
	rt := reflect.TypeFor[*V]()
	rv, err := m.getValue(rt)
	if err == nil {
		value = rv.Interface().(*V)
	}
	return
}

// New creates a manager for you.
func New() *Manager {
	return &Manager{
		cache: make(map[string]reflect.Value),
	}
}
