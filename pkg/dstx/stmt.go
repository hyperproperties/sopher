package dstx

import (
	"github.com/dave/dst"
)

func Statements[S ~[]E, E dst.Stmt](statements S) []dst.Stmt {
	stmts := make([]dst.Stmt, len(statements))
	for idx := range statements {
		stmts[idx] = statements[idx]
	}
	return stmts
}
