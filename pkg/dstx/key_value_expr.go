package dstx

import "github.com/dave/dst"

type KeyValueBuilder struct {
	Key dst.Expr
}

func Key(key dst.Expr) *KeyValueBuilder {
	return &KeyValueBuilder{
		Key: key,
	}
}

func KeyS(key string) *KeyValueBuilder {
	return Key(Ident(key))
}

func (builder *KeyValueBuilder) Value(value dst.Expr) *dst.KeyValueExpr {
	return &dst.KeyValueExpr{
		Key:   builder.Key,
		Value: value,
	}
}

func (builder *KeyValueBuilder) ValueS(value string) *dst.KeyValueExpr {
	return builder.Value(Ident(value))
}
