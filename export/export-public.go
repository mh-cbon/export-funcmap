package export

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"golang.org/x/tools/go/loader"
)

// PublicIdents exports
// a map of pkg import path and function ident
// for every values of the target funcMap defined
// as an ident or a selector expression.
// For a funcmap defined such
// package y
// var x := map[string]interface{}{
//  "f1": template.HTMLEscaper,
//  "f2": PublicFunc,
// }
// PublicIdents exports their information
// var yy = []map[string][string]{
//  map[string]string{
//   "FuncName": "f1",
//   "Sel": "template.HTMLEscaper",
//   "Pkg": "html/template",
//  },
//   "FuncName": "f2",
//   "Sel": "y.PublicFunc",
//   "Pkg": "some/package/path",
//  },
//}
func PublicIdents(targetPackagePaths Targets, outvarname string, prog *loader.Program, destFile *ast.File) (ast.Decl, error) {

	var err error
	var res []map[string]string

	for _, targetPackagePath := range targetPackagePaths {

		ourpkg := prog.Package(targetPackagePath.PkgPath)

		for _, searchIdent := range targetPackagePath.Idents {
			found := false
			for id := range findMapStringInterface(ourpkg.Pkg, ourpkg.Defs) {
				if id.Name == searchIdent {
					found = true
					valueSpec := id.Obj.Decl.(*ast.ValueSpec)
					if len(valueSpec.Values) > 0 {
						keyValues := valueSpec.Values[0].(*ast.CompositeLit).Elts
						for _, e := range keyValues {
							expr := e.(*ast.KeyValueExpr).Value
							key := e.(*ast.KeyValueExpr).Key.(*ast.BasicLit)
							switch node := expr.(type) {
							case *ast.Ident:
								if ast.IsExported(node.Name) {
									res = append(res, map[string]string{
										"FuncName": key.Value[1 : len(key.Value)-1], // remove quotes
										"Sel":      ourpkg.Pkg.Name() + "." + node.Name,
										"Pkg":      ourpkg.Pkg.Path(),
									})
								}
							case *ast.SelectorExpr:
								if ast.IsExported(node.Sel.Name) {

									pkgg := ourpkg.Uses[node.X.(*ast.Ident)]
									// dirty way :x
									// str will look like
									// package alias ("html/template")
									// or
									// package fmt
									str := pkgg.String()
									str = str[8:] // get ride of package
									if strings.Index(str, " ") > -1 {
										str = strings.Split(str, " ")[1]
										str = str[2 : len(str)-2] // get ride of parenthesis and quotes
									}
									importPath := str
									pkgName := filepath.Base(importPath)

									res = append(res, map[string]string{
										"FuncName": key.Value[1 : len(key.Value)-1], // remove quotes
										"Sel":      pkgName + "." + node.Sel.Name,
										"Pkg":      importPath,
									})
								}
							case *ast.FuncLit:
								//pass
							default:
								panic(
									fmt.Errorf("export.PublicIdents: unhandled ast node type %v\n%#v",
										node, node),
								)
							}
						}
					}
				}
			}
			if found == false {
				return nil, fmt.Errorf(
					"variable %v not found in %v",
					searchIdent,
					targetPackagePath.PkgPath,
				)
			}
		}
	}

	gocode := `package y
  var ` + outvarname + ` [] map[string]string = ` + fmt.Sprintf("%#v", res)
	astNode := stringToAst(gocode)

	return astNode.Decls[0], err
}

func stringToAst(gocode string) *ast.File {
	f, err := parser.ParseFile(token.NewFileSet(), "", gocode, 0)
	if err != nil {
		err := fmt.Errorf(
			"stringToAst: Failed to convert string to ast\n%v",
			gocode)
		panic(err)
	}
	return f
}
