package dstx

import "github.com/dave/dst"

func NewLineBefore[T dst.Node](node T) T {
	node.Decorations().Before = dst.NewLine
	return node
}

func NewLineAfter[T dst.Node](node T) T {
	node.Decorations().After = dst.NewLine
	return node
}

func NewLineAround[T dst.Node](node T) T {
	NewLineBefore(node)
	NewLineAfter(node)
	return node
}

func AppendStart[T dst.Node](node T, comments ...string) T {
	node.Decorations().Start.Append(comments...)
	return node
}


func PrependStart[T dst.Node](node T, comments ...string) T {
	node.Decorations().Start.Prepend(comments...)
	return node
}

func AppendEnd[T dst.Node](node T, comments ...string) T {
	node.Decorations().End.Append(comments...)
	return node
}

func PrependEnd[T dst.Node](node T, comments ...string) T {
	node.Decorations().End.Prepend(comments...)
	return node
}