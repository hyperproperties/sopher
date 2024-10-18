package language

import "github.com/hyperproperties/sopher/pkg/iterx"

// IncrementalMonitor represents an interface for tracking and evaluating the state of
// assignments over a sequence of elements, typically in the context of
// quantifiers in formal verification or monitoring systems.
type IncrementalMonitor[T any] interface {

	// Increment processes an update to the monitored sequence of elements,
	// evaluating the impact of the latest increment on the overall state.
	//
	// Parameters:
	// - assignments: A slice containing all current valuations for each element
	//   variable used in at least this quantifier and the preceding ones.
	//   This represents the complete state of element assignments at the current
	//   step.
	//
	// - slice: A slice representing the sequence of elements being monitored,
	//   excluding the most recently added elements (which will be added after the
	//   increment operation). It effectively captures the elements before the
	//   most recent addition in the current sequence.
	//
	// - added: An integer specifying the number of new elements added to the
	//   sequence `slice` from right to left. These are the elements not present
	//   in the previous step but introduced during this increment.
	//
	// Returns:
	// - LiftedBoolean: The result of evaluating the updated state based on the
	//   latest assignments and the incremented sequence of elements. This value
	//   reflects the outcome of the monitoring logic after the increment.
	Increment(assignments []T, slice []T, added int) LiftedBoolean
}

// Compile time checking interface implementations:
var (
	_ IncrementalMonitor[any] = (*PredicateMonitor[any])(nil)
	_ IncrementalMonitor[any] = (*UniversalMonitor[any])(nil)
	_ IncrementalMonitor[any] = (*ExistentialMonitor[any])(nil)
)

func IncrementalMonitorFromAST[T any](node Node) IncrementalMonitor[T] {
	var recurse func(node Node, variables int) IncrementalMonitor[T]
	recurse = func(node Node, variables int) IncrementalMonitor[T] {
		switch cast := node.(type) {
		case Universal:
			body := recurse(cast.assertion, variables+len(cast.variables))
			monitor := NewUniversalMonitor(variables, len(cast.variables), body)
			return &monitor
		case Existential:
			body := recurse(cast.assertion, variables+len(cast.variables))
			monitor := NewExistentialMonitor(variables, len(cast.variables), body)
			return &monitor
		case PredicateExpression[T]:
			monitor := NewPredicateMonitor(cast.predicate)
			return &monitor
		}
		panic("unknown or unsupported AST node for the incremental monitor")
	}
	return recurse(node, 0)
}

type PredicateMonitor[T any] struct {
	predicate func(assignments []T) bool
}

func NewPredicateMonitor[T any](predicate func(assignments []T) bool) PredicateMonitor[T] {
	return PredicateMonitor[T]{
		predicate: predicate,
	}
}

func (monitor *PredicateMonitor[T]) Increment(assignments []T, _ []T, _ int) LiftedBoolean {
	return LiftBoolean(monitor.predicate(assignments))
}

type UniversalMonitor[T any] struct {
	offset, size int
	body         IncrementalMonitor[T]
	result       LiftedBoolean
}

func NewUniversalMonitor[T any](offset, size int, body IncrementalMonitor[T]) UniversalMonitor[T] {
	return UniversalMonitor[T]{
		offset: offset,
		size:   size,
		body:   body,
		result: LiftedTrue,
	}
}

func (monitor *UniversalMonitor[T]) Increment(assignments []T, slice []T, added int) LiftedBoolean {
	iterator := iterx.Permutations(monitor.size, len(slice))
	for permutation := range iterx.Map(slice, iterator) {
		for idx, element := range permutation {
			assignments[monitor.offset+idx] = element
		}

		if monitor.body.Increment(assignments, slice, added).IsFalse() {
			return LiftedFalse
		}
	}

	return LiftedTrue
}

type ExistentialMonitor[T any] struct {
	offset, size int
	body         IncrementalMonitor[T]
}

func NewExistentialMonitor[T any](offset, size int, body IncrementalMonitor[T]) ExistentialMonitor[T] {
	return ExistentialMonitor[T]{
		offset: offset,
		size:   size,
		body:   body,
	}
}

func (monitor *ExistentialMonitor[T]) Increment(assignments []T, slice []T, added int) LiftedBoolean {
	iterator := iterx.Permutations(monitor.size, len(slice))
	for permutation := range iterx.Map(slice, iterator) {
		for idx, element := range permutation {
			assignments[monitor.offset+idx] = element
		}

		if monitor.body.Increment(assignments, slice, added).IsTrue() {
			return LiftedTrue
		}
	}

	return LiftedFalse
}
