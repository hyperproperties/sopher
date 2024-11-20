package dstx

import "github.com/dave/dst"

type CallExprBuilder struct {
	function dst.Expr
}

func Call(function dst.Expr) *CallExprBuilder {
	return &CallExprBuilder{
		function: function,
	}
}

func (builder *CallExprBuilder) PassS(arguments ...string) *dst.CallExpr {
	return builder.Pass(StringsToExprs(arguments...)...)
}

func (builder *CallExprBuilder) Pass(arguments ...dst.Expr) *dst.CallExpr {
	return &dst.CallExpr{
		Fun:  builder.function,
		Args: arguments,
	}
}
