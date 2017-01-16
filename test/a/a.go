package a

import (
	"github.com/mh-cbon/export-funcmap/test/b"
	"github.com/mh-cbon/export-funcmap/test/d"
)

type Tomate interface{}

type unexported struct{}

func SomeFn(s string) b.SomeType {
	return b.SomeType{}
}

func SomeFnUnexported(s string) unexported {
	return unexported{}
}

func SomeFnPointer(s string) *b.SomeType {
	return &b.SomeType{}
}

func SomeFnInterface(s string) Tomate {
	return nil
}

var SomeOtherFn = d.FNd
