package export

import (
	"go/ast"
)

// Export exports symbolic and public idents information of targets.
func Export(targets Targets, outfilename, outpackage, outvarname string) (*ast.File, error) {

	// gather all targeted packages
	targetPackages := targets.GetPackagePaths()
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

	return destFile, nil
}
