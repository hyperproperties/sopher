package language

import "iter"

type AGHyperContractTester[T any] struct{}

func NewAGHyperContractTester[T any]() AGHyperContractTester[T] {
	return AGHyperContractTester[T]{}
}

// Returns an iterator which yields the result of each execution satisfying the assumption being tested on the guaratee.
func (filter AGHyperContractTester[T]) Test(inputs iter.Seq[T], contract AGHyperContract[T], model ...T) iter.Seq2[[]T, LiftedBoolean] {
	return func(yield func([]T, LiftedBoolean) bool) {
		// Create the incremental interpreter with the existing model.
		explorer := NewIncrementalExplorer(model, 0)
		interpreter := NewInterpreter(&explorer)

		for input := range inputs {
			idx := explorer.Increment(input)

			satAssume := LiftedTrue
			for _, assumption := range contract.assumptions {
				satAssume = satAssume.And(
					interpreter.Satisfies(assumption),
				)
			}

			// The assumption was satisfied and therfore we call the function and update the entry in the model.
			if satAssume.IsTrue() {
				output := contract.call(input)
				// It is okay to update because the call does not change any state used to check the assumptions.
				explorer.Update(idx, output)

				satGuarantee := LiftedTrue
				for _, guarantee := range contract.guarantees {
					satGuarantee = satGuarantee.And(
						interpreter.Satisfies(guarantee),
					)
				}

				// If true then we know this subset of the "model"
				// satisfies all assumptions and guarantees.
				// For this reason we reset the incremental interpreter.
				// Such that no elements is marked as added and should be tested.
				if satGuarantee.IsTrue() {
					explorer.Decrement(explorer.Added())
				}

				if !yield(explorer.Model(), satGuarantee) {
					return
				}
			} else {
				explorer.Decrement(1)
			}
		}
	}
}
