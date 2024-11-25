package dstx

import "github.com/dave/dst"

type CompositeLitBuilder struct {
	ident    *dst.Ident
	elements []dst.Expr
}

func Compose(ident *dst.Ident) *CompositeLitBuilder {
	return &CompositeLitBuilder{
		ident: ident,
	}
}

func ComposeS(ident string) *CompositeLitBuilder {
	return Compose(Ident(ident))
}

func (builder *CompositeLitBuilder) Elements(elements ...dst.Expr) *CompositeLitBuilder {
	builder.elements = append(builder.elements, elements...)
	return builder
}

func (builder *CompositeLitBuilder) ElementsS(elements ...string) *CompositeLitBuilder {
	for _, ident := range elements {
		builder.elements = append(builder.elements, Ident(ident))
	}
	return builder
}

func (builder *CompositeLitBuilder) Final() *dst.CompositeLit {
	return &dst.CompositeLit{
		Type: builder.ident,
		Elts: builder.elements,
	}
}
