package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mh-cbon/export-funcmap/export"
)

var version = "0.0.0"

func main() {

	var help = flag.Bool("help", false, "Show help")
	var shelp = flag.Bool("h", false, "Show help")
	var sver = flag.Bool("v", false, "Show version")

	flag.Parse()

	if *help || *shelp {
		showHelp()
		return
	} else if *sver {
		showVersion()
		return
	}

	args := os.Args[1:]

	// small trick for go run,
	// it needs -- to separate arguments for go run and the runned program
	// but for the final built executable it does not exists.
	// lets detect it and remove it.
	if args[0] == "--" {
		args = args[1:]
	}

	if len(args) < 4 {
		showHelp()
		fmt.Println()
		fmt.Println("Not enough arguments.")
		return
	}
	outfilename := args[0]
	outpackage := args[1]
	outvarname := args[2]

	targets := export.Targets{}
	if err := targets.Parse(args[3:]); err != nil {
		showHelp()
		fmt.Println()
		fmt.Println(err)
		return
	}

	destFile, err := export.Export(targets, outfilename, outpackage, outvarname)
	if err != nil {
		panic(err)
	}

	// print the result.
	export.PrintAstFile(os.Stdout, destFile)
}

func showHelp() {
	fmt.Println(`export-funcmap - ` + version + `
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

Example
	export-funcmap gen.go gen export text/template:builtins
	export-funcmap gen.go gen export text/template:builtins:builtins
	export-funcmap gen.go gen export text/template:builtins text/template:builtins
`)
}
func showVersion() {
	fmt.Println(`export-funcmap - ` + version)
}
