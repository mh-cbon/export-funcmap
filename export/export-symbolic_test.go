package export_test

import (
	"bytes"
	"go/format"
	"regexp"
	"testing"

	"golang.org/x/tools/go/loader"

	"github.com/mh-cbon/export-funcmap/export"
)

type testData struct {
	pkg            string
	varnames       []string
	expectErr      bool
	expectKeyCount int
	expectContents string
	// expectContents []*regexp.Regexp
}

var keyMatch = regexp.MustCompile(`"[^"]+":`)

func TestBool(t *testing.T) {

	export.EnableCache = false

	tpkg := "github.com/mh-cbon/export-funcmap/export/test"

	datas := []testData{
		testData{
			pkg:            tpkg,
			varnames:       []string{"boolfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g bool) bool {
return false
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"bytefn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g byte) byte {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"intfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g int) int {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int8fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g int8) int8 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int16fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g int16) int16 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int32fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g int32) int32 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int64fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g int64) int64 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uintfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g uint) uint {
return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint8fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g uint8) uint8 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint16fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g uint16) uint16 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint32fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g uint32) uint32 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint64fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g uint64) uint64 {
  return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"float32fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g float32) float32 {
return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"float64fn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g float64) float64 {
return 0
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"stringfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g string) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"interfacefn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g interface {
}) interface {
} {
return nil
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"ellipsisfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g ...string) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyArgsEllipsisfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(k int, g ...string) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"errorTypefn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g string) error {
return nil
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyArgsfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g string, v string) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyResultsfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g string) (string, error) {
return "", nil
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyArgsManyResultsfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g string, v string) (string, error) {
return "", nil
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"typeAliasToBasicfn"},
			expectKeyCount: 1,
			expectContents: `package gen

import (
"html/template"
)

var tomate = map[string]interface {
}{"fn": func(g template.HTML) template.HTML {
return template.HTML("")
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"structfn"},
			expectKeyCount: 1,
			expectContents: `package gen

import (
"github.com/mh-cbon/export-funcmap/export/test"
)

var tomate = map[string]interface {
}{"fn": func(g a.SomeStruct) a.SomeStruct {
return a.SomeStruct{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"typeInterfacefn"},
			expectKeyCount: 1,
			expectContents: `package gen

import (
"github.com/mh-cbon/export-funcmap/export/test"
)

var tomate = map[string]interface {
}{"fn": func(g a.SomeInterface) a.SomeInterface {
return nil
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicefn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g []string) []string {
return []string{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicePointerBasicfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g []*string) []*string {
return []*string{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicePointerStructfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g []*a.SomeStruct) []*a.SomeStruct {
return []*a.SomeStruct{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"sliceslicefn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g [][]string) [][]string {
return [][]string{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicesliceStructfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g [][]a.SomeStruct) [][]a.SomeStruct {
return [][]a.SomeStruct{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"mapfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g map[string]interface {
}) map[string]interface {
} {
return map[string]interface {
}{}
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"mapfn", "stringfn"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"fn": func(g string) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"otherstringfn", "stringfn"},
			expectKeyCount: 2,
			expectContents: `package gen

var tomate = map[string]interface {
}{"otherfn": func(o string) string {
return ""
}, "fn": func(g string) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"notbytesfn"},
			expectKeyCount: 1,
			expectContents: `package gen

import (
"bytes"
)

var tomate = map[string]interface {
}{"fn": func(o bytes.Buffer) string {
return ""
}}`,
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"funcMap"},
			expectKeyCount: 1,
			expectContents: `package gen

var tomate = map[string]interface {
}{"_html_template_attrescaper": func() {
}}`,
		},
		testData{
			pkg:       tpkg,
			varnames:  []string{"typeUnexportedfn"},
			expectErr: true,
		},
		testData{
			pkg:       tpkg,
			varnames:  []string{"nosuchvar"},
			expectErr: true,
		},
	}

	progs := map[string]*loader.Program{}

	for _, data := range datas {
		if _, ok := progs[data.pkg]; ok == false {
			prog, err := export.GetProgram([]string{data.pkg})
			if err != nil {
				panic(err)
			}
			progs[data.pkg] = prog
		}
		if execTest(data, t, progs[data.pkg]) == false {
			break
		}
		//note, add concurrency does not enhance execution time
		// (guess, the loader is already highly concurrent)
	}
}

func execTest(data testData, t *testing.T, prog *loader.Program) bool {

	targets := []export.Target{
		export.Target{
			PkgPath: data.pkg,
			Idents:  data.varnames,
		},
	}
	_, destFile := export.NewPkg("gen.go", "gen")
	mapVar, imported, err := export.Symbolic(targets, "tomate", prog, destFile)

	// create and inject the import statement
	export.AddImportDecl(destFile, imported)

	// add the new var to the file.
	destFile.Decls = append(destFile.Decls, mapVar)

	if data.expectErr {
		if err == nil {
			t.Errorf(
				"Test %v: Expected an error, got=%v",
				data.varnames, err,
			)
		}
		return false
	} else if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	export.PrintAstFile(&b, destFile)
	str := b.String()

	keysFound := keyMatch.FindAllString(str, -1)
	if len(keysFound) != data.expectKeyCount {
		t.Errorf(
			"Test %v: Invalid count of keys in map results, expected=%v got=%v",
			data.varnames, data.expectKeyCount, len(keysFound),
		)
	}

	if formatGoCode(data.expectContents) != formatGoCode(str) {
		t.Errorf(
			"Test %v: Invalid content did not match,\nexpected=\n%v\n\ngot=\n%v",
			data.varnames, data.expectContents, str,
		)
		return false
	}
	return true

	// for _, r := range data.expectContents {
	// 	if r.MatchString(str) == false {
	// 		t.Errorf(
	// 			"Test %v: Invalid content did not match the regexp, expected=%v got=%v",
	// 			data.varnames, r, str,
	// 		)
	// 	}
	// }
}

func TestParse(t *testing.T) {
	targets := export.Targets{}
	err := targets.Parse([]string{"package:var", "path/package:var2"})
	if err != nil {
		t.Error(err)
	}
	if targets[0].PkgPath != "package" {
		t.Errorf("Expected targets[0].PkgPath=package, got=%v", targets[0].PkgPath)
	}
	if targets[0].Idents[0] != "var" {
		t.Errorf("Expected targets[0].Idents=var, got=%v", targets[0].Idents)
	}
	if targets[1].PkgPath != "path/package" {
		t.Errorf("Expected targets[0].PkgPath=path/package, got=%v", targets[0].PkgPath)
	}
	if targets[1].Idents[0] != "var2" {
		t.Errorf("Expected targets[0].Idents=var2, got=%v", targets[0].Idents)
	}
}

func formatGoCode(s string) string {
	fmtExpected, err := format.Source([]byte(s))
	if err != nil {
		panic(err)
	}
	return string(fmtExpected)
}
