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
	assert.Equal(t, uint64(934562412326081415), a)
	assert.Equal(t, uint64(98071240265192299), b)
}
