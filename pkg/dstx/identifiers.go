package dstx

import "github.com/dave/dst"

// Converts a slice of strings to a slice of dst.Expr which is []*dst.Ident.
func StringsToExprs(strings ...string) []dst.Expr {
	identifiers := make([]dst.Expr, len(strings))
	for idx := range strings {
		identifiers[idx] = dst.NewIdent(strings[idx])
	}
	return identifiers
}
