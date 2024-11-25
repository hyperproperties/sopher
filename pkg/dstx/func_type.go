package dstx

import "github.com/dave/dst"

type FunctionTypeBuilder struct {
	parameters *dst.FieldList
}

func Taking(parameters *dst.FieldList) *FunctionTypeBuilder {
	return &FunctionTypeBuilder{
		parameters: parameters,
	}
}

func TakingN(paramters ...*dst.Field) *FunctionTypeBuilder {
	return Taking(Fields(paramters...))
}

func (builder *FunctionTypeBuilder) Results(results *dst.FieldList) *dst.FuncType {
	return &dst.FuncType{
		Params:  builder.parameters,
		Results: results,
	}
}

func (builder *FunctionTypeBuilder) ResultsN(results ...*dst.Field) *dst.FuncType {
	return builder.Results(Fields(results...))
}
