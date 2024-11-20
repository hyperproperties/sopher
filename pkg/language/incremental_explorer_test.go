package language

import (
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewIncrementalExplorer(t *testing.T) {
	assert.NotPanics(t, func() { NewIncrementalExplorer[bool](nil, nil) })
	assert.NotPanics(t, func() { NewIncrementalExplorer[bool](nil, []bool{}) })
	assert.NotPanics(t, func() { NewIncrementalExplorer[bool]([]bool{}, nil) })
	assert.NotPanics(t, func() { NewIncrementalExplorer[bool]([]bool{}, []bool{}) })
	assert.NotPanics(t, func() { NewIncrementalExplorer[bool]([]bool{false}, []bool{}) })
	assert.NotPanics(t, func() { NewIncrementalExplorer[bool]([]bool{}, []bool{true}) })
}

func TestIncrementalExplorerIncrement(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		index       int
	}{
		{
			description: "No initial model and increment",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			index:       0,
		},
		{
			description: "Model with one element and no initial increment",
			explorer:    NewIncrementalExplorer[int]([]int{0}, nil),
			index:       1,
		},
		{
			description: "No initial model but an initial increment with one element",
			explorer:    NewIncrementalExplorer[int](nil, []int{0}),
			index:       1,
		},
		{
			description: "Initial model and increment has two elements",
			explorer:    NewIncrementalExplorer[int]([]int{212, -2}, []int{0, 2}),
			index:       4,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			// The size of the increment does not matter to the index.
			// The index should be for the first element.
			index := tt.explorer.Increment(rand.Int(), rand.Int(), rand.Int())
			assert.Equal(t, tt.index, index)
		})
	}
}

func TestIncrementalExplorerEmptyIncrement(t *testing.T) {
	assert.Panics(t, func() {
		explorer := NewIncrementalExplorer[int](nil, nil)
		explorer.Increment()
	})
}

func TestIncrementalExplorerDecrement(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		amount      int
		length      int
	}{
		{
			description: "Decrement by zero in an empty increment",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			amount:      0,
			length:      0,
		},
		{
			description: "Decrement by one which empties the whole increment",
			explorer:    NewIncrementalExplorer[int](nil, []int{1}),
			amount:      1,
			length:      0,
		},
		{
			description: "Decrement by two leaving one element in the increment",
			explorer:    NewIncrementalExplorer[int](nil, []int{1, 2, 3}),
			amount:      2,
			length:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.explorer.Decrement(tt.amount)
			assert.Len(t, tt.explorer.increment, tt.length)
		})
	}
}

func TestIncrementalExplorerDecrementPanics(t *testing.T) {
	assert.Panics(t, func() {
		explorer := NewIncrementalExplorer[int](nil, nil)
		explorer.Decrement(-1)
	})
	assert.Panics(t, func() {
		explorer := NewIncrementalExplorer[int](nil, []int{})
		explorer.Decrement(1)
	})
}

func TestIncrementalExplorerCommit(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		length      int
	}{
		{
			description: "Committing to a nil model and increment",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			length:      0,
		},
		{
			description: "Committing a nil increment to a model",
			explorer:    NewIncrementalExplorer[int]([]int{1}, nil),
			length:      1,
		},
		{
			description: "Comitting an increment to a nil model",
			explorer:    NewIncrementalExplorer[int](nil, []int{1}),
			length:      1,
		},
		{
			description: "Committing an increment to a model",
			explorer:    NewIncrementalExplorer[int]([]int{0}, []int{1, 2}),
			length:      3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.explorer.Commit()
			assert.Len(t, tt.explorer.increment, 0)
			assert.Len(t, tt.explorer.model, tt.length)
		})
	}
}

func TestIncrementalExplorerRollback(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		length      int
	}{
		{
			description: "Rollback a nil increment",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			length:      0,
		},
		{
			description: "Rollback a nil increment but with a model",
			explorer:    NewIncrementalExplorer[int]([]int{1}, nil),
			length:      1,
		},
		{
			description: "Rollback a nil increment",
			explorer:    NewIncrementalExplorer[int](nil, []int{1}),
			length:      0,
		},
		{
			description: "Rollback an increment to a model",
			explorer:    NewIncrementalExplorer[int]([]int{0}, []int{1, 2}),
			length:      1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.explorer.Rollback()
			assert.Len(t, tt.explorer.increment, 0)
			assert.Len(t, tt.explorer.model, tt.length)
		})
	}
}

func TestIncrementalExplorerUpdate(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		index       int
		element     int
	}{
		{
			description: "",
			explorer:    NewIncrementalExplorer[int]([]int{0}, nil),
			index:       0,
			element:     1,
		},
		{
			description: "",
			explorer:    NewIncrementalExplorer[int]([]int{0, -1}, nil),
			index:       1,
			element:     1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			tt.explorer.Update(tt.index, tt.element)
			assert.Equal(t, tt.explorer.model[tt.index], tt.element)
		})
	}
}

