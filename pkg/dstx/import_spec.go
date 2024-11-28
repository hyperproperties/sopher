package dstx

import (
	"go/token"
	"strconv"

	"github.com/dave/dst"
)

func ImportS(name, path string) *dst.ImportSpec {
	return Import(Ident(name), BasicString(strconv.Quote(path)))
}

func AnonymousImportS(path string) *dst.ImportSpec {
	return Import(nil, BasicString(strconv.Quote(path)))
}

func Import(name *dst.Ident, path *dst.BasicLit) *dst.ImportSpec {
	return &dst.ImportSpec{
		Name: name,
		Path: path,
	}
}

func DeclareImports(imports ...dst.Spec) *dst.GenDecl {
	return &dst.GenDecl{
		Tok:   token.IMPORT,
		Specs: imports,
	}
}

func FindImportDeclaration(declarations []dst.Decl) *dst.GenDecl {
	for _, declaration := range declarations {
		if cast, ok := declaration.(*dst.GenDecl); ok && cast.Tok == token.IMPORT {
			return cast
		}
	}
	return nil
}

func ImportIntoS(file *dst.File, imports map[string]string) {
	specs := make([]dst.Spec, 0)
	for name, path := range imports {
		importSpec := ImportS(name, path)
		file.Imports = append(file.Imports, importSpec)
		specs = append(specs, importSpec)
	}

	if declaration := FindImportDeclaration(file.Decls); declaration != nil {
		declaration.Specs = append(declaration.Specs, specs...)
	} else {
		declaration := DeclareImports(specs...)
		file.Decls = append([]dst.Decl{declaration}, file.Decls...)
	}
}
