package language

import "iter"

type RandomExplorer[T any] struct {
	generator iter.Seq[T]
	assignments []T
}

func NewRandomExplorer[T any](generator iter.Seq[T]) RandomExplorer[T] {
	return RandomExplorer[T]{
		generator: generator,
	}
}

func (explorer *RandomExplorer[T]) UniversalHyperAssertion(assertion UniversalHyperAssertion[T]) {
	explorer.assignments = append(explorer.assignments, make([]T, assertion.size)...)
}

func (explorer *RandomExplorer[T]) ExistentialHyperAssertion(assertion ExistentialHyperAssertion[T]) {
	explorer.assignments = append(explorer.assignments, make([]T, assertion.size)...)
}

func (explorer *RandomExplorer[T]) Explore(assertion PredicateHyperAssertion[T]) iter.Seq[LiftedBoolean]  {
	return func(yield func(LiftedBoolean) bool) {
		pull, stop := iter.Pull(explorer.generator)
		defer stop()

		for {
			for idx := range explorer.assignments {
				value, ok := pull()
				if !ok {
					return
				}
				explorer.assignments[idx] = value
			}
	
			result := assertion.predicate(explorer.assignments)
			if !yield(LiftBoolean(result)) {
				return
			}
		}
	}
}