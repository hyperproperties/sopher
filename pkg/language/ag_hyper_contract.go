package language

import "github.com/hyperproperties/sopher/pkg/quick"

type AGHyperContract[T any] struct {
	assumption HyperAssertion[T]
	call       func(input T) (output T)
	guarantee  HyperAssertion[T]
	explorer   IncrementalExplorer[T]
}

func NewAGHyperContract[T any](
	assumptions HyperAssertion[T],
	call func(input T) (output T),
	guarantees HyperAssertion[T],
) AGHyperContract[T] {
	return AGHyperContract[T]{
		assumption: assumptions,
		call:       call,
		guarantee:  guarantees,
		explorer:   NewIncrementalExplorer[T](nil, nil),
	}
}

func (contract *AGHyperContract[T]) Assume(executions ...T) (satisfied LiftedBoolean) {
	satisfied = LiftedTrue

	interpreter := NewInterpreter(&contract.explorer)
	start := contract.explorer.IncrementLength()
	contract.explorer.Increment(executions...)

	// TODO: Add a contract configuration that allows a user defined loop-termination criteria.
	for counter := 0; counter < 1000; counter++ {
		if satisfied = interpreter.Satisfies(contract.assumption); satisfied.IsTrue() {
			break
		}
		contract.explorer.Increment(quick.New[T]())
	}

	end := contract.explorer.IncrementLength()
	contract.explorer.Decrement(end-start)

	return satisfied
}

func (contract *AGHyperContract[T]) Call(input T) (assumption LiftedBoolean, output T, guarantee LiftedBoolean) {
	if assumption = contract.Assume(input); assumption.IsFalse() {
		return
	}

	output = contract.call(input)

	if guarantee = contract.Guarantee(output); guarantee.IsTrue() {
		contract.explorer.Increment(output)
		contract.explorer.Commit()
	}

	return assumption, output, guarantee
}

func (contract *AGHyperContract[T]) Guarantee(executions ...T) (satisfied LiftedBoolean) {
	// FIXME: Implement a test here.
	return LiftedTrue
}
