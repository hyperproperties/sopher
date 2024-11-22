package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreterV2(t *testing.T) {
	interpreter := NewInterpreterV2(NewIncrementalExplorerV2([]int{}, []int{1}))
	assertion := UniversalHyperAssertion[int]{
		size: 2,
		body: BinaryHyperAssertion[int]{
			lhs: ExistentialHyperAssertion[int]{
				size: 1,
				body: TrueHyperAssertion[int]{},
			},
			operator: LogicalConjunction,
			rhs:      TrueHyperAssertion[int]{},
		},
	}
	satisfied := interpreter.Satisfies(assertion)
	assert.Equal(t, LiftedTrue, satisfied)
}
