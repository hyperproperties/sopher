package dstx

import "github.com/dave/dst"

func ArrayType(length dst.Expr, element dst.Expr) *dst.ArrayType {
	return &dst.ArrayType{
		Len: length,
		Elt: element,
	}
}

func ArrayTypeI(length int, element dst.Expr) *dst.ArrayType {
	return ArrayType(BasicInt(length), element)
}

func ArrayTypeIS(length int, element string) *dst.ArrayType {
	return ArrayTypeI(length, Ident(element))
}

func ArrayTypeS(length dst.Expr, element string) *dst.ArrayType {
	return ArrayType(length, BasicString(element))
}

func SliceType(element dst.Expr) *dst.ArrayType {
	return &dst.ArrayType{
		Elt: element,
	}
}

func SliceTypeS(element string) *dst.ArrayType {
	return SliceType(BasicString(element))
}
