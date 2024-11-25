package dstx

import "github.com/dave/dst"

func Ident(name string) *dst.Ident {
	return dst.NewIdent(name)
}

func Idents(names ...string) []*dst.Ident {
	idents := make([]*dst.Ident, len(names))
	for idx := range names {
		idents[idx] = Ident(names[idx])
	}
	return idents
}
