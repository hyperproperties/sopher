package dstx

import "github.com/dave/dst"

func Function(t *dst.FuncType, body *dst.BlockStmt) *dst.FuncLit {
	return &dst.FuncLit{
		Type: t,
		Body: body,
	}
}

func FunctionOf(function *dst.FuncDecl) *dst.FuncLit {
	var TypeParams *dst.FieldList = nil
	if function.Type.TypeParams != nil {
		TypeParams = dst.Clone(function.Type.TypeParams).(*dst.FieldList)
	}

	var Params *dst.FieldList = nil
	if function.Type.Params != nil {
		Params = dst.Clone(function.Type.Params).(*dst.FieldList)
	}

	var Results *dst.FieldList = nil
	if function.Type.Results != nil {
		Results = dst.Clone(function.Type.Results).(*dst.FieldList)
	}

	return &dst.FuncLit{
		Type: &dst.FuncType{
			TypeParams: TypeParams,
			Params:     Params,
			Results:    Results,
		},
		Body: dst.Clone(function.Body).(*dst.BlockStmt),
	}
}
