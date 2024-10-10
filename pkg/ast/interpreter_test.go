package ast

import (
	"math/rand"
	"testing"

	"github.com/hyperproperties/sopher/pkg/iterx"
	"github.com/stretchr/testify/assert"
)

func TestNonInterference(t *testing.T) {
	// ∀ e1, e2. (e1.low == e2.low ∧ ¬(e1.high = e2.high)) -> (e1.ret0 = e2.ret0)
	// or equivalently
	// ∀ e1, e2. ¬(e1.low = e2.low ∧ ¬(e1.high = e2.high)) ∨ (e1.ret0 = e2.ret0)
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

	// Construct the hyperproperty.
	predicate := NewPredicateExpression(
		func(assignments []Execution) bool {
			e1, e2 := assignments[0], assignments[1]
			return !(e1.low == e2.low && e1.high != e2.high) || (e1.ret0 == e2.ret0)
		},
	)
	universal := NewUniversalQuantifier(0, 2, predicate)

	// Collect executions to check sat for hyperproperty.
	executions := iterx.Collect(executor)
	interpreter := NewInterpreter(2, executions)

	// Check hyperproperty sat.
	satisfied, counter := interpreter.Model(universal)
	assert.True(t, satisfied, "Expected non-interference to be satisfied")
	assert.Nil(t, counter, "Expected non-interference to be satisfied and not have a counter example")
}

func TestGeneralisedNonInterference(t *testing.T) {
	// ∀ e1, e2.∃ e3. ¬(e3.low == e2.low ∧ e3.low == e1.low) ∧
	//						e3.high == e1.high ∧ e3.ret0 == e2.ret0
	// "!(e3.low == e2.low ∧ e3.low == e1.low)" is added to say that e3
	//		is not the same as e1 and e2 in terms of low inputs.
	//		Oterwise, if e1 == e2 == e3 then GNI is trivially satisfied.
	GNI := func(low, high bool) bool {
		if high {
			return !low
		}
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

	// Construct the hyperproperty.
	predicate := NewPredicateExpression(
		func(assignments []Execution) bool {
			e1, e2, e3 := assignments[0], assignments[1], assignments[2]
			return !(e3.low == e1.low && e3.low == e2.low) &&
				e3.high == e1.high && e3.ret0 == e2.ret0
		},
	)
	existential := NewExistentialQuantifier(2, 1, predicate)
	universal := NewUniversalQuantifier(0, 2, existential)

	// Collect executions to check sat for hyperproperty.
	executions := iterx.Collect(executor)
	interpreter := NewInterpreter(3, executions)

	// Check hyperproperty sat.
	satisfied, counter := interpreter.Model(universal)
	assert.False(t, satisfied, "Expected generalised non-interference to not be satisfied")
	assert.NotNil(t, counter, "Expected generalised non-interference to not be satisfied and have a counter example")
}

func TestProbabilisticObservationalDeterminism(t *testing.T) {
	// P_{π, π'}(π =^L_{in} π' | π =^L_{out} π') > 0.8
	POBS := func(input int) int {
		if rand.Float32() < 0.1 {
			return -input
		}
		return input
	}

	type Execution struct {
		input int
		ret0  int
	}

	executor := func(yield func(Execution) bool) {
		for {
			input := rand.Intn(100)
			ret0 := POBS(input)
			execution := Execution{input, ret0}

			if !yield(execution) {
				return
			}
		}
	}

	// Construct the hyperproperty.
	event := NewPredicateExpression(
		func(assignments []Execution) bool {
			e1, e2 := assignments[0], assignments[1]
			return e1.ret0 == e2.ret0
		},
	)
	given := NewPredicateExpression(
		func(assignments []Execution) bool {
			e1, e2 := assignments[0], assignments[1]
			return e1.input == e2.input
		},
	)
	probability := NewConditionalProbabilityQuantifier(0, 2, event, given)
	inequality := GreaterThanOrEqual(probability, Number(0.8))

	// Collect executions to check sat for hyperproprety.
	executions := iterx.CollectN(executor, 25_000)
	interpreter := NewInterpreter(2, executions)

	// Check hyperproperty sat.
	satisfied, counter := interpreter.Model(inequality)
	assert.True(t, satisfied, "Expected probabilistic observational determinsm to be satisfied")
	assert.Nil(t, counter, "Expected probabilistic observational determinsm to be satisfied and not have a counter example")
}