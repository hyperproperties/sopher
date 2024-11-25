package language

var _ Explorer[any] = (*IncrementalExplorer[any])(nil)

type IncrementalExplorer[T any] struct {
	domain *IncrementalDomain[T]
}

func NewIncrementalExplorer[T any](domain *IncrementalDomain[T]) IncrementalExplorer[T] {
	return IncrementalExplorer[T]{
		domain: domain,
	}
}

// Checks if the assertion is satisfied given a fixed model where a prefix of it has already been tested.
// It works by assigning a incremental permutation of the model to the universally quantified variables.
// Then it attempts to recursively find a solution to the existentially quantified variables by
// using permutations (not incrementally) of the mode. If no solution for the existentially quantified
// variables is found then a counter-example for the universal quantifiers and we return that. Otherwise, true.
func (explorer *IncrementalExplorer[T]) Explore(
	scope Scope, interpreter *Interpreter[T], assertion HyperAssertion[T],
) LiftedBoolean {
	from, _ := interpreter.assignments.Expand(scope.Size())
	defer interpreter.assignments.Shrink(scope.Size())

	// If there are no increment then we dont test it.  However, we know
	// per the invariant that the model must be a set which satifies the
	// assertion and therefore we return true.
	if explorer.domain.IncrementLength() == 0 {
		return LiftedTrue
	}

	var Recursive func(depth, offset int) LiftedBoolean
	Recursive = func(depth, offset int) LiftedBoolean {
		// We have reached the depth of the scope so we check the assertion.
		if depth >= scope.Depth() {
			return interpreter.Satisfies(assertion)
		}

		// Generate all assignments for the existential quantifier's variables.
		quantifier := scope.quantifiers[depth]
		if quantifier.quantification != ExistentialQuantification {
			return Recursive(depth+1, offset+quantifier.size)
		}

		// For every permutation is tried and tested against the assertion.
		for permutation := range explorer.domain.Permutations(quantifier.Size()) {
			interpreter.assignments.Assign(from+offset, permutation...)

			// Get the result and check for short circuit.
			result := Recursive(depth+1, offset+quantifier.size)
			if result.IsTrue() {
				return result
			}
		}

		// No satisfied exitential example was found.
		return LiftedFalse
	}

	if scope.OnlyExistential() {
		return Recursive(0, 0)
	}

	// Generate all assignments for universal variables.
	for permutation := range explorer.domain.IncrementalPermutations(scope.UniversalSize()) {
		// For each assignment of universal variables:
		local := 0
		for offset, quantifier := range scope.Quantifiers() {
			if quantifier.quantification == UniversalQuantification {
				interpreter.assignments.Assign(from+offset, permutation[local:local+quantifier.Size()]...)
				local += quantifier.size
			}
		}

		// Check if there exists a satisfying assignment for the existential variables.
		// Use early exit if no such assignment exists.
		result := Recursive(0, from)
		if result.IsFalse() {
			return result
		}
	}

	// If all universal assignments are satisfied, return true; otherwise, return false.
	return LiftedTrue
}
