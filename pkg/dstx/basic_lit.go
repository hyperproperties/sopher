package dstx

import (
	"go/token"

	"github.com/dave/dst"
)

func BasicString(str string) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.STRING,
		Value: str,
	}
}
