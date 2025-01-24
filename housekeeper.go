package housekeeper

import (
	"io"
	"reflect"
)

const (
	_MethodInit = "Init"
)

type Housekeeper struct {
	cache map[string]reflect.Value
}

/*
Gets finds and returns value from housekeeper.

The value parameters should be a pointer to pointer to target type.

Exmaple:

	type Foo struct {}

	hk := New()
	var foo *Foo
	if err := hk.Get(&foo); err != nil {
		log.Fatal(err)
	}
	// TODO: Works with foo
*/
func (h *Housekeeper) Get(value any) (err error) {
	rv := reflect.ValueOf(value)
	rt := rv.Type()
	if err = checkTypeForGet(rt); err != nil {
		return
	}
	if v, err := h.getValue(rt.Elem()); err == nil {
		rv.Set(v.Addr())
	}
	return
}

// Put puts value to Housekeeper cache.
func (h *Housekeeper) Put(value any) (err error) {
	rt := reflect.TypeOf(value)
	if err = checkType(rt); err != nil {
		return
	}
	cacheKey := getCacheKey(rt.Elem())
	h.cache[cacheKey] = reflect.ValueOf(value)
	return
}

// getValue gets value from cache.
//
// The input t should be a pointer type.
func (h *Housekeeper) getValue(t reflect.Type) (v reflect.Value, err error) {
	if err = checkType(t); err != nil {
		return
	}
	cacheKey := getCacheKey(t)
	var found bool
	if v, found = h.cache[cacheKey]; !found {
		if v, err = h.newValue(t); err == nil {
			h.cache[cacheKey] = v
		}
	}
	return
}

func (h *Housekeeper) newValue(t reflect.Type) (v reflect.Value, err error) {
	v = reflect.New(t.Elem())
	if err = h.initValue(t, v); err != nil {
		v.SetZero()
	}
	return
}

func (h *Housekeeper) initValue(t reflect.Type, v reflect.Value) (err error) {
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
		args[i], err = h.getValue(ft.In(i))
		if err != nil {
			return
		}
	}
	// Call init method
	return findError(im.Func.Call(args))
}

func (h *Housekeeper) Close() (err error) {
	for _, rv := range h.cache {
		v := rv.Interface()
		if c, ok := v.(io.Closer); ok {
			c.Close()
		}
	}
	clear(h.cache)
	return
}

// GetFrom gets value from housekeeper with generic support.
func GetFrom[V any](hk *Housekeeper) (value *V, err error) {
	rt := reflect.TypeFor[*V]()
	if v, err := hk.getValue(rt); err == nil {
		value = v.Interface().(*V)
	}
	return
}

// New creates a housekeeper for you.
func New() *Housekeeper {
	return &Housekeeper{
		cache: make(map[string]reflect.Value),
	}
}
