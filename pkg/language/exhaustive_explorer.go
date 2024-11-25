package language

import (
	"iter"

	"github.com/hyperproperties/sopher/pkg/iterx"
)

type ExhaustiveExplorer[T any] struct {
	domain []T
}

func NewExhaustiveExplorer[T any](domain ...T) ExhaustiveExplorer[T] {
	return ExhaustiveExplorer[T]{
		domain: domain,
	}
}

func (explorer ExhaustiveExplorer[T]) permutations(sub int) iter.Seq[[]T] {
	return iterx.Map(explorer.domain, iterx.Permutations(sub, len(explorer.domain)))
}

func (explorer *ExhaustiveExplorer[T]) Explore(
	scope Scope, interpreter *Interpreter[T], assertion HyperAssertion[T],
) LiftedBoolean {
	if scope.Size() > 0 && len(explorer.domain) == 0 {
		panic("explorer with an empty domain cannot explore the scope")
	}

	from, _ := interpreter.assignments.Expand(scope.Size())
	defer interpreter.assignments.Shrink(scope.Size())

	var Recursive func(depth, offset int) LiftedBoolean
	Recursive = func(depth, offset int) LiftedBoolean {
		if depth >= scope.Depth() {
			return interpreter.Satisfies(assertion)
		}

		switch quantifier := scope.quantifiers[depth]; quantifier.quantification {
		// Look for a negative example in the universal quantification.
		case UniversalQuantification:
			for permutation := range explorer.permutations(quantifier.Size()) {
				interpreter.assignments.Assign(offset, permutation...)
				satisfied := Recursive(depth+1, offset+quantifier.Size())
				if !satisfied.IsTrue() {
					return LiftedFalse
				}
			}
			return LiftedTrue
		// Look for a positive example in the existential quantification.
		case ExistentialQuantification:
			for permutation := range explorer.permutations(quantifier.Size()) {
				interpreter.assignments.Assign(offset, permutation...)
				satisfied := Recursive(depth+1, offset+quantifier.Size())
				if satisfied.IsTrue() {
					return LiftedTrue
				}
			}
			return LiftedFalse
		}
		panic("unknown quantifier")
	}

	return Recursive(0, from)
}
