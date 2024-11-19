package language

var _ HyperAssertionVisitor[any] = (*Interpreter[any])(nil)

type Interpreter[T any] struct {
	explorer Explorer[T]
	scope    Scope
	stack    Stack[LiftedBoolean]
}

func NewInterpreter[T any](explorer Explorer[T]) Interpreter[T] {
	return Interpreter[T]{
		explorer: explorer,
	}
}

func (interpreter *Interpreter[T]) Satisfies(assertion HyperAssertion[T]) LiftedBoolean {
	assertion.Accept(interpreter)
	return interpreter.stack.Pop()
}

func (interpreter *Interpreter[T]) AllAssertion(assertions AllAssertion[T]) {
	for _, assertion := range assertions {
		if satisfied := interpreter.Satisfies(assertion); satisfied.IsFalse() {
			interpreter.stack.Push(LiftedFalse)
			return
		}
	}

	interpreter.stack.Push(LiftedTrue)
}

func (interpreter *Interpreter[T]) AnyAssertion(assertions AnyAssertion[T]) {
	for _, assertion := range assertions {
		if satisfied := interpreter.Satisfies(assertion); satisfied.IsTrue() {
			interpreter.stack.Push(LiftedTrue)
			return
		}
	}

	interpreter.stack.Push(LiftedFalse)
}

func (interpreter *Interpreter[T]) UniversalHyperAssertion(assertion UniversalHyperAssertion[T]) {
	quantifier := NewUniversalQuantifier(assertion.size)
	interpreter.scope.Push(quantifier)
	assertion.body.Accept(interpreter)
	interpreter.scope.Pop()
}

func (interpreter *Interpreter[T]) ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T]) {
	quantifier := NewExistentialQuantifier(assertion.size)
	interpreter.scope.Push(quantifier)
	assertion.body.Accept(interpreter)
	interpreter.scope.Pop()
}

func (interpreter *Interpreter[T]) UnaryHyperAssertion(assertion UnaryHyperAssertion[T]) {
	assertion.operand.Accept(interpreter)
	operand := interpreter.stack.Pop()

	switch assertion.operator {
	case LogicalNegation:
		interpreter.stack.Push(operand.Not())
	default:
		panic("unknown unary operator")
	}
}

func (interpreter *Interpreter[T]) BinaryHyperAssertion(assertion BinaryHyperAssertion[T]) {
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

func (interpreter *Interpreter[T]) PredicateHyperAssertion(assertion PredicateHyperAssertion[T]) {
	// TODO: To support nested quantifiers the assignments must be saved for the nested quantifiers.
	result := interpreter.explorer.Explore(interpreter.scope, assertion)
	interpreter.stack.Push(result)
}

func (interpreter *Interpreter[T]) TrueHyperAssertion(assertion TrueHyperAssertion[T]) {
	interpreter.stack.Push(LiftedTrue)
}
