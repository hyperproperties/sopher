package dstx

import (
	"github.com/dave/dst"
)

func Expressions[S ~[]E, E dst.Expr](expressions S) []dst.Expr {
	exprs := make([]dst.Expr, len(expressions))
	for idx := range expressions {
		exprs[idx] = expressions[idx]
	}
	return exprs
}
