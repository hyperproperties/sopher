package dstx

import "github.com/dave/dst"

func Clone[T dst.Node](node T) T {
	return dst.Clone(node).(T)
}

func Clones[S ~[]E, E dst.Node](slice S) S {
	clones := make(S, len(slice))
	for idx, node := range slice {
		clones[idx] = Clone(node)
	}
	return clones
}
