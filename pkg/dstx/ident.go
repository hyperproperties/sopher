package dstx

import "github.com/dave/dst"

func Ident(name string) *dst.Ident {
	return dst.NewIdent(name)
}
