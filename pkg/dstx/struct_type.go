package dstx

import "github.com/dave/dst"

func StructType(fields *dst.FieldList) *dst.StructType {
	return &dst.StructType{
		Fields: fields,
	}
}

func StructTypeN(fields ...*dst.Field) *dst.StructType {
	return StructType(Fields(fields...))
}
