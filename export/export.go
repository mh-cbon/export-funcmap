package export

import (
	"fmt"
	"go/ast"
	"go/format"
	"go/token"
	"go/types"
	"io"
	"strings"

	"golang.org/x/tools/go/loader"
)

// Target defines a target package and idents to export.
type Target struct {
	PkgPath string
	Idents  []string
}

// Targets is an alias of []Target
type Targets []Target

// Parse a string of package:var
func (t *Targets) Parse(s []string) error {
	for i := 0; i < len(s); i++ {
		parts := strings.Split(s[i], ":")
		if len(parts) < 2 {
			return fmt.Errorf("Invalid package target: %v", s[i])
		}
		*t = append(*t, Target{
			PkgPath: parts[0],
			Idents:  parts[1:],
		})
	}
	return nil
}

// GetPackagePaths returns the list of package path
func (t *Targets) GetPackagePaths() []string {
	var ret []string
	for _, target := range *t {
		ret = append(ret, target.PkgPath)
	}
	return ret
}

// Export a symbolic map of given target package and ther idents.
func Export(targetPackagePaths Targets, outvarname string, prog *loader.Program, destFile *ast.File) (*ast.GenDecl, []string, error) {

	var err error
	var imported []string

	// Add a varDecl, var xx = map[string]interface{}{}
	mapStrIntDecl, elts := newMapStringInterfaceDelc(outvarname)

	for _, targetPackagePath := range targetPackagePaths {

		ourpkg := prog.Package(targetPackagePath.PkgPath)

		for _, searchIdent := range targetPackagePath.Idents {
			found := false
			for id := range findMapStringInterface(ourpkg.Pkg, ourpkg.Defs) {
				if id.Name == searchIdent {
					found = true
					keyValues := id.Obj.Decl.(*ast.ValueSpec).Values[0].(*ast.CompositeLit).Elts
					for i := range keyValues {

						key := keyValues[i].(*ast.KeyValueExpr).Key.(*ast.BasicLit)

						// Create a key on the map, "x":func(){}
						kv, fn := newKeyValueStringFuncLit(key.Value)
						// Add the kvalue on the map
						injectKvIntoMapStringInterface(kv, elts)

						identLike := keyValues[i].(*ast.KeyValueExpr).Value
						signature := ourpkg.Types[identLike].Type.(*types.Signature)

						var err2 error
						// Define func parameters func(p string...) {}
						in := signature.Params()
						// printTuple(in)
						fn.Type.Params, err2 = newFuncParams(in, signature.Variadic())
						if err2 != nil {
							return nil, nil, err2
						}

						// Define func returns func(...) string... {}
						out := signature.Results()
						// printTuple(out)
						fn.Type.Results, err2 = newFuncResults(out)
						if err2 != nil {
							return nil, nil, err2
						}

						// Define func body func(...) ... { return ""...}
						fn.Body, err2 = newFuncBodyZeroValue(out)
						if err2 != nil {
							return nil, nil, err2
						}

						// extracts imports from the func
						imported = append(imported, extractImports(in)...)
						imported = append(imported, extractImports(out)...)
					}
				}
			}
			if found == false {
				return nil, nil, fmt.Errorf(
					"variable %v not found in %v",
					searchIdent,
					targetPackagePath.PkgPath,
				)
			}
		}
	}

	return mapStrIntDecl, imported, err
}

// GetVarDecl returns the ast node of the variable declaration.
func GetVarDecl(p *ast.File) *ast.GenDecl {
	for _, v := range p.Decls {
		if n, ok := v.(*ast.GenDecl); ok {
			if n.Tok == token.VAR {
				return n
			}
		}
	}
	return nil
}

// GetImportDecl returns the ast node of the import declaration.
func GetImportDecl(p *ast.File) *ast.GenDecl {
	for _, v := range p.Decls {
		if n, ok := v.(*ast.GenDecl); ok {
			if n.Tok == token.IMPORT {
				return n
			}
		}
	}
	return nil
}

// MustGetImportDecl returns the ast node of the import declaration.
func MustGetImportDecl(p *ast.File) *ast.GenDecl {
	i := GetImportDecl(p)
	if i == nil {
		i = NewImportDecl()
	}
	return i
}

//PrintAstFile prints an ast to given writer.
func PrintAstFile(w io.Writer, node interface{}) error {
	fset := token.NewFileSet()
	return format.Node(w, fset, node)
}

