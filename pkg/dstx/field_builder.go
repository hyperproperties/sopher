package dstx

import "github.com/dave/dst"

type FieldBuilder struct {
	names []*dst.Ident
}

func Field(names ...*dst.Ident) *FieldBuilder {
	return &FieldBuilder{
		names: names,
	}
}

func FieldS(names ...string) *FieldBuilder {
	return Field(Idents(names...)...)
}

func (builder *FieldBuilder) Type(t *dst.Ident) *dst.Field {
	return &dst.Field{
		Names: builder.names,
		Type:  t,
	}
}

func (builder *FieldBuilder) TypeS(name string) *dst.Field {
	return builder.Type(Ident(name))
}

func Fields(fields ...*dst.Field) *dst.FieldList {
	return &dst.FieldList{
		List: fields,
	}
}
