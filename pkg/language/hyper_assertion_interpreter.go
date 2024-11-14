package language

// TODO: Iterative evaluation should be a lot faster.

import "github.com/hyperproperties/sopher/pkg/iterx"

type HyperAssertionInterpreter[T any] struct {
	elements []T
	assignments []T
	assertion HyperAssertion[T]
	satisfied bool
}

func NewHyperAssertionInterpreter[T any]() HyperAssertionInterpreter[T] {
	return HyperAssertionInterpreter[T]{}
}

func (interpreter *HyperAssertionInterpreter[T]) Satisfies(assertion HyperAssertion[T], elements []T) bool {
	interpreter.elements = elements
	interpreter.assignments = make([]T, assertion.Size())
	interpreter.assertion = assertion
	interpreter.assertion.Accept(interpreter)
	return interpreter.satisfied
}

func (interpreter *HyperAssertionInterpreter[T]) UniversalHyperAssertion(assertion UniversalHyperAssertion[T]) {
	permutations := iterx.Permutations(assertion.size, len(interpreter.elements))
	for permutation := range iterx.Map(interpreter.elements, permutations) {
		for idx := 0; idx < assertion.size; idx++ {
			interpreter.assignments[assertion.offset + idx] = permutation[idx]
		}

		assertion.body.Accept(interpreter)
		if !interpreter.satisfied {
			return
		}
	}
}

func (interpreter *HyperAssertionInterpreter[T]) ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T]) {
	permutations := iterx.Permutations(assertion.size, len(interpreter.elements))
	for permutation := range iterx.Map(interpreter.elements, permutations) {
		for idx := 0; idx < assertion.size; idx++ {
			interpreter.assignments[assertion.offset + idx] = permutation[idx]
		}

		assertion.body.Accept(interpreter)
		if interpreter.satisfied {
			return
		}
	}
}

func (interpreter *HyperAssertionInterpreter[T]) PredicateHyperAssertion(assertion PredicateHyperAssertion[T]) {
	interpreter.satisfied = assertion.predicate(interpreter.assignments)
}

func (interpreter *HyperAssertionInterpreter[T]) TrueHyperAssertion(assertion TrueHyperAssertion[T]) {
	interpreter.satisfied = true
}
