package language

import "iter"

type IncrementalTester[T any] struct{}

func NewIncrementalTester[T any]() IncrementalTester[T] {
	return IncrementalTester[T]{}
}

func (tester IncrementalTester[T]) Test(
	inputs iter.Seq[T], call func(input T) (output T), contract AGHyperContract[T], model ...T,
) iter.Seq2[[]T, LiftedBoolean] {
	return func(yield func([]T, LiftedBoolean) bool) {
		domain := NewIncrementalDomain(model, []T{})
		explorer := NewIncrementalExplorer(&domain)
		interpreter := NewInterpreter(&explorer)

		for input := range inputs {
			idx := domain.Increment(input)

			assumption := interpreter.Satisfies(contract.assumption)

			// The assumption was satisfied and therfore we call the function and update the entry in the model.
			if assumption.IsTrue() {
				output := call(input)

				// It is okay to update because the call does not change any state used to check the assumptions.
				domain.Update(idx, output)

				guarantee := interpreter.Satisfies(contract.guarantee)

				if guarantee.IsFalse() {
					if !yield(append(domain.set, domain.increment...), guarantee) {
						return
					}
				} else {
					if !yield(domain.set, guarantee) {
						return
					}
				}
			} else {
				domain.Decrement(1)
			}
		}
	}
}
