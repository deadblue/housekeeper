# Housekeeper

A lightweight, reflection-based dependency-injection and lifecycle management 
framework.

## Example

```golang
import (
    "log"

    "github.com/deadblue/housekeeper"
)

type FooService struct {}

// Close is used to release resources, it will be automatically called by 
// housekeeper during manager closing.
func (f *FooService) Close() error {
    // TODO: Release used resource.
    return nil
}

type BarService struct {}

// newBarService is custom BarService provide function.
func newBarService() (*BarService, error) {
    // TODO: Initlaize bar, return error when failed.
    bar := &BarService{}
    return bar, nil
}

type App struct {
    // foo field will be assigned in Init method
    foo *FooService
    // Bar field will be wired by housekeeper.
    Bar *BarService `autowire:""`
}

// App initialize method, it will be automatically called by housekeeper, the 
// arguments will be resolved by housekeeper.
func (a *App) Init(
    foo *FooService,
) error {
    a.foo = foo
    // TODO: Initialize app, return error when failed.
    return nil
}

func (a *App) Run() {
    // Run application
}

func main() {
    mgr := housekeeper.New()
    defer mgr.Close()

    if err := mgr.Provide(newBarService); err != nil {
        log.Printf("Register privoder failed: %s", err)
    }
    // Or call MustProvide when you want to register several providers or 
    // ignore the errors.
    // mgr.MustProvide(newBarService)

    var app *App
    if err := mgr.Get(&app); err != nil {
        log.Fatalf("Initialize app failed: %s", err)
    }
    app.Run()
}
```

## Mechanism

### Initialize

Housekeeper manager resolve values follow this steps:

1. Retrieve value by type from cache. When found, return it.
2. When not found, follow this steps to make it:
    1. Call provider for type when present, or simply new it.
    2. Automatically wire field values which has `autowire` tag.
    3. Call `Init` method on target value when defined.
    4. Cache made value to internal cache.
3. Repeat these steps for arguments of provider function or `Init` method, and 
values for`autowire` field.

> You can custom `Init` method name through `InitMethodOption` when making manager.

### Fianlize

During manager closing, all managed values that implements `io.Closer` will be 
closed, developer can release resources in its `Close` method.

## Limitation

1. All values managed by housekeeper will be cached, in other words, they are SINGLETON.
2. You SHOULD declare dependent value in pointer type.

## License

MIT
