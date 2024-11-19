package language

import "iter"

type Quantification uint8

const (
	UniversalQuantification = Quantification(iota)
	ExistentialQuantification
)

type Quantifier struct {
	quantification Quantification
	size           int
}

func NewUniversalQuantifier(size int) Quantifier {
	return Quantifier{
		quantification: UniversalQuantification,
		size:           size,
	}
}

func NewExistentialQuantifier(size int) Quantifier {
	return Quantifier{
		quantification: ExistentialQuantification,
		size:           size,
	}
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
		size += quantifier.size
	}
	return size
}

func (scope Scope) UniversalSize() (size int) {
	for _, quantifier := range scope.quantifiers {
		if quantifier.quantification == UniversalQuantification {
			size += quantifier.size
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

			offset += quantifier.size
		}
	}
}

func (scope *Scope) Push(quantifier Quantifier) {
	scope.quantifiers.Push(quantifier)
}

func (scope *Scope) Pop() Quantifier {
	return scope.quantifiers.Pop()
}