func TestIncrementalExplorerUpdateOutOfBounds(t *testing.T) {
	assert.Panics(t, func() {
		explorer := NewIncrementalExplorer[int]([]int{0}, nil)
		explorer.Update(-1, 0)
	})
	assert.Panics(t, func() {
		explorer := NewIncrementalExplorer[int]([]int{0}, nil)
		explorer.Update(1, 0)
	})
	assert.Panics(t, func() {
		explorer := NewIncrementalExplorer[int](nil, nil)
		explorer.Update(0, 0)
	})
}

func TestIncrementalExplorerIncrementLength(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		increment   int
	}{
		{
			description: "Empty increment",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			increment:   0,
		},
		{
			description: "Empty increment",
			explorer:    NewIncrementalExplorer[int](nil, []int{}),
			increment:   0,
		},
		{
			description: "Empty increment",
			explorer:    NewIncrementalExplorer[int](nil, []int{1}),
			increment:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			assert.Equal(t, tt.increment, len(tt.explorer.increment))
			assert.Len(t, tt.explorer.increment, tt.increment)
		})
	}
}

func TestIncrementalExplorerExplore(t *testing.T) {
	tests := []struct {
		description string
		explorer    IncrementalExplorer[int]
		scope       Scope
		assertion   PredicateHyperAssertion[int]
		satisfied   LiftedBoolean
	}{
		{
			description: "No model or increment but the assertion returns constant true",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			scope:       NewScope(),
			assertion: NewPredicateHyperAssertion(func(_ []int) bool {
				return true
			}),
			satisfied: LiftedTrue,
		},
		{
			description: "No model or increment but the assertion returns constant false",
			explorer:    NewIncrementalExplorer[int](nil, nil),
			scope:       NewScope(),
			assertion: NewPredicateHyperAssertion(func(_ []int) bool {
				return false
			}),
			satisfied: LiftedFalse,
		},
		{
			description: "Has a model but no increment so the assertion is not tested",
			explorer:    NewIncrementalExplorer[int]([]int{0}, nil),
			scope:       NewScope(NewUniversalQuantifierScope(1)),
			assertion: NewPredicateHyperAssertion(func(assignments []int) bool {
				// forall e0. e0 > 0
				return assignments[0] > 0
			}),
			satisfied: LiftedTrue,
		},
		{
			description: "Has an increment but it does not satisfy the assertion",
			explorer:    NewIncrementalExplorer[int](nil, []int{0}),
			scope:       NewScope(NewUniversalQuantifierScope(1)),
			assertion: NewPredicateHyperAssertion(func(assignments []int) bool {
				// forall e0. e0 > 0
				return assignments[0] > 0
			}),
			satisfied: LiftedFalse,
		},
		{
			description: "Has an increment and it does satisfy the assertion",
			explorer:    NewIncrementalExplorer[int](nil, []int{1}),
			scope:       NewScope(NewUniversalQuantifierScope(1)),
			assertion: NewPredicateHyperAssertion(func(assignments []int) bool {
				// forall e0. e0 > 0
				return assignments[0] > 0
			}),
			satisfied: LiftedTrue,
		},
		{
			description: "Existential quantifiers also considers the model and not jsut the increment thereby satisfying the assertion",
			explorer:    NewIncrementalExplorer[int]([]int{0}, nil),
			scope:       NewScope(NewExistentialQuantifierScope(1)),
			assertion: NewPredicateHyperAssertion(func(assignments []int) bool {
				// exists e0. eo == 0
				return assignments[0] == 0
			}),
			satisfied: LiftedTrue,
		},
		{
			description: "Universal quantifiers does not consider the model but since the increment is empty then we assume it to be satisfied",
			explorer:    NewIncrementalExplorer[int]([]int{0}, nil),
			scope:       NewScope(NewUniversalQuantifierScope(1)),
			assertion: NewPredicateHyperAssertion(func(assignments []int) bool {
				// forall e0. eo == 0
				return assignments[0] == 0
			}),
			satisfied: LiftedTrue,
		},
		{
			description: "Universal quantifiers considers solely the increment which does not satisfy the assertion even though the model would alone",
			explorer:    NewIncrementalExplorer[int]([]int{0}, []int{-1}),
			scope:       NewScope(NewUniversalQuantifierScope(1)),
			assertion: NewPredicateHyperAssertion(func(assignments []int) bool {
				// forall e0. eo == 0
				return assignments[0] == 0
			}),
			satisfied: LiftedFalse,
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			satisfied := tt.explorer.Explore(tt.scope, tt.assertion)
			assert.Equal(t, tt.satisfied, satisfied)
		})
	}
}
