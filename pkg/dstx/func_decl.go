package dstx

import "github.com/dave/dst"

type FuncDelcBuilder struct {
	name     *dst.Ident
	receiver *dst.FieldList
	body     *dst.BlockStmt
}

func DeclareFunction(name string) *FuncDelcBuilder {
	return &FuncDelcBuilder{
		name: Ident(name),
	}
}

func (builder *FuncDelcBuilder) For(receiver *dst.FieldList) *FuncDelcBuilder {
	builder.receiver = receiver
	return builder
}

func (builder *FuncDelcBuilder) With(body *dst.BlockStmt) *FuncDelcBuilder {
	builder.body = body
	return builder
}

func (builder *FuncDelcBuilder) As(t *dst.FuncType) *dst.FuncDecl {
	return &dst.FuncDecl{
		Recv: builder.receiver,
		Name: builder.name,
		Type: t,
		Body: builder.body,
	}
}

func HasNamedOutputs(function *dst.FuncDecl) bool {
	for _, output := range function.Type.Results.List {
		if len(output.Names) > 0 {
			return true
		}
	}
	return false
}

func HasReceiver(function *dst.FuncDecl) bool {
	return function.Recv != nil || len(function.Recv.List) > 0
}
