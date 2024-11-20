package language

import "iter"

type AGHyperContractTester[T any] struct{}

func NewAGHyperContractTester[T any]() AGHyperContractTester[T] {
	return AGHyperContractTester[T]{}
}

// Q: Should this behaviour be an explorer in and of itself? (See comment at top of increment_explorer.go)
// Returns an iterator which yields the result of each execution satisfying the assumption being tested on the guaratee.
func (filter AGHyperContractTester[T]) Test(inputs iter.Seq[T], call func(input T) (output T), contract AGHyperContract[T], model ...T) iter.Seq2[[]T, LiftedBoolean] {
	return func(yield func([]T, LiftedBoolean) bool) {
		// Create the incremental interpreter with the existing model.
		explorer := NewIncrementalExplorer(model, nil)
		interpreter := NewInterpreter(&explorer)

		for input := range inputs {
			idx := explorer.Increment(input)

			satAssume := interpreter.Satisfies(contract.assumption)

			// The assumption was satisfied and therfore we call the function and update the entry in the model.
			if satAssume.IsTrue() {
				output := call(input)
				// It is okay to update because the call does not change any state used to check the assumptions.
				explorer.Update(idx, output)

				satGuarantee := interpreter.Satisfies(contract.guarantee)

				// If true then we know this subset of the "model"
				// satisfies all assumptions and guarantees.
				// For this reason we reset the incremental interpreter.
				// Such that no elements is marked as added and should be tested.
				if satGuarantee.IsTrue() {
					explorer.Commit()
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
