package housekeeper

var (
	// Default manager for one who just need one manager for application.
	defaultManager = New()
)

// Get gets value from default manager.
func Get(ptrptr any) error {
	return defaultManager.Get(ptrptr)
}

// GetFor gets a pointer value of the type V, from default manager.
func GetFor[V any]() (*V, error) {
	return GetFrom[V](defaultManager)
}

// Put puts value to default manager, the ptr argument should be a pointer to value.
func Put(ptr any) error {
	return defaultManager.Put(ptr)
}

// MustPut puts several values to default manager, and ignores all errors.
func MustPut(ptrs ...any) {
	defaultManager.MustProvide(ptrs...)
}

// Provide registers type provider to default manager.
func Provide(provider any) error {
	return defaultManager.Provide(provider)
}

// MustProvide registers several type providers to default manager, and ignores
// all errors.
func MustProvide(providers ...any) {
	defaultManager.MustProvide(providers...)
}

// Close closes default manager.
func Close() error {
	return defaultManager.Close()
}