func findMapStringInterface(pkg *types.Package, defs map[*ast.Ident]types.Object) map[*ast.Ident]types.Object {
	res := make(map[*ast.Ident]types.Object)

	for id, obj := range defs {
		if obj != nil {
			inner := pkg.Scope().Innermost(id.Pos())
			if inner.Parent() == pkg.Scope() { // find only package level declarations
				if isMapStringInterface(obj.Type()) {
					res[id] = obj
				} else if m, ok := obj.Type().(*types.Named); ok {
					if isMapStringInterface(m.Underlying()) {
						res[id] = obj
					}
				}
			}
		}
	}
	return res
}

func newMapStringInterfaceDelc(varName string) (*ast.GenDecl, *ast.CompositeLit) {
	s := &ast.GenDecl{}
	s.Tok = token.VAR

	spec := &ast.ValueSpec{}
	spec.Names = append(spec.Names, &ast.Ident{Name: varName})

	value := &ast.CompositeLit{}
	spec.Values = append(spec.Values, value)

	mapType := &ast.MapType{}
	mapType.Key = &ast.Ident{Name: "string"}
	interfaceType := &ast.InterfaceType{Methods: &ast.FieldList{}}
	mapType.Value = interfaceType
	value.Type = mapType

	value.Elts = make([]ast.Expr, 0)

	s.Specs = append(s.Specs, spec)

	return s, value
}

func newKeyValueStringFuncLit(keyName string) (*ast.KeyValueExpr, *ast.FuncLit) {
	s := &ast.KeyValueExpr{}
	s.Key = &ast.BasicLit{Kind: token.STRING, Value: keyName}

	f := &ast.FuncLit{}
	f.Type = &ast.FuncType{}
	s.Value = f

	return s, f
}

func newFuncParams(tuple *types.Tuple, isVariadic bool) (*ast.FieldList, error) {
	return typesTupleToAstFieldList(tuple, isVariadic, true)
}

func newFuncResults(tuple *types.Tuple) (*ast.FieldList, error) {
	return typesTupleToAstFieldList(tuple, false, false)
}

func newFuncBodyZeroValue(tuple *types.Tuple) (*ast.BlockStmt, error) {
	s := &ast.BlockStmt{}
	if tuple.Len() > 0 {
		ret := &ast.ReturnStmt{}
		s.List = append(s.List, ret)

		for i := 0; i < tuple.Len(); i++ {
			t, err := typesTypeToAstZeroValue(tuple.At(i).Type(), false)
			if err != nil {
				return nil, err
			}
			ret.Results = append(ret.Results, t)
		}
	}
	return s, nil
}

