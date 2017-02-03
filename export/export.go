package export

import (
	"go/ast"
)

// Export exports symbolic and public idents information of targets.
func Export(targets Targets, outfilename, outpackage, outvarname string) (*ast.File, error) {

	// gather all targeted packages
	targetPackages := targets.GetPackagePaths()

	// is it already processed ?
	if f := getCached(targetPackages); f != nil {
		// yup.
		return f, nil
	}

	// make a program of them
	prog, err := GetProgram(targetPackages)
	if err != nil {
		return nil, err
	}

	// create a new file of a package.
	_, destFile := NewPkg(outfilename, outpackage)

	// generate the symbolic expression of the funcmap as a declaration
	// as a var xx map[string]interface{} = map[string]interface{}{...}
	mapVar, imported, err := Symbolic(targets, outvarname, prog, destFile)
	if err != nil {
		return nil, err
	}

	publicIdents, err := PublicIdents(targets, outvarname+"Public", prog, destFile)
	if err != nil {
		return nil, err
	}

	// create and inject the import statement
	AddImportDecl(destFile, imported)

	// add the new var to the file.
	destFile.Decls = append(destFile.Decls, mapVar)
	destFile.Decls = append(destFile.Decls, publicIdents)

	storeCached(targetPackages, destFile)

	return destFile, nil
}

// EnableCache set cache status
var EnableCache = true

var cached = map[string]*ast.File{}

func getCacheKey(targetPackages []string) string {
	n := ""
	for _, t := range targetPackages {
		n += t + "\n"
	}
	return n
}
func getCached(targetPackages []string) *ast.File {
	if !EnableCache {
		return nil
	}
	key := getCacheKey(targetPackages)
	if f, ok := cached[key]; ok {
		// always return a copy.
		return stringToAst(astNodeToString(f))
	}
	return nil
}

func storeCached(targetPackages []string, f *ast.File) {
	if EnableCache {
		key := getCacheKey(targetPackages)
		cached[key] = f
	}
}
