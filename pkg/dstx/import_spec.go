package dstx

import (
	"go/token"
	"strconv"

	"github.com/dave/dst"
)

func ImportS(name, path string) *dst.ImportSpec {
	return Import(Ident(name), BasicString(strconv.Quote(path)))
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
