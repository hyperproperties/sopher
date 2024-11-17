package language

type AGHyperContract[T any] struct {
	assumptions []HyperAssertion[T]
	call        func(input T) (output T)
	guarantees  []HyperAssertion[T]
}

func NewAGHyperContract[T any](
	assumptions []HyperAssertion[T],
	call func(input T) (output T),
	guarantees []HyperAssertion[T],
) AGHyperContract[T] {
	return AGHyperContract[T]{
		assumptions: assumptions,
		call:        call,
		guarantees:  guarantees,
	}
}

func (contract *AGHyperContract[T]) Assume(executions ...T) LiftedBoolean {
	/*interpreter := NewHyperAssertionInterpreter[T]()
	for _, assumption := range contract.assumptions {
		satisfied := interpreter.Satisfies(assumption, executions)
		if !satisfied {
			return LiftedFalse
		}
	}
	return LiftedTrue*/
	panic("not implemented yet")
}

func (contract *AGHyperContract[T]) Call(execution T) T {
	return contract.call(execution)
}

func (contract *AGHyperContract[T]) Guarantee(executions ...T) LiftedBoolean {
	/*interpreter := NewHyperAssertionInterpreter[T]()
	for _, guarantee := range contract.guarantees {
		satisfied := interpreter.Satisfies(guarantee, executions)
		if !satisfied {
			return LiftedFalse
		}
	}
	return LiftedTrue*/
	panic("not implemented yet")
}
