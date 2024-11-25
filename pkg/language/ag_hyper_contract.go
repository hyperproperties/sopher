package language

type AGHyperContract[T any] struct {
	assumption HyperAssertion[T]
	guarantee  HyperAssertion[T]
	models     map[uint64][]T
}

func NewAGHyperContract[T any](
	assumptions, guarantees HyperAssertion[T],
) AGHyperContract[T] {
	return AGHyperContract[T]{
		assumption: assumptions,
		guarantee:  guarantees,
		models:     map[uint64][]T{},
	}
}

// Checks if a set of executions from a caller satisfies the assertion given its existing model.
func (contract *AGHyperContract[T]) Satisfies(caller uint64, assertion HyperAssertion[T], executions ...T) LiftedBoolean {
	model := contract.models[caller]
	domain := NewIncrementalDomain(model, executions)
	explorer := NewIncrementalExplorer(&domain)
	interpreter := NewInterpreter(&explorer)
	satisfied := interpreter.Satisfies(assertion)
	return satisfied
}

// Checks if the caller and the executions (as an increment) satisfies the contract's assumption.
func (contract *AGHyperContract[T]) Assume(caller uint64, executions ...T) LiftedBoolean {
	return contract.Satisfies(caller, contract.assumption, executions...)
}

// Calls the function (call) with the input as the caller. If the guarantee is
// not not-satisfied then the output is appended to the model of the caller.
func (contract *AGHyperContract[T]) Call(
	caller uint64, input T, call func(input T) (output T),
) (assumption LiftedBoolean, output T, guarantee LiftedBoolean) {
	assumption, guarantee = LiftedFalse, LiftedFalse

	if assumption = contract.Assume(caller, input); assumption.IsFalse() {
		return
	}

	output = call(input)

	if guarantee = contract.Guarantee(caller, output); guarantee.IsFalse() {
		return
	}

	if model, exists := contract.models[caller]; exists {
		contract.models[caller] = append(model, output)
	} else {
		contract.models[caller] = []T{output}
	}

	return assumption, output, guarantee
}

// Checks if the caller and the executions (as an increment) satisfies the contract's guarantee.
func (contract *AGHyperContract[T]) Guarantee(caller uint64, executions ...T) LiftedBoolean {
	return contract.Satisfies(caller, contract.guarantee, executions...)
}
