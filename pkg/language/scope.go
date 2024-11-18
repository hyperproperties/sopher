package language

import "iter"

type QuantifierVisitor interface {
	UniversalQuantifier(quantifier UniversalQuantifier)
	ExistentialQuantifier(quantifier ExistentialQuantifier)
}

type Quantifier interface {
	Size() int
	Accept(visitor QuantifierVisitor)
}

type UniversalQuantifier struct {
	size int
}

func NewUniversalQuantifier(size int) UniversalQuantifier {
	return UniversalQuantifier{
		size: size,
	}
}

func (quantifier UniversalQuantifier) Size() int {
	return quantifier.size
}

func (quantifier UniversalQuantifier) Accept(visitor QuantifierVisitor) {
	visitor.UniversalQuantifier(quantifier)
}

type ExistentialQuantifier struct {
	size int
}

func NewExistentialQuantifier(size int) ExistentialQuantifier {
	return ExistentialQuantifier{
		size: size,
	}
}

func (quantifier ExistentialQuantifier) Size() int {
	return quantifier.size
}

func (quantifier ExistentialQuantifier) Accept(visitor QuantifierVisitor) {
	visitor.ExistentialQuantifier(quantifier)
}

type Scope struct {
	quantifiers Stack[Quantifier]
}

func NewScope() Scope {
	return Scope{}
}

func (scope Scope) Depth() int {
	return len(scope.quantifiers)
}

func (scope Scope) Size() (size int) {
	for _, quantifier := range scope.quantifiers {
		size += quantifier.Size()
	}
	return size
}

func (scope Scope) UniversalSize() (size int) {
	for idx := range scope.quantifiers {
		switch cast := scope.quantifiers[idx].(type) {
		case UniversalQuantifier:
			size += cast.Size()
		case ExistentialQuantifier:
		default:
			panic("unknown quantifier")
		}
	}
	return size
}

func (scope Scope) Quantifiers() iter.Seq2[int, Quantifier] {
	return func(yield func(int, Quantifier) bool) {
		offset := 0
		for _, quantifier := range scope.quantifiers {
			if !yield(offset, quantifier) {
				return
			}

			offset += quantifier.Size()
		}
	}
}

func (scope *Scope) Push(quantifier Quantifier) {
	scope.quantifiers.Push(quantifier)
}

func (scope *Scope) Pop() Quantifier {
	return scope.quantifiers.Pop()
}
