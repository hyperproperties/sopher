package dstx

import "github.com/dave/dst"

func ExprStmt(expression dst.Expr) dst.Stmt {
	return &dst.ExprStmt{
		X: expression,
	}
}