// typesTypeToAstZeroValue transforms a types.Type
// into an ast expression of its zero value
func typesTypeToAstZeroValue(t types.Type, asIdent bool) (ast.Expr, error) {
	var ret ast.Expr
	switch m := t.(type) {
	case *types.Basic:

		if asIdent == false {
			switch m.Kind() {
			case types.String:
				ret = &ast.BasicLit{Kind: token.STRING, Value: "\"\""}
			case types.Bool:
				ret = &ast.BasicLit{Kind: token.STRING, Value: "false"}
			case types.Int, types.Int8, types.Int16, types.Int32, types.Int64:
				ret = &ast.BasicLit{Kind: token.STRING, Value: "0"}
			case types.Uint, types.Byte /*==Uint8*/, types.Uint16, types.Uint32, types.Uint64:
				ret = &ast.BasicLit{Kind: token.STRING, Value: "0"}
			case types.Float32, types.Float64:
				ret = &ast.BasicLit{Kind: token.STRING, Value: "0"}
			default:
				return nil, fmt.Errorf("Unhandled basic zero value %v", m)
			}
		} else {
			ret = &ast.Ident{Name: m.Name()}
		}

	case *types.Named:
		switch u := t.Underlying().(type) {
		case *types.Basic:
			call := &ast.CallExpr{}
			fun := &ast.SelectorExpr{}
			if m.Obj().Exported() == false {
				return nil, fmt.Errorf("Cannot use unexported type to create a zero value %v", m)
			}
			fun.X = &ast.Ident{Name: m.Obj().Pkg().Name()}
			fun.Sel = &ast.Ident{Name: m.Obj().Name()}
			call.Fun = fun
			t, err := typesTypeToAstZeroValue(u, false)
			if err != nil {
				return nil, err
			}
			call.Args = append(call.Args, t)
			ret = call

		case *types.Struct:
			sel := &ast.SelectorExpr{}
			if m.Obj().Exported() == false {
				return nil, fmt.Errorf("Cannot use unexported type to create a zero value %v", m)
			}
			sel.X = &ast.Ident{Name: m.Obj().Pkg().Name()}
			sel.Sel = &ast.Ident{Name: m.Obj().Name()}
			if asIdent == false {
				ret = &ast.CompositeLit{Type: sel}
			} else {
				ret = sel
			}
		case *types.Interface:
			ret = &ast.Ident{Name: "nil"}

		default:
			fmt.Printf("%#v\n", m)
			return nil, fmt.Errorf("Unhandled named arg zero value %v", m)
		}

	case *types.Pointer:
		if asIdent == false {
			ret = &ast.Ident{Name: "nil"}
		} else {
			x, err := typesTypeToAstZeroValue(m.Elem(), true)
			if err != nil {
				return nil, err
			}
			ret = &ast.StarExpr{X: x}
		}

	case *types.Interface:
		if asIdent == false {
			ret = &ast.Ident{Name: "nil"}
		} else {
			ret = &ast.InterfaceType{Methods: &ast.FieldList{}}
		}

	case *types.Slice:
		x, err := typesTypeToAstZeroValue(m.Elem(), true)
		if err != nil {
			return nil, err
		}
		ret = &ast.ArrayType{Elt: x}
		if asIdent == false {
			ret = &ast.CompositeLit{Type: ret}
		}

	case *types.Map:
		x, err := typesTypeToAstZeroValue(m.Key(), true)
		if err != nil {
			return nil, err
		}
		x2, err2 := typesTypeToAstZeroValue(m.Elem(), true)
		if err2 != nil {
			return nil, err2
		}
		ret = &ast.MapType{Key: x, Value: x2}
		if asIdent == false {
			ret = &ast.CompositeLit{Type: ret}
		}

	default:
		fmt.Printf("%#v\n", m)
		return nil, fmt.Errorf("Unhandled named arg zero value %v", m)
	}
	return ret, nil
}

func extractImports(tuple *types.Tuple) []string {
	ret := make([]string, 0)
	for i := 0; i < tuple.Len(); i++ {
		importPath := typesTupleToImportPath(tuple.At(i).Type())
		if importPath != "" {
			ret = append(ret, importPath)
		}
	}
	return ret
}

func typesTupleToImportPath(t types.Type) string {
	ret := ""
	switch m := t.(type) {
	case *types.Named:
		if m.Obj().Pkg() != nil {
			ret = m.Obj().Pkg().Path()
		}
	case *types.Pointer:
		ret = typesTupleToImportPath(m.Elem())
	}
	return ret
}

func typesTupleToAstFieldList(tuple *types.Tuple, isVariadic, withNames bool) (*ast.FieldList, error) {
	if tuple.Len() == 0 {
		return nil, nil
	}

	var err error
	s := &ast.FieldList{List: make([]*ast.Field, 0)}

	for i := 0; i < tuple.Len(); i++ {
		field := &ast.Field{}
		if withNames {
			name := &ast.Ident{Name: tuple.At(i).Name()}
			field.Names = append(field.Names, name)
		}
		field.Type, err = typesTypeToAstExpr(tuple.At(i).Type(), isVariadic && i == tuple.Len()-1)
		if err != nil {
			return nil, err
		}
		s.List = append(s.List, field)
	}

	return s, err
}

