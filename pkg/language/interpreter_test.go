package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpreterV2(t *testing.T) {
	domain := NewIncrementalDomain([]int{}, []int{1})
	explorer := NewIncrementalExplorer(&domain)
	interpreter := NewInterpreter(&explorer)
	assertion := ExistentialHyperAssertion[int]{
		size: 1,
		body: UnaryHyperAssertion[int]{
			operator: LogicalNegation,
			operand:  TrueHyperAssertion[int]{},
		},
	}
	satisfied := interpreter.Satisfies(assertion)
	assert.Equal(t, LiftedFalse, satisfied)
}
