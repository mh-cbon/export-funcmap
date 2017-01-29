package a

import (
	notbytes "bytes"
	"html/template"
	text "text/template"
)

var stringfn = map[string]interface{}{
	"fn": func(g string) string { return "" },
}
var boolfn = map[string]interface{}{
	"fn": func(g bool) bool { return true },
}
var bytefn = map[string]interface{}{
	"fn": func(g byte) byte { return 0 },
}
var intfn = map[string]interface{}{
	"fn": func(g int) int { return 0 },
}
var int8fn = map[string]interface{}{
	"fn": func(g int8) int8 { return 0 },
}
var int16fn = map[string]interface{}{
	"fn": func(g int16) int16 { return 0 },
}
var int32fn = map[string]interface{}{
	"fn": func(g int32) int32 { return 0 },
}
var int64fn = map[string]interface{}{
	"fn": func(g int64) int64 { return 0 },
}
var uintfn = map[string]interface{}{
	"fn": func(g uint) uint { return 0 },
}
var uint8fn = map[string]interface{}{
	"fn": func(g uint8) uint8 { return 0 },
}
var uint16fn = map[string]interface{}{
	"fn": func(g uint16) uint16 { return 0 },
}
var uint32fn = map[string]interface{}{
	"fn": func(g uint32) uint32 { return 0 },
}
var uint64fn = map[string]interface{}{
	"fn": func(g uint64) uint64 { return 0 },
}
var float32fn = map[string]interface{}{
	"fn": func(g float32) float32 { return 0 },
}
var float64fn = map[string]interface{}{
	"fn": func(g float64) float64 { return 0.0 },
}
var interfacefn = map[string]interface{}{
	"fn": func(g interface{}) interface{} { return nil },
}
var ellipsisfn = map[string]interface{}{
	"fn": func(g ...string) string { return "" },
}
var manyArgsEllipsisfn = map[string]interface{}{
	"fn": func(k int, g ...string) string { return "" },
}
var errorTypefn = map[string]interface{}{
	"fn": func(g string) error { return nil },
}
var manyArgsfn = map[string]interface{}{
	"fn": func(g string, v string) string { return "" },
}
var manyResultsfn = map[string]interface{}{
	"fn": func(g string) (string, error) { return "", nil },
}
var manyArgsManyResultsfn = map[string]interface{}{
	"fn": func(g string, v string) (string, error) { return "", nil },
}
var typeAliasToBasicfn = map[string]interface{}{
	"fn": func(g template.HTML) template.HTML { return template.HTML("") },
}
var structfn = map[string]interface{}{
	"fn": func(g SomeStruct) SomeStruct { return SomeStruct{} },
}
var typeInterfacefn = map[string]interface{}{
	"fn": func(g SomeInterface) SomeInterface { return nil },
}
var typeUnexportedfn = map[string]interface{}{
	"fn": func(g string) unexportedType { return nil },
}
var slicefn = map[string]interface{}{
	"fn": func(g []string) []string { return []string{} },
}
var sliceslicefn = map[string]interface{}{
	"fn": func(g [][]string) [][]string { return [][]string{} },
}
var slicePointerBasicfn = map[string]interface{}{
	"fn": func(g []*string) []*string { return []*string{} },
}
var slicePointerStructfn = map[string]interface{}{
	"fn": func(g []*SomeStruct) []*SomeStruct { return []*SomeStruct{} },
}
var slicesliceStructfn = map[string]interface{}{
	"fn": func(g [][]SomeStruct) [][]SomeStruct { return [][]SomeStruct{} },
}
var mapfn = map[string]interface{}{
	"fn": func(g map[string]interface{}) map[string]interface{} { return map[string]interface{}{} },
}

var otherstringfn = map[string]interface{}{
	"otherfn": func(o string) string { return "" },
}
var notbytesfn = map[string]interface{}{
	"fn": func(o notbytes.Buffer) string { return "" },
}

var funcMap = text.FuncMap{
	"_html_template_attrescaper": func() {},
}

// SomeStruct with a comment.
type SomeStruct struct{}

// SomeInterface with a comment.
type SomeInterface interface{}
type unexportedType interface{}
