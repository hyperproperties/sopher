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

func TestLiftedBooleanString(t *testing.T) {
	assert.Equal(t, "true", LiftedTrue.String())
	assert.Equal(t, "false", LiftedFalse.String())
	assert.Equal(t, "unknown", LiftedUnknown.String())
}

func TestLiftedBooleanIs(t *testing.T) {
	assert.True(t, LiftedTrue.IsTrue())
	assert.False(t, LiftedTrue.IsFalse())
	assert.False(t, LiftedTrue.IsUnknown())
	assert.False(t, LiftedFalse.IsTrue())
	assert.True(t, LiftedFalse.IsFalse())
	assert.False(t, LiftedFalse.IsUnknown())
	assert.False(t, LiftedUnknown.IsTrue())
	assert.False(t, LiftedUnknown.IsFalse())
	assert.True(t, LiftedUnknown.IsUnknown())
}

func TestLiftedBooleanNot(t *testing.T) {
	assert.Equal(t, LiftedTrue, LiftedFalse.Not(), "not false")
	assert.Equal(t, LiftedUnknown, LiftedUnknown.Not(), "not unknown")
	assert.Equal(t, LiftedFalse, LiftedTrue.Not(), "not true")
}

func TestLiftedBooleanlIf(t *testing.T) {
	assert.Equal(t, LiftedTrue, LiftedTrue.If(LiftedTrue))
	assert.Equal(t, LiftedFalse, LiftedTrue.If(LiftedFalse))
	assert.Equal(t, LiftedUnknown, LiftedTrue.If(LiftedUnknown))
	assert.Equal(t, LiftedTrue, LiftedUnknown.If(LiftedTrue))
	assert.Equal(t, LiftedUnknown, LiftedUnknown.If(LiftedFalse))
	assert.Equal(t, LiftedUnknown, LiftedUnknown.If(LiftedUnknown))
	assert.Equal(t, LiftedTrue, LiftedFalse.If(LiftedTrue))
	assert.Equal(t, LiftedTrue, LiftedFalse.If(LiftedFalse))
	assert.Equal(t, LiftedTrue, LiftedFalse.If(LiftedUnknown))
}

func TestLiftedBooleanIff(t *testing.T) {
	assert.Equal(t, LiftedTrue, LiftedTrue.Iff(LiftedTrue))
	assert.Equal(t, LiftedFalse, LiftedTrue.Iff(LiftedFalse))
	assert.Equal(t, LiftedUnknown, LiftedTrue.Iff(LiftedUnknown))
	assert.Equal(t, LiftedUnknown, LiftedUnknown.Iff(LiftedTrue))
	assert.Equal(t, LiftedUnknown, LiftedUnknown.Iff(LiftedFalse))
	assert.Equal(t, LiftedUnknown, LiftedUnknown.Iff(LiftedUnknown))
	assert.Equal(t, LiftedFalse, LiftedFalse.Iff(LiftedTrue))
	assert.Equal(t, LiftedTrue, LiftedFalse.Iff(LiftedFalse))
	assert.Equal(t, LiftedUnknown, LiftedFalse.Iff(LiftedUnknown))
}

func TestLiftBoolean(t *testing.T) {
	assert.Equal(t, LiftedTrue, LiftBoolean(true))
	assert.Equal(t, LiftedFalse, LiftBoolean(false))
}