// typesTypeToAstExpr transform a types.Type
// into an ast.Expr suitable for func params/results
func typesTypeToAstExpr(t types.Type, ellisped bool) (ast.Expr, error) {
	switch m := t.(type) {
	case *types.Basic:
		return &ast.Ident{Name: m.Name()}, nil

	case *types.Named:
		isUniverse := m.Obj().Parent().Parent() == nil
		isError := m.Obj().Name() == "error"
		if isUniverse && isError {
			return &ast.Ident{Name: "error"}, nil
		}
		if m.Obj().Exported() == false {
			return nil, fmt.Errorf("Cannot use unexported type %v", m)
		}
		sel := &ast.SelectorExpr{}
		sel.X = &ast.Ident{Name: m.Obj().Pkg().Name()}
		sel.Sel = &ast.Ident{Name: m.Obj().Name()}
		return sel, nil

	case *types.Pointer:
		var err error
		ret := &ast.StarExpr{}
		ret.X, err = typesTypeToAstExpr(m.Elem(), false)
		return ret, err

	case *types.Interface:
		return &ast.InterfaceType{Methods: &ast.FieldList{}}, nil

	case *types.Slice:
		if ellisped {
			t, err := typesTypeToAstExpr(m.Elem(), false)
			return &ast.Ellipsis{Elt: t}, err
		}
		t, err := typesTypeToAstExpr(m.Elem(), false)
		return &ast.ArrayType{Elt: t}, err

	case *types.Map:
		ret := &ast.MapType{}
		t, err := typesTypeToAstExpr(m.Key(), false)
		if err != nil {
			return nil, err
		}
		ret.Key = t
		t2, err2 := typesTypeToAstExpr(m.Elem(), false)
		if err2 != nil {
			return nil, err2
		}
		ret.Value = t2
		return ret, nil
	}
	return nil, fmt.Errorf("Unhandled param type %v", t)
}

func newImportSpec(importPath string, importName string) *ast.ImportSpec {

	s := &ast.ImportSpec{}

	importLit := &ast.BasicLit{}
	importLit.Kind = 9 // need to check this :x
	importLit.Value = `"` + importPath + `"`
	s.Path = importLit

	if importName != "" {
		importIdent := &ast.Ident{}
		importIdent.Name = importName
		s.Name = importIdent
	}

	return s
}

// NewImportDecl creates a new import(...) declaration
func NewImportDecl() *ast.GenDecl {
	importGenDecl := &ast.GenDecl{}
	importGenDecl.Tok = token.IMPORT
	return importGenDecl
}

// InjectImportPaths injects given import paths into the provided decl.
func InjectImportPaths(importPaths []string, decl *ast.GenDecl) {
	for _, importPath := range importPaths {
		duplicate := false
		for _, importSpec := range decl.Specs {
			if i, ok := importSpec.(*ast.ImportSpec); ok {
				if i.Path.Value == "\""+importPath+"\"" {
					duplicate = true
					break
				}
			}
		}
		if duplicate == false {
			spec := newImportSpec(importPath, "")
			decl.Specs = append(decl.Specs, spec)
		}
	}
}

func isMapStringInterface(t types.Type) bool {
	if t == nil {
		return false
	}
	if m, ok := t.(*types.Map); ok { // find only map declaration
		// look for map[string]interface{}
		if b, okk := m.Key().(*types.Basic); okk && b.Kind() == types.String {
			if _, okkk := m.Elem().(*types.Interface); okkk {
				return true
			}
		}
	}
	return false
}

func injectKvIntoMapStringInterface(kv *ast.KeyValueExpr, elts *ast.CompositeLit) {
	newkeyName := kv.Key.(*ast.BasicLit).Value
	for i, e := range elts.Elts {
		x := e.(*ast.KeyValueExpr).Key.(*ast.BasicLit).Value
		if x == newkeyName {
			elts.Elts[i] = kv
			return
		}
	}
	elts.Elts = append(elts.Elts, kv)
}

// GetProgram creates a new Program of a list of packages.
func GetProgram(pkgs []string) (*loader.Program, error) {

	var conf loader.Config
	_, err := conf.FromArgs(pkgs, false)
	if err != nil {
		return nil, err
	}

	prog, err2 := conf.Load()
	if err2 != nil {
		return nil, err2
	}

	return prog, err
}

// NewPkg creates a new go package.
func NewPkg(fileName string, pkgName string) (*ast.Package, *ast.File) {

	f := &ast.File{}
	f.Name = &ast.Ident{Name: pkgName}

	p := &ast.Package{}
	p.Name = pkgName

	p.Files = map[string]*ast.File{fileName: f}

	return p, f
}

// AddImportDecl creates and add an import statement to the file.
func AddImportDecl(file *ast.File, imports []string) {
	if len(imports) > 0 {
		// Add a GenDecl import statement node to the file tree.
		// literaly: import("pkg"...)
		importGenDecl := NewImportDecl()
		// append the declaration to the file
		file.Decls = append(file.Decls, importGenDecl)

		// inject package path consumed by the new map variable into the file.
		InjectImportPaths(imports, importGenDecl)

		// Make a mutiline import declaration.
		importGenDecl.Lparen = token.Pos(1)
	}
}
