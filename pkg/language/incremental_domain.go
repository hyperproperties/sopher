package language

import (
	"iter"
	"slices"

	"github.com/hyperproperties/sopher/pkg/iterx"
)

type IncrementalDomain[T any] struct {
	set       []T
	increment []T
}

func NewIncrementalDomain[T any](set []T, increment []T) IncrementalDomain[T] {
	return IncrementalDomain[T]{
		set:       set,
		increment: increment,
	}
}

// Pushes the executions to the increment.
func (domain *IncrementalDomain[T]) Increment(executions ...T) int {
	if len(executions) == 0 {
		panic("cannot increment with no executions")
	}
	index := len(domain.set) + len(domain.increment)
	domain.increment = append(domain.increment, executions...)
	return index
}

// Pops the executions from the increment.
func (domain *IncrementalDomain[T]) Decrement(amount int) {
	if amount < 0 {
		panic("cannot decrement by a negative number")
	}
	if amount == 0 {
		return
	}

	length := len(domain.increment)
	domain.increment = slices.Delete(domain.increment, length-amount, length)
}

// Finishes the increment and moves the executions to the model.
func (domain *IncrementalDomain[T]) Commit() {
	domain.set = append(domain.set, domain.increment...)
	domain.increment = nil
}

// Finishes the increment by deleting the increment.
func (domain *IncrementalDomain[T]) Rollback() {
	domain.increment = nil
}

// Updates an entry in the model. This does not work for
// elements in the increment which has not been committed yet.
func (domain *IncrementalDomain[T]) Update(index int, value T) {
	domain.set[index] = value
}

// Returns the model of the incremental domain which is the set of executions
// that should not solely be in an exploration permutation.
func (domain *IncrementalDomain[T]) Model() []T {
	return domain.set
}

// Returns the increment of the incremental domain which is the set where
// for every tested assertion every assignment must have had atleast one assignment
// to an element in the increment.
func (domain *IncrementalDomain[T]) Incremental() []T {
	return domain.increment
}

// The number of elements in the current increment.
func (domain *IncrementalDomain[T]) IncrementLength() int {
	return len(domain.increment)
}

// Returns an iterator over the incremental permutations guranteed to include atleast one
// element from the domain's increment. This is used by universal quantifiers.
func (domain *IncrementalDomain[T]) IncrementalPermutations(sub int) iter.Seq[[]T] {
	return iterx.Map(
		append(domain.set, domain.increment...), iterx.IncrementalPermutations(
			sub, len(domain.set)+len(domain.increment), len(domain.increment),
		),
	)
}

// Returns an iterators over all permutations including the increment.
func (domain *IncrementalDomain[T]) Permutations(sub int) iter.Seq[[]T] {
	return iterx.Map(
		append(domain.set, domain.increment...), iterx.Permutations(
			sub, len(domain.set)+len(domain.increment),
		),
	)
}
