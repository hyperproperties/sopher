package language

import (
	"testing"

	"github.com/hyperproperties/sopher/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestMonitorNonInterference(t *testing.T) {
	NI := func(low, high bool) bool {
		return !low
	}

	type Execution struct {
		low, high bool
		ret0      bool
	}

	executor := func(yield func(Execution) bool) {
		for _, execution := range [4]Execution{
			{low: false, high: false},
			{low: false, high: true},
			{low: true, high: false},
			{low: true, high: true},
		} {
			execution.ret0 = NI(execution.low, execution.high)
			if !yield(execution) {
				return
			}
		}
	}

	predicate := NewPredicateMonitor(
		func(assignments []Execution) bool {
			e1, e2 := assignments[0], assignments[1]
			return !(e1.low == e2.low && e1.high != e2.high) || (e1.ret0 == e2.ret0)
		},
	)
	universal := NewUniversalMonitor(0, 2, predicate)

	executions := iterx.Collect(executor)
	satisfied := universal.Increment(make([]Execution, 2), executions, len(executions))

	assert.True(t, satisfied.IsTrue())
}

func TestMonitorInterference(t *testing.T) {
	NI := func(low, high bool) bool {
		if low {
			return high
		}
		return !low
	}

	type Execution struct {
		low, high bool
		ret0      bool
	}

	executor := func(yield func(Execution) bool) {
		for _, execution := range [4]Execution{
			{low: false, high: false},
			{low: false, high: true},
			{low: true, high: false},
			{low: true, high: true},
		} {
			execution.ret0 = NI(execution.low, execution.high)
			if !yield(execution) {
				return
			}
		}
	}

	predicate := NewPredicateMonitor(
		func(assignments []Execution) bool {
			e1, e2 := assignments[0], assignments[1]
			return !(e1.low == e2.low && e1.high != e2.high) || (e1.ret0 == e2.ret0)
		},
	)
	universal := NewUniversalMonitor(0, 2, predicate)

	executions := iterx.Collect(executor)
	satisfied := universal.Increment(make([]Execution, 2), executions, len(executions))

	assert.True(t, satisfied.IsFalse())
}

func TestMonitorGeneralisedNonInterference(t *testing.T) {
	GNI := func(low, high bool) bool {
		return low
	}

	type Execution struct {
		low, high bool
		ret0      bool
	}

	executor := func(yield func(Execution) bool) {
		for _, execution := range [4]Execution{
			{low: false, high: false},
			{low: false, high: true},
			{low: true, high: false},
			{low: true, high: true},
		} {
			execution.ret0 = GNI(execution.low, execution.high)
			if !yield(execution) {
				return
			}
		}
	}

	predicate := NewPredicateMonitor(
		func(assignments []Execution) bool {
			e1, e2, e3 := assignments[0], assignments[1], assignments[2]
			return e3.high == e1.high && e3.ret0 == e2.ret0
		},
	)
	existential := NewExistentialMonitor(2, 1, predicate)
	universal := NewUniversalMonitor(0, 2, existential)

	executions := iterx.Collect(executor)
	satisfied := universal.Increment(make([]Execution, 3), executions, len(executions))

	assert.True(t, satisfied.IsTrue())
}

func TestMonitorGeneralisedInterference(t *testing.T) {
	GNI := func(low, high bool) bool {
		if high && low {
			return !low
		}
		return high
	}

	type Execution struct {
		low, high bool
		ret0      bool
	}

	executor := func(yield func(Execution) bool) {
		for _, execution := range [4]Execution{
			{low: false, high: false},
			{low: false, high: true},
			{low: true, high: false},
			{low: true, high: true},
		} {
			execution.ret0 = GNI(execution.low, execution.high)
			if !yield(execution) {
				return
			}
		}
	}

	predicate := NewPredicateMonitor(
		func(assignments []Execution) bool {
			e1, e2, e3 := assignments[0], assignments[1], assignments[2]
			return e3.high == e1.high && e3.ret0 == e2.ret0
		},
	)
	existential := NewExistentialMonitor(2, 1, predicate)
	universal := NewUniversalMonitor(0, 2, existential)

	executions := iterx.Collect(executor)
	satisfied := universal.Increment(make([]Execution, 3), executions, len(executions))

	assert.True(t, satisfied.IsFalse())
}
