# Housekeeper

Lightweight dependency-injection and lifecycle manager for your go project.

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

    var app *App
    if err := mgr.Get(&app); err != nil {
        log.Fatalf("Initialize app failed: %s", err)
    }
    app.Run()
}
```

## Mechanism

### Initialize

Housekeeper manager will automatically make dependent values follow this steps:

1. Call provider when present, or simply call builtin `new()` function.
2. Automatically wire field values which has `autowire` tag.
3. Call `Init` method on target value when defined.

The arguments of provider or `Init` method, and wired field values will be 
automatically resolved by manager.

The made values will be cached by manager, in other words, all of them are 
singleton.

### Fianlize

During manager closing, all managed values that implements `io.Closer` will be 
closed, developer can release resources in its `Close` method.

## License

MIT