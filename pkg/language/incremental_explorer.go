package language

import "github.com/hyperproperties/sopher/pkg/iterx"

type IncrementalExplorer[T any] struct {
	model []T
	added int
}

func NewIncrementalExplorer[T any](model []T, added int) IncrementalExplorer[T] {
	return IncrementalExplorer[T]{
		model: model,
		added: added,
	}
}

// Appends executions to the model and assumed as ground truth.
func (explorer *IncrementalExplorer[T]) Append(executions ...T) {
	explorer.model = append(explorer.model, executions...)
}

// Increments the model and marks them as added. The index of the first element is
// returned and the subsequent to to the len of the executions are the indices
// of the executions in the model
func (explorer *IncrementalExplorer[T]) Increment(executions ...T) int {
	explorer.Append(executions...)
	explorer.added += len(executions)
	return len(explorer.model) - len(executions) - 1
}

// Decrements the number of newly added executions.
func (explorer *IncrementalExplorer[T]) Decrement(amount int) {
	if amount > explorer.added {
		panic("cannot decrement by more than incremented")
	}
	explorer.added -= amount
}

// Updates an execution in the model.
func (explorer *IncrementalExplorer[T]) Update(index int, value T) {
	explorer.model[index] = value
}

// Returns the number of executions in the model which are untested.
func (explorer IncrementalExplorer[T]) Added() int {
	return explorer.added
}

// Returns the full model with tested as well as untested executions.
func (explorer IncrementalExplorer[T]) Model() []T {
	return explorer.model
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
		satisfied := assertion.predicate([]T{})
		return LiftBoolean(satisfied)
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
		existential, ok := quantifier.(ExistentialQuantifier)
		if !ok {
			return Exists(depth+1, offset+quantifier.Size())
		}

		// For every permutation is tried and tested against the assertion.
		// TODO: Convert the ranging over assignments for the quantifier to a strategy.
		for permutation := range iterx.Map(
			explorer.model, iterx.Permutations(
				existential.Size(), len(explorer.model),
			),
		) {
			// TODO: copy?
			for idx := 0; idx < existential.Size(); idx++ {
				assignments[offset+idx] = permutation[idx]
			}

			// Get the result and check for short circuit.
			result := Exists(depth+1, offset+existential.Size())
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
		explorer.model, iterx.IncrementalPermutations(
			scope.UniversalSize(), len(explorer.model), explorer.added,
		),
	) {
		//For each assignment of universal variables:
		local := 0
		for offset, quantifier := range scope.Quantifiers() {
			if cast, ok := quantifier.(UniversalQuantifier); ok {
				// TODO: copy?
				for idx := 0; idx < cast.Size(); idx++ {
					assignments[offset+idx] = permutation[local+idx]
				}
				local += cast.Size()
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
