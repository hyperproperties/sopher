package language

import "iter"

type HyperAssertionQuantifier[T any] interface {}
type HyperAssertionQualitative[T any] interface {}

type HyperAssertionQuantifierInterpreter[T any] interface {
	Assignments(domain iter.Seq[T]) iter.Seq[[]T]
}

type HyperAssertionQualitativeInterpreter[T any] interface {
	Satisfies(assignments []T) LiftedBoolean
}

type HyperAssertionNode[T any] struct {
	quantifier HyperAssertionQuantifier[T]
	qualitative HyperAssertionQualitative[T]
}











/*type HyperAssertionNode[T any] struct {
	quantifiers []HyperAssertionQuantifier[T]
	qualitative HyperAssertionQualitative[T]
}

func NewHyperAssertionNode[T any](
	quantifiers []HyperAssertionQuantifier[T],
	qualitative HyperAssertionQualitative[T],
) HyperAssertionNode[T] {
	return HyperAssertionNode[T]{
		quantifiers: quantifiers,
		qualitative: qualitative,
	}
}

type HyperAssertionQualitativeInterpreter[T any] struct {
	assignments []T
	stack Stack[LiftedBoolean]
}

func NewHyperAssertionQualitativeInterpreter[T any]() HyperAssertionQualitativeInterpreter[T] {
	return HyperAssertionQualitativeInterpreter[T]{}
}

func (interpreter *HyperAssertionQualitativeInterpreter[T]) Satisfies(qualitative HyperAssertionQualitative[T], assignments []T) LiftedBoolean {
	interpreter.assignments = assignments
	qualitative.Accept(interpreter)
	return interpreter.stack.Pop()
}

func (interpreter *HyperAssertionQualitativeInterpreter[T]) PredicateHyperAssertion(assertion PredicateHyperAssertion[T]) {
	satisfied := assertion.predicate(interpreter.assignments)
	interpreter.stack.Push(LiftBoolean(satisfied))
}

func (interpreter *HyperAssertionQualitativeInterpreter[T]) TrueHyperAssertion(assertion TrueHyperAssertion[T])	 {
	interpreter.stack.Push(LiftedTrue)
}*/