package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreterExhaustiveExplorer(t *testing.T) {
	tautology := NewPredicateHyperAssertion(func(_ []int) bool { return true })
	contradiction := NewPredicateHyperAssertion(func(_ []int) bool { return false })

	type Domain struct {
		description string
		set         []int
		satisfied   LiftedBoolean
	}

	tests := []struct {
		description string
		assertion   HyperAssertion[int]
		domains     []Domain
	}{
		{
			description: "true",
			assertion:   tautology,
			domains: []Domain{
				{
					description: "{}",
					satisfied:   LiftedTrue,
				},
			},
		},
		{
			description: "false",
			assertion:   contradiction,
			domains: []Domain{
				{
					description: "{}",
					satisfied:   LiftedFalse,
				},
			},
		},
		{
			description: "{1, 2, 3} :: forall n. n > 0",
			assertion: NewUniversalHyperAssertion(
				1, NewPredicateHyperAssertion(func(assignments []int) bool {
					n := assignments[0]
					return n > 0
				}),
			),
			domains: []Domain{
				{
					description: "{1, 2, 3}",
					set:         []int{1, 2, 3},
					satisfied:   LiftedTrue,
				},
				{
					description: "{-1}",
					set:         []int{-1},
					satisfied:   LiftedFalse,
				},
			},
		},
		{
			description: "!forall n. n > 0",
			assertion: NewUnaryHyperAssertion(
				LogicalNegation, NewUniversalHyperAssertion(
					1, NewPredicateHyperAssertion(func(assignments []int) bool {
						n := assignments[0]
						return n > 0
					}),
				),
			),
			domains: []Domain{
				{
					description: "{-1, -2, -3}",
					set:         []int{-1, -2, -3},
					satisfied:   LiftedTrue,
				},
				{
					description: "{1, 3, 2}",
					set:         []int{1, 3, 2},
					satisfied:   LiftedFalse,
				},
			},
		},
		{
			description: "exists n. n == 1",
			assertion: NewExistentialHyperAssertion(
				1, NewPredicateHyperAssertion(func(assignments []int) bool {
					n := assignments[0]
					return n == 1
				}),
			),
			domains: []Domain{
				{
					description: "{1}",
					set:         []int{1},
					satisfied:   LiftedTrue,
				},
				{
					description: "{-1}",
					set:         []int{-1},
					satisfied:   LiftedFalse,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.description, func(t *testing.T) {
			for _, domain := range tt.domains {
				t.Run(domain.description, func(t *testing.T) {
					explorer := NewExhaustiveExplorer(domain.set...)
					interpreter := NewInterpreter(&explorer)
					satisfied := interpreter.Satisfies(tt.assertion)
					assert.Equal(t, domain.satisfied, satisfied)
					assert.True(t, interpreter.stack.IsEmpty())
				})
			}
		})
	}
}
