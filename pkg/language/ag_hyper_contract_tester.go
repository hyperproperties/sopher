package language

import "iter"

type AGHyperContractTest[T any] struct{}

func NewAGHyperContractTest[T any]() AGHyperContractTest[T] {
	return AGHyperContractTest[T]{}
}

// Returns an iterator which yields the result of each execution satisfying the assumption being tested on the guaratee.
func (filter AGHyperContractTest[T]) Test(inputs iter.Seq[T], call func(input T) (output T), contract AGHyperContract[T], model ...T) {
	// Create the incremental interpreter with the existing model.
	domain := NewIncrementalDomain(model, nil)
	explorer := NewIncrementalExplorer(&domain)
	interpreter := NewInterpreter(&explorer)

	for input := range inputs {
		idx := domain.Increment(input)

		satAssume := interpreter.Satisfies(contract.assumption)

		// The assumption was satisfied and therfore we call the function and update the entry in the model.
		if satAssume.IsTrue() {
			output := call(input)
			// It is okay to update because the call does not change any state used to check the assumptions.
			domain.Update(idx, output)

			satGuarantee := interpreter.Satisfies(contract.guarantee)

			// If true then we know this subset of the "model"
			// satisfies all assumptions and guarantees.
			// For this reason we reset the incremental interpreter.
			// Such that no elements is marked as added and should be tested.
			if satGuarantee.IsTrue() {
				domain.Commit()
			}
		} else {
			domain.Decrement(1)
		}
	}
}
