package dstx

import "github.com/dave/dst"

func Star(expr dst.Expr) *dst.StarExpr {
	return &dst.StarExpr{
		X: expr,
	}
}
