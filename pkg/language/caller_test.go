package language

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaller(t *testing.T) {
	foo := func() uint64 {
		return Caller()
	}
	// assert.NotEqual(t, foo(), foo()) Very difficult to get working.
	a := foo()
	b := foo()
	assert.NotEqual(t, a, b)
}
