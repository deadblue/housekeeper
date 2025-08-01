package housekeeper_test

import (
	"fmt"
	"log"

	"github.com/deadblue/housekeeper"
)

type FooService struct{}

func (f *FooService) GetName() string {
	return "world"
}

// Close method will be called during housekeeper manager closing.
func (f *FooService) Close() error {
	fmt.Println("Shutdown foo service.")
	return nil
}

type BarService struct{}

func (b *BarService) Greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Custom provider for BarService.
func newBarService() (bar *BarService, err error) {
	fmt.Println("Provide BarService through custom provider.")
	bar = &BarService{}
	// You can return error when provide BarService failed
	return
}

type App struct {
	// Foo service will be assigned in Init method.
	foo *FooService
	// Bar field with autowre tag will be automatically assigned by housekeeper.
	// Autowire field should be exported.
	Bar *BarService `autowire:""`
}

// Init method will be called by housekeeper during making App value.
// Init method should be exported.
func (a *App) Init(foo *FooService) (err error) {
	fmt.Println("Initialize app with foo service.")
	a.foo = foo
	// You can return error init failed
	return nil
}

func (a *App) Run() {
	name := a.foo.GetName()
	a.Bar.Greet(name)
}

func Example() {
	mgr := housekeeper.New()
	// Close manager
	defer mgr.Close()

	// Register custom provider for BarService
	if err := mgr.Provide(newBarService); err != nil {
		log.Fatalf("Register provider failed: %s", err)
	}

	// Assemble app and run it
	var app *App
	if err := mgr.Get(&app); err != nil {
		log.Fatalf("Get app value failed: %s\n", err)
	} else {
		app.Run()
	}

	// Output:
	// Provide BarService through custom provider.
	// Initialize app with foo service.
	// Hello, world!
	// Shutdown foo service.
}
