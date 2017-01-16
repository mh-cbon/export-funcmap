package a

import (
	"html/template"

	"github.com/mh-cbon/export-funcmap/test/a"
)

var k = map[string]interface{}{
	"a":  template.JSEscapeString,
	"b":  rr,
	"c":  func() {},
	"yy": a.SomeFn,
	"zz": a.SomeOtherFn,
	"g": func() template.HTML {
		return template.HTML("")
	},
	"e":   func(g bool) bool { return bool(true) },
	"e0":  func(g byte) byte { return 0 },
	"e1":  func(g int) int { return 0 },
	"e2":  func(g int8) int8 { return 0 },
	"e3":  func(g int16) int16 { return 0 },
	"e4":  func(g int32) int32 { return 0 },
	"e5":  func(g int64) int64 { return 0 },
	"e6":  func(g uint) uint { return 0 },
	"e7":  func(g uint8) uint8 { return 0 },
	"e8":  func(g uint16) uint16 { return 0 },
	"e9":  func(g uint32) uint32 { return 0 },
	"e10": func(g uint64) uint64 { return 0 },
	"e11": func(g float32) float32 { return 0 },
	"e12": func(g float64) float64 { return 0.0 },
	"dup": func(g int8) int8 { return 0 },
}

var k2 = template.FuncMap{
	"a":  template.JSEscapeString,
	"b":  rr,
	"c":  func() {},
	"yy": a.SomeFn,
	// "tt": a.SomeFnUnexported,
	"zz": a.SomeOtherFn,
	"ii": a.SomeFnInterface,
	"uu": a.SomeFnPointer,
	"g": func() template.HTML {
		return template.HTML("")
	},
	"dup": func(g string) string { return "" },
}

func rr() {
	if true {
		b := 1
		c := 2
		if b+c > 3 {

		}
	}
}
