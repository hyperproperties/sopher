package language

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLiftedBooleanOr(t *testing.T) {
	values := [3]LiftedBoolean{
		LiftedFalse, LiftedTrue, LiftedUnknown,
	}

	for _, lhs := range values {
		for _, rhs := range values {
			t.Run(fmt.Sprintf("%s || %s", lhs, rhs), func(t *testing.T) {
				resultant := lhs.Or(rhs)
				if lhs == LiftedFalse && rhs == LiftedFalse {
					assert.Equal(t, LiftedFalse, resultant)
				} else if lhs == LiftedTrue || rhs == LiftedTrue {
					assert.Equal(t, LiftedTrue, resultant)
				} else {
					assert.Equal(t, LiftedUnknown, resultant)
				}
			})
		}
	}
}

func TestLiftedBooleanAnd(t *testing.T) {
	values := [3]LiftedBoolean{
		LiftedFalse, LiftedTrue, LiftedUnknown,
	}

	for _, lhs := range values {
		for _, rhs := range values {
			t.Run(fmt.Sprintf("%s && %s", lhs, rhs), func(t *testing.T) {
				resultant := lhs.And(rhs)
				if lhs == LiftedFalse || rhs == LiftedFalse {
					assert.Equal(t, LiftedFalse, resultant)
				} else if lhs == LiftedTrue && rhs == LiftedTrue {
					assert.Equal(t, LiftedTrue, resultant)
				} else {
					assert.Equal(t, LiftedUnknown, resultant)
				}
			})
		}
	}
}

func TestFunctionName(t *testing.T) {
	assert.Equal(t, LiftedTrue, LiftedFalse.Not(), "not false")
	assert.Equal(t, LiftedUnknown, LiftedUnknown.Not(), "not unknown")
	assert.Equal(t, LiftedFalse, LiftedTrue.Not(), "not true")
}
