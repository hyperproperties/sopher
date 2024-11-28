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

func (builder *FunctionTypeBuilder) Taking(parameters *dst.FieldList) *FunctionTypeBuilder {
	builder.parameters.List = append(builder.parameters.List, parameters.List...)
	return builder
}

func (builder *FunctionTypeBuilder) TakingN(paramters ...*dst.Field) *FunctionTypeBuilder {
	builder.Taking(Fields(paramters...))
	return builder
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

func (builder *FunctionTypeBuilder) Void() *dst.FuncType {
	return builder.Results(nil)
}
