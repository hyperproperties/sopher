package language

import (
	"slices"

	"github.com/hyperproperties/sopher/pkg/iterx"
)

type IncrementalExplorer[T any] struct {
	model     []T
	increment []T
}

func NewIncrementalExplorer[T any](model []T, increment []T) IncrementalExplorer[T] {
	return IncrementalExplorer[T]{
		model:     model,
		increment: increment,
	}
}

// Pushes the executions to the increment.
func (explorer *IncrementalExplorer[T]) Increment(executions ...T) int {
	if len(executions) == 0 {
		panic("cannot increment with no executions")
	}
	index := len(explorer.model) + len(explorer.increment)
	explorer.increment = append(explorer.increment, executions...)
	return index
}

// Pops the executions from the increment.
func (explorer *IncrementalExplorer[T]) Decrement(amount int) {
	if amount < 0 {
		panic("cannot decrement by a negative number")
	}
	if amount == 0 {
		return
	}

	length := len(explorer.increment)
	explorer.increment = slices.Delete(explorer.increment, length-amount, length)
}

// Finishes the increment and moves the executions to the model.
func (explorer *IncrementalExplorer[T]) Commit() {
	explorer.model = append(explorer.model, explorer.increment...)
	explorer.increment = nil
}

// Finishes the increment by deleting the increment.
func (explorer *IncrementalExplorer[T]) Rollback() {
	explorer.increment = nil
}

// Updates an entry in the model. This does not work for
// elements in the increment which has not been committed yet.
func (explorer *IncrementalExplorer[T]) Update(index int, value T) {
	explorer.model[index] = value
}

// Returns the model of the incremental explorer which is the set of executions
// that should not solely be in an exploration permutation.
func (explorer *IncrementalExplorer[T]) Model() []T {
	return explorer.model
}

// Returns the increment of the incremental explorer which is the set where
// for every tested assertion every assignment must have had atleast one assignment
// to an element in the increment.
func (explorer *IncrementalExplorer[T]) Incremental() []T {
	return explorer.increment
}

// The number of elements in the current increment.
func (explorer *IncrementalExplorer[T]) IncrementLength() int {
	return len(explorer.increment)
}

// Checks if the assertion is satisfied given a fixed model where a prefix of it has already been tested.
// It works by assigning a incremental permutation of the model to the universally quantified variables.
// Then it attempts to recursively find a solution to the existentially quantified variables by
// using permutations (not incrementally) of the mode. If no solution for the existentially quantified
// variables is found then a counter-example for the universal quantifiers and we return that. Otherwise, true.
func (explorer *IncrementalExplorer[T]) Explore(scope Scope, assertion PredicateHyperAssertion[T]) LiftedBoolean {
	// Initialise all assignments for all quantifiers within the scope.
	assignments := make([]T, scope.Size())
	if len(assignments) == 0 {
		satisfied := assertion.predicate(nil)
		return LiftBoolean(satisfied)
	}

	// If there are no increment then we dont test it.  However, we know
	// per the invariant that the model must be a set which satifies the
	// assertion and therefore we return true.
	if explorer.IncrementLength() == 0 {
		return LiftedTrue
	}

	var Exists func(depth, offset int) LiftedBoolean
	Exists = func(depth, offset int) LiftedBoolean {
		// We have reached the depth of the scope so we check the assertion.
		if depth >= scope.Depth() {
			result := assertion.predicate(assignments)
			return LiftBoolean(result)
		}

		// Generate all assignments for the existential quantifier's variables.
		quantifier := scope.quantifiers[depth]
		if quantifier.quantification != ExistentialQuantification {
			return Exists(depth+1, offset+quantifier.size)
		}

		// For every permutation is tried and tested against the assertion.
		// TODO: Convert the ranging over assignments for the quantifier to a strategy.
		for permutation := range iterx.Map(
			append(explorer.model, explorer.increment...), iterx.Permutations(
				quantifier.size, len(explorer.model),
			),
		) {
			// TODO: copy?
			for idx := 0; idx < quantifier.size; idx++ {
				assignments[offset+idx] = permutation[idx]
			}

			// Get the result and check for short circuit.
			result := Exists(depth+1, offset+quantifier.size)
			if result.IsTrue() {
				return result
			}
		}

		// No satisfied exitential example was found.
		return LiftedFalse
	}

	// Generate all assignments for universal variables.
	// TODO: Convert the ranging over assignments for the quantifier to a strategy.
	for permutation := range iterx.Map(
		append(explorer.model, explorer.increment...), iterx.IncrementalPermutations(
			scope.UniversalSize(), len(explorer.model), len(explorer.increment),
		),
	) {
		//For each assignment of universal variables:
		local := 0
		for offset, quantifier := range scope.Quantifiers() {
			if quantifier.quantification == UniversalQuantification {
				// TODO: copy?
				for idx := 0; idx < quantifier.size; idx++ {
					assignments[offset+idx] = permutation[local+idx]
				}
				local += quantifier.size
			}
		}

		/* Check if there exists a satisfying assignment for the existential variables.
		   Use early exit if no such assignment exists. */
		result := Exists(0, 0)
		if result.IsFalse() {
			return result
		}
	}

	// If all universal assignments are satisfied, return true; otherwise, return false.
	return LiftedTrue
}
