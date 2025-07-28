package housekeeper

import (
	"log"
)

const (
	_DefaultInitMethodName = "Init"
)

type Option interface {
	isOption()
}

// InitMethodOption for user who wants to custom init method name.
type InitMethodOption string

func (p InitMethodOption) isOption() {}

type options struct {
	InitMethodName string
}

func mergeOptions(opts ...Option) options {
	result := options{
		InitMethodName: _DefaultInitMethodName,
	}
	for _, opt := range opts {
		switch opt := opt.(type) {
		case InitMethodOption:
			if name := string(opt); name != "" {
				result.InitMethodName = name
			}
		default:
			log.Printf("Unsupported option: %v", opt)
		}
	}
	return result
}
