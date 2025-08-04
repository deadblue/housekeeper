# Housekeeper

A lightweight, reflection-based dependency-injection and lifecycle management 
framework.

[Example](https://pkg.go.dev/github.com/deadblue/housekeeper#example-package).

## Mechanism

### Initialize

Housekeeper manager resolve values follow this steps:

- Retrieve value by type from cache, return it when found.
- When not found, follow this steps to make it:
    1. Call provider for type when present, or simply new it.
    2. Wire value to fields which has `autowire` tag.
    3. Call `Init` method on target value when defined.
    4. Put made value to cache.
- Repeat these steps for parameters of provider function or `Init` method, and 
values for`autowire` field.

> You can customize `Init` method name through `InitMethodOption` when making manager.

### Fianlize

During manager closing, all managed values that implements `io.Closer` will be 
closed, developer can release resources in its `Close` method.

## Specification

### Provider

Provider is a function that returns a pointer value to target type. It should 
follow these specifications:

* Provider function should have at least one result.
* The type of first result should be a pointer type, which will be treated as 
  target type.
* The type of second or later result can be an error type, which will be treated
  as providing error.
* All parameters of provider function, should be pointer type, except for
  `context.Context` parameter.
* Provider function can NOT have variadic parameter.

### Init Method

Init method is the method to initialize value, It should follow these 
specifications:

* Init method should be exportable.
* The receiver of init method should be the pointer to value.
* Init method can have a result of error type, which will be treated as 
  initializing error.
* All parameters of init method, should be pointer type, except for 
  `context.Context` parameter.
* Init method can NOT have variadic parameter.


## Limitation

1. All values managed by housekeeper will be cached, in other words, they are 
   SINGLETON.
2. You SHOULD declare dependent value in pointer type.

## License

MIT
