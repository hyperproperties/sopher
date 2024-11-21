package dstx

import "github.com/dave/dst"

func Return(expressions ...dst.Expr) *dst.ReturnStmt {
	return &dst.ReturnStmt{
		Results: expressions,
	}
}
