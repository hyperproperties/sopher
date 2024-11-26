package dstx

import (
	"fmt"
	"go/token"

	"github.com/dave/dst"
)

func BasicString(str string) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.STRING,
		Value: str,
	}
}

func BasicInt(integer int) *dst.BasicLit {
	return &dst.BasicLit{
		Kind:  token.INT,
		Value: fmt.Sprintf("%v", integer),
	}
}
