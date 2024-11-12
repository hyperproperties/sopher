package language

type AGContract[T any] struct {
	assumptions, guarantees []IncrementalMonitor[T]
	executions              []T
}

func NewAGContract[T any](
	assumptions, guarantees []IncrementalMonitor[T],
) AGContract[T] {
	return AGContract[T]{
		assumptions: assumptions,
		guarantees:  guarantees,
	}
}

func (contract *AGContract[T]) Assume(executions ...T) LiftedBoolean {
	// Assume everything to be fine initially.
	result := LiftedTrue

	// We must check these executions and all other executions.
	combined := append(contract.executions, executions...)

	// For every assumption we must _NOT_ have a counter example.
	for _, assumption := range contract.assumptions {
		assigments := make([]T, assumption.Size())
		interpretation := assumption.Increment(assigments, combined, len(executions))

		// Combine the result and check if we have found a counter example.
		result = result.And(interpretation)
		if result.IsFalse() {
			break
		}
	}

	// For a lifted boolean it can either be "True" or "Unknown".
	return result
}

func (contract *AGContract[T]) Guarantee(executions ...T) LiftedBoolean {
	// Assume everything to be fine initially.
	result := LiftedTrue

	// The executions have return values so we add them to the set.
	contract.executions = append(contract.executions, executions...)

	for _, guarantee := range contract.guarantees {
		assignments := make([]T, guarantee.Size())
		interpretation := guarantee.Increment(assignments, contract.executions, len(executions))

		// Combine the result and check if we have found a counter example.
		result = result.And(interpretation)
		if result.IsFalse() {
			break
		}
	}

	// For a lifted boolean it can either be "True" or "Unknown".
	return result
}
