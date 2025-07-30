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

func (f *FooService) Close() error {
	fmt.Println("Shutdown foo service.")
	return nil
}

type BarService struct{}

func (b *BarService) Greet(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Custom provider for BarService.
// provider function can be unexported.
func newBarService() (bar *BarService, err error) {
	fmt.Println("Provide BarService through custom provider.")
	bar = &BarService{}
	// You can return error when provide BarService failed
	return
}

type App struct {
	// FooService will be auto-wired by housekeeper.
	// Autowire field should be exported.
	Foo *FooService `autowire:""`
	// Bar service will be injected in init method
	bar *BarService
}

// Init method will be called by housekeeper.
func (a *App) Init(bar *BarService) (err error) {
	fmt.Println("Initialize app.")
	a.bar = bar
	// You can return error init failed
	return nil
}

func (a *App) Run() {
	name := a.Foo.GetName()
	a.bar.Greet(name)
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
	// Initialize app.
	// Hello, world!
	// Shutdown foo service.
}
