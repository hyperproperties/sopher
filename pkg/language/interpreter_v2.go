package language

// --- FIXME ---
// The main idea behind nested quantifiers which are not grouped
// is to recursively interpret sub-trees and build upoin existing
// assignments from the previous sub-trees such that exploration
// remembers the assignments to arrive at the tree but can expand
// upon them for the sub-tree quantifiers.

var (
	_ HyperAssertionVisitor[any]             = (*InterpreterV2[any])(nil)
	_ HyperAssertionQuantitativeVisitor[any] = (*InterpreterV2[any])(nil)
	_ HyperAssertionQualitativeVisitor[any]  = (*InterpreterV2[any])(nil)
)

type InterpreterV2[T any] struct {
	scopes      Stack[Scope]
	stack       Stack[LiftedBoolean]
	assignments Assignments[T]
	explorer    Explorer[T]
}

func NewInterpreterV2[T any](explorer Explorer[T]) InterpreterV2[T] {
	return InterpreterV2[T]{
		scopes:   make(Stack[Scope], 1),
		explorer: explorer,
	}
}

func (interpreter *InterpreterV2[T]) Satisfies(assertion HyperAssertion[T]) LiftedBoolean {
	interpreter.HyperAssertion(assertion)
	return LiftedFalse
}

func (interpreter *InterpreterV2[T]) HyperAssertion(assertion HyperAssertion[T]) {
	switch cast := assertion.(type) {
	case HyperAssertionQuantitative[T]:
		cast.Quantitatively(interpreter)
	case HyperAssertionQualitative[T]:
		scope := interpreter.scopes.Peek()

		// If there are any quantifiers in the current scope then we need to start an
		// exploration of assignments based on the scope's sequence of quantifiers.
		if len(scope.quantifiers) > 0 {
			// When recursively exploring we create a new scope such that proceeding
			// quantifiers will add variables to that scope and not the current one.
			interpreter.scopes.Push(NewScope())

			// Recursively explore whether given the explore the assertion is satisfiable.
			satisfied := interpreter.explorer.Explore(scope, interpreter, assertion)
			interpreter.stack.Push(satisfied)

			interpreter.scopes.Pop()
		} else {
			cast.Qualitatively(interpreter)
		}
	default:
		panic("unknown hyper assertion")
	}
}

func (interpreter *InterpreterV2[T]) UniversalHyperAssertion(assertion UniversalHyperAssertion[T]) {
	top := interpreter.scopes.Top()
	interpreter.scopes[top].Push(NewUniversalQuantifierScope(assertion.size))
	interpreter.HyperAssertion(assertion.body)
	interpreter.scopes[top].Pop()
}

func (interpreter *InterpreterV2[T]) ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T]) {
	top := interpreter.scopes.Top()
	interpreter.scopes[top].Push(NewExistentialQuantifierScope(assertion.size))
	interpreter.HyperAssertion(assertion.body)
	interpreter.scopes[top].Pop()
}

func (interpreter *InterpreterV2[T]) UnaryHyperAssertion(assertion UnaryHyperAssertion[T]) {
	operand := interpreter.Satisfies(assertion.operand)

	switch assertion.operator {
	case LogicalNegation:
		interpreter.stack.Push(operand.Not())
	default:
		panic("unknown unary operator")
	}
}

func (interpreter *InterpreterV2[T]) BinaryHyperAssertion(assertion BinaryHyperAssertion[T]) {
	lhs := interpreter.Satisfies(assertion.lhs)

	switch assertion.operator {
	case LogicalConjunction:
		if lhs.IsFalse() {
			interpreter.stack.Push(lhs)
		} else {
			rhs := interpreter.Satisfies(assertion.rhs)
			interpreter.stack.Push(lhs.And(rhs))
		}
	case LogicalDisjunction:
		if lhs.IsTrue() {
			interpreter.stack.Push(lhs)
		} else {
			rhs := interpreter.Satisfies(assertion.rhs)
			interpreter.stack.Push(lhs.Or(rhs))
		}
	case LogicalBiimplication:
		rhs := interpreter.Satisfies(assertion.rhs)
		interpreter.stack.Push(LiftBoolean(lhs == rhs))
	case LogicalImplication:
		if lhs.IsTrue() {
			rhs := interpreter.Satisfies(assertion.rhs)
			interpreter.stack.Push(rhs)
		} else {
			interpreter.stack.Push(LiftedTrue)
		}
	default:
		panic("unknown binary operator")
	}
}

func (interpreter *InterpreterV2[T]) PredicateHyperAssertion(assertion PredicateHyperAssertion[T]) {
	panic("not supported yet")
}

func (interpreter *InterpreterV2[T]) TrueHyperAssertion(assertion TrueHyperAssertion[T]) {
	interpreter.stack.Push(LiftedTrue)
}

func (interpreter *InterpreterV2[T]) AllAssertion(assertions AllAssertion[T]) {
	for _, assertion := range assertions {
		if satisfied := interpreter.Satisfies(assertion); satisfied.IsFalse() {
			interpreter.stack.Push(LiftedFalse)
			return
		}
	}

	interpreter.stack.Push(LiftedTrue)
}

func (interpreter *InterpreterV2[T]) AnyAssertion(assertions AnyAssertion[T]) {
	for _, assertion := range assertions {
		if satisfied := interpreter.Satisfies(assertion); satisfied.IsTrue() {
			interpreter.stack.Push(LiftedTrue)
			return
		}
	}

	interpreter.stack.Push(LiftedFalse)
}
