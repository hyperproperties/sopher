package language

type HyperAssertionIncrementalInterpreter[T any] struct {
	explorer *IncrementalExplorer[T]
	interpreter HyperAssertionInterpreter[T]
}

func NewHyperAssertionIncrementalInterpreter[T any]() HyperAssertionIncrementalInterpreter[T] {
	explorer := NewIncrementalExplorer[T](nil, 0)
	return HyperAssertionIncrementalInterpreter[T]{
		explorer: &explorer,
		interpreter: NewHyperAssertionInterpreter(&explorer),
	}
}

func (incremental *HyperAssertionIncrementalInterpreter[T]) Increment(executions ...T) {
	incremental.explorer.model = append(incremental.explorer.model, executions...)
	incremental.explorer.added += len(executions)
}

func (incremental *HyperAssertionIncrementalInterpreter[T]) Satisfies(assertion HyperAssertion[T]) LiftedBoolean {
	verdict := incremental.interpreter.Satisfies(assertion)
	incremental.explorer.added = 0
	return verdict
}

func (incremental *HyperAssertionIncrementalInterpreter[T]) Decrement(amount int) {
	if amount < incremental.explorer.added {
		incremental.explorer.added -= amount
	} else {
		incremental.explorer.added = 0
	}
	incremental.explorer.model = incremental.explorer.model[:len(incremental.explorer.model)-amount]
}

func (incremental *HyperAssertionIncrementalInterpreter[T]) Reset() {
	incremental.explorer.model = nil
	incremental.explorer.added = 0
}

func (incremental *HyperAssertionIncrementalInterpreter[T]) Model() []T {
	return incremental.explorer.model
}