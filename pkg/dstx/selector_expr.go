package dstx

import "github.com/dave/dst"

type SelectorExprBuilder struct {
	identifier *dst.Ident
}

func Select(identifier *dst.Ident) *SelectorExprBuilder {
	return &SelectorExprBuilder{
		identifier: identifier,
	}
}

func SelectS(name string) *SelectorExprBuilder {
	return Select(dst.NewIdent(name))
}

func (builder *SelectorExprBuilder) FromS(name string) *dst.SelectorExpr {
	return builder.From(dst.NewIdent(name))
}

func (builder *SelectorExprBuilder) From(lhs dst.Expr) *dst.SelectorExpr {
	return &dst.SelectorExpr{
		X:   lhs,
		Sel: builder.identifier,
	}
}
