package language

type HyperAssertionInterpreter[T any] struct {
	explorer Explorer[T]
	scope    Scope
	stack    Stack[LiftedBoolean]
}

func NewHyperAssertionInterpreter[T any](explorer Explorer[T]) HyperAssertionInterpreter[T] {
	return HyperAssertionInterpreter[T]{
		explorer: explorer,
	}
}

func (interpreter *HyperAssertionInterpreter[T]) Satisfies(assertion HyperAssertion[T]) LiftedBoolean {
	assertion.Accept(interpreter)
	return interpreter.stack.Pop()
}

func (interpreter *HyperAssertionInterpreter[T]) UniversalHyperAssertion(assertion UniversalHyperAssertion[T]) {
	quantifier := NewUniversalQuantifier(assertion.size)
	interpreter.scope.Push(quantifier)
	assertion.body.Accept(interpreter)
	interpreter.scope.Pop()
}

func (interpreter *HyperAssertionInterpreter[T]) ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T]) {
	quantifier := NewExistentialQuantifier(assertion.size)
	interpreter.scope.Push(quantifier)
	assertion.body.Accept(interpreter)
	interpreter.scope.Pop()
}

func (interpreter *HyperAssertionInterpreter[T]) UnaryHyperAssertion(assertion UnaryHyperAssertion[T]) {
	assertion.operand.Accept(interpreter)
	operand := interpreter.stack.Pop()

	switch assertion.operator {
	case LogicalNegation:
		interpreter.stack.Push(operand.Not())
	default:
		panic("unknown unary operator")
	}
}

func (interpreter *HyperAssertionInterpreter[T]) BinaryHyperAssertion(assertion BinaryHyperAssertion[T]) {
	assertion.lhs.Accept(interpreter)
	lhs := interpreter.stack.Pop()

	switch assertion.operator {
	case LogicalConjunction:
		if lhs.IsFalse() {
			interpreter.stack.Push(lhs)
		} else {
			assertion.rhs.Accept(interpreter)
			rhs := interpreter.stack.Pop()
			interpreter.stack.Push(lhs.And(rhs))
		}
	case LogicalDisjunction:
		if lhs.IsTrue() {
			interpreter.stack.Push(lhs)
		} else {
			assertion.rhs.Accept(interpreter)
			rhs := interpreter.stack.Pop()
			interpreter.stack.Push(lhs.Or(rhs))
		}
	case LogicalBiimplication:
		assertion.rhs.Accept(interpreter)
		rhs := interpreter.stack.Pop()
		interpreter.stack.Push(LiftBoolean(lhs == rhs))
	case LogicalImplication:
		if lhs.IsTrue() {
			assertion.rhs.Accept(interpreter)
			rhs := interpreter.stack.Pop()
			interpreter.stack.Push(rhs)
		} else {
			interpreter.stack.Push(LiftedTrue)
		}
	default:
		panic("unknown binary operator")
	}
}

func (interpreter *HyperAssertionInterpreter[T]) PredicateHyperAssertion(assertion PredicateHyperAssertion[T]) {
	result := interpreter.explorer.Explore(
		interpreter.scope, assertion,
	)
	interpreter.stack.Push(result)
}

func (interpreter *HyperAssertionInterpreter[T]) TrueHyperAssertion(assertion TrueHyperAssertion[T]) {
	interpreter.stack.Push(LiftedTrue)
}
