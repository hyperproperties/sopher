package dstx

import (
	"go/token"

	"github.com/dave/dst"
)

type DeclareVarBuilder struct {
	names []*dst.Ident
	t     dst.Expr
}

func DeclareVariable(names ...*dst.Ident) *DeclareVarBuilder {
	return &DeclareVarBuilder{
		names: names,
	}
}

func DeclareVariableS(names ...string) *DeclareVarBuilder {
	return DeclareVariable(Idents(names...)...)
}

func (builder *DeclareVarBuilder) Type(t dst.Expr) *DeclareVarBuilder {
	builder.t = t
	return builder
}

func (builder *DeclareVarBuilder) TypeS(t string) *DeclareVarBuilder {
	return builder.Type(Ident(t))
}

func (builder *DeclareVarBuilder) Values(exprs ...dst.Expr) *dst.GenDecl {
	return &dst.GenDecl{
		Tok: token.VAR,
		Specs: []dst.Spec{
			&dst.ValueSpec{
				Names:  builder.names,
				Type:   builder.t,
				Values: exprs,
			},
		},
	}
}

func (builder *DeclareVarBuilder) ValuesS(exprs ...string) *dst.GenDecl {
	return builder.Values(Expressions(Idents(exprs...))...)
}

func DeclareType(name *dst.Ident, t dst.Expr) *dst.GenDecl {
	return &dst.GenDecl{
		Tok: token.TYPE,
		Specs: []dst.Spec{
			&dst.TypeSpec{
				Name: name,
				Type: t,
			},
		},
	}
}

func DeclareTypeS(name string, t dst.Expr) *dst.GenDecl {
	return DeclareType(Ident(name), t)
}

func DeclareStructType(name *dst.Ident, fields *dst.FieldList) *dst.GenDecl {
	return DeclareType(name, StructType(fields))
}

func DeclareStructTypeS(name string, fields *dst.FieldList) *dst.GenDecl {
	return DeclareStructType(Ident(name), fields)
}

func DeclareStructTypeN(name *dst.Ident, fields ...*dst.Field) *dst.GenDecl {
	return DeclareStructType(name, Fields(fields...))
}

func DeclareStructTypeSN(name string, fields ...*dst.Field) *dst.GenDecl {
	return DeclareStructTypeN(Ident(name), fields...)
}
