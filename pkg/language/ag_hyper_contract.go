package language

import "github.com/hyperproperties/sopher/pkg/quick"

type AGHyperContract[T any] struct {
	assumptions []HyperAssertion[T]
	guarantees  []HyperAssertion[T]
	// A set known to pass both the assumptions and guarantees.
	model []T
}

func NewAGHyperContract[T any](
	assumptions, guarantees []HyperAssertion[T],
) AGHyperContract[T] {
	return AGHyperContract[T]{
		assumptions: assumptions,
		guarantees:  guarantees,
	}
}

func (contract *AGHyperContract[T]) Model(call func(input T) T) {
	if len(contract.model) > 0 {
		return
	}

	model := make([]T, 10)

	for {
		// Find a model which satisfies the assumptions.
		for {
			for idx := range model {
				model[idx] = quick.New[T]()
			}
	
			if contract.Assume(model...).IsTrue() {
				break
			}
		}

		// The model satisfies the assumptions and
		// should then satisfy the guarantees.
		for idx, execution := range model {
			model[idx] = call(execution)
		}

		if contract.Guarantee(model...).IsTrue() {
			break
		} else {
			panic("model does not satisfy guarantee")
		}
	}

	contract.model = model
}

func (contract *AGHyperContract[T]) Assume(executions ...T) LiftedBoolean {
	interpreter := NewHyperAssertionInterpreter[T]()
	for _, assumption := range contract.assumptions {
		satisfied := interpreter.Satisfies(assumption, append(contract.model, executions...))
		if !satisfied {
			return LiftedFalse
		}
	}
	return LiftedTrue
}

func (contract *AGHyperContract[T]) Guarantee(executions ...T) LiftedBoolean {
	interpreter := NewHyperAssertionInterpreter[T]()
	for _, guarantee := range contract.guarantees {
		satisfied := interpreter.Satisfies(guarantee, append(contract.model, executions...))
		if !satisfied {
			return LiftedFalse
		}
	}
	return LiftedTrue
}
