package dstx

import "github.com/dave/dst"

func Function(t *dst.FuncType, body *dst.BlockStmt) *dst.FuncLit {
	return &dst.FuncLit{
		Type: t,
		Body: body,
	}
}

func FunctionOf(function *dst.FuncDecl) *dst.FuncLit {
	var typeParams *dst.FieldList = nil
	if function.Type.TypeParams != nil {
		typeParams = dst.Clone(function.Type.TypeParams).(*dst.FieldList)
	}

	var params *dst.FieldList = nil
	if function.Type.Params != nil {
		params = dst.Clone(function.Type.Params).(*dst.FieldList)
	}

	var results *dst.FieldList = nil
	if function.Type.Results != nil {
		results = dst.Clone(function.Type.Results).(*dst.FieldList)
	}

	return &dst.FuncLit{
		Type: &dst.FuncType{
			TypeParams: typeParams,
			Params:     params,
			Results:    results,
		},
		Body: dst.Clone(function.Body).(*dst.BlockStmt),
	}
}
