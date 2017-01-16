package export_test

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/mh-cbon/export-funcmap/export"
)

type testData struct {
	pkg            string
	varnames       []string
	expectErr      bool
	expectKeyCount int
	expectContents []*regexp.Regexp
}

var keyMatch = regexp.MustCompile(`"[^"]+":`)

func TestBool(t *testing.T) {
	tpkg := "github.com/mh-cbon/export-funcmap/export/test"
	datas := []testData{
		testData{
			pkg:            tpkg,
			varnames:       []string{"boolfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+bool\)\s+bool\s+\{\s+return\s+false\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"bytefn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+byte\)\s+byte\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"intfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+int\)\s+int\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int8fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+int8\)\s+int8\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int16fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+int16\)\s+int16\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int32fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+int32\)\s+int32\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"int64fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+int64\)\s+int64\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uintfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+uint\)\s+uint\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint8fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+uint8\)\s+uint8\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint16fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+uint16\)\s+uint16\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint32fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+uint32\)\s+uint32\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"uint64fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+uint64\)\s+uint64\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"float32fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+float32\)\s+float32\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"float64fn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+float64\)\s+float64\s+\{\s+return\s+0\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"stringfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"interfacefn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+interface\s+\{\s+\}\)\s+interface\s+\{\s+\}\s+\{\s+return\s+nil\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"ellipsisfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+\.\.\.string\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyArgsEllipsisfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(k int, g\s+\.\.\.string\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"errorTypefn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string\)\s+error\s+\{\s+return\s+nil\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyArgsfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string,\s+v\s+string\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyResultsfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string\)\s+\(string,\s+error\)\s+\{\s+return\s+"",\s+nil\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"manyArgsManyResultsfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string,\s+v\s+string\)\s+\(string,\s+error\)\s+\{\s+return\s+"",\s+nil\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"typeAliasToBasicfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+template\.HTML\)\s+template\.HTML\s+\{\s+return\s+template\.HTML\(""\)\s+\}`),
				regexp.MustCompile(`import "html/template"`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"structfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+a\.SomeStruct\)\s+a\.SomeStruct\s+\{\s+return\s+a.SomeStruct\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"typeInterfacefn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+a\.SomeInterface\)\s+a\.SomeInterface\s+\{\s+return\s+nil\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicefn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+\[\]string\)\s+\[\]string\s+\{\s+return\s+\[\]string\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicePointerBasicfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+\[\]\*string\)\s+\[\]\*string\s+\{\s+return\s+\[\]\*string\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicePointerStructfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+\[\]\*a\.SomeStruct\)\s+\[\]\*a\.SomeStruct\s+\{\s+return\s+\[\]\*a\.SomeStruct\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"sliceslicefn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+\[\]\[\]string\)\s+\[\]\[\]string\s+\{\s+return\s+\[\]\[\]string\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"slicesliceStructfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+\[\]\[\]a\.SomeStruct\)\s+\[\]\[\]a\.SomeStruct\s+\{\s+return\s+\[\]\[\]a\.SomeStruct\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"mapfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+map\[string\]interface\s+\{\s+\}\)\s+map\[string\]interface\s+\{\s+\}\s+\{\s+return\s+map\[string\]interface\s+\{\s+\}\{\}\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"mapfn", "stringfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"otherstringfn", "stringfn"},
			expectKeyCount: 2,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(g\s+string\)\s+string\s+\{\s+return\s+""\s+\}`),
				regexp.MustCompile(`func\(o\s+string\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
		},
		testData{
			pkg:            tpkg,
			varnames:       []string{"notbytesfn"},
			expectKeyCount: 1,
			expectContents: []*regexp.Regexp{
				regexp.MustCompile(`func\(o\s+bytes\.Buffer\)\s+string\s+\{\s+return\s+""\s+\}`),
			},
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

	for _, data := range datas {
		execTest(data, t) //note, add concurrency does not enhance execution time
		// (guess, the loader is already highly concurrent)
	}
}

func execTest(data testData, t *testing.T) {

	targets := []export.Target{
		export.Target{
			PkgPath: data.pkg,
			Idents:  data.varnames,
		},
	}
	newPkg, err := export.Export(targets, "gen.go", "gen", "tomate")

	if data.expectErr {
		if err == nil {
			t.Errorf(
				"Test %v: Expected an error, got=%v",
				data.varnames, err,
			)
		}
		return
	} else if err != nil {
		panic(err)
	}

	var b bytes.Buffer
	export.PrintAstFile(&b, newPkg.Files["gen.go"])
	str := b.String()

	keysFound := keyMatch.FindAllString(str, -1)
	if len(keysFound) != data.expectKeyCount {
		t.Errorf(
			"Test %v: Invalid count of keys in map results, expected=%v got=%v",
			data.varnames, data.expectKeyCount, len(keysFound),
		)
	}

	for _, r := range data.expectContents {
		if r.MatchString(str) == false {
			t.Errorf(
				"Test %v: Invalid content did not match the regexp, expected=%v got=%v",
				data.varnames, r, str,
			)
		}
	}

}
