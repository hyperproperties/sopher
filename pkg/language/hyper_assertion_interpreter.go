package language

import "iter"

type HyperAssertionInterpreter[T any] struct {
	pull        func() (T, bool)
	assignments []T
	stack       Stack[LiftedBoolean]
}

func NewHyperAssertionInterpreter[T any]() HyperAssertionInterpreter[T] {
	return HyperAssertionInterpreter[T]{}
}

func (interpreter *HyperAssertionInterpreter[T]) Satisfies(assertion HyperAssertion[T]) LiftedBoolean {
	assertion.Accept(interpreter)
	return interpreter.stack.Pop()
}

func (interpreter *HyperAssertionInterpreter[T]) Model() iter.Seq[T] {
	// TODO: Seperate the model and assignments.
	return func(yield func(T) bool) {
		for idx := range interpreter.assignments {
			if !yield(interpreter.assignments[idx]) {
				return
			}
		}
	}
}

func (interpreter *HyperAssertionInterpreter[T]) Instantiate(offset, size int) {
	for idx := 0; idx < size; idx++ {
		element, _ := interpreter.pull()
		interpreter.assignments[offset+idx] = element
	}
}

// TODO: Make to a strategy such that we can either re-assign or in-/de-crease the model.
func (interpreter *HyperAssertionInterpreter[T]) FindInstantiation(body HyperAssertion[T], limit, offset, size int) iter.Seq[LiftedBoolean] {
	return func(yield func(LiftedBoolean) bool) {
		interpreter.Instantiate(offset, size)

		for attempt := 0; attempt < limit; attempt++ {
			body.Accept(interpreter)
			satisfied := interpreter.stack.Peek()

			if !yield(satisfied) {
				return
			}

			interpreter.Instantiate(offset + (attempt % size), 1)
		}
	}
}

func (interpreter *HyperAssertionInterpreter[T]) UniversalHyperAssertion(assertion UniversalHyperAssertion[T]) {
	// TODO: Include scopes for quantifiers instead of this.
	oldLength := len(interpreter.assignments)
	interpreter.assignments = append(interpreter.assignments, make([]T, assertion.size)...)

	for satisfied := range interpreter.FindInstantiation(
		assertion.body, 1000, assertion.offset, assertion.size,
		) {
		if satisfied.IsFalse() {
			return
		}
	}

	interpreter.assignments = interpreter.assignments[:oldLength]
}

func (interpreter *HyperAssertionInterpreter[T]) ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T]) {
	// TODO: Include scopes for quantifiers instead of this.
	oldLength := len(interpreter.assignments)
	interpreter.assignments = append(interpreter.assignments, make([]T, assertion.size)...)

	for satisfied := range interpreter.FindInstantiation(
		assertion.body, 1000, assertion.offset, assertion.size,
	) {
		if satisfied.IsTrue() {
			return
		}
	}

	interpreter.assignments = interpreter.assignments[:oldLength]
}

func (interpreter *HyperAssertionInterpreter[T]) BinaryHyperAssertion(assertion BinaryHyperAssertion[T]) {
	// TODO: Short circuit.
	assertion.lhs.Accept(interpreter)
	lhs := interpreter.stack.Pop()
	assertion.rhs.Accept(interpreter)
	rhs := interpreter.stack.Pop()

	switch assertion.operator {
	case LogicalConjunction:
		interpreter.stack.Push(lhs.And(rhs))
	case LogicalDisjunction:
		interpreter.stack.Push(lhs.Or(rhs))
	case LogicalBiimplication:
		interpreter.stack.Push(LiftBoolean(lhs == rhs))
	case LogicalImplication:
		interpreter.stack.Push(lhs.Not().Or(rhs))
	default:
		panic("unknown binary operator")
	}
}

func (interpreter *HyperAssertionInterpreter[T]) PredicateHyperAssertion(assertion PredicateHyperAssertion[T]) {
	var explorer Explorer[T]
	for verdict := range explorer.Explore(assertion) {
		if !verdict.IsUnknown() {
			interpreter.stack.Push(verdict)
		}
	}
}

func (interpreter *HyperAssertionInterpreter[T]) TrueHyperAssertion(assertion TrueHyperAssertion[T]) {
	interpreter.stack.Push(LiftedTrue)
}
