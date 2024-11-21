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

func (builder *DefineStmtBuilder) AsS(rhs ...string) *dst.AssignStmt {
	return builder.As(StringsToExprs(rhs...)...)
}

func (builder *DefineStmtBuilder) As(rhs ...dst.Expr) *dst.AssignStmt {
	return &dst.AssignStmt{
		Lhs: builder.lhs,
		Tok: token.DEFINE,
		Rhs: rhs,
	}
}

type AssignStmtBuilder struct {
	lhs []dst.Expr
}

func Assign(lhs ...dst.Expr) *AssignStmtBuilder {
	return &AssignStmtBuilder{
		lhs: lhs,
	}
}

func AssignS(lhs ...string) *AssignStmtBuilder {
	return Assign(StringsToExprs(lhs...)...)
}

func (builder *AssignStmtBuilder) ToS(rhs ...string) *dst.AssignStmt {
	return builder.To(StringsToExprs(rhs...)...)
}

func (builder *AssignStmtBuilder) To(rhs ...dst.Expr) *dst.AssignStmt {
	return &dst.AssignStmt{
		Lhs: builder.lhs,
		Tok: token.ASSIGN,
		Rhs: rhs,
	}
}
