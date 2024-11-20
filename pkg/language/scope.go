package language

import "iter"

// The type of quantifier (universal / existential)
type Quantification uint8

const (
	// Represents a universal quantifier.
	UniversalQuantification = Quantification(iota)
	// Represents an existential quantifier.
	ExistentialQuantification
)

// The type of quantifier and the number of declared variables in the quantifier.
type QuantifierScope struct {
	quantification Quantification
	size           int
}

// Creates a new universal quantifier scope with "size" variables.
func NewUniversalQuantifierScope(size int) QuantifierScope {
	return QuantifierScope{
		quantification: UniversalQuantification,
		size:           size,
	}
}

// Creates a new existential quantifier scope with "size" variables.
func NewExistentialQuantifierScope(size int) QuantifierScope {
	return QuantifierScope{
		quantification: ExistentialQuantification,
		size:           size,
	}
}

// Returns the quantification of the quantifier's scope.
func (scope QuantifierScope) Quantification() Quantification {
	return scope.quantification
}

// Returns the number of variables in the quantifier.
func (scope QuantifierScope) Size() int {
	return scope.size
}

// A scope represents, in a fifo structure, the active quantifiers.
type Scope struct {
	quantifiers Stack[QuantifierScope]
}

// Creates a new an empty scope.
func NewScope(quantifiers ...QuantifierScope) Scope {
	return Scope{
		quantifiers: quantifiers,
	}
}

// Returns the number of quantifiers in the scope.
func (scope Scope) Depth() int {
	return len(scope.quantifiers)
}

// Returns the number of declared variables.
func (scope Scope) Size() (size int) {
	for _, quantifier := range scope.quantifiers {
		size += quantifier.size
	}
	return size
}

// Returns the number of declared universally quantified variables.
func (scope Scope) UniversalSize() (size int) {
	for _, quantifier := range scope.quantifiers {
		if quantifier.quantification == UniversalQuantification {
			size += quantifier.size
		}
	}
	return size
}

// Returns an iterator from bottom to top of the quantifiers with the variable offset as key.
func (scope Scope) Quantifiers() iter.Seq2[int, QuantifierScope] {
	return func(yield func(int, QuantifierScope) bool) {
		offset := 0
		for _, quantifier := range scope.quantifiers {
			if !yield(offset, quantifier) {
				return
			}

			offset += quantifier.size
		}
	}
}

// Pushes the quantifier to the scope.
func (scope *Scope) Push(quantifier QuantifierScope) {
	scope.quantifiers.Push(quantifier)
}

// Pops the latest quantifier from the scope.
func (scope *Scope) Pop() QuantifierScope {
	return scope.quantifiers.Pop()
}
