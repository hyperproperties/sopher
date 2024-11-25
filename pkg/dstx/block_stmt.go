package dstx

import "github.com/dave/dst"

func Block(stmts ...dst.Stmt) *dst.BlockStmt {
	return &dst.BlockStmt{
		List: stmts,
	}
}

type BlockBuilder struct {
	statements []dst.Stmt
}

func Sequence(statements ...dst.Stmt) *BlockBuilder {
	return &BlockBuilder{
		statements: statements,
	}
}

func (builder *BlockBuilder) Append(statements ...dst.Stmt) *BlockBuilder {
	builder.statements = append(builder.statements, statements...)
	return builder
}

func (builder *BlockBuilder) Terminate(statements ...dst.Stmt) *dst.BlockStmt {
	return Block(append(builder.statements, statements...)...)
}
