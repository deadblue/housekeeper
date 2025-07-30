package housekeeper

import (
	"log"
	"unicode"
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
				if isExportable(name) {
					result.InitMethodName = name
				} else {
					log.Printf("Can not use unexported method \"%s\" as init method", name)
				}
			}
		default:
			log.Printf("Unsupported option: %v", opt)
		}
	}
	return result
}

func isExportable(s string) bool {
	var chs = []rune(s)
	return len(chs) > 0 && unicode.IsUpper(chs[0])
}
