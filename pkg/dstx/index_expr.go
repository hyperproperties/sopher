package dstx

import "github.com/dave/dst"

type IndexExprBuilder struct {
	index dst.Expr
}

func Index(index dst.Expr) *IndexExprBuilder {
	return &IndexExprBuilder{
		index: index,
	}
}

func IndexS(index string) *IndexExprBuilder {
	return Index(Ident(index))
}

func (builder *IndexExprBuilder) Of(of dst.Expr) *dst.IndexExpr {
	return &dst.IndexExpr{
		X:     of,
		Index: builder.index,
	}
}

func (builder *IndexExprBuilder) OfS(of string) *dst.IndexExpr {
	return builder.Of(Ident(of))
}
