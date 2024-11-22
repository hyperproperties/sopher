package dstx

import "github.com/dave/dst"

func HasNamedOutputs(function *dst.FuncDecl) bool {
	for _, output := range function.Type.Results.List {
		if len(output.Names) > 0 {
			return true
		}
	}
	return false
}
