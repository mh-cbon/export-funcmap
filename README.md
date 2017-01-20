# export-funcmap

Export a funcmap variable declaration to
- its symbolic version
- its public identifiers, when possible

The output is useful to perform source code analysis of go templates,
it will be similar to this,

```go
package gen

var export = map[string]interface {}
{
  "and": func(arg0 interface {}, args ...interface {}) interface {} {
  	return nil
  },
  "call": func(fn interface {}, args ...interface {}) (interface {}, error) {
  	return nil, nil
  },
  "html": func(args ...interface {}) string {
   return ""
  },
}
var exportPublic []map[string]string = []map[string]string{
  map[string]string{
    "FuncName": "html",
    "Sel": "template.HTMLEscaper",
    "Pkg": "text/template",
  },
}
```

# Install

```sh
go get -u github.com/mh-cbon/export-funcmap
```

# Cli

```sh
export-funcmap - 0.0.0
Export a funcmap variable declaration to its symbolic version.

Usage

	export-funcmap <outfilename> <outpackage> <outvarname> <pkgpath:var...>....

	outfilename
		The output filepath of the export result.
		required.

	outpackage
		The output package name of the export result.
		required.

	outvarname
		The output variable name of the export result.
		required.

	pkgpath:var...
		A repeatable argument of package path (for example html/template),
		followed by at least one semi-colon to indicate the desired variable
		to export.
		Each package path can be followed by multiple semi-colon variable if
		multiple variable needs to be extracted from the same package.
		required.

	-v
		Show version

	-h|--help
		Show help

Example
	export-funcmap gen.go gen export text/template:builtins
	export-funcmap gen.go gen export text/template:builtins:builtins
	export-funcmap gen.go gen export text/template:builtins text/template:builtins
```

# Usage

```go
package main

import(
  "github.com/mh-cbon/export-funcmap/export"
)

func main () {
  targets := []string{
    "text/template:builtins",
    "text/template:builtins",
  }
  outfilename := "gen.go"
  outpackage := "gen"
  outvarname := "funcsMap"

  file, err := export.Export(targets, outfilename, outpackage, outvarname)
  if err != nil {
  	panic(err)
  }

  // print the result.
  export.PrintAstFile(os.Stdout, file)
}
```
