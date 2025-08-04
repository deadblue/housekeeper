package housekeeper_test

import (
	"context"
	"fmt"
	"log"

	"github.com/deadblue/housekeeper"
)

type contextKey struct{}

var configKey = contextKey{}

type AppConfig struct {
	DbUri string
}

type FooService struct{}

func (f *FooService) Init(ctx context.Context) error {
	config := ctx.Value(configKey).(*AppConfig)
	fmt.Printf("[Foo] Connect tot database: %s.\n", config.DbUri)
	return nil
}

func (f *FooService) GetName() string {
	return "world"
}

// Close method will be called during housekeeper manager closing.
func (f *FooService) Close() error {
	fmt.Println("[Foo] Close database connection.")
	return nil
}

type BarService struct{}

func (b *BarService) Greet(name string) {
	fmt.Printf("[Bar] Hello, %s!\n", name)
}

// Custom provider for BarService.
func newBarService() *BarService {
	fmt.Println("[Bar] Provide via newBarService.")
	return &BarService{}
}

type App struct {
	// Fields with autowre tag will be assigned by housekeeper.
	// Autowire fields should be exported.

	Foo *FooService `autowire:""`
	Bar *BarService `autowire:""`
}

func (a *App) Run() {
	a.Bar.Greet(a.Foo.GetName())
}

func Example() {
	mgr := housekeeper.New()
	// Close manager
	defer mgr.Close()

	// Register custom provider for BarService
	mgr.MustProvide(newBarService)

	// Prepare config
	ctx := context.WithValue(context.TODO(), configKey, &AppConfig{
		DbUri: "mysql://127.0.0.1:3306/mydb",
	})

	// Assemble app and run it
	var app *App
	if err := mgr.GetWithContext(ctx, &app); err != nil {
		log.Fatalf("Get app value failed: %s\n", err)
	} else {
		app.Run()
	}

	// Output:
	// [Foo] Connect tot database: mysql://127.0.0.1:3306/mydb.
	// [Bar] Provide via newBarService.
	// [Bar] Hello, world!
	// [Foo] Close database connection.
}
