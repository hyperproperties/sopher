package dstx

import "github.com/dave/dst"

type FileBuilder struct {
	pkg          *dst.Ident
	imports      []*dst.ImportSpec
	declarations []dst.Decl
}

func File(pkg *dst.Ident) *FileBuilder {
	return &FileBuilder{
		pkg: pkg,
		declarations: []dst.Decl{
			DeclareImports(),
		},
	}
}

func FileS(pkg string) *FileBuilder {
	return File(Ident(pkg))
}

func (builder *FileBuilder) Import(specs ...*dst.ImportSpec) *FileBuilder {
	if declaration := FindImportDeclaration(builder.declarations); declaration != nil {
		for _, spec := range specs {
			declaration.Specs = append(declaration.Specs, spec)
		}
		builder.imports = append(builder.imports, specs...)
	}
	return builder
}

func (builder *FileBuilder) Declare(declaration dst.Decl) *FileBuilder {
	builder.declarations = append(builder.declarations, declaration)
	return builder
}

func (builder *FileBuilder) EOF() *dst.File {
	return &dst.File{
		Name:    builder.pkg,
		Decls:   builder.declarations,
		Imports: builder.imports,
	}
}
