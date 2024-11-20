package dstx

import (
	"go/token"

	"github.com/dave/dst"
)

type DefineStmtBuilder struct {
	lhs []dst.Expr
}

func Define(lhs ...dst.Expr) *DefineStmtBuilder {
	return &DefineStmtBuilder{
		lhs: lhs,
	}
}

func DefineS(lhs ...string) *DefineStmtBuilder {
	return Define(StringsToExprs(lhs...)...)
}

func (builder *DefineStmtBuilder) As(rhs ...dst.Expr) *dst.AssignStmt {
	return &dst.AssignStmt{
		Lhs: builder.lhs,
		Tok: token.DEFINE,
		Rhs: rhs,
	}
}
