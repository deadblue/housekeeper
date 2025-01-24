package housekeeper

import (
	"errors"
	"log"
)

type Bar struct {
	state string
}

func (b *Bar) Init() {
	b.state = "bar is inited"
}

func (b *Bar) Close() (err error) {
	println("close bar")
	return
}

type Foo struct {
	bar *Bar
}

func (f *Foo) Init(bar *Bar) (err error) {
	if bar == nil {
		err = errors.New("Bar is nil!")
	} else {
		f.bar = bar
	}
	return
}

func (f *Foo) PrintBarState() {
	println(f.bar.state)
}

func ExmapleGetFrom() {
	mgr := New()
	defer mgr.Close()

	foo, err := GetFrom[Foo](mgr)
	if err != nil {
		log.Fatalf("Get foo failed: %s", err)
	}
	foo.PrintBarState()
}

func ExmapleManager_Get() {
	mgr := New()
	defer mgr.Close()

	var foo *Foo
	if err := mgr.Get(&foo); err != nil {
		log.Fatalf("Get foo failed: %s", err)
	}
	foo.PrintBarState()
}
