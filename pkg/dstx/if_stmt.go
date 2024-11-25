package dstx

import (
	"github.com/dave/dst"
)

type IfStmtBuilder struct {
	init        dst.Stmt
	condition   dst.Expr
	consequence *dst.BlockStmt
	alternative dst.Stmt
}

func If(condition dst.Expr) *IfStmtBuilder {
	return &IfStmtBuilder{
		condition: condition,
	}
}

func (builder *IfStmtBuilder) Initalise(init dst.Stmt) *IfStmtBuilder {
	builder.init = init
	return builder
}

func (builder *IfStmtBuilder) Then(consequence *dst.BlockStmt) *IfStmtBuilder {
	builder.consequence = consequence
	return builder
}

func (builder *IfStmtBuilder) ThenN(stmts ...dst.Stmt) *IfStmtBuilder {
	return builder.Then(Block(stmts...))
}

func (builder *IfStmtBuilder) Else(alternative dst.Stmt) *IfStmtBuilder {
	builder.alternative = alternative
	return builder
}

func (builder *IfStmtBuilder) End() *dst.IfStmt {
	return &dst.IfStmt{
		Init: builder.init,
		Cond: builder.condition,
		Body: builder.consequence,
		Else: builder.alternative,
	}
}
